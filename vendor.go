package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// vendor は弁当会社を表す.
type vendor struct {
	name               string
	iconEmoji          string
	location           string
	fetch              fetch
	parseDateAndDishes func(io.Reader) (time.Time, []dish, error)
}

// menu は1食分のメニュー.
type menu struct {
	date   time.Time
	dishes []dish
}

// dish はメニューの1品目を表す.
type dish string

// menu は弁当会社のサイトから今日のメニューを取得する.
func (v *vendor) menu() (m *menu, err error) {
	body, fErr := v.fetch(v.location)
	if fErr != nil {
		err = fmt.Errorf("fetching %v: %v", v.location, fErr)
		return
	}

	defer func() {
		if cErr := body.Close(); cErr != nil && err == nil {
			err = cErr
		}
	}()

	date, ds, pErr := v.parseDateAndDishes(body)
	if pErr != nil {
		err = fmt.Errorf("parsing body: %v", pErr)
		return
	}

	return &menu{
		date:   date,
		dishes: ds,
	}, nil
}

func newTamagoya() *vendor {
	return &vendor{
		name:               "玉子屋",
		iconEmoji:          "hatching_chick",
		location:           "http://www.tamagoya.co.jp/",
		fetch:              fetchHTTPBody,
		parseDateAndDishes: parseTamagoya,
	}
}

// 2019年以前のあづま給食.
func newAzuma() *vendor {
	return &vendor{
		name:               "あづま給食",
		iconEmoji:          "bento",
		location:           "http://azuma-catering.com/lunch.php",
		fetch:              fetchHTTPBody,
		parseDateAndDishes: parseAzuma,
	}
}

// 2020年以降のあづま給食.
func newAzuma2020() *vendor {
	return &vendor{
		name:               "あづま給食",
		iconEmoji:          "bento",
		location:           "http://azuma-catering.co.jp/calendar/menu.txt",
		fetch:              fetchHTTPBody,
		parseDateAndDishes: (azuma2020Parser(time.Now)).parse,
	}
}

// azuma2020Parser は、あづま給食月間メニューから、特定の日付のメニューを返す.
type azuma2020Parser func() time.Time

func parseTamagoya(body io.Reader) (time.Time, []dish, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return time.Time{}, nil, err
	}

	mn := doc.Find("div[class=cnt] > div[class=text]").First()
	date, main, err := parseDateAndMainDish(mn.Find("h3").First().Text())
	if err != nil {
		return time.Time{}, nil, err
	}

	sides := parseSideDishes(mn.Find("p").First().Text())
	ds := []dish{main}
	return date, append(ds, sides...), nil
}

func parseDateAndMainDish(dateAndMain string) (date time.Time, main dish, err error) {
	dateAndMain = strings.ReplaceAll(dateAndMain, "\n", "")
	dateAndMain = strings.ReplaceAll(dateAndMain, " ", "")
	runes := []rune(dateAndMain)
	i := index(runes, '(')
	if i == -1 {
		err = fmt.Errorf("date not found")
		return
	}

	date, err = parseDate(runes[:i])
	if err != nil {
		return
	}

	runes = runes[i:]
	i = index(runes, ')')
	if i == -1 {
		err = fmt.Errorf("parences for day of week not found")
		return
	}

	main = dish(runes[i+1:])
	return
}

func parseSideDishes(sides string) []dish {
	ss := strings.Split(sides, "\n")
	ds := make([]dish, 0, len(ss))
	for i := 0; i < len(ss); i++ {
		ds = append(ds, dish(strings.TrimSpace(ss[i])))
	}

	return ds
}

func parseAzuma(body io.Reader) (time.Time, []dish, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return time.Time{}, nil, err
	}

	date, err := parseDate([]rune(doc.Find("h2").Text()))
	if err != nil {
		return time.Time{}, nil, err
	}

	var ds []dish
	doc.Find("ul[class=lunch_menu] > li").Each(func(_ int, sel *goquery.Selection) {
		ds = append(ds, dish(sel.Text()))
	})

	return date, ds, nil
}

func (a azuma2020Parser) parse(body io.Reader) (date time.Time, dishes []dish, err error) {
	date = a()
	// csv としてのパースはエラー. カラム数が行ごとに異なる
	r := bufio.NewReader(body)
	for {
		l, rErr := r.ReadString('\n')
		if rErr == io.EOF {
			err = fmt.Errorf("no menu for %s found", date.Format("20060102"))
			return
		}

		if rErr != nil {
			err = fmt.Errorf("reading a line: %v", rErr)
			return
		}

		ss := strings.Split(l, ",")
		if len(ss) < 3 { // date,dishes(at least one),calorie
			err = fmt.Errorf("unexpected column number in %q", l)
			return
		}

		if ss[0] == date.Format("20060102") {
			for _, s := range ss[1 : len(ss)-1] {
				dishes = append(dishes, dish(s))
			}

			return
		}
	}
}

func index(rs []rune, r rune) int {
	for i := 0; i < len(rs); i++ {
		if rs[i] == r {
			return i
		}
	}

	return -1
}

// parseDate は "2019年7月15日" の形式の日付を解釈する.
func parseDate(date []rune) (d time.Time, err error) {
	i := index(date, '年')
	if i == -1 {
		err = fmt.Errorf("year not found")
		return
	}

	year, err := strconv.ParseInt(string(date[:i]), 10, 0)
	if err != nil {
		return
	}

	date = date[i+1:]
	i = index(date, '月')
	if i == -1 {
		err = fmt.Errorf("month not found")
		return
	}

	month, err := strconv.ParseInt(string(date[:i]), 10, 0)
	if err != nil {
		return
	}

	date = date[i+1:]
	i = index(date, '日')
	if i == -1 {
		err = fmt.Errorf("day not found")
		return
	}

	day, err := strconv.ParseInt(string(date[:i]), 10, 0)
	if err != nil {
		return
	}

	d = time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0, time.Local)
	return
}

// String はデバッグ用文字列を返す.
func (m *menu) String() string {
	b := new(strings.Builder)
	for i, d := range m.dishes {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(string(d))
	}

	return fmt.Sprintf("%s: [%s]", m.date.Format("2006/01/02"), b.String())
}

func (d dish) isSouplike() bool {
	return strings.HasSuffix(string(d), "カレー") ||
		strings.HasSuffix(string(d), "シチュー") ||
		strings.HasSuffix(string(d), "麻婆豆腐")
}

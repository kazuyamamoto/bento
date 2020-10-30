package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

// slackHandler は Slack にメニューを通知する http.Handler.
type slackHandler struct {
	shouldNotify func() (bool, error) // 通知すべきかどうか
	url          string               // Slack integration の URL
	vendors      []*vendor            // メニューを得る弁当会社
}

func (s *slackHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	ok, err := s.shouldNotify()
	if err != nil {
		m := fmt.Sprintf("failed to decide whether to notify today's menu: %v", err)
		log.Println(m)
		http.Error(w, m, http.StatusInternalServerError)
		return
	}

	if ok {
		s.fetchMenu(w)
	}
}

func (s *slackHandler) fetchMenu(w http.ResponseWriter) {
	// App Engine がいつインスタンスを落とすかわからないので、
	// 確実に goroutine を終わらせる
	wg := new(sync.WaitGroup)
	defer wg.Wait()

	for _, v := range s.vendors {
		wg.Add(1)
		go func(v *vendor) {
			if err := s.notify(v); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log.Println(err)
			}

			wg.Done()
		}(v)
	}
}

func (s *slackHandler) notify(v *vendor) (err error) {
	var p *payload
	m, err := v.menu()
	if err == nil {
		p = menuPayload(v, m)
	} else {
		p = errorPayload(v, fmt.Errorf("bento: failed to fetch a menu: %v", err))
	}

	bs, err := p.marshal()
	if err != nil {
		err = fmt.Errorf("marshalling payload: %v", err)
		return
	}

	vs := url.Values{
		"payload": {string(bs)},
	}

	res, err := http.PostForm(s.url, vs)
	if err != nil {
		err = fmt.Errorf("posting payload: %v", err)
		return
	}

	defer func() {
		if errClose := res.Body.Close(); errClose != nil && err == nil {
			err = fmt.Errorf("closing response body: %v", errClose)
		}
	}()

	return
}

// payload は Slack integration に送信するペイロード.
type payload struct {
	Text      string `json:"text,omitempty"`
	Name      string `json:"username,omitempty"`
	IconEmoji string `json:"icon_emoji,omitempty"`
}

// menuPayload はメニューのペイロードを返す.
func menuPayload(v *vendor, m *menu) *payload {
	b := new(strings.Builder)
	for _, d := range m.dishes {
		b.WriteString("> - ")
		b.WriteString(string(d))
		if d.isSouplike() {
			b.WriteString(" :spoon:")
		}
		b.WriteRune('\n')
	}

	t := fmt.Sprintf(`%s
%s`,
		m.date.Format("*2006/01/02*"),
		b.String())

	return &payload{
		Text:      t,
		Name:      v.name,
		IconEmoji: v.iconEmoji,
	}
}

// errorPayload はエラーメッセージのペイロードを返す.
func errorPayload(v *vendor, err error) *payload {
	return &payload{
		Text:      err.Error(),
		Name:      v.name,
		IconEmoji: v.iconEmoji,
	}
}

func (p *payload) marshal() ([]byte, error) {
	return json.Marshal(p)
}

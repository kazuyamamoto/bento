package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

// 平日かどうかを Web API を呼んで判断する.
type weekday struct {
	location string
	fetch    fetch
	date     func() time.Time
}

func newWeekday() *weekday {
	return &weekday{
		location: "https://holidays-jp.github.io/api/v1/date.json",
		fetch:    fetchHTTPBody,
		date:     time.Now,
	}
}

func (h *weekday) isToday() (bool, error) {
	b, err := h.fetch(h.location)
	if err != nil {
		return false, fmt.Errorf("fetching holiday data from holiday API: %v", err)
	}

	defer func() {
		_ = b.Close()
	}()

	bs, err := ioutil.ReadAll(b)
	if err != nil {
		return false, fmt.Errorf("reading fetched body from holiday API: %v", err)
	}

	holidays := make(map[string]string)
	if err := json.Unmarshal(bs, &holidays); err != nil {
		return false, fmt.Errorf("unmarshalling int JSON data %q: %v", string(bs), err)
	}

	_, ok := holidays[h.date().Format("2006-01-02")]
	return !ok, nil
}

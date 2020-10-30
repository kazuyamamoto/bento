package main

import (
	"testing"
	"time"
)

func TestNewSlackPayload(t *testing.T) {
	v := &vendor{
		name:      "弁当や",
		iconEmoji: "bento",
	}

	m := &menu{
		date:   time.Date(2020, 1, 2, 3, 4, 5, 6, time.UTC),
		dishes: []dish{"うどん", "シチュー", "そば"},
	}

	p := menuPayload(v, m)
	got := p.Text
	want := `*2020/01/02*
> - うどん
> - シチュー :spoon:
> - そば
`
	if want != got {
		t.Fatalf("want\n%q\ngot\n%q", want, got)
	}
}

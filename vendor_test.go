package main

import (
	"testing"
	"time"
)

func dateFunc(date time.Time) func() time.Time {
	return func() time.Time {
		return date
	}
}

func TestMenu_String(t *testing.T) {
	sut := &menu{
		date:   time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC),
		dishes: []dish{"うどん", "そば"},
	}

	want := "2020/01/02: [うどん, そば]"
	got := sut.String()
	if want != got {
		t.Fatalf("want %s, got %s", want, got)
	}
}

func TestDish_isSouplike(t *testing.T) {
	tests := map[dish]bool{
		"カレー":     true,
		"シチュー":    true,
		"麻婆豆腐":    true,
		"カレーライス":  false,
		"タイカレー":   true,
		"ビーフシチュー": true,
		"麻婆茄子":    false,
		"豆腐":      false,
	}

	for d, want := range tests {
		t.Run(string(d), func(t *testing.T) {
			if d.isSouplike() != want {
				t.Fatalf("%v should return %v", d, want)
			}
		})
	}
}

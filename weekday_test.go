package main

import (
	"testing"
	"time"
)

func TestWeekday(t *testing.T) {
	date := time.Date(2020, 8, 1, 0, 0, 0, 0, time.Local)
	sut := &weekday{
		location: "testdata/date.json",
		fetch:    fetchFile,
		date:     dateFunc(date),
	}

	weekday, err := sut.isToday()
	if err != nil {
		t.Fatal(err)
	}

	if !weekday {
		t.Fatalf("%v should be a weekday:", date)
	}
}

func TestHoliday(t *testing.T) {
	date := time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local)
	sut := &weekday{
		location: "testdata/date.json",
		fetch:    fetchFile,
		date:     dateFunc(date),
	}

	weekday, err := sut.isToday()
	if err != nil {
		t.Fatal(err)
	}

	if weekday {
		t.Fatalf("%v should be a holiday:", date)
	}
}

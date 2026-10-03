package progress

import (
	"errors"
	"testing"
)

func TestSheetsPercent(t *testing.T) {
	Reset()
	Begin()
	Sheets(0, 0)
	if got := Get(); got.Percent != 0 || got.Status != InProgress {
		t.Fatalf("reset %+v", got)
	}
	Sheets(0, 1)
	if got := Get(); got.Done != 0 || got.Total != 1 || got.Percent != 0 {
		t.Fatalf("one sheet start %+v", got)
	}
	Sheets(1, 1)
	if got := Get(); got.Percent != 100 {
		t.Fatalf("one sheet %+v", got)
	}
	Sheets(1, 4)
	if got := Get(); got.Done != 1 || got.Total != 4 || got.Percent != 25 {
		t.Fatalf("many sheets %+v", got)
	}
	Sheets(4, 4)
	Finish()
	if got := Get(); got.Percent != 100 || got.Status != Done {
		t.Fatalf("many sheets done %+v", got)
	}
}

func TestFailKeepsTheError(t *testing.T) {
	boom := errors.New("boom")
	Reset()
	Begin()
	Fail(boom)
	if got := Get(); got.Status != Error || got.Err != boom {
		t.Fatalf("after Fail %+v", got)
	}
	Reset()
	if got := Get(); got.Err != nil || got.Status != Empty {
		t.Fatalf("after Reset %+v", got)
	}
	Begin()
	Fail(boom)
	Begin()
	if Get().Err != nil {
		t.Fatal("Begin starts a run without the old error")
	}
}

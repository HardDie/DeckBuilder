package grid

import "testing"

func TestSize(t *testing.T) {
	t.Parallel()
	cases := []struct {
		n, cols, rows int
	}{
		{1, 2, 2},
		{4, 2, 2},
		{5, 3, 2},
		{69, 10, 7},
		{70, 10, 7},
	}
	for _, tt := range cases {
		cols, rows := Size(tt.n)
		if cols != tt.cols || rows != tt.rows {
			t.Fatalf("Size(%d) = %dx%d, want %dx%d", tt.n, cols, rows, tt.cols, tt.rows)
		}
	}
}

func TestSlot(t *testing.T) {
	t.Parallel()
	col, row := Slot(0, 10)
	if col != 0 || row != 0 {
		t.Fatalf("first slot = %d,%d", col, row)
	}
	col, row = Slot(10, 10)
	if col != 0 || row != 1 {
		t.Fatalf("eleventh slot = %d,%d", col, row)
	}
	col, row = Slot(69, 10)
	if col != 9 || row != 6 {
		t.Fatalf("last 10x7 slot = %d,%d", col, row)
	}
}

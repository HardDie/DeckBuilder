package fit

import "testing"

func TestApply(t *testing.T) {
	t.Parallel()

	scale, got := Apply(20, 30, 1, Size{})
	if scale != 1 || got != (Size{Inner: 1, Width: 20, Height: 30, Set: true}) {
		t.Fatalf("plain: scale %d size %+v", scale, got)
	}

	scale, got = Apply(20, 30, 0, Size{})
	if scale != 1 || got.Width != 20 || got.Height != 30 {
		t.Fatalf("zero scale: scale %d size %+v", scale, got)
	}

	_, got = Apply(20, 30, 2, Size{})
	if got.Width != 10 || got.Height != 15 {
		t.Fatalf("divisor: %+v", got)
	}

	_, got = Apply(10, 2000, 1, Size{})
	if got.Width != 10 || got.Height != 2000 || got.Inner != 1 {
		t.Fatalf("tall face stays full size: %+v", got)
	}

	_, got = Apply(1001, 8, 1, Size{})
	if !got.Set || got.Width >= 1001 || got.Width*10 > 10_000 {
		t.Fatalf("wide face should shrink under the 10000px page limit: %+v", got)
	}

	_, got = Apply(1001, 8, 1, Size{Inner: 1, Width: 4, Height: 6, Set: true})
	if got.Width != 4 || got.Height != 6 {
		t.Fatalf("locked cell: %+v", got)
	}
}

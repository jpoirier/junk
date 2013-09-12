package clock

import "testing"

func TestClockMonotonic(t *testing.T) {
	now := Monotonic.Now()
	if now.IsZero() {
		t.Fatal("time was zero, expecting non zero")
	}
}

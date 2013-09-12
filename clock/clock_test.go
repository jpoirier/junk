package clock

import "testing"
import "time"

func TestClockMonotonic(t *testing.T) {
	now := Monotonic.Now()
	unixzero := time.Unix(0, 0)
	if now == unixzero {
		t.Fatal("time was zero, expecting non zero")
	}
}

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

func TestClockRealtime(t *testing.T) {
	now := Realtime.Now()
	unixzero := time.Unix(0, 0)
	if now == unixzero {
		t.Fatal("time was zero, expecting non zero")
	}
}

func TestClockProcess(t *testing.T) {
	now := Process.Now()
	unixzero := time.Unix(0, 0)
	if now == unixzero {
		t.Fatal("time was zero, expecting non zero")
	}
}

func TestTimerUptime(t *testing.T) {
	bt := Uptime.Elapsed()
	if bt == 0 {
		t.Fatal("time was zero, expecting non zero")
	}
	t.Log("Boottime", bt)
}

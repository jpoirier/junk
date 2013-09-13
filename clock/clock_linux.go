package clock

import (
	"syscall"
	"time"
	"unsafe"
)

const (
	// from /usr/include/linux/time.h

	CLOCK_REALTIME = iota
	CLOCK_MONOTONIC
	CLOCK_PROCESS_CPUTIME_ID
	CLOCK_THREAD_CPUTIME_ID
	CLOCK_MONOTONIC_RAW
	CLOCK_REALTIME_COARSE
	CLOCK_MONOTONIC_COARSE
	CLOCK_BOOTTIME
	CLOCK_REALTIME_ALARM
	CLOCK_BOOTTIME_ALARM
)

// Available Clocks
var (
	// System-wide clock that measures real (i.e., wall-clock) time.
	// This clock is affected by discontinuous jumps in the system time (e.g., if
	// the system administrator manually changes the clock), and by the incremental
	// adjustments performed by adjtime(3) and NTP.
	Realtime Clock = &clock{CLOCK_REALTIME}

	// Clock that cannot be set and represents monotonic time since some unspecified
	// starting point. This clock is not affected by discontinuous jumps in the
	// system time (e.g., if the system administrator manually changes the clock),
	// but is affected by the incremental adjustments performed by adjtime(3) and NTP.
	Monotonic Clock = &clock{CLOCK_MONOTONIC}

	// High-resolution per-process timer from the CPU.
	Process Clock = &clock{CLOCK_PROCESS_CPUTIME_ID}
)

// Available Timers
var (
	// The amount of elapsed time since the bootup of this host.
	Uptime Timer = &timer{clock: clock{CLOCK_BOOTTIME}}
)

type clock struct {
	clockid uintptr
}

func (c *clock) Now() time.Time {
	var ts syscall.Timespec
	syscall.Syscall(syscall.SYS_CLOCK_GETTIME, c.clockid, uintptr(unsafe.Pointer(&ts)), 0)
	sec, nsec := ts.Unix()
	return time.Unix(sec, nsec)
}

type timer struct {
	clock
}

var epoch = time.Unix(0, 0)

func (t *timer) Elapsed() time.Duration {
	return t.Now().Sub(epoch)
}

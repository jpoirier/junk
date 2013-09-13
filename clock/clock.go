// Package clock provides various Clock implementations that mesure time according to various
// cronometers available.
//
// See also https://groups.google.com/d/msg/golang-nuts/D11F4zMs-E0/Nnh6W6nN-3YJ
package clock

import "time"

// A Clock represents a cronometer which can be queried for a time.Time value.
type Clock interface {

	// Now returns the current time according to this Clock
	Now() time.Time
}

// A Timer represents an amount of time that elapsed according to a Clock.
type Timer interface {

	// Elapsed returns the amount of time that has elapsed as measured
	// by the Timers underlying clock.
	Elapsed() time.Duration

	// Clock returns the underlying clock which powers this timer.
	// Clock() Clock
}

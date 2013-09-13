// Package poller allows readiness notification
package poller

import "io"

type Pollable interface {
	io.ReadWriteCloser
	Fd() uintptr
}

type Poller interface {
	// Register
	Register(Pollable) io.ReadWriteCloser
}

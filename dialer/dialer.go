// Package dialer provides a Dialer interface and implementation which
// provide a mechanism to reused net.Conn connections. 
// 
// The semantics around reuse of a Dialer provided Conn are defined
// by the consumer of the connection.
package dialer

import "net"

// Dialer represents a type which can Dial a remote network server. Dialer 
// implementations may return existing connections if safe to do so.
type Dialer interface {

	// Dial connects to the address on the named network.
	// Implementations may return an existing Conn if one is available.
	// Multiple calls to Dial may block if the implementations define
	// global or per remote connection limits.
	Dial(network, addr string) (Conn, error)
}

// Conn extends the net.Conn interface with a method of making connections 
// available for reuse, Release.
type Conn interface {
	net.Conn
	
	// Release returns the connection to the Dialer without closing it.
	// Callers must only call Release when it is known that the previous
	// request on the Conn has been fully consumed.
	Release()
}

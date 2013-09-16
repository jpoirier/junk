// Package dialer provides a Dialer interface and implementation which
// provide a mechanism to reused net.Conn connections.
//
// The semantics around reuse of a Dialer provided Conn are defined
// by the consumer of the connection.
package dialer

import (
	"net"
	"sync"
)

// Dialer represents a type which can Dial a remote network server. Dialer
// implementations may return existing connections if safe to do so.
type Dialer interface {

	// Dial connects to the address on the named network.
	// Implementations may return an existing Conn if one is available.
	// Multiple calls to Dial may block if the implementations define
	// global or per remote connection limits.
	Dial(network, addr string) (Conn, error)

	// Shutdown shuts down the Dialer.
	Shutdown()
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

type dialer struct {
	shutdown chan struct{} // closed when dialer is closed

	sync.Mutex // protects remaining fields
	dw         map[tupple]*dialworker
}

// New returns a Dialer implementation.
func New() Dialer {
	return &dialer{
		shutdown: make(chan struct{}),
		dw:       make(map[tupple]*dialworker),
	}
}

func (d *dialer) Shutdown() {
	close(d.shutdown)
}

func (d *dialer) Dial(network, addr string) (Conn, error) {
	t := tupple{network, addr}
	dw, ok := d.dw[t]
	if !ok {
		dw = &dialworker{
			tupple:   t,
			dialer:   d,
			shutdown: make(chan struct{}),
			pool:     make(chan Conn, 8),
			dial:     make(chan chan result),
		}
		go dw.loop()
		d.dw[t] = dw
	}
	return dw.Dial()
}

// A tupple represents an endpoint that can be dialed.
type tupple struct {
	network, addr string
}

type result struct {
	Conn
	error
}

// A dialworker manages a pool of connections that are not assigned to callers.
type dialworker struct {
	tupple
	*dialer
	shutdown chan struct{}
	pool     chan Conn
	dial     chan chan result
}

func (d *dialworker) Dial() (Conn, error) {
	r := make(chan result)
	d.dial <- r
	c := <-r
	return c.Conn, c.error
}

func (d *dialworker) loop() {
	for {
		select {
		case <-d.dialer.shutdown:
			// global shutdown
		case <-d.shutdown:
			// local shutdown
		case r := <-d.dial:
			// request to dial
			select {
			case c := <-d.pool:
				r <- result{c, nil}
			default:
				c, err := net.Dial(d.network, d.addr)
				r <- result{
					Conn: &conn{
						Conn: c,
						pool: d.pool,
					},
					error: err,
				}
			}
		}
	}
}

type conn struct {
	net.Conn
	pool chan Conn
}

func (c *conn) Release() {
	select {
	case c.pool <- c:
		// released
	default:
		// discard
	}
}

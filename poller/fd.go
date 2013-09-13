package poller

import (
	"errors"
	"io"
	"sync"
)

// ReadWriteCloser implementation that supports concurrent Read/Write and Close operations.
type rwc struct {
	io.ReadWriteCloser
	sync.Mutex // protects refcount and closing
	refcount   int
	closing    bool
}

var errClosing = errors.New("closing")

func (c *rwc) incRef(closing bool) error {
	c.Lock()
	defer c.Unlock()
	if c.closing {
		return errClosing
	}
	c.refcount++
	if closing {
		c.closing = true
	}
	return nil
}

func (c *rwc) decRef() {
	c.Lock()
	defer c.Unlock()
	c.refcount--
	if c.closing && c.refcount == 0 {
		c.ReadWriteCloser.Close()
	}
}

func (c *rwc) Close() error {
	if err := c.incRef(true); err != nil {
		return err
	}
	c.decRef()
	return nil
}

func (c *rwc) Read(b []byte) (int, error) {
	if err := c.incRef(false); err != nil {
		return 0, err
	}
	defer c.decRef()
	return c.ReadWriteCloser.Read(b)
}

func (c *rwc) Write(b []byte) (int, error) {
	if err := c.incRef(false); err != nil {
		return 0, err
	}
	defer c.decRef()
	return c.ReadWriteCloser.Write(b)
}

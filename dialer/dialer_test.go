package dialer

import (
	"net"
	"testing"
)

func server(t *testing.T) (net.Addr, func()) {
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatal(err)
	}
	return l.Addr(), func() {
		if err := l.Close(); err != nil {
			t.Fatal(err)
		}
	}
}

func TestNewDialer(t *testing.T) {
	d := New().(*dialer)
	d.Shutdown()
	select {
	case _, _ = <-d.shutdown:
		// cool
	default:
		t.Fatal("d.shutdown was not closed")
	}
}

func TestDialTwice(t *testing.T) {
	_, shutdown := server(t)
	defer shutdown()

}

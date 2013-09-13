package poller

import (
	"os"
	"testing"
	"time"
)

func TestPollerLoop(t *testing.T) {
	pr, pw, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer pr.Close()
	defer pw.Close()
	p := &poller{
		pr: pr,
		pw: pw,
	}
	if err := p.loop(time.Millisecond); err != nil {
		t.Fatal(err)
	}
}

func TestPollerLoop2(t *testing.T) {
	pr, pw, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer pr.Close()
	defer pw.Close()
	p := &poller{
		pr: pr,
		pw: pw,
	}
	pw.Write([]byte{0})
	if err := p.loop(time.Second); err != nil {
		t.Fatal(err)
	}
}

func TestNewPoller(t *testing.T) {
	p, err := newPoller()
	if err != nil {
		t.Fatal(err)
	}
	if err := p.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestPollerWakeup(t *testing.T) {
	p, err := newPoller()
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()
	p.wakeup()
}

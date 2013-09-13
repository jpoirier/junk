package poller

import (
	"log"
	"os"
	"syscall"
	"time"
	"unsafe"
)

type poller struct {
	pr, pw *os.File
}

func newPoller() (*poller, error) {
	pr, pw, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	p := poller{
		pr: pr,
		pw: pw,
	}
	go p.run()
	return &p, nil
}

func (p *poller) Close() error {
	err1 := p.pr.Close()
	err2 := p.pw.Close()
	return firstErr(err1, err2)
}

func (p *poller) run() {
	for {
		if err := p.loop(time.Second); err != nil {
			log.Fatal(err)
		}
	}
}

var wakebuf [1]byte

func (p *poller) wakeup() error {
	_, err := p.pw.Write(wakebuf[:])
	return err
}

func (p *poller) loop(timeout time.Duration) error {
	var rset, wset syscall.FdSet
	var numfd int
	set(&rset, p.pr.Fd(), &numfd)
	tv := toTimeval(timeout)
	n, err := syscall.Select(numfd+1, &rset, &wset, nil, &tv)
	if err != nil {
		return err
	}
	for fd := uintptr(0); n > 0 && fd <= uintptr(numfd); fd++ {
		if isset(&rset, fd) {
			n--
			log.Println(fd, "read")
		}
		if isset(&wset, fd) {
			n--
			log.Println(fd, "write")
		}
	}
	return nil
}

func set(set *syscall.FdSet, fd uintptr, n *int) {
	width := unsafe.Sizeof(set.Bits[0])
	index := fd / width
	offset := fd % width
	set.Bits[index] |= 1 << offset
	*n = max(int(fd), *n)
}

func isset(set *syscall.FdSet, fd uintptr) bool {
	width := unsafe.Sizeof(set.Bits[0])
	index := fd / width
	offset := fd % width
	return 1<<offset&set.Bits[index] > 0
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func toTimeval(d time.Duration) syscall.Timeval { return syscall.NsecToTimeval(int64(d)) }

func firstErr(err ...error) error {
	for _, err := range err {
		if err != nil {
			return err
		}
	}
	return nil
}

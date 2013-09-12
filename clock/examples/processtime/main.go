package main

import (
	"fmt"
	"time"

	"github.com/davecheney/junk/clock"
)

func main() {
	rt, pt := clock.Realtime.Now(), clock.Process.Now()
	time.Sleep(time.Second)
	fmt.Printf("Wall clock time: %v, process CPU time: %v\n", clock.Realtime.Now().Sub(rt), clock.Process.Now().Sub(pt))
}

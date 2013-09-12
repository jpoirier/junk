package main

import (
	"fmt"
	"github.com/davecheney/junk/clock"
	"time"
)

func main() {
	t := time.NewTicker(2 * time.Second)
	rt, mt := clock.Realtime.Now(), clock.Monotonic.Now()
	var rtt, mtt time.Duration
	for {
		<-t.C
		rt1, mt1 := clock.Realtime.Now(), clock.Monotonic.Now()
		rtt += rt1.Sub(rt)
		mtt += mt1.Sub(mt)
		rt, mt = rt1, mt1
		fmt.Printf("CLOCK_REALTIME: %v\t CLOCK_MONOTONIC: %v\n", rtt, mtt)
	}
}

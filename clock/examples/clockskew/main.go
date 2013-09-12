package main

import (
	"fmt"
	"github.com/davecheney/junk/clock"
	"time"
)

func main() {
	t := time.NewTicker(2 * time.Second)
	for {
		<-t.C
		fmt.Printf("CLOCK_REALTIME: %v\t CLOCK_MONOTONIC: %v\n", clock.Realtime.Now(), clock.Monotonic.Now())
	}
}

package dynamic

import (
	"fmt"
	"time"
)

// WatchProcess creates a go routine that watches for the dynamicd process and restarts if stopped
func WatchProcess(stopchan chan struct{}, stoppedchan chan struct{}, sleepSecs uint16) {
	go func(stopchan chan struct{}, stoppedchan chan struct{}) {
		i := 1
		defer close(stoppedchan)
		//fmt.Println("WatchProcess chan", stopchan, stoppedchan)
		//defer func() {
		// tear down here
		//}()
		for {
			select {
			default:
				proc, _ := FindDynamicdProcess()
				//fmt.Println("WatchProcess FindDynamicdProcess", i)
				if proc == nil {
					fmt.Println("WatchProcess restarting dynamicd process", i)
					dynamicd, err := LoadRPCDynamicd()
					if err != nil {
						fmt.Println("WatchProcess error restarting dynamicd process", i, err)
					} else {
						fmt.Println("WatchProcess error restarting dynamicd process", i, dynamicd)
					}
					// TODO: unlock wallet if locked.
					i++
				}
				time.Sleep(time.Duration(sleepSecs) * time.Second)
			case <-stopchan:
				fmt.Println("WatchProcess stopped")
				return
			}
		}
	}(stopchan, stoppedchan)
}

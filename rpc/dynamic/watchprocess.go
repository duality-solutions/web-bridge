package dynamic

import (
	"fmt"
	"time"
)

// WatchProcess creates a go routine that watches for the dynamicd process and restarts if stopped
func WatchProcess(stopchan chan struct{}, stoppedchan chan struct{}, sleepSecs uint16, walletpassphrase string) {
	go func(stopchan chan struct{}, stoppedchan chan struct{}) {
		//i := 1
		restarts := 0
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
					restarts++
					fmt.Println("WatchProcess restarting dynamicd process", restarts)
					dynamicd, err := LoadRPCDynamicd()
					if err != nil {
						fmt.Println("WatchProcess error restarting dynamicd process", restarts, err)
						continue
					}
					time.Sleep(time.Duration(sleepSecs) * time.Second)
					// unlock wallet if locked.
					if len(walletpassphrase) > 0 {
						dynamicd.UnlockWallet(walletpassphrase)
					}
				}
				time.Sleep(time.Duration(sleepSecs) * time.Second)
				//i++
			case <-stopchan:
				fmt.Println("WatchProcess stopped")
				return
			}
		}
	}(stopchan, stoppedchan)
}

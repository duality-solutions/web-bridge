package dynamic

import (
	"time"

	"github.com/duality-solutions/web-bridge/internal/util"
)

// WatchProcess creates a go routine that watches for the dynamicd process and restarts if stopped
func WatchProcess(stopchan chan struct{}, sleepSecs uint16, walletpassphrase string) {
	go func(stopchan chan struct{}) {
		//i := 1
		restarts := 0
		//util.Info.Println("WatchProcess chan", stopchan, stoppedchan)
		//defer func() {
		// tear down here
		//}()
		for {
			select {
			default:
				_, proc, _ := FindDynamicdProcess(false, time.Second*1)
				if proc == nil {
					restarts++
					util.Info.Println("WatchProcess restarting dynamicd process", restarts)
					dynamicd, _, err := FindDynamicdProcess(true, time.Second*30)
					if err != nil {
						util.Error.Println("WatchProcess error restarting dynamicd process", restarts, err)
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
				util.Info.Println("WatchProcess stopped")
				return
			}
		}
	}(stopchan)
	util.Info.Println("WatchProcess started")
}

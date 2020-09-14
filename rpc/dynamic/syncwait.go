package dynamic

import (
	"time"

	"github.com/duality-solutions/web-bridge/internal/util"
)

// WaitForSync waits for the Dynamic blockchain to fully sync
func (d *Dynamicd) WaitForSync(stopchan *chan struct{}, checkDelaySeconds, endDelaySeconds uint16) bool {
	status, _ := d.GetSyncStatus()
	for status.SyncProgress < 1 {
		select {
		case <-time.After(time.Duration(checkDelaySeconds) * time.Second):
			status, _ = d.GetSyncStatus()
		case <-*stopchan:
			util.Info.Println("WaitForSync stopped")
			return false
		}
	}
	time.Sleep(time.Duration(endDelaySeconds) * time.Second)
	return true
}

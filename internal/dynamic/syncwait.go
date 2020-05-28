package dynamic

import (
	"time"
)

// WaitForSync waits for the Dynamic blockchain to fully sync
func WaitForSync(d *Dynamicd, checkDelaySeconds, endDelaySeconds uint16) {
	status, _ := d.GetSyncStatus()
	for status.SyncProgress < 1 {
		time.Sleep(time.Duration(checkDelaySeconds) * time.Second)
		status, _ = d.GetSyncStatus()
	}
	time.Sleep(time.Duration(endDelaySeconds) * time.Second)
}

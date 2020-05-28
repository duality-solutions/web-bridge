package dynamic

import "time"

// WaitForSync waits for the Dynamic blockchain to fully sync
func WaitForSync(d *Dynamicd) {
	status, _ := d.GetSyncStatus()
	for status.SyncProgress < 1 {
		time.Sleep(time.Second * 30)
		status, _ = d.GetSyncStatus()
	}
	time.Sleep(time.Second * 10)
}

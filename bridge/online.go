package bridge

import "sync"

// OnlineNotification stores start and end WebBridge session time
type OnlineNotification struct {
	StartTime int64 `json:"start_time"`
	EndTime   int64 `json:"end_time"`
	mut       *sync.RWMutex
}

func newOnlineNotification() OnlineNotification {
	on := OnlineNotification{
		StartTime: 0,
		EndTime:   0,
	}
	on.mut = new(sync.RWMutex)
	return on
}

func (on *OnlineNotification) updateDates(start int64, end int64) {
	on.mut.Lock()
	defer on.mut.Unlock()
	on.StartTime = start
	on.EndTime = end
}

// GetStartTime returns the online notification last epoch start time
func (on *OnlineNotification) GetStartTime() int64 {
	on.mut.RLock()
	defer on.mut.RUnlock()
	return on.StartTime
}

// GetEndTime returns the online notification last epoch end time
func (on *OnlineNotification) GetEndTime() int64 {
	on.mut.RLock()
	defer on.mut.RUnlock()
	return on.EndTime
}

func (on *OnlineNotification) getCurrentNotification() OnlineNotification {
	on.mut.RLock()
	defer on.mut.RUnlock()
	return *on
}

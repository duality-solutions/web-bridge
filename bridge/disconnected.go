package bridge

import (
	"time"

	util "github.com/duality-solutions/web-bridge/internal/utilities"
)

// StopDisconnected looks for disconnected bridges and resets them
func StopDisconnected(stopchan chan struct{}) bool {
	util.Info.Println("StopDisconnected started")
	currentEpoch := time.Now().Unix()
	for _, link := range linkBridges.connected {
		if (link.State == StateWaitForRTC || link.State == StateEstablishRTC) && link.RTCState == "failed" && (currentEpoch-link.OnStateChangeEpoch > 180) {
			util.Info.Println("StopDisconnected found", link.LinkParticipants(), link.LinkID())
			link.State = StateInit
			delete(linkBridges.connected, link.LinkID())
			linkBridges.unconnected[link.LinkID()] = link
		}
	}
	return true
}

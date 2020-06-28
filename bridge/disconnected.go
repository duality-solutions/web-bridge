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
		if link.State == StateWaitForRTC || link.State == StateEstablishRTC {
			if link.PeerConnection != nil {
				failedPeerConnection := (link.PeerConnection.ConnectionState().String() == "failed")
				failedICEConnection := (link.PeerConnection.ICEConnectionState().String() == "failed")
				if failedPeerConnection || failedICEConnection {
					if failedPeerConnection && failedICEConnection {
						util.Info.Println("StopDisconnected failed peer and ICE connections", link.LinkParticipants(), link.LinkID())
					} else if failedPeerConnection {
						util.Info.Println("StopDisconnected failed peer connection", link.LinkParticipants(), link.LinkID())
					} else if failedICEConnection {
						util.Info.Println("StopDisconnected failed ICE connection", link.LinkParticipants(), link.LinkID())
					}
					if failedICEConnection {
						util.Info.Println("StopDisconnected close peer connection", link.LinkParticipants(), link.LinkID())
						link.PeerConnection.Close()
					}
					link.State = StateInit
					delete(linkBridges.connected, link.LinkID())
					linkBridges.unconnected[link.LinkID()] = link
					continue
				}
			}
		} else if (link.State == StateWaitForRTC || link.State == StateEstablishRTC) && link.RTCState == "failed" && (currentEpoch-link.OnStateChangeEpoch > 360) {
			util.Info.Println("StopDisconnected failed state found", link.LinkParticipants(), link.LinkID())
			failedICEConnection := (link.PeerConnection.ICEConnectionState().String() == "failed")
			if failedICEConnection {
				util.Info.Println("StopDisconnected close peer connection", link.LinkParticipants(), link.LinkID())
				link.PeerConnection.Close()
			}
			link.State = StateInit
			delete(linkBridges.connected, link.LinkID())
			linkBridges.unconnected[link.LinkID()] = link
		}
	}
	return true
}

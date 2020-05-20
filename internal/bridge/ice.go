package bridge

import (
	"fmt"

	"github.com/duality-solutions/web-bridge/internal/settings"
	webrtc "github.com/pion/webrtc/v2"
)

// newIceSetting create a new WebRTC ICE Server setting
func newIceSetting(config settings.Configuration) (*webrtc.ICEServer, error) {
	if (len(config.IceServers)) == 0 {
		return nil, fmt.Errorf("No ICE service URL found")
	}
	urls := []string{config.IceServers[0].URL}
	iceSettings := webrtc.ICEServer{
		URLs:           urls,
		Username:       config.IceServers[0].UserName,
		Credential:     config.IceServers[0].Credential,
		CredentialType: webrtc.ICECredentialTypePassword,
	}
	return &iceSettings, nil
}

// ConnectToIceServices uses the configuration settings to establish a connection with ICE servers
func ConnectToIceServices(config settings.Configuration) (*webrtc.PeerConnection, error) {
	iceSettings, err := newIceSetting(config)
	if err != nil {
		return nil, fmt.Errorf("NewIceSetting", err)
	}
	configICE := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{*iceSettings},
	}

	peerConnection, err := webrtc.NewPeerConnection(configICE)
	if err != nil {
		return nil, fmt.Errorf("NewPeerConnection", err)
	}
	return peerConnection, nil
}

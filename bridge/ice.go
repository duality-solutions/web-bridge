package bridge

import (
	"fmt"

	"github.com/duality-solutions/web-bridge/init/settings"
	util "github.com/duality-solutions/web-bridge/internal/utilities"
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

func connectToIceServicesOption(config settings.Configuration, detached bool) (*webrtc.PeerConnection, error) {
	iceSettings, err := newIceSetting(config)
	if err != nil {
		return nil, fmt.Errorf("NewIceSetting %v", err)
	}
	configICE := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{*iceSettings},
	}

	s := webrtc.SettingEngine{}
	if detached {
		s.DetachDataChannels()
	}

	// Create an API object with the engine
	api := webrtc.NewAPI(webrtc.WithSettingEngine(s))
	peerConnection, err := api.NewPeerConnection(configICE)
	if err != nil {
		return nil, fmt.Errorf("NewPeerConnection %v", err)
	}
	return peerConnection, nil
}

// ConnectToIceServices uses the configuration settings to establish a connection with ICE servers
func ConnectToIceServices(config settings.Configuration) (*webrtc.PeerConnection, error) {
	return connectToIceServicesOption(config, false)
}

// ConnectToIceServicesDetached uses the configuration settings to establish a connection with ICE servers with detached channels
func ConnectToIceServicesDetached(config settings.Configuration) (*webrtc.PeerConnection, error) {
	return connectToIceServicesOption(config, true)
}

// DisconnectBridgeIceServices calls the close method for a WebRTC peer connection
func DisconnectBridgeIceServices(bridges *Bridges) error {
	for i, v := range bridges.connected {
		util.Info.Println("DisconnectBridgeIceServices", i, v)
		DisconnectIceService(v.PeerConnection)
	}
	for i, v := range bridges.unconnected {
		util.Info.Println("DisconnectBridgeIceServices", i, v)
		err := DisconnectIceService(v.PeerConnection)
		if err != nil {
			util.Error.Println("DisconnectBridgeIceServices error", i, err)
		}
	}
	return nil
}

// DisconnectIceService calls the close method for a WebRTC peer connection
func DisconnectIceService(pc *webrtc.PeerConnection) error {
	return pc.Close()
}

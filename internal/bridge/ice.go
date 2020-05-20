package bridge

import (
	"fmt"

	"github.com/duality-solutions/web-bridge/internal/settings"
	webrtc "github.com/pion/webrtc/v2"
)

// NewIceSetting create a new WebRTC ICE Server setting
func NewIceSetting(config settings.Configuration) (*webrtc.ICEServer, error) {
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

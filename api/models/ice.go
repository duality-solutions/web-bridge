package models

// IceServerConfig stores the ICE server configuration information needed for WebRTC connections
// swagger:parameters models.IceServerConfig
type IceServerConfig struct {
	// The ICE server's full URL with protocol prefix and port: turn:ice.bdap.io:3478
	URL string `json:"URL"`
	// The ICE server's user name. Leave blank if it doesn't apply
	UserName string `json:"UserName"`
	// The ICE server's credentials. Leave blank if it doesn't apply
	Credential string `json:"Credential"`
}

package models

type BridgeInfo struct {
	SessionID          uint16 `json:"session_id"`
	LinkID             string `json:"link_id"`
	State              string `json:"state"`
	MyAccount          string `json:"my_account"`
	LinkAccount        string `json:"link_account"`
	OnOpenEpoch        int64  `json:"on_open_epoch"`
	OnStateChangeEpoch int64  `json:"on_state_changed_epoch"`
	OnLastDataEpoch    int64  `json:"on_last_data_epoch"`
	OnErrorEpoch       int64  `json:"on_error_epoch"`
	RTCState           string `json:"rtc_status"`
	HTTPListenPort     uint16 `json:"http_listen_port"`
	HTTPSListenPort    uint16 `json:"https_listen_port"`
}

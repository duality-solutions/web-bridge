package dynamic

import (
	"encoding/json"
	"fmt"
)

// MessageReturnJSON stores dynamic RPC send message response
type MessageReturnJSON struct {
	TimestampEpoch int    `json:"timestamp_epoch"`
	SharedPubkey   string `json:"shared_pubkey"`
	SubjectID      string `json:"subject_id"`
	MessageID      string `json:"message_id"`
	MessageHash    string `json:"message_hash"`
	MessageSize    int    `json:"message_size"`
	SignatureSize  int    `json:"signature_size"`
	KeepLast       string `json:"keep_last"`
}

// SendLinkMessage calls the DHT put command to add an encrypted record for the given account link
func (d *Dynamicd) SendLinkMessage(sender, receiver, answer string) (*MessageReturnJSON, error) {
	var ret MessageReturnJSON
	cmd := "dynamic-cli link sendmessage " + sender + " " + receiver + " bridge " + `"` + answer + `" 1`
	req, err := NewRequest(cmd)
	if err != nil {
		return nil, err
	}
	res := <-d.ExecCmdRequest(req)
	err = json.Unmarshal([]byte(res), &ret)
	if err != nil {
		return nil, err
	}
	fmt.Println(cmd, &ret)
	return &ret, nil
}

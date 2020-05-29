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

// GetMessageReturnJSON stores dynamic RPC get message response
type GetMessageReturnJSON struct {
	Type           string `json:"type"`
	Message        string `json:"message"`
	MessageID      string `json:"message_id"`
	MessageSize    int    `json:"message_size"`
	TimestampEpoch int    `json:"timestamp_epoch"`
	RecordNum      int    `json:"record_num"`
}

// SendLinkMessage calls the DHT put command to add an encrypted record for the given account link
func (d *Dynamicd) SendLinkMessage(sender, receiver, message string) (*MessageReturnJSON, error) {
	var ret MessageReturnJSON
	cmd := "dynamic-cli link sendmessage " + sender + " " + receiver + " bridge " + `"` + message + `"`
	req, err := NewRequest(cmd)
	if err != nil {
		return nil, err
	}
	res := <-d.ExecCmdRequest(req)
	err = json.Unmarshal([]byte(res), &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

// GetLinkMessages calls the local VGP instant message queue
func (d *Dynamicd) GetLinkMessages(receiver, sender string) (*[]GetMessageReturnJSON, error) {
	var ret []GetMessageReturnJSON
	cmd := "dynamic-cli link getaccountmessages " + receiver + " " + sender + " bridge"
	req, err := NewRequest(cmd)
	if err != nil {
		return nil, err
	}
	res := <-d.ExecCmdRequest(req)
	var messagesGeneric map[string]interface{}
	err = json.Unmarshal([]byte(res), &messagesGeneric)
	if err != nil {
		fmt.Println("GetLinkMessages messagesGeneric error", err)
		return nil, err
	}
	for _, v := range messagesGeneric {
		b, err := json.Marshal(v)
		if err != nil {
			fmt.Println("GetLinkMessages json.Marshal error", err)
		} else {
			var message GetMessageReturnJSON
			err := json.Unmarshal(b, &message)
			if err != nil {
				fmt.Println("GetLinkMessages json.Unmarshal error", err)
			} else {
				ret = append(ret, message)
			}
		}
	}
	return &ret, nil
}

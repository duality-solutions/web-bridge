package dynamic

import (
	"encoding/json"

	util "github.com/duality-solutions/web-bridge/internal/utilities"
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

// GetVGPMessageReturn stores dynamic RPC get message response with the from and to
type GetVGPMessageReturn struct {
	LinkID   string
	From     string
	To       string
	Type     string
	Messages []GetMessageReturnJSON
}

// SendNotificationMessage sends an encrypted message to the the given account link using VGP IM
func (d *Dynamicd) SendNotificationMessage(sender, receiver, message string) (*MessageReturnJSON, error) {
	var ret MessageReturnJSON
	cmd := "dynamic-cli link sendmessage " + sender + " " + receiver + " online " + `"` + message + `"`
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

// GetNotificationMessages calls the local VGP instant message queue
func (d *Dynamicd) GetNotificationMessages(receiver, sender string) (*[]GetMessageReturnJSON, error) {
	var ret []GetMessageReturnJSON
	cmd := "dynamic-cli link getaccountmessages " + receiver + " " + sender + " online"
	req, err := NewRequest(cmd)
	if err != nil {
		return nil, err
	}
	res := <-d.ExecCmdRequest(req)
	var messagesGeneric map[string]interface{}
	err = json.Unmarshal([]byte(res), &messagesGeneric)
	if err != nil {
		util.Error.Println("GetLinkMessages messagesGeneric error", err)
		return nil, err
	}
	for _, v := range messagesGeneric {
		b, err := json.Marshal(v)
		if err != nil {
			util.Error.Println("GetLinkMessages json.Marshal error", err)
		} else {
			var message GetMessageReturnJSON
			err := json.Unmarshal(b, &message)
			if err != nil {
				util.Error.Println("GetLinkMessages json.Unmarshal error", err)
			} else {
				ret = append(ret, message)
			}
		}
	}
	return &ret, nil
}

// SendLinkMessage calls the VGP IM send message command to add an encrypted record for the given account link
func (d *Dynamicd) SendLinkMessage(sender, receiver, message, msgtype string) (*MessageReturnJSON, error) {
	util.Info.Println("SendLinkMessage", sender, receiver, len(message), msgtype)
	var ret MessageReturnJSON
	cmd := "dynamic-cli link sendmessage " + sender + " " + receiver + " " + msgtype + " " + `"` + message + `"`
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
func (d *Dynamicd) GetLinkMessages(receiver, sender, msgtype string) (*[]GetMessageReturnJSON, error) {
	var ret []GetMessageReturnJSON
	cmd := "dynamic-cli link getaccountmessages " + receiver + " " + sender + " " + msgtype
	req, err := NewRequest(cmd)
	if err != nil {
		return nil, err
	}
	res := <-d.ExecCmdRequest(req)
	var messagesGeneric map[string]interface{}
	err = json.Unmarshal([]byte(res), &messagesGeneric)
	if err != nil {
		util.Error.Println("GetLinkMessages messagesGeneric error", err)
		return nil, err
	}
	for _, v := range messagesGeneric {
		b, err := json.Marshal(v)
		if err != nil {
			util.Error.Println("GetLinkMessages json.Marshal error", err)
		} else {
			var message GetMessageReturnJSON
			err := json.Unmarshal(b, &message)
			if err != nil {
				util.Error.Println("GetLinkMessages json.Unmarshal error", err)
			} else {
				ret = append(ret, message)
			}
		}
	}
	return &ret, nil
}

// SendNotificationMessageAsync sends an encrypted message to the the given account link using VGP IM
func (d *Dynamicd) SendNotificationMessageAsync(sender, receiver, message string, out chan<- MessageReturnJSON) {
	go func() {
		var ret MessageReturnJSON
		cmd := "dynamic-cli link sendmessage " + sender + " " + receiver + " online " + `"` + message + `"`
		req, err := NewRequest(cmd)
		if err != nil {
			util.Error.Println("SendNotificationMessageAsync NewRequest error", err)
			return
		}
		res := <-d.ExecCmdRequest(req)
		err = json.Unmarshal([]byte(res), &ret)
		if err != nil {
			util.Error.Println("SendNotificationMessageAsync Unmarshal error", err)
			return
		}
		util.Info.Println("SendNotificationMessageAsync sent notification message from", sender, "to", receiver)
		out <- ret
	}()
}

// GetLinkMessagesAsync asynchronously executes the getaccountmessages RPC command and returns the results to a channel
func (d *Dynamicd) GetLinkMessagesAsync(id, receiver, sender, msgtype string, out chan<- GetVGPMessageReturn) {
	go func() {
		var ret GetVGPMessageReturn
		ret.LinkID = id
		ret.To = receiver
		ret.From = sender
		ret.Type = msgtype
		cmd := "dynamic-cli link getaccountmessages " + receiver + " " + sender + " " + msgtype
		req, err := NewRequest(cmd)
		if err != nil {
			util.Error.Println("GetLinkMessagesAsync NewRequest error", err)
			return
		}
		res := <-d.ExecCmdRequest(req)
		var messagesGeneric map[string]interface{}
		err = json.Unmarshal([]byte(res), &messagesGeneric)
		if err != nil {
			util.Error.Println("GetLinkMessagesAsync Unmarshal error", err)
			return
		}
		for _, v := range messagesGeneric {
			b, err := json.Marshal(v)
			if err != nil {
				util.Error.Println("GetLinkMessages json.Marshal error", err)
			} else {
				var message GetMessageReturnJSON
				err := json.Unmarshal(b, &message)
				if err != nil {
					util.Error.Println("GetLinkMessages json.Unmarshal error", err)
				} else {
					ret.Messages = append(ret.Messages, message)
				}
			}
		}
		out <- ret
	}()
}

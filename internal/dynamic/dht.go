package dynamic

import (
	"encoding/json"
	"fmt"
)

// DHTPutReturn used to store information returned by dht putlinkrecord dev test01 pshare "<offer>"
type DHTPutReturn struct {
	LinkRequestor string `json:"link_requestor"`
	LinkAcceptor  string `json:"link_acceptor"`
	PutPubkey     string `json:"put_pubkey"`
	PutOperation  string `json:"put_operation"`
	PutSeq        int    `json:"put_seq"`
	PutDataSize   int    `json:"put_data_size"`
}

// DHTGetReturn used to store information returned by dht putlinkrecord dev test01 pshare "<offer>"
type DHTGetReturn struct {
	LinkRequestor   string `json:"link_requestor"`
	LinkAcceptor    string `json:"link_acceptor"`
	GetPubkey       string `json:"get_pubkey"`
	GetOperation    string `json:"get_operation"`
	GetSeq          int    `json:"get_seq"`
	DataEncrypted   string `json:"data_encrypted"`
	DataVersion     int    `json:"data_version"`
	DataChunks      int    `json:"data_chunks"`
	GetValue        string `json:"get_value"`
	GetValueSize    int    `json:"get_value_size"`
	GetMilliseconds uint32 `json:"get_milliseconds"`
}

// PutLinkRecord calls the DHT put command to add an encrypted record for the given account link
func (d *Dynamicd) PutLinkRecord(sender, receiver, offer string) (*DHTPutReturn, error) {
	var ret DHTPutReturn
	cmd := "dynamic-cli dht putlinkrecord " + sender + " " + receiver + " pshare " + `"` + offer + `"`
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

// GetLinkRecord calls the DHT get record command to fetch an encrypted record for the given account link
func (d *Dynamicd) GetLinkRecord(sender, receiver string) (*DHTGetReturn, error) {
	var ret DHTGetReturn
	cmd := "dynamic-cli dht getlinkrecord " + sender + " " + receiver + " pshare"
	req, err := NewRequest(cmd)
	if err != nil {
		return nil, err
	}
	fmt.Println("Request:", req)
	res := <-d.ExecCmdRequest(req)
	err = json.Unmarshal([]byte(res), &ret)
	if err != nil {
		return nil, err
	}
	fmt.Println(cmd, &ret)
	return &ret, nil
}

// ClearLinkRecord clears an encrypted record for the given account link
func (d *Dynamicd) ClearLinkRecord(sender, receiver string) (*DHTPutReturn, error) {
	var ret DHTPutReturn
	cmd := "dynamic-cli dht clearlinkrecord " + sender + " " + receiver + " pshare"
	req, _ := NewRequest(cmd)
	err := json.Unmarshal([]byte(<-d.ExecCmdRequest(req)), &ret)
	if err != nil {
		return nil, err
	}
	fmt.Println(cmd, &ret)
	return &ret, nil
}

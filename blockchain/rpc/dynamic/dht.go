package dynamic

import (
	"encoding/json"
	"time"

	"github.com/duality-solutions/web-bridge/internal/util"
)

// DHTPutJSON used to store information returned by dht putlinkrecord dev test01 pshare "<offer>"
type DHTPutJSON struct {
	LinkRequestor string `json:"link_requestor"`
	LinkAcceptor  string `json:"link_acceptor"`
	PutPubkey     string `json:"put_pubkey"`
	PutOperation  string `json:"put_operation"`
	PutSeq        int    `json:"put_seq"`
	PutDataSize   int    `json:"put_data_size"`
	Timestamp     int64  `json:"timestamp"`
}

// DHTPutReturn used to store information returned by dht putlinkrecord dev test01 pshare "<offer>"
type DHTPutReturn struct {
	Sender   string
	Receiver string
	DHTPutJSON
}

// DHTGetJSON used to store information returned by dht putlinkrecord dev test01 pshare "<offer>"
type DHTGetJSON struct {
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
	NullRecord      string `json:"null_record"`
	Timestamp       int64  `json:"timestamp"`
}

// Minutes from DHT timestamp to now.
func (d DHTGetJSON) Minutes() int64 {
	if d.Timestamp > 0 {
		return (time.Now().Unix() - d.Timestamp) / 60
	}
	return -1
}

// Minutes from DHT timestamp to now.
func (d DHTPutJSON) Minutes() int64 {
	if d.Timestamp > 0 {
		return (time.Now().Unix() - d.Timestamp) / 60
	}
	return -1
}

// DHTGetReturn used to store information returned by dht putlinkrecord dev test01 pshare "<offer>"
type DHTGetReturn struct {
	Sender   string
	Receiver string
	DHTGetJSON
}

// PutLinkRecord calls the DHT put command to add an encrypted record for the given account link
func (d *Dynamicd) PutLinkRecord(sender, receiver, offer string, out chan<- DHTPutReturn) {
	go func() {
		var ret = DHTPutReturn{
			Sender:   sender,
			Receiver: receiver,
		}
		cmd := "dynamic-cli dht putlinkrecord " + sender + " " + receiver + " pshare " + `"` + offer + `"`
		req, err := NewRequest(cmd)
		if err != nil {
			util.Error.Println("PutLinkRecord error", err)
			out <- ret
		} else {
			var r DHTPutJSON
			json.Unmarshal([]byte(<-d.ExecCmdRequest(req)), &r)
			ret = DHTPutReturn{
				Sender:     sender,
				Receiver:   receiver,
				DHTPutJSON: r,
			}
			out <- ret
		}
	}()
}

// ClearLinkRecord clears an encrypted record for the given account link
func (d *Dynamicd) ClearLinkRecord(sender, receiver string, out chan<- DHTPutReturn) {
	go func() {
		var ret = DHTPutReturn{
			Sender:   sender,
			Receiver: receiver,
		}
		cmd := "dynamic-cli dht clearlinkrecord " + sender + " " + receiver + " pshare"
		req, err := NewRequest(cmd)
		if err != nil {
			util.Error.Println("ClearLinkRecord error", err)
			out <- ret
		} else {
			var r DHTPutJSON
			json.Unmarshal([]byte(<-d.ExecCmdRequest(req)), &r)
			ret = DHTPutReturn{
				Sender:     sender,
				Receiver:   receiver,
				DHTPutJSON: r,
			}
			out <- ret
		}
	}()
}

// GetLinkRecord calls the DHT get link record command to fetch an encrypted record for the given account link
func (d *Dynamicd) GetLinkRecord(sender, receiver string, out chan<- DHTGetReturn) {
	go func() {
		var ret = DHTGetReturn{
			Sender:   sender,
			Receiver: receiver,
		}
		cmd := "dynamic-cli dht getlinkrecord " + sender + " " + receiver + " pshare"
		req, err := NewRequest(cmd)
		if err != nil {
			util.Error.Println("GetLinkRecord error", err)
			out <- ret
		} else {
			var r DHTGetJSON
			json.Unmarshal([]byte(<-d.ExecCmdRequest(req)), &r)
			ret = DHTGetReturn{
				Sender:     sender,
				Receiver:   receiver,
				DHTGetJSON: r,
			}
			recordEpoch := time.Unix(ret.DHTGetJSON.Timestamp, 0).Unix()
			if recordEpoch > 0 && r.NullRecord == "false" {
				currentEpoch := time.Now().Unix()
				util.Info.Println("GetLinkRecord:", sender, "record timestamp", (currentEpoch-recordEpoch)/60, "minutes ago")
			}
			out <- ret
		}
	}()
}

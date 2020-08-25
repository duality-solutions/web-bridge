package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/duality-solutions/web-bridge/rpc/dynamic"
	"github.com/gin-gonic/gin"
)

type Link struct {
	LinkStatus             string
	RequestorFQDN          string `json:"requestor_fqdn"`
	RecipientFQDN          string `json:"recipient_fqdn"`
	SharedRequestPubkey    string `json:"shared_request_pubkey"`
	SharedAcceptPubkey     string `json:"shared_accept_pubkey"`
	RequestorWalletAddress string `json:"requestor_wallet_address"`
	RecipientWalletAddress string `json:"recipient_wallet_address"`
	RequestorLinkPubkey    string `json:"requestor_link_pubkey"`
	TxID                   string `json:"txid"`
	Time                   int64  `json:"time"`
	ExpiresOn              int64  `json:"expires_on"`
	Expired                bool   `json:"expired"`
	RecipientLinkPubkey    string `json:"recipient_link_pubkey"`
	AcceptTxID             string `json:"accept_txid"`
	AcceptTime             int64  `json:"accept_time"`
	AcceptExpiresOn        int64  `json:"accept_expires_on"`
	AcceptExpired          bool   `json:"accept_expired"`
	LinkMessage            string `json:"link_message"`
}

func (l *Link) SetValue(fieldname, value string) {
	switch fieldname {
	case "requestor_fqdn":
		l.RequestorFQDN = value
	case "recipient_fqdn":
		l.RecipientFQDN = value
	case "shared_request_pubkey":
		l.SharedRequestPubkey = value
	case "shared_accept_pubkey":
		l.SharedAcceptPubkey = value
	case "requestor_wallet_address":
		l.RequestorWalletAddress = value
	case "recipient_wallet_address":
		l.RecipientWalletAddress = value
	case "requestor_link_pubkey":
		l.RequestorLinkPubkey = value
	case "txid":
		l.TxID = value
	case "time":
		i, _ := strconv.ParseInt(value, 10, 64)
		l.Time = i
	case "expires_on":
		i, _ := strconv.ParseInt(value, 10, 64)
		l.ExpiresOn = i
	case "Expired":
		b, _ := strconv.ParseBool(value)
		l.Expired = b
	case "recipient_link_pubkey":
		l.RecipientLinkPubkey = value
	case "accept_txid":
		l.AcceptTxID = value
	case "accept_time":
		i, _ := strconv.ParseInt(value, 10, 64)
		l.AcceptTime = i
	case "accept_expires_on":
		i, _ := strconv.ParseInt(value, 10, 64)
		l.AcceptExpiresOn = i
	case "accept_expired":
		b, _ := strconv.ParseBool(value)
		l.AcceptExpired = b
	case "link_message":
		l.LinkMessage = value
	}
}

func (w *WebBridgeRunner) links(c *gin.Context) {
	cmd, _ := dynamic.NewRequest(`dynamic-cli link complete`)
	response, _ := <-w.dynamicd.ExecCmdRequest(cmd)
	complete := map[string]interface{}{}
	err := json.Unmarshal([]byte(response), &complete)
	if err != nil {
		strErrMsg := fmt.Sprintf("Results JSON unmarshal error %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": strErrMsg})
		return
	}

	myLinks := make(map[string]Link)
	for key, linkInterface := range complete {
		if key != "locked_links" {
			linkObj := Link{LinkStatus: "Complete"}
			linkVal := reflect.ValueOf(linkInterface)
			for _, lk := range linkVal.MapKeys() {
				link := linkVal.MapIndex(lk)
				linkObj.SetValue(lk.String(), fmt.Sprintf("%v", link))
			}
			myLinks[key] = linkObj
		}
	}

	cmd, _ = dynamic.NewRequest(`dynamic-cli link pending`)
	response, _ = <-w.dynamicd.ExecCmdRequest(cmd)
	pending := map[string]interface{}{}
	err = json.Unmarshal([]byte(response), &pending)
	if err != nil {
		strErrMsg := fmt.Sprintf("Results JSON unmarshal error %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": strErrMsg})
		return
	}

	for key, linkInterface := range pending {
		if key != "locked_links" {
			linkObj := Link{LinkStatus: "Pending"}
			linkVal := reflect.ValueOf(linkInterface)
			for _, lk := range linkVal.MapKeys() {
				link := linkVal.MapIndex(lk)
				linkObj.SetValue(lk.String(), fmt.Sprintf("%v", link))
			}
			myLinks[key] = linkObj
		}
	}
	c.JSON(http.StatusOK, gin.H{"result": myLinks})
}

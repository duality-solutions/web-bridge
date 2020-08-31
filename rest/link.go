package rest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"

	util "github.com/duality-solutions/web-bridge/internal/utilities"
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
		i, _ := util.ScientificNotationToInt64(value)
		l.Time = i
	case "expires_on":
		i, _ := util.ScientificNotationToInt64(value)
		l.ExpiresOn = i
	case "Expired":
		b, _ := strconv.ParseBool(value)
		l.Expired = b
	case "recipient_link_pubkey":
		l.RecipientLinkPubkey = value
	case "accept_txid":
		l.AcceptTxID = value
	case "accept_time":
		i, _ := util.ScientificNotationToInt64(value)
		l.AcceptTime = i
	case "accept_expires_on":
		i, _ := util.ScientificNotationToInt64(value)
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

type linkRequest struct {
	RequestorFQDN string `json:"requestor_fqdn"`
	RecipientFQDN string `json:"recipient_fqdn"`
	LinkMessage   string `json:"link_message"`
}

type linkAccept struct {
	RecipientFQDN string `json:"recipient_fqdn"`
	RequestorFQDN string `json:"requestor_fqdn"`
}

type linkRequestResponse struct {
	Expired                bool   `json:"expired"`
	ExpiresOn              int    `json:"expires_on"`
	LinkMessage            string `json:"link_message"`
	RecipientFQDN          string `json:"recipient_fqdn"`
	RecipientWalletAddress string `json:"recipient_wallet_address"`
	RequestorFQDN          string `json:"requestor_fqdn"`
	RequestorLinkPubkey    string `json:"requestor_link_pubkey"`
	RequestorWalletAddress string `json:"requestor_wallet_address"`
	SignatureProof         string `json:"signature_proof"`
	Time                   int    `json:"time"`
	TxID                   string `json:"txid"`
}

func (w *WebBridgeRunner) linkrequest(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		strErrMsg := fmt.Sprintf("Request body read all error %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": strErrMsg})
		return
	}
	if len(body) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request body is empty"})
		return
	}
	req := linkRequest{}
	err = json.Unmarshal(body, &req)
	if err != nil {
		strErrMsg := fmt.Sprintf("Request body JSON unmarshal error %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": strErrMsg})
		return
	}

	cmd, _ := dynamic.NewRequest(`dynamic-cli link request "` + req.RequestorFQDN + `" "` + req.RecipientFQDN + `" "` + req.LinkMessage + `"`)
	response, _ := <-w.dynamicd.ExecCmdRequest(cmd)
	var result interface{}
	err = json.Unmarshal([]byte(response), &result)
	if err != nil {
		strErrMsg := fmt.Sprintf("Results JSON unmarshal error %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": strErrMsg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}

func (w *WebBridgeRunner) linkaccept(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		strErrMsg := fmt.Sprintf("Request body read all error %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": strErrMsg})
		return
	}
	if len(body) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request body is empty"})
		return
	}
	req := linkAccept{}
	err = json.Unmarshal(body, &req)
	if err != nil {
		strErrMsg := fmt.Sprintf("Request body JSON unmarshal error %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": strErrMsg})
		return
	}

	cmd, _ := dynamic.NewRequest(`dynamic-cli link accept "` + req.RecipientFQDN + `" "` + req.RequestorFQDN + `"`)
	response, _ := <-w.dynamicd.ExecCmdRequest(cmd)
	var result interface{}
	err = json.Unmarshal([]byte(response), &result)
	if err != nil {
		strErrMsg := fmt.Sprintf("Results JSON unmarshal error %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": strErrMsg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}

type sendMessageRequest struct {
	SenderFQDN    string `json:"sender_fqdn"`
	RecipientFQDN string `json:"recipient_fqdn"`
	MessageType   string `json:"message_type"`
	Message       string `json:"message"`
	keepLast      bool   `json:"keep_last"`
}

func (w *WebBridgeRunner) sendlinkmessage(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		strErrMsg := fmt.Sprintf("Request body read all error %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": strErrMsg})
		return
	}
	if len(body) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request body is empty"})
		return
	}
	reqBody := sendMessageRequest{}
	err = json.Unmarshal(body, &reqBody)
	if err != nil {
		strErrMsg := fmt.Sprintf("Request body JSON unmarshal error %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": strErrMsg})
		return
	}
	var keepLast = "1"
	if reqBody.keepLast == false {
		keepLast = "0"
	}
	// Set dynamic CLI command
	cmd := `dynamic-cli link sendmessage "` + reqBody.SenderFQDN + `" "` + reqBody.RecipientFQDN + `" "` +
		reqBody.MessageType + `" "` + reqBody.Message + `" "` + keepLast + `"`
	// Create new dynamic CLI request from command
	req, err := dynamic.NewRequest(cmd)
	if err != nil {
		strErrMsg := fmt.Sprintf("Dynamic CLI new request error %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": strErrMsg})
		return
	}
	// Execute dynamic CLI request
	res := <-w.dynamicd.ExecCmdRequest(req)
	var ret dynamic.MessageReturnJSON
	err = json.Unmarshal([]byte(res), &ret)
	if err != nil {
		strErrMsg := fmt.Sprintf("Dynamic CLI response JSON unmarshal error %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": strErrMsg})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": ret})
}

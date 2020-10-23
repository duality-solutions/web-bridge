package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/duality-solutions/web-bridge/api/models"
	"github.com/duality-solutions/web-bridge/blockchain/rpc/dynamic"
	"github.com/gin-gonic/gin"
)

// LinkStatus is the current link status
type LinkStatus string

const (
	// Complete link
	Complete LinkStatus = "Complete"
	// Pending link
	Pending LinkStatus = "Pending"
)

// GetLinks returns a list of pending or complete links from the Dynamic RPC server
func (w *WebBridgeRunner) GetLinks(status LinkStatus) (*models.LinksResponse, error) {
	var links models.LinksResponse
	var cmdStr = ""
	if status == Complete {
		cmdStr = `dynamic-cli link complete`
	} else if status == Pending {
		cmdStr = `dynamic-cli link pending`
	}
	cmd, _ := dynamic.NewRequest(cmdStr)
	response, _ := <-w.dynamicd.ExecCmdRequest(cmd)
	results := map[string]interface{}{}
	err := json.Unmarshal([]byte(response), &results)
	if err != nil {
		strErrMsg := fmt.Sprintf("Results JSON unmarshal error %v", err)
		return nil, errors.New(strErrMsg)
	}
	myLinks := make(map[string]models.Link)
	for key, linkInterface := range results {
		if key != "locked_links" {
			linkObj := models.Link{LinkStatus: string(status)}
			linkVal := reflect.ValueOf(linkInterface)
			for _, lk := range linkVal.MapKeys() {
				link := linkVal.MapIndex(lk)
				linkObj.SetValue(lk.String(), fmt.Sprintf("%v", link))
			}
			myLinks[key] = linkObj
		} else {
			links.LockedLinks = linkInterface.(float64)
		}
	}
	links.Links = myLinks
	return &links, nil
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

	myLinks := make(map[string]models.Link)
	for key, linkInterface := range complete {
		if key != "locked_links" {
			linkObj := models.Link{LinkStatus: "Complete"}
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
			linkObj := models.Link{LinkStatus: "Pending"}
			linkVal := reflect.ValueOf(linkInterface)
			for _, lk := range linkVal.MapKeys() {
				link := linkVal.MapIndex(lk)
				linkObj.SetValue(lk.String(), fmt.Sprintf("%v", link))
			}
			myLinks[key] = linkObj
		}
	}
	c.JSON(http.StatusOK, myLinks)
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
	req := models.LinkRequest{}
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

	c.JSON(http.StatusOK, result)
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
	req := models.LinkAccept{}
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

	c.JSON(http.StatusOK, result)
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
	reqBody := models.SendMessageRequest{}
	err = json.Unmarshal(body, &reqBody)
	if err != nil {
		strErrMsg := fmt.Sprintf("Request body JSON unmarshal error %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": strErrMsg})
		return
	}
	var keepLast = "1"
	if reqBody.KeepLast == false {
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
	c.JSON(http.StatusOK, ret)
}

func (w *WebBridgeRunner) getlinkmessages(c *gin.Context) {
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
	reqBody := models.GetMessageRequest{}
	err = json.Unmarshal(body, &reqBody)
	if err != nil {
		strErrMsg := fmt.Sprintf("Request body JSON unmarshal error %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": strErrMsg})
		return
	}
	// Set dynamic CLI command
	cmd := `dynamic-cli link getaccountmessages "` + reqBody.RecipientFQDN + `" "` + reqBody.SenderFQDN + `"`
	if len(reqBody.MessageType) > 0 {
		cmd += ` "` + reqBody.MessageType + `"`
	}
	// Create new dynamic CLI request from command
	req, err := dynamic.NewRequest(cmd)
	if err != nil {
		strErrMsg := fmt.Sprintf("Dynamic CLI new request error %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": strErrMsg})
		return
	}
	// Execute dynamic CLI request
	res := <-w.dynamicd.ExecCmdRequest(req)
	var ret interface{}
	err = json.Unmarshal([]byte(res), &ret)
	if err != nil {
		strErrMsg := fmt.Sprintf("Dynamic CLI response JSON unmarshal error %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": strErrMsg})
		return
	}
	c.JSON(http.StatusOK, ret)
}

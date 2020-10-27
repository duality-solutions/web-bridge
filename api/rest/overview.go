package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/duality-solutions/web-bridge/api/models"
	"github.com/duality-solutions/web-bridge/blockchain/rpc/dynamic"
	"github.com/duality-solutions/web-bridge/bridge"
	"github.com/gin-gonic/gin"
)

// GetBlockchainOverview returns the current blockchain overview status
func (w *WebBridgeRunner) GetBlockchainOverview() (*models.BlockchainOverview, int, error) {
	cmd, _ := dynamic.NewRequest("dynamic-cli syncstatus")
	response, _ := <-w.dynamicd.ExecCmdRequest(cmd)
	if strings.Contains(response, "Error:") {
		result := models.RPCError{}
		err := json.Unmarshal([]byte(response), &result)
		if err != nil {
			strErrMsg := fmt.Sprintf("GetBlockchainOverview() (syncstatus) JSON unmarshal error %v", err)
			return nil, http.StatusInternalServerError, errors.New(strErrMsg)
		}
		return nil, http.StatusInternalServerError, errors.New(result.Error.Message)
	}
	var status models.SyncStatus
	err := json.Unmarshal([]byte(response), &status)
	if err != nil {
		strErrMsg := fmt.Sprintf("GetBlockchainOverview() results (syncstatus) JSON unmarshal error %v", err)
		return nil, http.StatusInternalServerError, errors.New(strErrMsg)
	}

	overview := models.BlockchainOverview{
		Network:           status.ChainName,
		ClientVersion:     status.ClientVersion,
		Peers:             status.Peers,
		Blocks:            status.Blocks,
		TotalBlocks:       status.CurrentBlockHeight,
		SyncProgress:      status.SyncProgress,
		StatusDescription: status.StatusDescription,
		FullySynced:       status.FullySynced,
	}
	return &overview, http.StatusOK, nil
}

// GetWalletBalanceOverview returns the current wallet overview status
func (w *WebBridgeRunner) GetWalletBalanceOverview() (*models.WalletOverview, int, error) {
	cmd, _ := dynamic.NewRequest(`dynamic-cli getwalletinfo`)
	response, _ := <-w.dynamicd.ExecCmdRequest(cmd)
	if strings.Contains(response, "Error:") {
		result := models.RPCError{}
		err := json.Unmarshal([]byte(response), &result)
		if err != nil {
			strErrMsg := fmt.Sprintf("Response (getwalletinfo) JSON unmarshal error %v", err)
			return nil, http.StatusInternalServerError, errors.New(strErrMsg)
		}
		return nil, http.StatusInternalServerError, errors.New(result.Error.Message)
	}
	var encrypted bool = false
	if strings.Contains(response, `"unlocked_until"`) {
		encrypted = true
	}
	var walletinfo models.WalletInfoResponse
	err := json.Unmarshal([]byte(response), &walletinfo)
	if err != nil {
		strErrMsg := fmt.Sprintf("Results JSON unmarshal error %v", err)
		return nil, http.StatusInternalServerError, errors.New(strErrMsg)
	}
	cmd, _ = dynamic.NewRequest("dynamic-cli getcredits")
	response, _ = <-w.dynamicd.ExecCmdRequest(cmd)
	if strings.Contains(response, "Error:") {
		result := models.RPCError{}
		err := json.Unmarshal([]byte(response), &result)
		if err != nil {
			strErrMsg := fmt.Sprintf("GetBlockchainOverview() (getcredits) JSON unmarshal error %v", err)
			return nil, http.StatusInternalServerError, errors.New(strErrMsg)
		}
		return nil, http.StatusInternalServerError, errors.New(result.Error.Message)
	}
	results := map[string]interface{}{}
	err = json.Unmarshal([]byte(response), &results)
	if err != nil {
		strErrMsg := fmt.Sprintf("GetBlockchainOverview() results (getcredits) JSON unmarshal error %v", err)
		return nil, http.StatusInternalServerError, errors.New(strErrMsg)
	}
	credits := models.CreditsResponse{}
	creditTxs := make(map[string]models.CreditTransaction)
	for key, credittx := range results {
		if key != "total_credits" && key != "total_deposits" && key != "total_dynamic" {
			linkObj := models.CreditTransaction{}
			linkVal := reflect.ValueOf(credittx)
			for _, lk := range linkVal.MapKeys() {
				credit := linkVal.MapIndex(lk)
				linkObj.SetValue(lk.String(), fmt.Sprintf("%v", credit))
			}
			creditTxs[key] = linkObj
		} else if key == "total_credits" {
			credits.TotalCredits = credittx.(float64)
		} else if key == "total_deposits" {
			credits.TotalDeposits = credittx.(float64)
		} else if key == "total_dynamic" {
			credits.TotalDynamic = credittx.(string)
		}
	}
	overview := models.WalletOverview{
		Transactions:     walletinfo.TxCount,
		Encrypted:        encrypted,
		UnlockedEpoch:    walletinfo.UnlockedUntil,
		AvailableBalance: walletinfo.Balance - walletinfo.UnconfirmedBalance,
		PendingBalance:   walletinfo.UnconfirmedBalance,
		TotalBalance:     walletinfo.Balance,
		Credits:          credits.TotalCredits,
		Deposits:         credits.TotalDeposits,
	}
	return &overview, http.StatusOK, nil
}

// GetAccountOverview returns the current account overview status
func (w *WebBridgeRunner) GetAccountOverview() (*models.AccountOverview, int, error) {
	cmd, _ := dynamic.NewRequest("dynamic-cli mybdapaccounts")
	response, _ := <-w.dynamicd.ExecCmdRequest(cmd)
	if strings.Contains(response, "Error:") {
		result := models.RPCError{}
		err := json.Unmarshal([]byte(response), &result)
		if err != nil {
			strErrMsg := fmt.Sprintf("GetAccountOverview() (mybdapaccounts) JSON unmarshal error %v", err)
			return nil, http.StatusInternalServerError, errors.New(strErrMsg)
		}
		return nil, http.StatusInternalServerError, errors.New(result.Error.Message)
	}
	myAccounts := make(map[string]models.Account)
	err := json.Unmarshal([]byte(response), &myAccounts)
	if err != nil {
		strErrMsg := fmt.Sprintf("GetAccountOverview() results (mybdapaccounts) JSON unmarshal error %v", err)
		return nil, http.StatusInternalServerError, errors.New(strErrMsg)
	}
	myCompleteLinks, err := w.GetLinks(Complete)
	myPendingLinks, err := w.GetLinks(Pending)
	overview := models.AccountOverview{
		Users:         len(myAccounts),
		CompleteLinks: len(myCompleteLinks.Links) + int(myCompleteLinks.LockedLinks),
		PendingLinks:  len(myPendingLinks.Links) + int(myPendingLinks.LockedLinks),
		Certificates:  0,
		Audits:        0,
	}
	return &overview, http.StatusOK, nil
}

// GetBridgeOverview returns the current bridge overview status
func (w *WebBridgeRunner) GetBridgeOverview() (*models.BridgeOverview, int, error) {
	controler, err := bridge.Controler()
	if err != nil {
		strErrMsg := fmt.Sprintf("GetBridgeOverview() (bridge.Controler) error %v", err)
		return nil, http.StatusInternalServerError, errors.New(strErrMsg)
	}
	bridges := controler.AllBridges()
	var connected, connecting, idle, disabled, stopped int
	for _, bridge := range bridges {
		if bridge.State().String() == "StateOpenConnection" {
			connected++
		} else if bridge.State().String() == "StateNew" {
			idle++
		} else if bridge.State().String() == "StateWaitForOffer" {
			idle++
		} else if bridge.State().String() == "StateWaitForAnswer" {
			idle++
		} else if bridge.State().String() == "StateSendAnswer" {
			connecting++
		} else if bridge.State().String() == "StateWaitForRTC" {
			connecting++
		} else if bridge.State().String() == "StateEstablishRTC" {
			connecting++
		} else if bridge.State().String() == "StateDisconnected" {
			stopped++
		} else if bridge.State().String() == "StateShutdown" {
			stopped++
		}
	}
	overview := models.BridgeOverview{
		Total:      len(bridges),
		Connected:  connected,
		Connecting: connecting,
		Idle:       idle,
		Disabled:   disabled,
		Stopped:    stopped,
	}
	return &overview, http.StatusOK, nil
}

//
// @Description Returns the current WebBridge overview status
// @Accept  json
// @Produce  json
// @Success 200 {object} models.WalletSetupStatus "ok"
// @Failure 400 {object} string "Bad request"
// @Failure 500 {object} string "Internal error"
// @Router /api/v1/overview [get]
func (w *WebBridgeRunner) overview(c *gin.Context) {
	response := models.OverviewResponse{}
	blockchain, httpStatus, err := w.GetBlockchainOverview()
	if err != nil {
		strErrMsg := fmt.Sprintf("%v", err)
		c.JSON(httpStatus, gin.H{"error": strErrMsg})
		return
	}
	wallet, httpStatus, err := w.GetWalletBalanceOverview()
	if err != nil {
		strErrMsg := fmt.Sprintf("%v", err)
		c.JSON(httpStatus, gin.H{"error": strErrMsg})
		return
	}
	accounts, httpStatus, err := w.GetAccountOverview()
	if err != nil {
		strErrMsg := fmt.Sprintf("%v", err)
		c.JSON(httpStatus, gin.H{"error": strErrMsg})
		return
	}
	if accounts.CompleteLinks > 0 {
		bridges, _, err := w.GetBridgeOverview()
		if err != nil {
			response = models.OverviewResponse{
				Chain:    *blockchain,
				Wallet:   *wallet,
				Accounts: *accounts,
				Bridges:  models.DefaultBridgeOverview(),
			}
		} else {
			response = models.OverviewResponse{
				Chain:    *blockchain,
				Wallet:   *wallet,
				Accounts: *accounts,
				Bridges:  *bridges,
			}
		}
	} else {
		response = models.OverviewResponse{
			Chain:    *blockchain,
			Wallet:   *wallet,
			Accounts: *accounts,
			Bridges:  models.DefaultBridgeOverview(),
		}
	}
	c.JSON(httpStatus, response)
}

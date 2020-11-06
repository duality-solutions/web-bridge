package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/duality-solutions/web-bridge/api/models"
	"github.com/duality-solutions/web-bridge/blockchain/rpc/dynamic"
	"github.com/gin-gonic/gin"
)

// GetWalletSetupInfo gathers data configuration file and dynamicd to determine the current wallet setup status
func (w *WebBridgeRunner) GetWalletSetupInfo() (*models.WalletSetupStatus, int, error) {
	status := w.configuration.WalletSetupStatus()
	cmd := `dynamic-cli getwalletinfo`
	reqCnd, _ := dynamic.NewRequest(cmd)
	response, _ := <-w.dynamicd.ExecCmdRequest(reqCnd)
	if strings.Contains(response, "Error:") {
		result := models.RPCError{}
		err := json.Unmarshal([]byte(response), &result)
		if err != nil {
			strErrMsg := fmt.Sprintf("Response (getwalletinfo) JSON unmarshal error %v", err)
			return nil, http.StatusInternalServerError, errors.New(strErrMsg)
		}
		return nil, http.StatusInternalServerError, errors.New(result.Error.Message)
	}
	if strings.Contains(response, `"unlocked_until"`) {
		status.WalletEncrypted = true
	}
	var walletinfo models.WalletInfoResponse
	err := json.Unmarshal([]byte(response), &walletinfo)
	if err != nil {
		strErrMsg := fmt.Sprintf("Results JSON unmarshal error %v", err)
		return nil, http.StatusInternalServerError, errors.New(strErrMsg)
	}
	var currentStatus = *w.configuration.WalletSetupStatus()
	status.MnemonicBackup = currentStatus.MnemonicBackup
	if walletinfo.TxCount > 0 {
		status.HasTransactions = true
	}
	if status.WalletEncrypted {
		status.UnlockedUntil = walletinfo.UnlockedUntil
	}
	pending, err := w.GetLinks(Pending)
	if err != nil {
		strErrMsg := fmt.Sprintf("GetPendingLinks failed. error %v", err)
		return nil, http.StatusInternalServerError, errors.New(strErrMsg)
	}
	if len(pending.Links) > 0 || pending.LockedLinks > 0 {
		status.HasLinks = true
		status.HasAccounts = true // assuming accounts are needed to create links
	} else {
		complete, err := w.GetLinks(Complete)
		if err != nil {
			strErrMsg := fmt.Sprintf("GetCompleteLinks failed. error %v", err)
			return nil, http.StatusInternalServerError, errors.New(strErrMsg)
		}
		if len(complete.Links) > 0 || complete.LockedLinks > 0 {
			status.HasLinks = true
			status.HasAccounts = true // assuming accounts are needed to create links
		} else {
			// No links, check for accounts
			cmd = `dynamic-cli mybdapaccounts`
			reqCnd, _ = dynamic.NewRequest(cmd)
			response, _ = <-w.dynamicd.ExecCmdRequest(reqCnd)
			if strings.Contains(response, "error:") || strings.Contains(response, "Error:") {
				result := models.RPCError{}
				err := json.Unmarshal([]byte(response), &result)
				if err != nil {
					strErrMsg := fmt.Sprintf("Response (mybdapaccounts) JSON unmarshal error %v", err)
					return nil, http.StatusInternalServerError, errors.New(strErrMsg)
				}
				return nil, http.StatusInternalServerError, errors.New(result.Error.Message)
			}
			var myAccounts map[string]models.MyAccountResponse
			err = json.Unmarshal([]byte(response), &myAccounts)
			if err != nil {
				strErrMsg := fmt.Sprintf("Results (MyAccountsResponse) JSON unmarshal error %v", err)
				return nil, http.StatusInternalServerError, errors.New(strErrMsg)
			}
			if len(myAccounts) > 0 {
				status.HasAccounts = true
			}
		}
	}
	return status, http.StatusOK, nil
}

//
// @Description Returns the current wallet setup status
// @Accept  json
// @Produce  json
// @Success 200 {object} models.WalletSetupStatus "ok"
// @Failure 400 {object} string "Bad request"
// @Failure 500 {object} string "Internal error"
// @Router /api/v1/wallet/setup [get]
func (w *WebBridgeRunner) walletsetup(c *gin.Context) {
	status, httpStatus, err := w.GetWalletSetupInfo()
	if err != nil {
		strErrMsg := fmt.Sprintf("%v", err)
		c.JSON(httpStatus, gin.H{"error": strErrMsg})
		return
	}
	c.JSON(httpStatus, status)
}

//
// @Description Updates the wallet setup backup mnemonic configuration to true
// @Accept  json
// @Produce  json
// @Success 200 {object} models.WalletSetupStatus "ok"
// @Failure 400 {object} string "Bad request"
// @Failure 500 {object} string "Internal error"
// @Router /api/v1/wallet/setup/backup [post]
func (w *WebBridgeRunner) setupmnemonicbackup(c *gin.Context) {
	status, httpStatus, err := w.GetWalletSetupInfo()
	if err != nil {
		strErrMsg := fmt.Sprintf("%v", err)
		c.JSON(httpStatus, gin.H{"error": strErrMsg})
		return
	}
	status.MnemonicBackup = true
	w.configuration.UpdateWalletSetup(*status)
	c.JSON(httpStatus, *status)
}

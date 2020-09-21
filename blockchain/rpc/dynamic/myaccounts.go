package dynamic

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/duality-solutions/web-bridge/api/models"
	"github.com/duality-solutions/web-bridge/internal/util"
)

// Account stores a BDAP account object
type Account struct {
	OID                string `json:"oid"`
	Version            int    `json:"version"`
	DomainComponent    string `json:"domain_component"`
	CommonName         string `json:"common_name"`
	OrganizationalUnit string `json:"organizational_unit"`
	OrganizationName   string `json:"organization_name"`
	ObjectID           string `json:"object_id"`
	ObjectFullPath     string `json:"object_full_path"`
	ObjectType         string `json:"object_type"`
	WalletAddress      string `json:"wallet_address"`
	Public             int8   `json:"public"`
	DHTPublicKey       string `json:"dht_publickey"`
	LinkAddress        string `json:"link_address"`
	TxID               string `json:"txid"`
	Time               int    `json:"time"`
	ExpiresOn          int    `json:"expires_on"`
	Expired            bool   `json:"expired"`
}

func (d *Dynamicd) myAccounts() (*[]Account, error) {
	var accounts = []Account{}
	req, _ := NewRequest("dynamic-cli mybdapaccounts")
	rawResp := []byte(<-d.ExecCmdRequest(req))
	err := json.Unmarshal(rawResp, &accounts)
	if err != nil {
		var rpcErr models.RPCError
		err = json.Unmarshal(rawResp, &rpcErr)
		if err != nil {
			util.Error.Printf("Getting accounts from blockchain daemon failed (%v)\n", string(rawResp))
		}
	} else {
		return &accounts, nil
	}
	var accountsGeneric map[string]interface{}
	err = json.Unmarshal(rawResp, &accountsGeneric)
	for _, v := range accountsGeneric {
		b, err := json.Marshal(v)
		if err == nil {
			var account Account
			err = json.Unmarshal(b, &account)
			if err != nil {
				util.Error.Printf("Getting accounts from blockchain daemon failed (%v)\n", string(rawResp))
				return nil, err
			}
			accounts = append(accounts, account)
		}
	}
	return &accounts, nil
}

// GetMyAccounts returns all BDAP accounts from the wallet
func (d *Dynamicd) GetMyAccounts(timeout time.Duration) (*[]Account, error) {
	myAccounts, err := d.myAccounts()
	if err != nil {
		for {
			select {
			case <-time.After(time.Second * 5):
				myAccounts, err = d.myAccounts()
				if err == nil {
					return myAccounts, nil
				}
			case <-time.After(timeout):
				return nil, fmt.Errorf("GetMyAccounts failed after timeout")
			}
		}
	} else {
		return myAccounts, nil
	}
}

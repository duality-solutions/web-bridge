package models

// Account stores a BDAP user and group record
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
	Public             int    `json:"public"`
	DhtPublickey       string `json:"dht_publickey"`
	LinkAddress        string `json:"link_address"`
	TxID               string `json:"txid"`
	Time               int    `json:"time"`
	ExpiresOn          int    `json:"expires_on"`
	Expired            bool   `json:"expired"`
}

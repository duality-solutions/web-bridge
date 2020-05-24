package dynamic

import "strings"

// Link stores a BDAP link data that is returned by running dynamic link RPC commands
type Link struct {
	AcceptExpired          bool   `json:"accept_expired"`
	AcceptExpiresOn        int    `json:"accept_expires_on"`
	AcceptTime             int    `json:"accept_time"`
	AcceptTxID             string `json:"accept_txid"`
	Expired                bool   `json:"expired"`
	ExpiresOn              int    `json:"expires_on"`
	LinkMessage            string `json:"link_message"`
	RecipientFQDN          string `json:"recipient_fqdn"`
	RecipientLinkPubkey    string `json:"recipient_link_pubkey"`
	RecipientWalletAddress string `json:"recipient_wallet_address"`
	RequestorFQDN          string `json:"requestor_fqdn"`
	RequestorLinkPubkey    string `json:"requestor_link_pubkey"`
	RequestorWalletAddress string `json:"requestor_wallet_address"`
	SharedAcceptPubkey     string `json:"shared_accept_pubkey"`
	SharedRequestPubkey    string `json:"shared_request_pubkey"`
	Time                   int    `json:"time"`
	TxID                   string `json:"txid"`
}

//GetRecipientObjectID returns the recipient object id from the FQDN
func (l *Link) GetRecipientObjectID() string {
	p := strings.Split(l.RecipientFQDN, "@")
	if len(p) > 0 {
		return p[0]
	}
	return ""
}

//GetRequestorObjectID returns the requestor object id from the FQDN
func (l *Link) GetRequestorObjectID() string {
	p := strings.Split(l.RequestorFQDN, "@")
	if len(p) > 0 {
		return p[0]
	}
	return ""
}

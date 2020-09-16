package models

import (
	"strconv"

	"github.com/duality-solutions/web-bridge/internal/util"
)

type Link struct {
	LinkStatus             string `json:"link_status"`
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

type LinkRequest struct {
	RequestorFQDN string `json:"requestor_fqdn"`
	RecipientFQDN string `json:"recipient_fqdn"`
	LinkMessage   string `json:"link_message"`
}

type LinkAccept struct {
	RecipientFQDN string `json:"recipient_fqdn"`
	RequestorFQDN string `json:"requestor_fqdn"`
}

type SendMessageRequest struct {
	SenderFQDN    string `json:"sender_fqdn"`
	RecipientFQDN string `json:"recipient_fqdn"`
	MessageType   string `json:"message_type"`
	Message       string `json:"message"`
	KeepLast      bool   `json:"keep_last"`
}

type GetMessageRequest struct {
	RecipientFQDN string `json:"recipient_fqdn"`
	SenderFQDN    string `json:"sender_fqdn"`
	MessageType   string `json:"message_type"`
}

/*
type LinkRequestResponse struct {
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
*/

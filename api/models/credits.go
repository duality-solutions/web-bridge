package models

import "strconv"

// CreditTransaction stores a credit transaction
type CreditTransaction struct {
	Type          string `json:"type"`
	Operation     string `json:"operation"`
	Address       string `json:"address"`
	Pubkey        string `json:"pubkey"`
	SharedPubkey  string `json:"shared_pubkey"`
	DynamicAmount string `json:"dynamic_amount"`
	Credits       int    `json:"credits"`
}

// SetValue sets a field value for a CreditTransaction
func (ct *CreditTransaction) SetValue(fieldname, value string) {
	switch fieldname {
	case "type":
		ct.Type = value
	case "operation":
		ct.Operation = value
	case "address":
		ct.Address = value
	case "pubkey":
		ct.Pubkey = value
	case "shared_pubkey":
		ct.SharedPubkey = value
	case "dynamic_amount":
		ct.DynamicAmount = value
	case "credits":
		ct.Credits, _ = strconv.Atoi(value)
	}
}

// CreditsResponse response for getcredits JSON RPC method
type CreditsResponse struct {
	Credits       map[string]CreditTransaction
	TotalCredits  float64
	TotalDeposits float64
	TotalDynamic  string
}

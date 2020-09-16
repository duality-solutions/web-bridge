package models

type JSONRPC struct {
	JSONRPC interface{}   `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      interface{}   `json:"id"`
}

package models

type RPCError struct {
	Error struct {
		Code    int64  `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
	ID      interface{} `json:"id"`
	JSONRPC string      `json:"jsonrpc"`
}

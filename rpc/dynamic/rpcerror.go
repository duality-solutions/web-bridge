package dynamic

// RPCError contains error code and message
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// RPCErrorResponse contains the JSON error response message
type RPCErrorResponse struct {
	Result string   `json:"result"`
	Error  RPCError `json:"error"`
	ID     string   `json:"id"`
}

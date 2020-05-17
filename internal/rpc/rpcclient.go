// Copyright (c) 2014 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package rpcclient

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"

	socks "github.com/duality-solutions/web-bridge/internal/socks"
)

// RPCErrorCode represents an error code to be used as a part of an RPCError
// which is in turn used in a JSON-RPC Response object.
//
// A specific type is used to help ensure the wrong errors aren't used.
type RPCErrorCode int

// RPCError represents an error that is used as a part of a JSON-RPC Response
// object.
type RPCError struct {
	Code    RPCErrorCode `json:"code,omitempty"`
	Message string       `json:"message,omitempty"`
}

// Guarantee RPCError satisfies the builtin error interface.
var _, _ error = RPCError{}, (*RPCError)(nil)

// Error returns a string describing the RPC error.  This satisfies the
// builtin error interface.
func (e RPCError) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

// NewRPCError constructs and returns a new JSON-RPC error that is suitable
// for use in a JSON-RPC Response object.
func NewRPCError(code RPCErrorCode, message string) *RPCError {
	return &RPCError{
		Code:    code,
		Message: message,
	}
}

// Response is the general form of a JSON-RPC response.  The type of the Result
// field varies from one command to the next, so it is implemented as an
// interface.  The ID field has to be a pointer for Go to put a null in it when
// empty.
type Response struct {
	Result json.RawMessage `json:"result"`
	Error  *RPCError       `json:"error"`
	ID     *interface{}    `json:"id"`
}

// Config defines the configuration options for btcctl.
//
// See loadConfig for details on the configuration load process.
type Config struct {
	ConfigFile     string `short:"C" long:"configfile" description:"Path to configuration file"`
	ListCommands   bool   `short:"l" long:"listcommands" description:"List all of the supported commands and exit"`
	NoTLS          bool   `long:"notls" description:"Disable TLS"`
	Proxy          string `long:"proxy" description:"Connect via SOCKS5 proxy (eg. 127.0.0.1:9050)"`
	ProxyPass      string `long:"proxypass" default-mask:"-" description:"Password for proxy server"`
	ProxyUser      string `long:"proxyuser" description:"Username for proxy server"`
	RegressionTest bool   `long:"regtest" description:"Connect to the regression test network"`
	RPCCert        string `short:"c" long:"rpccert" description:"RPC server certificate chain for validation"`
	RPCPassword    string `short:"P" long:"rpcpass" default-mask:"-" description:"RPC password"`
	RPCServer      string `short:"s" long:"rpcserver" description:"RPC server to connect to"`
	RPCUser        string `short:"u" long:"rpcuser" description:"RPC username"`
	SimNet         bool   `long:"simnet" description:"Connect to the simulation test network"`
	TLSSkipVerify  bool   `long:"skipverify" description:"Do not verify tls certificates (not recommended!)"`
	TestNet3       bool   `long:"testnet" description:"Connect to testnet"`
	ShowVersion    bool   `short:"V" long:"version" description:"Display version information and exit"`
	Wallet         bool   `long:"wallet" description:"Connect to wallet"`
}

// newHTTPClient returns a new HTTP client that is configured according to the
// proxy and TLS settings in the associated connection configuration.
func newHTTPClient(cfg *Config) (*http.Client, error) {
	// Configure proxy if needed.
	var dial func(network, addr string) (net.Conn, error)
	if cfg.Proxy != "" {
		proxy := &socks.Proxy{
			Addr:     cfg.Proxy,
			Username: cfg.ProxyUser,
			Password: cfg.ProxyPass,
		}
		dial = func(network, addr string) (net.Conn, error) {
			c, err := proxy.Dial(network, addr)
			if err != nil {
				return nil, err
			}
			return c, nil
		}
	}

	// Configure TLS if needed.
	var tlsConfig *tls.Config
	if !cfg.NoTLS && cfg.RPCCert != "" {
		pem, err := ioutil.ReadFile(cfg.RPCCert)
		if err != nil {
			return nil, err
		}

		pool := x509.NewCertPool()
		pool.AppendCertsFromPEM(pem)
		tlsConfig = &tls.Config{
			RootCAs:            pool,
			InsecureSkipVerify: cfg.TLSSkipVerify,
		}
	}

	// Create and return the new HTTP client potentially configured with a
	// proxy and TLS.
	client := http.Client{
		Transport: &http.Transport{
			Dial:            dial,
			TLSClientConfig: tlsConfig,
		},
	}
	return &client, nil
}

// SendPostRequest sends the marshalled JSON-RPC command using HTTP-POST mode
// to the server described in the passed config struct.  It also attempts to
// unmarshal the response as a JSON-RPC response and returns either the result
// field or the error field depending on whether or not there is an error.
func SendPostRequest(marshalledJSON []byte, cfg *Config) ([]byte, error) {
	// Generate a request to the configured RPC server.
	protocol := "http"
	if !cfg.NoTLS {
		protocol = "https"
	}
	url := protocol + "://" + cfg.RPCServer
	bodyReader := bytes.NewReader(marshalledJSON)
	httpRequest, err := http.NewRequest("POST", url, bodyReader)
	if err != nil {
		return nil, err
	}
	httpRequest.Close = true
	httpRequest.Header.Set("Content-Type", "application/json")

	// Configure basic access authorization.
	httpRequest.SetBasicAuth(cfg.RPCUser, cfg.RPCPassword)

	// Create the new HTTP client that is configured according to the user-
	// specified options and submit the request.
	httpClient, err := newHTTPClient(cfg)
	if err != nil {
		return nil, err
	}
	httpResponse, err := httpClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}

	// Read the raw bytes and close the response.
	respBytes, err := ioutil.ReadAll(httpResponse.Body)
	httpResponse.Body.Close()
	if err != nil {
		err = fmt.Errorf("error reading json reply: %v", err)
		return nil, err
	}

	// Handle unsuccessful HTTP responses
	if httpResponse.StatusCode < 200 || httpResponse.StatusCode >= 300 {
		// Generate a standard error to return if the server body is
		// empty.  This should not happen very often, but it's better
		// than showing nothing in case the target server has a poor
		// implementation.
		if len(respBytes) == 0 {
			return nil, fmt.Errorf("%d %s", httpResponse.StatusCode,
				http.StatusText(httpResponse.StatusCode))
		}
		return nil, fmt.Errorf("%s", respBytes)
	}

	// Unmarshal the response.
	var resp Response
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		return nil, err
	}

	if resp.Error != nil {
		return nil, resp.Error
	}
	return resp.Result, nil
}

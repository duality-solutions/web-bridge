/*
From: https://github.com/bu/gin-access-limit/blob/main/middleware.go
MIT License

Copyright (c) 2018 Buwei Chiu

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
package rest

import (
	"log"
	"net"
	"strings"

	"github.com/gin-gonic/gin"
)

// DisableLogging set up logging. default is false (logging)
var DisableLogging bool

// TrustedHeaderField is a header field that developer trusted in their env.
// e.g. Upstream proxy server's special header that only server can setup
// need to avoid use common forgry-able header fields.
var TrustedHeaderField string

// AllowCIDR is a middleware that check given CIDR rules and return 403 Forbidden
// when user is not coming from allowed source. CIDRs accepts a list of CIDRs,
// separated by comma. (e.g. 127.0.0.1/32, ::1/128 )
func AllowCIDR(CIDRs string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// retreieve user's connection origin from request remote addr
		// need to split the host because original remoteAddr contains port
		remoteAddr, _, splitErr := net.SplitHostPort(c.Request.RemoteAddr)

		if splitErr != nil {
			c.AbortWithError(500, splitErr)
			return
		}

		// if we have Trusted Header Field, and it exists, use it
		if TrustedHeaderField != "" {
			if trustedRemoteAddr := c.GetHeader(TrustedHeaderField); trustedRemoteAddr != "" {
				remoteAddr = trustedRemoteAddr
			}
		}

		// parse it into IP type
		remoteIP := net.ParseIP(remoteAddr)

		// split CIDRs by comma, and we gonna check them one by one
		cidrSlices := strings.Split(CIDRs, ",")

		// under of CIDR we were in
		var matchCount uint

		// go over each CIDR and do the tests
		for _, cidr := range cidrSlices {
			// remove unwanted spaces
			cidr = strings.TrimSpace(cidr)

			// try to parse the CIDR
			_, cidrIPNet, parseCIDRErr := net.ParseCIDR(cidr)

			if parseCIDRErr != nil {
				c.AbortWithError(500, parseCIDRErr)
				return
			}

			// This is the core of this middleware,
			// it ask current CIDR network range to test if current IP is in
			if cidrIPNet.Contains(remoteIP) {
				matchCount = matchCount + 1
			}
		}

		// if no CIDR ranges contains our IP
		if matchCount == 0 {
			if DisableLogging == false {
				log.Printf("[LIMIT] Request from [" + remoteAddr + "] is not allow to access `" + c.Request.RequestURI + "`, only allow from: [" + CIDRs + "]")
			}

			c.AbortWithStatus(403)
			return
		}
	}
}

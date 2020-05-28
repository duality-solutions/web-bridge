/*
Copyright (c) 2012, Samuel Stauffer <samuel@descolada.com>
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright
  notice, this list of conditions and the following disclaimer.
* Redistributions in binary form must reproduce the above copyright
  notice, this list of conditions and the following disclaimer in the
  documentation and/or other materials provided with the distribution.
* Neither the name of the author nor the
  names of its contributors may be used to endorse or promote products
  derived from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL <COPYRIGHT HOLDER> BE LIABLE FOR ANY
DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package socks

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"time"
)

const (
	protocolVersion = 5

	authNone             = 0
	authGssAPI           = 1
	authUsernamePassword = 2
	authUnavailable      = 0xff

	commandTCPConnect = 1

	addressTypeIPv4   = 1
	addressTypeDomain = 3
	addressTypeIPv6   = 4

	statusRequestGranted          = 0
	statusGeneralFailure          = 1
	statusConnectionNotAllowed    = 2
	statusNetworkUnreachable      = 3
	statusHostUnreachable         = 4
	statusConnectionRefused       = 5
	statusTTLExpired              = 6
	statusCommandNotSupport       = 7
	statusAddressTypeNotSupported = 8
)

var (
	// ErrAuthFailed is an authorization error
	ErrAuthFailed = errors.New("authentication failed")
	// ErrInvalidProxyResponse is an invalid proxy error
	ErrInvalidProxyResponse = errors.New("invalid proxy response")
	// ErrNoAcceptableAuthMethod is thrown when there is no acceptable authentication method
	ErrNoAcceptableAuthMethod = errors.New("no acceptable authentication method")

	statusErrors = map[byte]error{
		statusGeneralFailure:          errors.New("general failure"),
		statusConnectionNotAllowed:    errors.New("connection not allowed by ruleset"),
		statusNetworkUnreachable:      errors.New("network unreachable"),
		statusHostUnreachable:         errors.New("host unreachable"),
		statusConnectionRefused:       errors.New("connection refused by destination host"),
		statusTTLExpired:              errors.New("TTL expired"),
		statusCommandNotSupport:       errors.New("command not supported / protocol error"),
		statusAddressTypeNotSupported: errors.New("address type not supported"),
	}
)

// ProxiedAddr stores proxy address information
type ProxiedAddr struct {
	Net  string
	Host string
	Port int
}

// Network stores network for a proxy address
func (a *ProxiedAddr) Network() string {
	return a.Net
}

// String prints the proxy address struct data to string
func (a *ProxiedAddr) String() string {
	return fmt.Sprintf("%s:%d", a.Host, a.Port)
}

type proxiedConn struct {
	conn       net.Conn
	remoteAddr *ProxiedAddr
	boundAddr  *ProxiedAddr
}

func (c *proxiedConn) Read(b []byte) (int, error) {
	return c.conn.Read(b)
}

func (c *proxiedConn) Write(b []byte) (int, error) {
	return c.conn.Write(b)
}

func (c *proxiedConn) Close() error {
	return c.conn.Close()
}

func (c *proxiedConn) LocalAddr() net.Addr {
	if c.boundAddr != nil {
		return c.boundAddr
	}
	return c.conn.LocalAddr()
}

func (c *proxiedConn) RemoteAddr() net.Addr {
	if c.remoteAddr != nil {
		return c.remoteAddr
	}
	return c.conn.RemoteAddr()
}

func (c *proxiedConn) SetDeadline(t time.Time) error {
	return c.conn.SetDeadline(t)
}

func (c *proxiedConn) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

func (c *proxiedConn) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}

// Proxy stores proxy information
type Proxy struct {
	Addr         string
	Username     string
	Password     string
	TorIsolation bool
}

// Dial connects to network and address
func (p *Proxy) Dial(network, addr string) (net.Conn, error) {
	return p.dial(network, addr, 0)
}

func (p *Proxy) dial(network, addr string, timeout time.Duration) (net.Conn, error) {
	host, strPort, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	port, err := strconv.Atoi(strPort)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialTimeout("tcp", p.Addr, timeout)
	if err != nil {
		return nil, err
	}

	var user, pass string
	if p.TorIsolation {
		var b [16]byte
		_, err := io.ReadFull(rand.Reader, b[:])
		if err != nil {
			conn.Close()
			return nil, err
		}
		user = hex.EncodeToString(b[0:8])
		pass = hex.EncodeToString(b[8:16])
	} else {
		user = p.Username
		pass = p.Password
	}
	buf := make([]byte, 32+len(host)+len(user)+len(pass))

	// Initial greeting
	buf[0] = protocolVersion
	if user != "" {
		buf = buf[:4]
		buf[1] = 2 // num auth methods
		buf[2] = authNone
		buf[3] = authUsernamePassword
	} else {
		buf = buf[:3]
		buf[1] = 1 // num auth methods
		buf[2] = authNone
	}

	_, err = conn.Write(buf)
	if err != nil {
		conn.Close()
		return nil, err
	}

	// Server's auth choice

	if _, err := io.ReadFull(conn, buf[:2]); err != nil {
		conn.Close()
		return nil, err
	}
	if buf[0] != protocolVersion {
		conn.Close()
		return nil, ErrInvalidProxyResponse
	}
	err = nil
	switch buf[1] {
	default:
		err = ErrInvalidProxyResponse
	case authUnavailable:
		err = ErrNoAcceptableAuthMethod
	case authGssAPI:
		err = ErrNoAcceptableAuthMethod
	case authUsernamePassword:
		buf = buf[:3+len(user)+len(pass)]
		buf[0] = 1 // version
		buf[1] = byte(len(user))
		copy(buf[2:], user)
		buf[2+len(user)] = byte(len(pass))
		copy(buf[3+len(user):], pass)
		if _, err = conn.Write(buf); err != nil {
			conn.Close()
			return nil, err
		}
		if _, err = io.ReadFull(conn, buf[:2]); err != nil {
			conn.Close()
			return nil, err
		}
		if buf[0] != 1 { // version
			err = ErrInvalidProxyResponse
		} else if buf[1] != 0 { // 0 = succes, else auth failed
			err = ErrAuthFailed
		}
	case authNone:
		// Do nothing
	}
	if err != nil {
		conn.Close()
		return nil, err
	}

	// Command / connection request

	buf = buf[:7+len(host)]
	buf[0] = protocolVersion
	buf[1] = commandTCPConnect
	buf[2] = 0 // reserved
	buf[3] = addressTypeDomain
	buf[4] = byte(len(host))
	copy(buf[5:], host)
	buf[5+len(host)] = byte(port >> 8)
	buf[6+len(host)] = byte(port & 0xff)
	if _, err := conn.Write(buf); err != nil {
		conn.Close()
		return nil, err
	}

	// Server response

	if _, err := io.ReadFull(conn, buf[:4]); err != nil {
		conn.Close()
		return nil, err
	}

	if buf[0] != protocolVersion {
		conn.Close()
		return nil, ErrInvalidProxyResponse
	}

	if buf[1] != statusRequestGranted {
		conn.Close()
		err := statusErrors[buf[1]]
		if err == nil {
			err = ErrInvalidProxyResponse
		}
		return nil, err
	}

	paddr := &ProxiedAddr{Net: network}

	switch buf[3] {
	default:
		conn.Close()
		return nil, ErrInvalidProxyResponse
	case addressTypeIPv4:
		if _, err := io.ReadFull(conn, buf[:4]); err != nil {
			conn.Close()
			return nil, err
		}
		paddr.Host = net.IP(buf).String()
	case addressTypeIPv6:
		if _, err := io.ReadFull(conn, buf[:16]); err != nil {
			conn.Close()
			return nil, err
		}
		paddr.Host = net.IP(buf).String()
	case addressTypeDomain:
		if _, err := io.ReadFull(conn, buf[:1]); err != nil {
			conn.Close()
			return nil, err
		}
		domainLen := buf[0]
		if _, err := io.ReadFull(conn, buf[:domainLen]); err != nil {
			conn.Close()
			return nil, err
		}
		paddr.Host = string(buf[:domainLen])
	}

	if _, err := io.ReadFull(conn, buf[:2]); err != nil {
		conn.Close()
		return nil, err
	}
	paddr.Port = int(buf[0])<<8 | int(buf[1])

	return &proxiedConn{
		conn:       conn,
		boundAddr:  paddr,
		remoteAddr: &ProxiedAddr{network, host, port},
	}, nil
}

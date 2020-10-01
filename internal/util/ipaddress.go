package util

import (
	"net"
	"strings"
)

// IP address lengths (bytes).
const (
	IPv4len = 4
	IPv6len = 16
	// Bigger than we need, not too big to worry about overflow
	bigInt = 0xFFFFFF
)

// ParseIP parses s and returns true if valid
// The string s can be in dotted decimal ("192.0.2.1")
// or IPv6 ("2001:db8::68") form.
// If s is not a valid textual representation of an IP address,
// IsValidIPAddress returns false.
func IsValidIPAddress(s string) bool {
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '.':
			return parseIPv4(s)
		case ':':
			return parseIPv6(s)
		}
	}
	return false
}

// IsValidCIDRList parses s which is a comma delimited list of CIDRs and returns true if all the CIDR in the list are valid
// The string s can "192.0.2.0/24" or "2001:db8::/32", as defined in
// RFC 4632 and RFC 4291.
func IsValidCIDRList(s string) bool {
	cidrList := strings.Split(s, ",")
	for _, raw := range cidrList {
		cidr := strings.Trim(raw, " ")
		_, _, err := net.ParseCIDR(cidr)
		if err != nil {
			return false
		}
	}
	return true
}

func parseIPv4(s string) bool {
	return IsIPv4(s)
	//return net.ParseIP(s) != nil
}

func parseIPv6(s string) bool {
	return net.ParseIP(s) != nil
}

// Decimal to integer.
// Returns number, characters consumed, success.
func dtoi(s string) (n int, i int, ok bool) {
	n = 0
	for i = 0; i < len(s) && '0' <= s[i] && s[i] <= '9'; i++ {
		n = n*10 + int(s[i]-'0')
		if n >= bigInt {
			return bigInt, i, false
		}
	}
	if i == 0 {
		return 0, 0, false
	}
	return n, i, true
}

// Parse IPv4 address (d.d.d.d).
func IsIPv4(s string) bool {
	var p [IPv4len]byte
	for i := 0; i < IPv4len; i++ {
		if len(s) == 0 {
			// Missing octets.
			return false
		}
		if i > 0 {
			if s[0] != '.' {
				return false
			}
			s = s[1:]
		}
		n, c, ok := dtoi(s)
		if !ok || n > 0xFF {
			return false
		}
		s = s[c:]
		p[i] = byte(n)
	}
	if len(s) != 0 {
		return false
	}
	return true
}

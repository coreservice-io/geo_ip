package data

import (
	"errors"
	"net"
)

var privateIPBlocks []*net.IPNet

func init() {
	for _, cidr := range []string{
		"127.0.0.0/8",    // IPv4 loopback
		"10.0.0.0/8",     // RFC1918
		"172.16.0.0/12",  // RFC1918
		"192.168.0.0/16", // RFC1918
		"169.254.0.0/16", // RFC3927 link-local
		"::1/128",        // IPv6 loopback
		"fe80::/10",      // IPv6 link-local
		"fc00::/7",       // IPv6 unique local addr
	} {
		_, block, _ := net.ParseCIDR(cidr)
		privateIPBlocks = append(privateIPBlocks, block)
	}
}

func IsLanIp(ip string) (bool, error) {
	netIp := net.ParseIP(ip)
	if netIp == nil {
		return false, errors.New("ip format error:" + ip)
	}

	if netIp.IsLoopback() || netIp.IsLinkLocalUnicast() || netIp.IsLinkLocalMulticast() {
		return true, nil
	}

	for _, block := range privateIPBlocks {
		if block.Contains(netIp) {
			return true, nil
		}
	}

	return false, nil
}

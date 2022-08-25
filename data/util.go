package data

import (
	"errors"
	"net"
)

func IsLanIp(ip string) (bool, error) {
	netIp := net.ParseIP(ip)
	if netIp == nil {
		return false, errors.New("ip error:" + ip)
	}

	ip4 := netIp.To4()
	if ip4 == nil {
		return false, errors.New("ip error:" + ip)
	}

	if ip4[0] == 10 ||
		(ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31) ||
		(ip4[0] == 192 && ip4[1] == 168) {
		return true, nil
	}
	return false, nil
}

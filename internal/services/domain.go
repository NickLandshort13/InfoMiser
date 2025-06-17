package services

import (
	"fmt"
	"net"
)

func GetIPFromDomain(domain string) (string, error) {
	ips, err := net.LookupIP(domain)
	if err != nil || len(ips) == 0 {
		return "", fmt.Errorf("no IP found")
	}

	for _, ip := range ips {
		v4 := ip.To4()
		if v4 != nil {
			return v4.String(), nil
		}
	}

	return "", fmt.Errorf("no IPv4 found")
}

package utils

import (
	"fmt"
	"net"
	"os"
)

// GetHosts
func GetHosts(cidr string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}
	// remove network address and broadcast address
	// return ips[1 : len(ips)-1], nil
	return ips, nil
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

// GetCIDR get CIDR from ip range
func GetCIDR(start net.IP, end net.IP) string {
	var cidrString string
	maxLen := 32

	for l := maxLen; l >= 0; l-- {
		mask := net.CIDRMask(l, maxLen)
		na := start.Mask(mask)
		n := net.IPNet{IP: na, Mask: mask}

		if n.Contains(end) {
			cidrString = fmt.Sprintf("%v/%v", na, l)
			break
		}
	}

	return cidrString
}

// ResolveName get neighbor ip
func ResolveName(ip string) ([]string, error) {
	// Try to get Neighbor DNS Names
	names, err := net.LookupAddr(ip)
	if err != nil {
		return nil, err
	}
	return names, nil
}

// IsValidPath - Check if a given path is valid
func IsValidPath(fp string) bool {
	// Check if file already exists
	if _, err := os.Stat(fp); err == nil {
		return true
	}

	// Attempt to create it
	var d []byte
	if err := os.WriteFile(fp, d, 0o600); err == nil {
		if err := os.Remove(fp); err != nil {
			return false
		}
		return true
	}

	return false
}

// IsValidIP Validate Input IP
func IsValidIP(ip string) (net.IP, error) {
	netIP := net.ParseIP(ip)
	if netIP == nil {
		return nil, fmt.Errorf("Invalid IP : %s", ip)
	}
	return netIP.To4(), nil
}

func SplitSlice(logs []string, num int) [][]string {
	var divided [][]string

	chunkSize := (len(logs) + num - 1) / num

	for i := 0; i < len(logs); i += chunkSize {
		end := i + chunkSize

		if end > len(logs) {
			end = len(logs)
		}

		divided = append(divided, logs[i:end])
	}

	return divided
}

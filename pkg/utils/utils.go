package utils

import (
	"fmt"
	"net"
	"os"
)

// GetHosts returns all IP addresses in a given CIDR range
func GetHosts(cidr string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []string
	for currentIP := ip.Mask(ipnet.Mask); ipnet.Contains(currentIP); inc(currentIP) {
		ips = append(ips, currentIP.String())
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

// GetCIDR calculates CIDR notation from an IP range
func GetCIDR(start, end net.IP) string {
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

// ResolveName performs reverse DNS lookup for an IP address
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
	if err := os.WriteFile(fp, d, 0644); err == nil {
		_ = os.Remove(fp) // And delete it (ignore error as file may not exist)
		return true
	}

	return false
}

// IsValidIP validates an input IP address
func IsValidIP(ip string) (net.IP, error) {
	netIP := net.ParseIP(ip)
	if netIP == nil {
		return nil, fmt.Errorf("invalid IP: %q", ip)
	}
	return netIP.To4(), nil
}

// SplitSlice divides a slice into num chunks
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

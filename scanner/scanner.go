package scanner

import (
	"fmt"
	"net"

	"bitbucket.org/aminebenseddik/reverse-scan/conf"
)

func Start(config *conf.Config) {

	fmt.Println(config)
}

func loopRange(ip net.IP, end string) {

}

// ResolveName get neighbor ip
func resolveName(ip string) ([]string, error) {
	// Try to get Neighbor DNS Names
	names, err := net.LookupAddr(ip)
	if err != nil {
		return nil, err
	}
	return names, nil
}

func Hosts(cidr string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}
	// remove network address and broadcast address
	return ips[1 : len(ips)-1], nil
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

// Package config handles configuration loading and validation for reverse-scan
package config

import (
	"fmt"
	"net"

	"github.com/amine7536/reverse-scan/pkg/utils"

	"github.com/spf13/cobra"
)

// Config the application's configuration
type Config struct {
	CIDR    string
	CSV     string
	StartIP net.IP
	EndIP   net.IP
	WORKERS int
}

// LoadConfig loads the config from a file if specified, otherwise from the environment
func LoadConfig(cmd *cobra.Command) (*Config, error) {

	start, err := cmd.Flags().GetString("start")
	if err != nil {
		return nil, err
	}

	end, err := cmd.Flags().GetString("end")
	if err != nil {
		return nil, err
	}

	cidr, err := cmd.Flags().GetString("cidr")
	if err != nil {
		return nil, err
	}

	output, err := cmd.Flags().GetString("output")
	if err != nil {
		return nil, err
	}

	workers, err := cmd.Flags().GetInt("workers")
	if err != nil {
		return nil, err
	}

	config, err := validateConfig(start, end, cidr, output, workers)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func validateConfig(start, end, cidr, output string, workers int) (*Config, error) {
	// Check that either CIDR or (start and end) are provided, but not both
	hasCIDR := cidr != ""
	hasStartEnd := start != "" || end != ""

	if !hasCIDR && !hasStartEnd {
		return nil, fmt.Errorf("must specify either --cidr or --start/--end range")
	}

	if hasCIDR && hasStartEnd {
		return nil, fmt.Errorf("cannot specify both --cidr and --start/--end range")
	}

	if output == "" {
		return nil, fmt.Errorf("must specify output file")
	}

	if !utils.IsValidPath(output) {
		return nil, fmt.Errorf("invalid output file: %q", output)
	}

	var startIP, endIP net.IP
	var cidrStr string

	if hasCIDR {
		// Parse CIDR notation
		ip, ipnet, err := net.ParseCIDR(cidr)
		if err != nil {
			return nil, fmt.Errorf("invalid CIDR notation %q: %w", cidr, err)
		}

		// Get the first IP in the range
		startIP = ip.Mask(ipnet.Mask)

		// Get the last IP in the range
		endIP = make(net.IP, len(startIP))
		copy(endIP, startIP)
		for i := range endIP {
			endIP[i] |= ^ipnet.Mask[i]
		}

		cidrStr = cidr
	} else {
		// Validate start and end IPs
		if start == "" {
			return nil, fmt.Errorf("must specify start range")
		}

		if end == "" {
			return nil, fmt.Errorf("must specify end range")
		}

		var err error
		startIP, err = utils.IsValidIP(start)
		if err != nil {
			return nil, err
		}

		endIP, err = utils.IsValidIP(end)
		if err != nil {
			return nil, err
		}

		if startIP[0] != endIP[0] {
			return nil, fmt.Errorf("invalid range: start and end IP must be in the same network")
		}

		if endIP[2] < startIP[2] {
			return nil, fmt.Errorf("invalid range: end IP must be greater than start IP")
		}

		cidrStr = utils.GetCIDR(startIP, endIP)
	}

	config := Config{
		StartIP: startIP,
		EndIP:   endIP,
		CIDR:    cidrStr,
		CSV:     output,
		WORKERS: workers,
	}

	return &config, nil
}

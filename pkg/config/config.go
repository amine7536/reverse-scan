package config

import (
	"fmt"
	"net"

	"github.com/amine7536/reverse-scan/pkg/utils"

	"github.com/spf13/cobra"
)

// Config the application's configuration
type Config struct {
	StartIP net.IP
	EndIP   net.IP
	CIDR    string
	CSV     string
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

	output, err := cmd.Flags().GetString("output")
	if err != nil {
		return nil, err
	}

	workers, err := cmd.Flags().GetInt("workers")
	if err != nil {
		return nil, err
	}

	config, err := validateConfig(start, end, output, workers)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func validateConfig(start, end, output string, workers int) (*Config, error) {
	if start == "" {
		return nil, fmt.Errorf("must specify start range")
	}

	if end == "" {
		return nil, fmt.Errorf("must specify end range")
	}

	if output == "" {
		return nil, fmt.Errorf("must specify output file")
	}

	startIP, err := utils.IsValidIP(start)
	if err != nil {
		return nil, err
	}

	endIP, err := utils.IsValidIP(end)
	if err != nil {
		return nil, err
	}

	if startIP[0] != endIP[0] {
		return nil, fmt.Errorf("invalid range: start and end IP must be in the same network")
	}

	if endIP[2] < startIP[2] {
		return nil, fmt.Errorf("invalid range: end IP must be greater than start IP")
	}

	if !utils.IsValidPath(output) {
		return nil, fmt.Errorf("invalid output file: %q", output)
	}

	config := Config{
		StartIP: startIP,
		EndIP:   endIP,
		CIDR:    utils.GetCIDR(startIP, endIP),
		CSV:     output,
		WORKERS: workers,
	}

	return &config, nil
}

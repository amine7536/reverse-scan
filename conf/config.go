package conf

import (
	"fmt"
	"net"

	"bitbucket.org/aminebenseddik/reverse-scan/utils"

	"github.com/spf13/cobra"
)

// Config the application's configuration
type Config struct {
	StartIP net.IP
	EndIP   net.IP
	CIDR    string
	CSV     string
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

	config, err := validateConfig(start, end, output)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func validateConfig(start string, end string, output string) (*Config, error) {
	if start == "" {
		return nil, fmt.Errorf("Must specify start range")
	}

	if end == "" {
		return nil, fmt.Errorf("Must specify end range")
	}

	if output == "" {
		return nil, fmt.Errorf("Must specify output file")
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
		return nil, fmt.Errorf("Invalid Range")
	}

	if endIP[2] < startIP[2] {
		return nil, fmt.Errorf("Invalid Range")
	}

	if !utils.IsValidPath(output) {
		return nil, fmt.Errorf("Invalid output file : %s", output)
	}

	config := Config{
		StartIP: startIP,
		EndIP:   endIP,
		CIDR:    utils.GetCIDR(startIP, endIP),
		CSV:     output,
	}

	return &config, nil
}

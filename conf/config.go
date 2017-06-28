package conf

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"

	"github.com/spf13/cobra"
)

// Config the application's configuration
type Config struct {
	StartIP net.IP
	EndIP   net.IP
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

	startIP, err := IsValidIP(start)
	if err != nil {
		return nil, err
	}

	endIP, err := IsValidIP(end)
	if err != nil {
		return nil, err
	}

	if startIP[0] != endIP[0] {
		return nil, fmt.Errorf("Invalid Range")
	}

	if startIP[1] != endIP[1] {
		return nil, fmt.Errorf("Invalid Range")
	}

	if endIP[2] < startIP[2] {
		return nil, fmt.Errorf("Invalid Range")
	}

	if !IsValidPath(output) {
		return nil, fmt.Errorf("Invalid output file : %s", output)
	}

	config := Config{
		StartIP: startIP,
		EndIP:   endIP,
		CSV:     output,
	}

	return &config, nil
}

// IsValidPath - Check if a given path is valid
func IsValidPath(fp string) bool {
	// Check if file already exists
	if _, err := os.Stat(fp); err == nil {
		return true
	}

	// Attempt to create it
	var d []byte
	if err := ioutil.WriteFile(fp, d, 0644); err == nil {
		os.Remove(fp) // And delete it
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

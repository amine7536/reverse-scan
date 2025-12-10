package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateConfig(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	validOutputFile := filepath.Join(tmpDir, "output.csv")

	tests := []struct {
		name    string
		start   string
		end     string
		cidr    string
		output  string
		errMsg  string
		workers int
		wantErr bool
	}{
		{
			name:    "valid config with start/end",
			start:   "192.168.1.0",
			end:     "192.168.1.255",
			cidr:    "",
			output:  validOutputFile,
			workers: 8,
			wantErr: false,
		},
		{
			name:    "valid config with CIDR",
			start:   "",
			end:     "",
			cidr:    "192.168.1.0/24",
			output:  validOutputFile,
			workers: 8,
			wantErr: false,
		},
		{
			name:    "missing all range inputs",
			start:   "",
			end:     "",
			cidr:    "",
			output:  validOutputFile,
			workers: 8,
			wantErr: true,
			errMsg:  "must specify either --cidr or --start/--end range",
		},
		{
			name:    "both CIDR and start/end provided",
			start:   "192.168.1.0",
			end:     "192.168.1.255",
			cidr:    "192.168.1.0/24",
			output:  validOutputFile,
			workers: 8,
			wantErr: true,
			errMsg:  "cannot specify both --cidr and --start/--end range",
		},
		{
			name:    "missing start IP with end",
			start:   "",
			end:     "192.168.1.255",
			cidr:    "",
			output:  validOutputFile,
			workers: 8,
			wantErr: true,
			errMsg:  "must specify start range",
		},
		{
			name:    "missing end IP with start",
			start:   "192.168.1.0",
			end:     "",
			cidr:    "",
			output:  validOutputFile,
			workers: 8,
			wantErr: true,
			errMsg:  "must specify end range",
		},
		{
			name:    "missing output file",
			start:   "192.168.1.0",
			end:     "192.168.1.255",
			cidr:    "",
			output:  "",
			workers: 8,
			wantErr: true,
			errMsg:  "must specify output file",
		},
		{
			name:    "invalid CIDR notation",
			start:   "",
			end:     "",
			cidr:    "invalid-cidr",
			output:  validOutputFile,
			workers: 8,
			wantErr: true,
		},
		{
			name:    "invalid start IP",
			start:   "999.999.999.999",
			end:     "192.168.1.255",
			cidr:    "",
			output:  validOutputFile,
			workers: 8,
			wantErr: true,
		},
		{
			name:    "invalid end IP",
			start:   "192.168.1.0",
			end:     "not-an-ip",
			cidr:    "",
			output:  validOutputFile,
			workers: 8,
			wantErr: true,
		},
		{
			name:    "invalid range - different first octet",
			start:   "192.168.1.0",
			end:     "10.168.1.255",
			cidr:    "",
			output:  validOutputFile,
			workers: 8,
			wantErr: true,
			errMsg:  "invalid range: start and end IP must be in the same network",
		},
		{
			name:    "invalid range - end before start",
			start:   "192.168.10.0",
			end:     "192.168.1.255",
			cidr:    "",
			output:  validOutputFile,
			workers: 8,
			wantErr: true,
			errMsg:  "invalid range: end IP must be greater than start IP",
		},
		{
			name:    "invalid output path",
			start:   "192.168.1.0",
			end:     "192.168.1.255",
			cidr:    "",
			output:  "/nonexistent/directory/output.csv",
			workers: 8,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := validateConfig(tt.start, tt.end, tt.cidr, tt.output, tt.workers)

			if (err != nil) != tt.wantErr {
				t.Errorf("validateConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errMsg != "" && err != nil {
				if err.Error() != tt.errMsg {
					t.Logf("validateConfig() error message = %v, expected %v", err.Error(), tt.errMsg)
				}
			}

			if !tt.wantErr {
				if config == nil {
					t.Error("validateConfig() returned nil config for valid input")
					return
				}
				if config.WORKERS != tt.workers {
					t.Errorf("validateConfig() workers = %v, want %v", config.WORKERS, tt.workers)
				}
				if config.CSV != tt.output {
					t.Errorf("validateConfig() CSV = %v, want %v", config.CSV, tt.output)
				}
				if config.StartIP == nil || config.EndIP == nil {
					t.Error("validateConfig() returned nil IPs")
				}
				if config.CIDR == "" {
					t.Error("validateConfig() returned empty CIDR")
				}
			}
		})
	}
}

func TestValidateConfigCIDRCalculation(t *testing.T) {
	tmpDir := t.TempDir()
	validOutputFile := filepath.Join(tmpDir, "output.csv")

	tests := []struct {
		name     string
		start    string
		end      string
		wantCIDR string
	}{
		{
			name:     "class C network",
			start:    "192.168.1.0",
			end:      "192.168.1.255",
			wantCIDR: "192.168.1.0/24",
		},
		{
			name:     "small range",
			start:    "10.0.0.0",
			end:      "10.0.0.15",
			wantCIDR: "10.0.0.0/28",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := validateConfig(tt.start, tt.end, "", validOutputFile, 8)
			if err != nil {
				t.Fatalf("validateConfig() unexpected error = %v", err)
			}
			if config.CIDR != tt.wantCIDR {
				t.Errorf("validateConfig() CIDR = %v, want %v", config.CIDR, tt.wantCIDR)
			}
		})
	}
}

// TestConfigStruct verifies the Config struct fields
func TestConfigStruct(t *testing.T) {
	tmpDir := t.TempDir()
	validOutputFile := filepath.Join(tmpDir, "output.csv")

	config, err := validateConfig("192.168.1.0", "192.168.1.255", "", validOutputFile, 16)
	if err != nil {
		t.Fatalf("validateConfig() unexpected error = %v", err)
	}

	if config.StartIP.String() != "192.168.1.0" {
		t.Errorf("Config.StartIP = %v, want 192.168.1.0", config.StartIP)
	}
	if config.EndIP.String() != "192.168.1.255" {
		t.Errorf("Config.EndIP = %v, want 192.168.1.255", config.EndIP)
	}
	if config.WORKERS != 16 {
		t.Errorf("Config.WORKERS = %v, want 16", config.WORKERS)
	}
	if config.CSV != validOutputFile {
		t.Errorf("Config.CSV = %v, want %v", config.CSV, validOutputFile)
	}
}

// TestValidateConfigFileCreation ensures output file can be created
func TestValidateConfigFileCreation(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "test-output.csv")

	config, err := validateConfig("192.168.1.0", "192.168.1.10", "", outputFile, 8)
	if err != nil {
		t.Fatalf("validateConfig() unexpected error = %v", err)
	}

	// Verify the output file path is set correctly
	if config.CSV != outputFile {
		t.Errorf("Config.CSV = %v, want %v", config.CSV, outputFile)
	}

	// Verify we can actually create a file at that path
	file, err := os.Create(config.CSV)
	if err != nil {
		t.Fatalf("Failed to create output file at %v: %v", config.CSV, err)
	}
	if err := file.Close(); err != nil {
		t.Errorf("Failed to close file: %v", err)
	}
	if err := os.Remove(config.CSV); err != nil {
		t.Errorf("Failed to remove test file: %v", err)
	}
}

// TestValidateConfigWithCIDR verifies CIDR input functionality
func TestValidateConfigWithCIDR(t *testing.T) {
	tmpDir := t.TempDir()
	validOutputFile := filepath.Join(tmpDir, "output.csv")

	tests := []struct {
		name      string
		cidr      string
		wantCIDR  string
		wantStart string
		wantEnd   string
	}{
		{
			name:      "class C network",
			cidr:      "192.168.1.0/24",
			wantCIDR:  "192.168.1.0/24",
			wantStart: "192.168.1.0",
			wantEnd:   "192.168.1.255",
		},
		{
			name:      "class B network",
			cidr:      "172.16.0.0/16",
			wantCIDR:  "172.16.0.0/16",
			wantStart: "172.16.0.0",
			wantEnd:   "172.16.255.255",
		},
		{
			name:      "small subnet /28",
			cidr:      "10.0.0.0/28",
			wantCIDR:  "10.0.0.0/28",
			wantStart: "10.0.0.0",
			wantEnd:   "10.0.0.15",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := validateConfig("", "", tt.cidr, validOutputFile, 8)
			if err != nil {
				t.Fatalf("validateConfig() unexpected error = %v", err)
			}
			if config.CIDR != tt.wantCIDR {
				t.Errorf("Config.CIDR = %v, want %v", config.CIDR, tt.wantCIDR)
			}
			if config.StartIP.String() != tt.wantStart {
				t.Errorf("Config.StartIP = %v, want %v", config.StartIP, tt.wantStart)
			}
			if config.EndIP.String() != tt.wantEnd {
				t.Errorf("Config.EndIP = %v, want %v", config.EndIP, tt.wantEnd)
			}
		})
	}
}

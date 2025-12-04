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
		output  string
		workers int
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid config",
			start:   "192.168.1.0",
			end:     "192.168.1.255",
			output:  validOutputFile,
			workers: 8,
			wantErr: false,
		},
		{
			name:    "missing start IP",
			start:   "",
			end:     "192.168.1.255",
			output:  validOutputFile,
			workers: 8,
			wantErr: true,
			errMsg:  "Must specify start range",
		},
		{
			name:    "missing end IP",
			start:   "192.168.1.0",
			end:     "",
			output:  validOutputFile,
			workers: 8,
			wantErr: true,
			errMsg:  "Must specify end range",
		},
		{
			name:    "missing output file",
			start:   "192.168.1.0",
			end:     "192.168.1.255",
			output:  "",
			workers: 8,
			wantErr: true,
			errMsg:  "Must specify output file",
		},
		{
			name:    "invalid start IP",
			start:   "999.999.999.999",
			end:     "192.168.1.255",
			output:  validOutputFile,
			workers: 8,
			wantErr: true,
		},
		{
			name:    "invalid end IP",
			start:   "192.168.1.0",
			end:     "not-an-ip",
			output:  validOutputFile,
			workers: 8,
			wantErr: true,
		},
		{
			name:    "invalid range - different first octet",
			start:   "192.168.1.0",
			end:     "10.168.1.255",
			output:  validOutputFile,
			workers: 8,
			wantErr: true,
			errMsg:  "Invalid Range",
		},
		{
			name:    "invalid range - end before start",
			start:   "192.168.10.0",
			end:     "192.168.1.255",
			output:  validOutputFile,
			workers: 8,
			wantErr: true,
			errMsg:  "Invalid Range",
		},
		{
			name:    "invalid output path",
			start:   "192.168.1.0",
			end:     "192.168.1.255",
			output:  "/nonexistent/directory/output.csv",
			workers: 8,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := validateConfig(tt.start, tt.end, tt.output, tt.workers)

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
			config, err := validateConfig(tt.start, tt.end, validOutputFile, 8)
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

	config, err := validateConfig("192.168.1.0", "192.168.1.255", validOutputFile, 16)
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

	config, err := validateConfig("192.168.1.0", "192.168.1.10", outputFile, 8)
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
	file.Close()
	os.Remove(config.CSV)
}

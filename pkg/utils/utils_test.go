package utils

import (
	"net"
	"os"
	"path/filepath"
	"testing"
)

func TestIsValidIP(t *testing.T) {
	tests := []struct {
		name    string
		ip      string
		wantErr bool
	}{
		{
			name:    "valid IPv4",
			ip:      "192.168.1.1",
			wantErr: false,
		},
		{
			name:    "valid IPv4 with zeros",
			ip:      "10.0.0.0",
			wantErr: false,
		},
		{
			name:    "invalid IP - empty",
			ip:      "",
			wantErr: true,
		},
		{
			name:    "invalid IP - malformed",
			ip:      "256.256.256.256",
			wantErr: true,
		},
		{
			name:    "invalid IP - not an IP",
			ip:      "not-an-ip",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ip, err := IsValidIP(tt.ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsValidIP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && ip == nil {
				t.Errorf("IsValidIP() returned nil IP for valid input")
			}
		})
	}
}

func TestGetCIDR(t *testing.T) {
	tests := []struct {
		name  string
		start string
		end   string
		want  string
	}{
		{
			name:  "class C network",
			start: "192.168.1.0",
			end:   "192.168.1.255",
			want:  "192.168.1.0/24",
		},
		{
			name:  "class B network",
			start: "172.16.0.0",
			end:   "172.16.255.255",
			want:  "172.16.0.0/16",
		},
		{
			name:  "small range",
			start: "10.0.0.0",
			end:   "10.0.0.15",
			want:  "10.0.0.0/28",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start := net.ParseIP(tt.start)
			end := net.ParseIP(tt.end)
			got := GetCIDR(start, end)
			if got != tt.want {
				t.Errorf("GetCIDR() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetHosts(t *testing.T) {
	tests := []struct {
		name      string
		cidr      string
		wantCount int
		wantErr   bool
	}{
		{
			name:      "class C /24",
			cidr:      "192.168.1.0/24",
			wantCount: 256,
			wantErr:   false,
		},
		{
			name:      "small range /30",
			cidr:      "10.0.0.0/30",
			wantCount: 4,
			wantErr:   false,
		},
		{
			name:      "single IP /32",
			cidr:      "192.168.1.1/32",
			wantCount: 1,
			wantErr:   false,
		},
		{
			name:      "invalid CIDR",
			cidr:      "invalid",
			wantCount: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hosts, err := GetHosts(tt.cidr)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetHosts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(hosts) != tt.wantCount {
				t.Errorf("GetHosts() returned %d hosts, want %d", len(hosts), tt.wantCount)
			}
		})
	}
}

func TestSplitSlice(t *testing.T) {
	tests := []struct {
		name      string
		logs      []string
		num       int
		wantParts int
	}{
		{
			name:      "split 10 items into 3 parts",
			logs:      []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"},
			num:       3,
			wantParts: 3,
		},
		{
			name:      "split 5 items into 2 parts",
			logs:      []string{"1", "2", "3", "4", "5"},
			num:       2,
			wantParts: 2,
		},
		{
			name:      "split 3 items into 5 parts",
			logs:      []string{"1", "2", "3"},
			num:       5,
			wantParts: 3,
		},
		{
			name:      "split empty slice",
			logs:      []string{},
			num:       2,
			wantParts: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SplitSlice(tt.logs, tt.num)
			if len(got) != tt.wantParts {
				t.Errorf("SplitSlice() returned %d parts, want %d", len(got), tt.wantParts)
			}

			// Verify all items are accounted for
			totalItems := 0
			for _, part := range got {
				totalItems += len(part)
			}
			if totalItems != len(tt.logs) {
				t.Errorf("SplitSlice() lost items: got %d total items, want %d", totalItems, len(tt.logs))
			}
		})
	}
}

func TestIsValidPath(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()

	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			name: "valid path in temp dir",
			path: filepath.Join(tmpDir, "test.txt"),
			want: true,
		},
		{
			name: "invalid path - nonexistent directory",
			path: "/nonexistent/directory/that/does/not/exist/file.txt",
			want: false,
		},
		{
			name: "valid path - existing file",
			path: filepath.Join(tmpDir, "existing.txt"),
			want: true,
		},
	}

	// Create an existing file for the test
	existingFile := filepath.Join(tmpDir, "existing.txt")
	if err := os.WriteFile(existingFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidPath(tt.path)
			if got != tt.want {
				t.Errorf("IsValidPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResolveName(t *testing.T) {
	// Test with localhost which should always resolve
	t.Run("resolve localhost", func(t *testing.T) {
		_, err := ResolveName("127.0.0.1")
		// We don't check the result because it depends on the system configuration
		// We just verify the function doesn't panic and returns properly
		if err != nil {
			// It's ok if it errors - not all systems have reverse DNS for localhost
			t.Logf("ResolveName returned error (expected on some systems): %v", err)
		}
	})

	t.Run("resolve invalid IP", func(t *testing.T) {
		names, err := ResolveName("999.999.999.999")
		if err == nil && len(names) > 0 {
			t.Errorf("ResolveName() should error on invalid IP")
		}
	})
}

func TestInc(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "increment simple IP",
			input:    "192.168.1.1",
			expected: "192.168.1.2",
		},
		{
			name:     "increment with overflow to next octet",
			input:    "192.168.1.255",
			expected: "192.168.2.0",
		},
		{
			name:     "increment zero",
			input:    "0.0.0.0",
			expected: "0.0.0.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ip := net.ParseIP(tt.input).To4()
			inc(ip)
			got := ip.String()
			if got != tt.expected {
				t.Errorf("inc() = %v, want %v", got, tt.expected)
			}
		})
	}
}

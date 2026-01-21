package common_test

import (
	"os"
	"runtime"
	"testing"

	"github.com/kaichao/scalebox/pkg/common"
)

func TestGetLocalIP(t *testing.T) {
	// Save original environment variables
	originalLocalIP := os.Getenv("LOCAL_IP")
	originalLocalIPIndex := os.Getenv("LOCAL_IP_INDEX")
	defer func() {
		// Restore environment variables
		os.Setenv("LOCAL_IP", originalLocalIP)
		os.Setenv("LOCAL_IP_INDEX", originalLocalIPIndex)
	}()

	testCases := []struct {
		name           string
		localIP        string
		localIPIndex   string
		expectedPrefix string
		skipOnOS       string
	}{
		{
			name:           "Test with LOCAL_IP environment variable",
			localIP:        "192.168.1.100",
			localIPIndex:   "",
			expectedPrefix: "192.168.1.100",
		},
		{
			name:           "Test without environment variables on current OS",
			localIP:        "",
			localIPIndex:   "",
			expectedPrefix: "", // Will be validated based on actual output
		},
		{
			name:           "Test with LOCAL_IP_INDEX=2 on Linux",
			localIP:        "",
			localIPIndex:   "2",
			expectedPrefix: "",       // Will be validated based on actual output
			skipOnOS:       "darwin", // Skip on macOS as it uses ifconfig
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Skip test if specified for current OS
			if tc.skipOnOS != "" && runtime.GOOS == tc.skipOnOS {
				t.Skipf("Skipping test on %s", runtime.GOOS)
			}

			// Set environment variables
			os.Setenv("LOCAL_IP", tc.localIP)
			os.Setenv("LOCAL_IP_INDEX", tc.localIPIndex)

			// Call the function
			ip := common.GetLocalIP()

			// Validate the result
			if tc.expectedPrefix != "" {
				if ip != tc.expectedPrefix {
					t.Errorf("GetLocalIP() = %v, want %v", ip, tc.expectedPrefix)
				}
			} else {
				// Basic validation for IP address format
				if ip == "" {
					t.Error("GetLocalIP() returned empty string")
				}
				// Check if it's a valid IP address format (basic check)
				// Should be in format x.x.x.x where x is 0-255
				parts := 0
				for i := 0; i < len(ip); i++ {
					if ip[i] == '.' {
						parts++
					}
				}
				if parts != 3 {
					t.Errorf("GetLocalIP() returned invalid IP format: %v, expected x.x.x.x", ip)
				}
				// Check it's not loopback (unless that's what we expect)
				if ip == "127.0.0.1" && tc.localIP == "" {
					// This might be acceptable in some test environments
					t.Logf("GetLocalIP() returned loopback address: %v", ip)
				}
			}
		})
	}
}

func TestGetLocalIP_EdgeCases(t *testing.T) {
	// Save original environment variables
	originalLocalIP := os.Getenv("LOCAL_IP")
	originalLocalIPIndex := os.Getenv("LOCAL_IP_INDEX")
	defer func() {
		// Restore environment variables
		os.Setenv("LOCAL_IP", originalLocalIP)
		os.Setenv("LOCAL_IP_INDEX", originalLocalIPIndex)
	}()

	// Test with invalid LOCAL_IP (should still return it as-is)
	os.Setenv("LOCAL_IP", "invalid-ip")
	ip := common.GetLocalIP()
	if ip != "invalid-ip" {
		t.Errorf("GetLocalIP() with invalid LOCAL_IP = %v, want 'invalid-ip'", ip)
	}

	// Test with empty environment variables
	os.Setenv("LOCAL_IP", "")
	os.Setenv("LOCAL_IP_INDEX", "")

	// This should return a valid IP or 127.0.0.1
	ip = common.GetLocalIP()
	if ip == "" {
		t.Error("GetLocalIP() returned empty string with empty env vars")
	}
}

func TestGetLocalIP_OSSpecific(t *testing.T) {
	// This test validates that the function works on different OSes
	// without actually checking the exact IP (which varies by system)

	originalLocalIP := os.Getenv("LOCAL_IP")
	originalLocalIPIndex := os.Getenv("LOCAL_IP_INDEX")
	defer func() {
		os.Setenv("LOCAL_IP", originalLocalIP)
		os.Setenv("LOCAL_IP_INDEX", originalLocalIPIndex)
	}()

	// Clear environment variables
	os.Setenv("LOCAL_IP", "")
	os.Setenv("LOCAL_IP_INDEX", "")

	ip := common.GetLocalIP()

	// Basic validation
	if ip == "" {
		t.Error("GetLocalIP() returned empty string")
	}

	// Log the result for debugging
	t.Logf("GetLocalIP() on %s returned: %s", runtime.GOOS, ip)

	// Check if it's a loopback address (acceptable in test environments)
	if ip == "127.0.0.1" {
		t.Log("GetLocalIP() returned loopback address, which may be acceptable in test environment")
	}
}

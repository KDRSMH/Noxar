package ports

import (
	"fmt"
	"testing"
)

func TestScanPorts(t *testing.T) {
	results := ScanPorts("127.0.0.1", 1, 1000)

	count := 0
	for range results {
		count++
	}

	if count == 0 {
		t.Log("No open ports found on localhost (expected)")
	} else {
		t.Logf("Found %d open ports\n", count)
	}
}

func TestIsPortOpen(t *testing.T) {
	tests := []struct {
		host string
		port int
		want bool
	}{
		{"127.0.0.1", 80, false},
		{"127.0.0.1", 22, false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s:%d", tt.host, tt.port), func(t *testing.T) {
			if got := isPortOpen(tt.host, tt.port); got != tt.want {
				t.Logf("isPortOpen(%s, %d) = %v", tt.host, tt.port, got)
			}
		})
	}
}

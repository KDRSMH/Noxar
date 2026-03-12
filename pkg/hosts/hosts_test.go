package hosts

import (
	"testing"
)

func TestGenerateIPRange(t *testing.T) {
	ips, err := GenerateIPRange("127.0.0.1/32")
	if err != nil {
		t.Fatalf("GenerateIPRange failed: %v", err)
	}

	if len(ips) == 0 {
		t.Fatal("Expected IPs, got empty list")
	}

	t.Logf("Generated %d IPs\n", len(ips))
}

func TestIsHostOnline(t *testing.T) {
	result := IsHostOnline("127.0.0.1", 8080)
	t.Logf("127.0.0.1:8080 online: %v", result)
}

func TestScanNetworkRange(t *testing.T) {
	results := ScanNetworkRange("127.0.0.1/32", 80)

	count := 0
	for range results {
		count++
	}

	t.Logf("Scanned network, found %d results\n", count)
}

package dns

import (
	"testing"
)

func TestReverseLookupParallel(t *testing.T) {
	ips := []string{"8.8.8.8"}
	results := ReverseLookupParallel(ips)

	count := 0
	for result := range results {
		count++
		t.Logf("Reverse lookup: %s -> %s (error: %v)\n", result.Query, result.Result, result.Err)
	}

	if count == 0 {
		t.Fatal("Expected at least 1 result")
	}
}

func TestForwardLookupsParallel(t *testing.T) {
	domains := []string{"google.com"}
	results := ForwardLookupsParallel(domains)

	count := 0
	for result := range results {
		count++
		t.Logf("Forward lookup: %s -> %s (error: %v)\n", result.Query, result.Result, result.Err)
	}

	if count == 0 {
		t.Fatal("Expected at least 1 result")
	}
}

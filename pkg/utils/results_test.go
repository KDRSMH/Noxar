package utils

import (
	"testing"
)

func TestFilterOnlineResults(t *testing.T) {
	results := []ReconResult{
		{Type: ResultTypePort, Data: "result1", Error: nil},
		{Type: ResultTypePort, Data: "result2", Error: nil},
		{Type: ResultTypeHost, Data: nil, Error: nil},
	}

	filtered := FilterOnlineResults(results)
	if len(filtered) != 3 {
		t.Fatalf("Expected 3 results, got %d", len(filtered))
	}
}

func TestToJSON(t *testing.T) {
	results := []ReconResult{
		{Type: ResultTypePort, Data: "8080", Error: nil},
		{Type: ResultTypeHost, Data: "192.168.1.1", Error: nil},
	}

	jsonStr, err := ToJSON(results)
	if err != nil {
		t.Fatalf("ToJSON failed: %v", err)
	}

	if len(jsonStr) == 0 {
		t.Fatal("Expected JSON string, got empty")
	}

	t.Logf("JSON output length: %d\n", len(jsonStr))
}

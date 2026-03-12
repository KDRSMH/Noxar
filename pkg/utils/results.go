package utils

import (
	"encoding/json"
	"fmt"
	"time"
)

type ResultType int

const (
	ResultTypePort ResultType = iota + 1
	ResultTypeHost
	ResultTypeDNS
)

type ReconResult struct {
	Type      ResultType
	Timestamp time.Time
	Data      interface{}
	Error     error
}

type ReconSummary struct {
	TotalResults int
	SuccessCount int
	FailureCount int
	Duration     time.Duration
	StartTime    time.Time
	EndTime      time.Time
}

func MergeResults(portChan, hostChan, dnsChan interface{}) []ReconResult {
	var results []ReconResult
	return results
}

func FilterOnlineResults(results []ReconResult) []ReconResult {
	var filtered []ReconResult
	for _, result := range results {
		if result.Error == nil {
			filtered = append(filtered, result)
		}
	}
	return filtered
}

func PrintResults(results []ReconResult) {
	for _, result := range results {
		if result.Error != nil {
			fmt.Printf("❌ [%s] Hata: %v\n", resultTypeString(result.Type), result.Error)
		} else {
			fmt.Printf("✓ [%s] %+v\n", resultTypeString(result.Type), result.Data)
		}
	}
}

func ToJSON(results []ReconResult) (string, error) {
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return "", fmt.Errorf("JSON marshalinde hata: %w", err)
	}
	return string(data), nil
}

func resultTypeString(t ResultType) string {
	switch t {
	case ResultTypePort:
		return "PORT"
	case ResultTypeHost:
		return "HOST"
	case ResultTypeDNS:
		return "DNS"
	default:
		return "UNKNOWN"
	}
}

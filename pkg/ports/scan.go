package ports

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type PortScanResult struct {
	Host string
	Port int
	Open bool
}

const MaxConcurrentScans = 100
const ScanTimeout = 1500 * time.Millisecond
const MaxRetries = 2

func ScanPorts(host string, startPort, endPort int) <-chan PortScanResult {
	results := make(chan PortScanResult, 100)
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, MaxConcurrentScans)

	for port := startPort; port <= endPort; port++ {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			open := isPortOpen(host, p)
			results <- PortScanResult{Host: host, Port: p, Open: open}
		}(port)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}

func isPortOpen(host string, port int) bool {
	address := net.JoinHostPort(host, fmt.Sprintf("%d", port))

	for attempt := 0; attempt < MaxRetries; attempt++ {
		conn, err := net.DialTimeout("tcp", address, ScanTimeout)
		if err == nil {
			conn.Close()
			return true
		}

		if attempt < MaxRetries-1 {
			time.Sleep(50 * time.Millisecond)
		}
	}

	return false
}

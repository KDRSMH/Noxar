package dns

import (
	"fmt"
	"net"
	"sync"
)

type DNsLookupResult struct {
	Query  string
	Result string
	Err    error
}

func ReverseLookupParallel(ips []string) <-chan DNsLookupResult {
	results := make(chan DNsLookupResult, len(ips)*10)
	var wg sync.WaitGroup

	for _, ip := range ips {
		wg.Add(1)
		go func(query string) {
			defer wg.Done()

			hostnames, err := net.LookupAddr(query)

			if err != nil {
				results <- DNsLookupResult{Query: query, Result: "", Err: err}
			} else if len(hostnames) > 0 {
				for _, hostname := range hostnames {
					results <- DNsLookupResult{Query: query, Result: hostname, Err: nil}
				}
			} else {
				results <- DNsLookupResult{Query: query, Result: "", Err: fmt.Errorf("no hostnames found")}
			}
		}(ip)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}

func ForwardLookupsParallel(domains []string) <-chan DNsLookupResult {
	results := make(chan DNsLookupResult, len(domains)*10)
	var wg sync.WaitGroup

	for _, domain := range domains {
		wg.Add(1)
		go func(query string) {
			defer wg.Done()

			ips, err := net.LookupHost(query)

			if err != nil {
				results <- DNsLookupResult{Query: query, Result: "", Err: err}
			} else if len(ips) > 0 {
				for _, ip := range ips {
					results <- DNsLookupResult{Query: query, Result: ip, Err: nil}
				}
			} else {
				results <- DNsLookupResult{Query: query, Result: "", Err: fmt.Errorf("no IPs found")}
			}
		}(domain)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}

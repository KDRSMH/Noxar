package hosts

import (
"fmt"
"net"
"sync"
"time"
)

type HostDiscoveryResult struct {
IP   string
Port int
Err  error
}

func GenerateIPRange(cidr string) ([]string, error) {
ip, ipnet, err := net.ParseCIDR(cidr)
if err != nil {
return nil, fmt.Errorf("invalid CIDR: %w", err)
}

var ips []string

for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); incrementIP(ip) {
ips = append(ips, ip.String())
}

return ips, nil
}

func incrementIP(ip net.IP) {
for j := len(ip) - 1; j >= 0; j-- {
ip[j]++
if ip[j] > 0 {
break
}
}
}

func IsHostOnline(host string, port int) bool {
address := fmt.Sprintf("%s:%d", host, port)
conn, err := net.DialTimeout("tcp", address, 1*time.Second)
if err != nil {
return false
}
defer conn.Close()
return true
}

func ScanNetworkRange(cidr string, port int) <-chan HostDiscoveryResult {
results := make(chan HostDiscoveryResult, 100)
var wg sync.WaitGroup

ips, err := GenerateIPRange(cidr)
if err != nil {
go func() {
results <- HostDiscoveryResult{IP: "", Port: port, Err: err}
close(results)
}()
return results
}

for _, ipAddr := range ips {
wg.Add(1)
go func(ip string) {
defer wg.Done()
if IsHostOnline(ip, port) {
results <- HostDiscoveryResult{IP: ip, Port: port, Err: nil}
}
}(ipAddr)
}

go func() {
wg.Wait()
close(results)
}()

return results
}

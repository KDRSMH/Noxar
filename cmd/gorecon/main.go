package main

import (
"bufio"
"fmt"
"net"
"os"
"strconv"
"strings"

"github.com/KDRSMH/noxar/pkg/dns"
"github.com/KDRSMH/noxar/pkg/hosts"
"github.com/KDRSMH/noxar/pkg/ports"
"github.com/KDRSMH/noxar/pkg/utils"
)

const (
ColorRed    = "\033[91m"
ColorGreen  = "\033[92m"
ColorYellow = "\033[93m"
ColorBlue   = "\033[94m"
ColorPurple = "\033[95m"
ColorCyan   = "\033[96m"
ColorReset  = "\033[0m"
ColorBold   = "\033[1m"
)

func main() {
reader := bufio.NewScanner(os.Stdin)

for {
showBanner()
showMenu()

fmt.Print(ColorCyan + "[>] " + ColorReset)
if !reader.Scan() {
break
}

choice := strings.TrimSpace(reader.Text())

switch choice {
case "1":
portScanInteractive(reader)
case "2":
hostDiscoveryInteractive(reader)
case "3":
dnsLookupInteractive(reader)
case "4":
fmt.Println(ColorGreen + "\n[+] Goodbye!" + ColorReset)
return
default:
fmt.Println(ColorRed + "[!] Invalid choice!" + ColorReset)
}
}
}

func clearScreen() {
fmt.Print("\033[2J\033[H")
}

func showBanner() {
clearScreen()
fmt.Println(ColorPurple + ColorBold)
fmt.Println("============================================================")
fmt.Println("                   N O X A R v1.0")
fmt.Println("           Network Reconnaissance Tool")
fmt.Println("          Advanced Network Security Scanner")
fmt.Println("============================================================")
fmt.Println(ColorReset)
}

func showMenu() {
fmt.Println(ColorCyan + "[*] SELECT SCAN TYPE" + ColorReset)
fmt.Println(ColorCyan + "-------------------------------------------" + ColorReset)
fmt.Println(ColorGreen + "[1] Port Scanning" + ColorReset)
fmt.Println(ColorGreen + "[2] Host Discovery" + ColorReset)
fmt.Println(ColorGreen + "[3] DNS Lookup" + ColorReset)
fmt.Println(ColorRed + "[4] Exit" + ColorReset)
fmt.Println(ColorCyan + "-------------------------------------------" + ColorReset)
}

func portScanInteractive(reader *bufio.Scanner) {
clearScreen()
fmt.Println(ColorBlue + ColorBold + "[*] PORT SCANNING MODULE" + ColorReset)
fmt.Println(ColorCyan + "===========================================" + ColorReset)

target := promptInput(reader, "Target IP")
if target == "" {
return
}

startPortStr := promptInput(reader, "Start Port (default: 1)")
startPort := 1
if startPortStr != "" {
if p, err := strconv.Atoi(startPortStr); err == nil {
startPort = p
}
}

endPortStr := promptInput(reader, "End Port (default: 1000)")
endPort := 1000
if endPortStr != "" {
if p, err := strconv.Atoi(endPortStr); err == nil {
endPort = p
}
}

fmt.Println(ColorYellow + "\n[+] Starting port scan..." + ColorReset)
fmt.Printf("    Target: %s, Ports: %d-%d\n\n", target, startPort, endPort)

portResults := ports.ScanPorts(target, startPort, endPort)
var allResults []utils.ReconResult
scannedCount := 0
totalPorts := endPort - startPort + 1

for result := range portResults {
scannedCount++
if scannedCount%100 == 0 {
fmt.Printf(ColorCyan+"[*] Progress: %d/%d ports scanned\r"+ColorReset, scannedCount, totalPorts)
}
allResults = append(allResults, utils.ReconResult{
Type:  utils.ResultTypePort,
Data:  result,
Error: nil,
})
}

fmt.Println()
printPortResults(allResults)

fmt.Println(ColorCyan + "Press Enter to continue..." + ColorReset)
reader.Scan()
}

func hostDiscoveryInteractive(reader *bufio.Scanner) {
clearScreen()
fmt.Println(ColorBlue + ColorBold + "[*] HOST DISCOVERY MODULE" + ColorReset)
fmt.Println(ColorCyan + "===========================================" + ColorReset)

cidr := promptInput(reader, "Network CIDR (e.g., 192.168.1.0/24)")
if cidr == "" {
return
}

portStr := promptInput(reader, "Port to Check (default: 80)")
port := 80
if portStr != "" {
if p, err := strconv.Atoi(portStr); err == nil {
port = p
}
}

fmt.Println(ColorYellow + "\n[+] Starting host discovery..." + ColorReset)
fmt.Printf("    Network: %s, Port: %d\n\n", cidr, port)

hostResults := hosts.ScanNetworkRange(cidr, port)
var allResults []utils.ReconResult
foundCount := 0

for result := range hostResults {
foundCount++
if result.Err == nil {
fmt.Printf(ColorGreen+"[+] Found: %v\n"+ColorReset, result.IP)
}
allResults = append(allResults, utils.ReconResult{
Type:  utils.ResultTypeHost,
Data:  result,
Error: result.Err,
})
}

fmt.Println()
printHostResults(allResults)

fmt.Println(ColorCyan + "Press Enter to continue..." + ColorReset)
reader.Scan()
}

func dnsLookupInteractive(reader *bufio.Scanner) {
clearScreen()
fmt.Println(ColorBlue + ColorBold + "[*] DNS LOOKUP MODULE" + ColorReset)
fmt.Println(ColorCyan + "===========================================" + ColorReset)

query := promptInput(reader, "Enter IP or Domain")
if query == "" {
return
}

fmt.Println(ColorYellow + "\n[+] Performing DNS lookup..." + ColorReset + "\n")

var dnsResults <-chan dns.DNsLookupResult

if net.ParseIP(query) != nil {
fmt.Println(ColorCyan + "[*] Detected as IP - Reverse DNS lookup" + ColorReset)
dnsResults = dns.ReverseLookupParallel([]string{query})
} else {
fmt.Println(ColorCyan + "[*] Detected as domain - Forward DNS lookup" + ColorReset)
dnsResults = dns.ForwardLookupsParallel([]string{query})
}

var allResults []utils.ReconResult
for result := range dnsResults {
if result.Err == nil {
fmt.Printf(ColorGreen+"[+] %s -> %s\n"+ColorReset, result.Query, result.Result)
} else {
fmt.Printf(ColorRed+"[!] Error: %v\n"+ColorReset, result.Err)
}
allResults = append(allResults, utils.ReconResult{
Type:  utils.ResultTypeDNS,
Data:  result,
Error: result.Err,
})
}

fmt.Println()
printDNSResults(allResults)

fmt.Println(ColorCyan + "Press Enter to continue..." + ColorReset)
reader.Scan()
}

func promptInput(reader *bufio.Scanner, prompt string) string {
fmt.Print(ColorCyan + "[>] " + ColorReset + prompt + ": ")
if !reader.Scan() {
return ""
}
return strings.TrimSpace(reader.Text())
}

func printPortResults(results []utils.ReconResult) {
var openPorts []utils.ReconResult
for _, r := range results {
if r.Error == nil {
portResult := r.Data.(ports.PortScanResult)
if portResult.Open {
openPorts = append(openPorts, r)
}
}
}

fmt.Println(ColorGreen + "===========================================" + ColorReset)
fmt.Println(ColorGreen + "           OPEN PORTS" + ColorReset)
fmt.Println(ColorGreen + "===========================================" + ColorReset)

if len(openPorts) == 0 {
fmt.Println(ColorRed + "[!] No open ports found" + ColorReset)
} else {
for _, r := range openPorts {
portResult := r.Data.(ports.PortScanResult)
fmt.Printf(ColorGreen+"[+] Port %-5d: OPEN\n"+ColorReset, portResult.Port)
}
}

fmt.Println(ColorGreen + "===========================================" + ColorReset)
}

func printHostResults(results []utils.ReconResult) {
var onlineHosts []utils.ReconResult
for _, r := range results {
if r.Error == nil {
onlineHosts = append(onlineHosts, r)
}
}

fmt.Println(ColorGreen + "===========================================" + ColorReset)
fmt.Println(ColorGreen + "           ONLINE HOSTS" + ColorReset)
fmt.Println(ColorGreen + "===========================================" + ColorReset)

if len(onlineHosts) == 0 {
fmt.Println(ColorRed + "[!] No online hosts found" + ColorReset)
} else {
for _, r := range onlineHosts {
hostResult := r.Data.(hosts.HostDiscoveryResult)
fmt.Printf(ColorGreen+"[+] %s:%d\n"+ColorReset, hostResult.IP, hostResult.Port)
}
}

fmt.Println(ColorGreen + "===========================================" + ColorReset)
}

func printDNSResults(results []utils.ReconResult) {
fmt.Println(ColorGreen + "===========================================" + ColorReset)
fmt.Println(ColorGreen + "           DNS RESULTS" + ColorReset)
fmt.Println(ColorGreen + "===========================================" + ColorReset)

if len(results) == 0 {
fmt.Println(ColorRed + "[!] No results found" + ColorReset)
} else {
for _, r := range results {
if r.Error == nil {
dnsResult := r.Data.(dns.DNsLookupResult)
fmt.Printf(ColorGreen+"[+] %s\n"+ColorReset, dnsResult.Result)
}
}
}

fmt.Println(ColorGreen + "===========================================" + ColorReset)
}

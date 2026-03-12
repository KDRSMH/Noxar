package main

import (
	"bufio"
	"flag"
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

// ── Renk sabitleri ────────────────────────────────────────────────────────────

const (
	colorRed    = "\033[91m"
	colorGreen  = "\033[92m"
	colorYellow = "\033[93m"
	colorBlue   = "\033[94m"
	colorPurple = "\033[95m"
	colorCyan   = "\033[96m"
	colorReset  = "\033[0m"
	colorBold   = "\033[1m"
)

// noColor: pipe/redirect algılanırsa ya da --no-color verilirse renkler kapatılır.
var noColor bool

func clr(code, text string) string {
	if noColor {
		return text
	}
	return code + text + colorReset
}

// ── isatty (terminal tespiti) ─────────────────────────────────────────────────

func isTTY() bool {
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

// ── Çıktı yazıcı ─────────────────────────────────────────────────────────────

type output struct {
	w       *bufio.Writer
	toFile  bool
}

var out *output

func newOutput(path string) (*output, error) {
	if path == "" {
		return &output{w: bufio.NewWriter(os.Stdout), toFile: false}, nil
	}
	f, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("çıktı dosyası açılamadı: %w", err)
	}
	return &output{w: bufio.NewWriter(f), toFile: true}, nil
}

func (o *output) printf(format string, a ...interface{}) {
	fmt.Fprintf(o.w, format, a...)
}

func (o *output) println(s string) {
	fmt.Fprintln(o.w, s)
}

func (o *output) flush() {
	o.w.Flush()
}

// ── clearScreen: sadece gerçek terminal varsa temizle ─────────────────────────

func clearScreen() {
	if isTTY() && !out.toFile {
		fmt.Print("\033[2J\033[H")
	}
}

// ── Doğrulama yardımcıları ────────────────────────────────────────────────────

func validatePort(p int) error {
	if p < 1 || p > 65535 {
		return fmt.Errorf("geçersiz port: %d (1-65535 aralığında olmalı)", p)
	}
	return nil
}

func validatePortRange(start, end int) error {
	if err := validatePort(start); err != nil {
		return err
	}
	if err := validatePort(end); err != nil {
		return err
	}
	if start > end {
		return fmt.Errorf("başlangıç portu (%d) bitiş portundan (%d) büyük olamaz", start, end)
	}
	return nil
}

func validateCIDR(cidr string) error {
	_, _, err := net.ParseCIDR(cidr)
	if err != nil {
		return fmt.Errorf("geçersiz CIDR formatı '%s': %w", cidr, err)
	}
	return nil
}

// ── Ana giriş noktası ─────────────────────────────────────────────────────────

func main() {
	// CLI bayrakları
	scanType   := flag.String("scan", "", "Tarama tipi: port | host | dns")
	target     := flag.String("target", "", "Hedef IP, domain ya da CIDR (ör: 192.168.1.0/24)")
	startPort  := flag.Int("start-port", 1, "Başlangıç portu")
	endPort    := flag.Int("end-port", 1000, "Bitiş portu")
	checkPort  := flag.Int("port", 80, "Host discovery için kontrol edilecek port")
	outputPath := flag.String("output", "", "Sonuçları kaydet (ör: --output sonuclar.txt)")
	noColorArg := flag.Bool("no-color", false, "ANSI renk kodlarını devre dışı bırak")
	flag.Parse()

	// Renk: terminal yoksa ya da --no-color verilmişse kapat
	noColor = *noColorArg || !isTTY()

	// Çıktı başlat
	var err error
	out, err = newOutput(*outputPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer out.flush()

	// CLI modu: --scan verilmişse doğrudan çalıştır, menüsüz
	if *scanType != "" {
		runCLIMode(*scanType, *target, *startPort, *endPort, *checkPort)
		return
	}

	// İnteraktif menü modu
	reader := bufio.NewScanner(os.Stdin)
	for {
		showBanner()
		showMenu()

		fmt.Print(clr(colorCyan, "[>] "))
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
			fmt.Println(clr(colorGreen, "\n[+] Goodbye!"))
			return
		default:
			fmt.Println(clr(colorRed, "[!] Geçersiz seçim!"))
		}
	}
}

// ── CLI modu (--scan ile) ─────────────────────────────────────────────────────

func runCLIMode(scanType, target string, startPort, endPort, checkPort int) {
	if target == "" {
		fmt.Fprintln(os.Stderr, clr(colorRed, "[!] --target parametresi gerekli"))
		os.Exit(1)
	}

	switch strings.ToLower(scanType) {
	case "port":
		if err := validatePortRange(startPort, endPort); err != nil {
			fmt.Fprintln(os.Stderr, clr(colorRed, "[!] "+err.Error()))
			os.Exit(1)
		}
		out.println(fmt.Sprintf("[*] Port taraması: %s, port %d-%d", target, startPort, endPort))
		portResults := ports.ScanPorts(target, startPort, endPort)
		var allResults []utils.ReconResult
		for result := range portResults {
			allResults = append(allResults, utils.ReconResult{
				Type:  utils.ResultTypePort,
				Data:  result,
				Error: nil,
			})
		}
		printPortResults(allResults)

	case "host":
		if err := validateCIDR(target); err != nil {
			fmt.Fprintln(os.Stderr, clr(colorRed, "[!] "+err.Error()))
			os.Exit(1)
		}
		if err := validatePort(checkPort); err != nil {
			fmt.Fprintln(os.Stderr, clr(colorRed, "[!] "+err.Error()))
			os.Exit(1)
		}
		out.println(fmt.Sprintf("[*] Host discovery: %s, port %d", target, checkPort))
		hostResults := hosts.ScanNetworkRange(target, checkPort)
		var allResults []utils.ReconResult
		for result := range hostResults {
			allResults = append(allResults, utils.ReconResult{
				Type:  utils.ResultTypeHost,
				Data:  result,
				Error: result.Err,
			})
		}
		printHostResults(allResults)

	case "dns":
		out.println(fmt.Sprintf("[*] DNS sorgusu: %s", target))
		var dnsResults <-chan dns.DNsLookupResult
		if net.ParseIP(target) != nil {
			dnsResults = dns.ReverseLookupParallel([]string{target})
		} else {
			dnsResults = dns.ForwardLookupsParallel([]string{target})
		}
		var allResults []utils.ReconResult
		for result := range dnsResults {
			allResults = append(allResults, utils.ReconResult{
				Type:  utils.ResultTypeDNS,
				Data:  result,
				Error: result.Err,
			})
		}
		printDNSResults(allResults)

	default:
		fmt.Fprintf(os.Stderr, clr(colorRed, "[!] Bilinmeyen tarama tipi: '%s' (port | host | dns)\n"), scanType)
		os.Exit(1)
	}
}

// ── Banner & menü ─────────────────────────────────────────────────────────────

func showBanner() {
	clearScreen()
	fmt.Println(clr(colorPurple+colorBold, "============================================================"))
	fmt.Println(clr(colorPurple+colorBold, "                   ") +
		clr(colorRed, "N") + clr(colorPurple+colorBold, " O X ") +
		clr(colorRed, "A") + clr(colorPurple+colorBold, " R v1.1"))
	fmt.Println(clr(colorPurple+colorBold, "           ") + clr(colorRed, "N") + clr(colorPurple+colorBold, "etwork Reconnaissance Tool"))
	fmt.Println(clr(colorPurple+colorBold, "          ") + clr(colorRed, "A") + clr(colorPurple+colorBold, "dvanced Network Security Scanner"))
	fmt.Println(clr(colorPurple+colorBold, "============================================================"))
	fmt.Println()
}

func showMenu() {
	fmt.Println(clr(colorCyan, "[*] SELECT SCAN TYPE"))
	fmt.Println(clr(colorCyan, "-------------------------------------------"))
	fmt.Println(clr(colorGreen, "[1] Port Scanning"))
	fmt.Println(clr(colorGreen, "[2] Host Discovery"))
	fmt.Println(clr(colorGreen, "[3] DNS Lookup"))
	fmt.Println(clr(colorRed, "[4] Exit"))
	fmt.Println(clr(colorCyan, "-------------------------------------------"))
}

// ── İnteraktif modlar ─────────────────────────────────────────────────────────

func portScanInteractive(reader *bufio.Scanner) {
	clearScreen()
	fmt.Println(clr(colorBlue+colorBold, "[*] PORT SCANNING MODULE"))
	fmt.Println(clr(colorCyan, "==========================================="))

	target := promptInput(reader, "Target IP veya domain")
	if target == "" {
		return
	}

	startPortStr := promptInput(reader, "Start Port (varsayılan: 1)")
	startPort := 1
	if startPortStr != "" {
		if p, err := strconv.Atoi(startPortStr); err == nil {
			startPort = p
		} else {
			fmt.Println(clr(colorRed, "[!] Geçersiz port girişi, varsayılan 1 kullanılıyor"))
		}
	}

	endPortStr := promptInput(reader, "End Port (varsayılan: 1000)")
	endPort := 1000
	if endPortStr != "" {
		if p, err := strconv.Atoi(endPortStr); err == nil {
			endPort = p
		} else {
			fmt.Println(clr(colorRed, "[!] Geçersiz port girişi, varsayılan 1000 kullanılıyor"))
		}
	}

	if err := validatePortRange(startPort, endPort); err != nil {
		fmt.Println(clr(colorRed, "[!] "+err.Error()))
		fmt.Println(clr(colorCyan, "Devam etmek için Enter'a bas..."))
		reader.Scan()
		return
	}

	fmt.Println(clr(colorYellow, "\n[+] Port taraması başlatılıyor..."))
	fmt.Printf("    Hedef: %s, Portlar: %d-%d\n\n", target, startPort, endPort)

	portResults := ports.ScanPorts(target, startPort, endPort)
	var allResults []utils.ReconResult
	scannedCount := 0
	totalPorts := endPort - startPort + 1

	for result := range portResults {
		scannedCount++
		if scannedCount%100 == 0 && isTTY() {
			fmt.Printf(clr(colorCyan, fmt.Sprintf("[*] İlerleme: %d/%d port tarandı\r", scannedCount, totalPorts)))
		}
		allResults = append(allResults, utils.ReconResult{
			Type:  utils.ResultTypePort,
			Data:  result,
			Error: nil,
		})
	}

	fmt.Println()
	printPortResults(allResults)

	fmt.Println(clr(colorCyan, "Devam etmek için Enter'a bas..."))
	reader.Scan()
}

func hostDiscoveryInteractive(reader *bufio.Scanner) {
	clearScreen()
	fmt.Println(clr(colorBlue+colorBold, "[*] HOST DISCOVERY MODULE"))
	fmt.Println(clr(colorCyan, "==========================================="))

	cidr := promptInput(reader, "Ağ CIDR (ör: 192.168.1.0/24)")
	if cidr == "" {
		return
	}

	if err := validateCIDR(cidr); err != nil {
		fmt.Println(clr(colorRed, "[!] "+err.Error()))
		fmt.Println(clr(colorCyan, "Devam etmek için Enter'a bas..."))
		reader.Scan()
		return
	}

	portStr := promptInput(reader, "Kontrol edilecek port (varsayılan: 80)")
	port := 80
	if portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		}
	}

	if err := validatePort(port); err != nil {
		fmt.Println(clr(colorRed, "[!] "+err.Error()))
		fmt.Println(clr(colorCyan, "Devam etmek için Enter'a bas..."))
		reader.Scan()
		return
	}

	fmt.Println(clr(colorYellow, "\n[+] Host discovery başlatılıyor..."))
	fmt.Printf("    Ağ: %s, Port: %d\n\n", cidr, port)

	hostResults := hosts.ScanNetworkRange(cidr, port)
	var allResults []utils.ReconResult

	for result := range hostResults {
		if result.Err == nil {
			fmt.Printf(clr(colorGreen, fmt.Sprintf("[+] Bulundu: %v\n", result.IP)))
		}
		allResults = append(allResults, utils.ReconResult{
			Type:  utils.ResultTypeHost,
			Data:  result,
			Error: result.Err,
		})
	}

	fmt.Println()
	printHostResults(allResults)

	fmt.Println(clr(colorCyan, "Devam etmek için Enter'a bas..."))
	reader.Scan()
}

func dnsLookupInteractive(reader *bufio.Scanner) {
	clearScreen()
	fmt.Println(clr(colorBlue+colorBold, "[*] DNS LOOKUP MODULE"))
	fmt.Println(clr(colorCyan, "==========================================="))

	query := promptInput(reader, "IP veya domain girin")
	if query == "" {
		return
	}

	fmt.Println(clr(colorYellow, "\n[+] DNS sorgusu yapılıyor...") + "\n")

	var dnsResults <-chan dns.DNsLookupResult

	if net.ParseIP(query) != nil {
		fmt.Println(clr(colorCyan, "[*] IP algılandı — Reverse DNS sorgusu"))
		dnsResults = dns.ReverseLookupParallel([]string{query})
	} else {
		fmt.Println(clr(colorCyan, "[*] Domain algılandı — Forward DNS sorgusu"))
		dnsResults = dns.ForwardLookupsParallel([]string{query})
	}

	var allResults []utils.ReconResult
	for result := range dnsResults {
		if result.Err == nil {
			fmt.Printf(clr(colorGreen, fmt.Sprintf("[+] %s -> %s\n", result.Query, result.Result)))
		} else {
			fmt.Printf(clr(colorRed, fmt.Sprintf("[!] Hata: %v\n", result.Err)))
		}
		allResults = append(allResults, utils.ReconResult{
			Type:  utils.ResultTypeDNS,
			Data:  result,
			Error: result.Err,
		})
	}

	fmt.Println()
	printDNSResults(allResults)

	fmt.Println(clr(colorCyan, "Devam etmek için Enter'a bas..."))
	reader.Scan()
}

// ── Yardımcı: girdi al ────────────────────────────────────────────────────────

func promptInput(reader *bufio.Scanner, prompt string) string {
	fmt.Print(clr(colorCyan, "[>] ") + prompt + ": ")
	if !reader.Scan() {
		return ""
	}
	return strings.TrimSpace(reader.Text())
}

// ── Sonuç yazdırıcılar ────────────────────────────────────────────────────────

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

	out.println(clr(colorGreen, "==========================================="))
	out.println(clr(colorGreen, "           AÇIK PORTLAR"))
	out.println(clr(colorGreen, "==========================================="))

	if len(openPorts) == 0 {
		out.println(clr(colorRed, "[!] Açık port bulunamadı"))
	} else {
		for _, r := range openPorts {
			portResult := r.Data.(ports.PortScanResult)
			out.println(clr(colorGreen, fmt.Sprintf("[+] Port %-5d: OPEN", portResult.Port)))
		}
	}
	out.println(clr(colorGreen, "==========================================="))
	out.flush()
}

func printHostResults(results []utils.ReconResult) {
	var onlineHosts []utils.ReconResult
	for _, r := range results {
		if r.Error == nil {
			onlineHosts = append(onlineHosts, r)
		}
	}

	out.println(clr(colorGreen, "==========================================="))
	out.println(clr(colorGreen, "           ONLINE HOSTLAR"))
	out.println(clr(colorGreen, "==========================================="))

	if len(onlineHosts) == 0 {
		out.println(clr(colorRed, "[!] Online host bulunamadı"))
	} else {
		for _, r := range onlineHosts {
			hostResult := r.Data.(hosts.HostDiscoveryResult)
			out.println(clr(colorGreen, fmt.Sprintf("[+] %s:%d", hostResult.IP, hostResult.Port)))
		}
	}
	out.println(clr(colorGreen, "==========================================="))
	out.flush()
}

func printDNSResults(results []utils.ReconResult) {
	out.println(clr(colorGreen, "==========================================="))
	out.println(clr(colorGreen, "           DNS SONUÇLARI"))
	out.println(clr(colorGreen, "==========================================="))

	if len(results) == 0 {
		out.println(clr(colorRed, "[!] Sonuç bulunamadı"))
	} else {
		for _, r := range results {
			if r.Error == nil {
				dnsResult := r.Data.(dns.DNsLookupResult)
				out.println(clr(colorGreen, "[+] "+dnsResult.Result))
			}
		}
	}
	out.println(clr(colorGreen, "==========================================="))
	out.flush()
}

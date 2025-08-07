package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/go-ping/ping"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]

	switch command {
	case "ping":
		if len(os.Args) < 3 {
			fmt.Println("Usage: network-scanner-cli ping <host>")
			return
		}
		host := os.Args[2]
		if pingHost(host) {
			fmt.Printf("Host %s: ALIVE\n", host)
		} else {
			fmt.Printf("Host %s: NOT REACHABLE\n", host)
		}

	case "portscan":
		if len(os.Args) < 5 {
			fmt.Println("Usage: network-scanner-cli portscan <host> <start_port> <end_port>")
			return
		}
		host := os.Args[2]
		startPort, err1 := strconv.Atoi(os.Args[3])
		endPort, err2 := strconv.Atoi(os.Args[4])

		if err1 != nil || err2 != nil {
			fmt.Println("Error: Invalid port numbers")
			return
		}

		scanPorts(host, startPort, endPort)

	case "netscan":
		if len(os.Args) < 3 {
			fmt.Println("Usage: network-scanner-cli netscan <network_cidr>")
			return
		}
		network := os.Args[2]
		scanNetwork(network)

	default:
		printUsage()
	}
}

func printUsage() {
	fmt.Println("Network Scanner CLI")
	fmt.Println("Usage:")
	fmt.Println("  network-scanner-cli ping <host>")
	fmt.Println("  network-scanner-cli portscan <host> <start_port> <end_port>")
	fmt.Println("  network-scanner-cli netscan <network_cidr>")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  network-scanner-cli ping google.com")
	fmt.Println("  network-scanner-cli portscan 192.168.1.1 1 1000")
	fmt.Println("  network-scanner-cli netscan 192.168.1.0/24")
}

func pingHost(host string) bool {
	pinger, err := ping.NewPinger(host)
	if err != nil {
		return false
	}

	pinger.SetPrivileged(false)
	pinger.Count = 1
	pinger.Timeout = 1 * time.Second

	err = pinger.Run()
	if err != nil {
		return false
	}

	stats := pinger.Statistics()
	return stats.PacketsRecv > 0
}

func scanPorts(host string, startPort, endPort int) {
	fmt.Printf("Scanning ports %d-%d on %s...\n", startPort, endPort, host)

	openPorts := []int{}
	totalPorts := endPort - startPort + 1
	scannedPorts := 0

	for port := startPort; port <= endPort; port++ {
		address := fmt.Sprintf("%s:%d", host, port)

		conn, err := net.DialTimeout("tcp", address, 1*time.Second)
		if err == nil {
			conn.Close()
			openPorts = append(openPorts, port)
			fmt.Printf("Port %d: OPEN\n", port)
		}

		scannedPorts++
		if scannedPorts%100 == 0 {
			fmt.Printf("Progress: %d/%d ports scanned\n", scannedPorts, totalPorts)
		}
	}

	fmt.Printf("\nScan complete. Found %d open ports out of %d scanned.\n", len(openPorts), totalPorts)
}

func scanNetwork(network string) {
	fmt.Printf("Scanning network %s...\n", network)

	_, ipNet, err := net.ParseCIDR(network)
	if err != nil {
		fmt.Printf("Error parsing network: %v\n", err)
		return
	}

	var ips []string
	for ip := ipNet.IP.Mask(ipNet.Mask); ipNet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}

	totalIPs := len(ips)
	aliveHosts := []string{}
	var mu sync.Mutex
	var wg sync.WaitGroup

	semaphore := make(chan struct{}, 50)

	for i, ip := range ips {
		wg.Add(1)
		go func(ip string, index int) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			if pingHost(ip) {
				mu.Lock()
				aliveHosts = append(aliveHosts, ip)
				fmt.Printf("Host %s: ALIVE\n", ip)
				mu.Unlock()
			}

			if (index+1)%50 == 0 {
				fmt.Printf("Progress: %d/%d hosts scanned\n", index+1, totalIPs)
			}
		}(ip, i)
	}

	wg.Wait()
	fmt.Printf("\nNetwork scan complete. Found %d alive hosts out of %d scanned.\n", len(aliveHosts), totalIPs)
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

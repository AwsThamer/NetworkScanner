package main

import (
	"fmt"
	"image/color"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/go-ping/ping"
)

// Custom theme for better colors
type myTheme struct{}

func (m myTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNamePrimary:
		return color.RGBA{0, 123, 255, 255} // Bootstrap blue
	case theme.ColorNameSuccess:
		return color.RGBA{40, 167, 69, 255} // Bootstrap green
	case theme.ColorNameWarning:
		return color.RGBA{255, 193, 7, 255} // Bootstrap yellow
	case theme.ColorNameError:
		return color.RGBA{220, 53, 69, 255} // Bootstrap red
	case theme.ColorNameBackground:
		return color.RGBA{248, 249, 250, 255} // Light gray background
	case theme.ColorNameForeground:
		return color.RGBA{33, 37, 41, 255} // Dark text
	}
	return theme.DefaultTheme().Color(name, variant)
}

func (m myTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (m myTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (m myTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNamePadding:
		return 8
	case theme.SizeNameInlineIcon:
		return 20
	case theme.SizeNameText:
		return 14
	}
	return theme.DefaultTheme().Size(name)
}

type Scanner struct {
	results      *widget.List
	resultData   []ScanResult
	progress     *widget.ProgressBar
	status       *widget.Label
	mu           sync.Mutex
	isScanning   bool
	scanningBtn  *widget.Button
	stopScanning chan bool
}

type ScanResult struct {
	Message string
	Type    string // "info", "success", "warning", "error"
	Time    string
}

func NewScanner() *Scanner {
	s := &Scanner{
		resultData:   []ScanResult{},
		status:       widget.NewLabelWithStyle("🚀 Ready to scan networks", fyne.TextAlignLeading, fyne.TextStyle{}),
		progress:     widget.NewProgressBar(),
		stopScanning: make(chan bool),
	}

	// Enhanced progress bar
	s.progress.TextFormatter = func() string {
		return fmt.Sprintf("%.1f%%", s.progress.Value*100)
	}

	s.results = widget.NewList(
		func() int {
			return len(s.resultData)
		},
		func() fyne.CanvasObject {
			icon := widget.NewIcon(theme.InfoIcon())
			icon.Resize(fyne.NewSize(16, 16))

			timeLabel := widget.NewLabelWithStyle("00:00:00", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})
			timeLabel.Resize(fyne.NewSize(70, 20))

			messageLabel := widget.NewLabel("Template message")
			messageLabel.Wrapping = fyne.TextWrapWord

			// Create a colored background rectangle
			bg := canvas.NewRectangle(color.RGBA{255, 255, 255, 255})
			bg.StrokeColor = color.RGBA{230, 230, 230, 255}
			bg.StrokeWidth = 1

			content := container.NewHBox(
				icon,
				timeLabel,
				widget.NewSeparator(),
				messageLabel,
			)

			return container.NewStack(bg, container.NewPadded(content))
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			if i >= len(s.resultData) {
				return
			}

			result := s.resultData[i]
			stack := o.(*fyne.Container)
			bg := stack.Objects[0].(*canvas.Rectangle)
			padded := stack.Objects[1].(*fyne.Container)
			content := padded.Objects[0].(*fyne.Container)

			icon := content.Objects[0].(*widget.Icon)
			timeLabel := content.Objects[1].(*widget.Label)
			messageLabel := content.Objects[3].(*widget.Label)

			timeLabel.SetText(result.Time)
			messageLabel.SetText(result.Message)

			switch result.Type {
			case "success":
				icon.SetResource(theme.ConfirmIcon())
				bg.FillColor = color.RGBA{212, 237, 218, 255} // Light green
			case "warning":
				icon.SetResource(theme.WarningIcon())
				bg.FillColor = color.RGBA{255, 243, 205, 255} // Light yellow
			case "error":
				icon.SetResource(theme.ErrorIcon())
				bg.FillColor = color.RGBA{248, 215, 218, 255} // Light red
			default:
				icon.SetResource(theme.InfoIcon())
				bg.FillColor = color.RGBA{217, 237, 247, 255} // Light blue
			}
			bg.Refresh()
		},
	)

	return s
}

func (s *Scanner) addResult(message, resultType string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().Format("15:04:05")
	s.resultData = append(s.resultData, ScanResult{
		Message: message,
		Type:    resultType,
		Time:    now,
	})
	s.results.Refresh()
	s.results.ScrollToBottom()
}

func (s *Scanner) clearResults() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.resultData = []ScanResult{}
	s.results.Refresh()
}

func (s *Scanner) updateStatus(status string) {
	s.status.SetText(status)
}

func (s *Scanner) updateProgress(value float64) {
	s.progress.SetValue(value)
}

func (s *Scanner) setScanning(scanning bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.isScanning = scanning
	if s.scanningBtn != nil {
		if scanning {
			s.scanningBtn.SetText("⏹️ Stop Scan")
			s.scanningBtn.Importance = widget.HighImportance
		} else {
			s.scanningBtn.SetText("▶️ Start Scan")
			s.scanningBtn.Importance = widget.MediumImportance
		}
	}
}

func (s *Scanner) scanPorts(host string, startPort, endPort int) {
	s.clearResults()
	s.setScanning(true)
	s.updateStatus("🔍 Scanning ports...")
	s.addResult(fmt.Sprintf("🎯 Starting port scan on %s (ports %d-%d)", host, startPort, endPort), "info")

	totalPorts := endPort - startPort + 1
	scannedPorts := 0
	openPorts := 0

	for port := startPort; port <= endPort; port++ {
		select {
		case <-s.stopScanning:
			s.addResult("⏹️ Scan stopped by user", "warning")
			s.setScanning(false)
			s.updateStatus("⏹️ Scan stopped")
			return
		default:
		}

		address := fmt.Sprintf("%s:%d", host, port)

		conn, err := net.DialTimeout("tcp", address, 1*time.Second)
		if err == nil {
			conn.Close()
			openPorts++
			s.addResult(fmt.Sprintf("✅ Port %d: OPEN", port), "success")
		}

		scannedPorts++
		progress := float64(scannedPorts) / float64(totalPorts)
		s.updateProgress(progress)

		if scannedPorts%25 == 0 {
			s.updateStatus(fmt.Sprintf("🔍 Scanning... %d/%d ports (%d open)", scannedPorts, totalPorts, openPorts))
		}
	}

	s.setScanning(false)
	s.addResult(fmt.Sprintf("🎉 Scan complete! Found %d open ports out of %d scanned", openPorts, totalPorts), "info")
	s.updateStatus(fmt.Sprintf("✅ Scan complete. %d open ports found.", openPorts))
}

func (s *Scanner) scanNetwork(network string) {
	s.clearResults()
	s.setScanning(true)
	s.updateStatus("🌐 Scanning network...")
	s.addResult(fmt.Sprintf("🌍 Starting network discovery on %s", network), "info")

	_, ipNet, err := net.ParseCIDR(network)
	if err != nil {
		s.addResult(fmt.Sprintf("❌ Error parsing network: %v", err), "error")
		s.setScanning(false)
		return
	}

	var ips []string
	for ip := ipNet.IP.Mask(ipNet.Mask); ipNet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}

	totalIPs := len(ips)
	scannedIPs := 0
	aliveHosts := 0

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 50)

	for _, ip := range ips {
		select {
		case <-s.stopScanning:
			s.addResult("⏹️ Scan stopped by user", "warning")
			s.setScanning(false)
			s.updateStatus("⏹️ Scan stopped")
			return
		default:
		}

		wg.Add(1)
		go func(ip string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			if s.pingHost(ip) {
				s.mu.Lock()
				aliveHosts++
				s.mu.Unlock()
				s.addResult(fmt.Sprintf("💚 Host %s: ALIVE", ip), "success")
			}

			s.mu.Lock()
			scannedIPs++
			progress := float64(scannedIPs) / float64(totalIPs)
			s.updateProgress(progress)
			if scannedIPs%10 == 0 {
				s.updateStatus(fmt.Sprintf("🌐 Scanning... %d/%d hosts (%d alive)", scannedIPs, totalIPs, aliveHosts))
			}
			s.mu.Unlock()
		}(ip)
	}

	wg.Wait()
	s.setScanning(false)
	s.addResult(fmt.Sprintf("🎉 Network scan complete! Found %d alive hosts out of %d scanned", aliveHosts, totalIPs), "info")
	s.updateStatus(fmt.Sprintf("✅ Network scan complete. %d hosts found.", aliveHosts))
}

func (s *Scanner) pingHost(host string) bool {
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

func (s *Scanner) pingNetwork(network string) {
	s.clearResults()
	s.setScanning(true)
	s.updateStatus("🌐 Pinging network range...")
	s.addResult(fmt.Sprintf("🌍 Starting ping sweep on %s", network), "info")

	_, ipNet, err := net.ParseCIDR(network)
	if err != nil {
		s.addResult(fmt.Sprintf("❌ Error parsing network: %v", err), "error")
		s.setScanning(false)
		return
	}

	var ips []string
	for ip := ipNet.IP.Mask(ipNet.Mask); ipNet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}

	totalIPs := len(ips)
	scannedIPs := 0
	aliveHosts := 0

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 50)

	for _, ip := range ips {
		select {
		case <-s.stopScanning:
			s.addResult("⏹️ Ping sweep stopped by user", "warning")
			s.setScanning(false)
			s.updateStatus("⏹️ Ping sweep stopped")
			return
		default:
		}

		wg.Add(1)
		go func(ip string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			if s.pingHost(ip) {
				s.mu.Lock()
				aliveHosts++
				s.mu.Unlock()
				s.addResult(fmt.Sprintf("🟢 %s: ALIVE (ping successful)", ip), "success")
			} else {
				s.addResult(fmt.Sprintf("🔴 %s: No response", ip), "error")
			}

			s.mu.Lock()
			scannedIPs++
			progress := float64(scannedIPs) / float64(totalIPs)
			s.updateProgress(progress)
			if scannedIPs%5 == 0 {
				s.updateStatus(fmt.Sprintf("🌐 Pinging... %d/%d hosts (%d responding)", scannedIPs, totalIPs, aliveHosts))
			}
			s.mu.Unlock()
		}(ip)
	}

	wg.Wait()
	s.setScanning(false)
	s.addResult(fmt.Sprintf("🎉 Ping sweep complete! %d hosts responded out of %d pinged", aliveHosts, totalIPs), "info")
	s.updateStatus(fmt.Sprintf("✅ Ping sweep complete. %d hosts responding.", aliveHosts))
}

func (s *Scanner) pingRange(rangeStr string) {
	s.clearResults()
	s.setScanning(true)
	s.updateStatus("🎯 Pinging custom range...")
	s.addResult(fmt.Sprintf("🎯 Starting ping sweep on range %s", rangeStr), "info")

	// Parse custom range (e.g., "192.168.1.1-192.168.1.50")
	parts := strings.Split(rangeStr, "-")
	if len(parts) != 2 {
		s.addResult("❌ Error: Invalid range format. Use: IP1-IP2 (e.g., 192.168.1.1-192.168.1.50)", "error")
		s.setScanning(false)
		return
	}

	startIP := net.ParseIP(strings.TrimSpace(parts[0]))
	endIP := net.ParseIP(strings.TrimSpace(parts[1]))

	if startIP == nil || endIP == nil {
		s.addResult("❌ Error: Invalid IP addresses in range", "error")
		s.setScanning(false)
		return
	}

	// Generate IP range
	var ips []string
	current := make(net.IP, len(startIP))
	copy(current, startIP)

	for {
		ips = append(ips, current.String())
		if current.Equal(endIP) {
			break
		}
		inc(current)
		// Safety check to prevent infinite loops
		if len(ips) > 1000 {
			s.addResult("⚠️ Warning: Range too large (max 1000 IPs), truncating", "warning")
			break
		}
	}

	totalIPs := len(ips)
	scannedIPs := 0
	aliveHosts := 0

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 50)

	for _, ip := range ips {
		select {
		case <-s.stopScanning:
			s.addResult("⏹️ Range ping stopped by user", "warning")
			s.setScanning(false)
			s.updateStatus("⏹️ Range ping stopped")
			return
		default:
		}

		wg.Add(1)
		go func(ip string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			if s.pingHost(ip) {
				s.mu.Lock()
				aliveHosts++
				s.mu.Unlock()
				s.addResult(fmt.Sprintf("🟢 %s: ALIVE (ping successful)", ip), "success")
			} else {
				s.addResult(fmt.Sprintf("🔴 %s: No response", ip), "error")
			}

			s.mu.Lock()
			scannedIPs++
			progress := float64(scannedIPs) / float64(totalIPs)
			s.updateProgress(progress)
			if scannedIPs%5 == 0 {
				s.updateStatus(fmt.Sprintf("🎯 Range ping... %d/%d IPs (%d responding)", scannedIPs, totalIPs, aliveHosts))
			}
			s.mu.Unlock()
		}(ip)
	}

	wg.Wait()
	s.setScanning(false)
	s.addResult(fmt.Sprintf("🎉 Range ping complete! %d hosts responded out of %d pinged", aliveHosts, totalIPs), "info")
	s.updateStatus(fmt.Sprintf("✅ Range ping complete. %d hosts responding.", aliveHosts))
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

// Create beautiful card with gradient background
func createStyledCard(title string, icon fyne.Resource, content fyne.CanvasObject) *fyne.Container {
	// Create gradient background
	bg := canvas.NewLinearGradient(
		color.RGBA{255, 255, 255, 255},
		color.RGBA{248, 249, 250, 255},
		90, // vertical gradient
	)

	// Create header with icon and title
	titleIcon := widget.NewIcon(icon)
	titleIcon.Resize(fyne.NewSize(24, 24))

	titleLabel := widget.NewLabelWithStyle(title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	titleLabel.TextStyle.Bold = true

	header := container.NewHBox(
		titleIcon,
		titleLabel,
		layout.NewSpacer(),
	)

	// Create separator
	separator := canvas.NewRectangle(color.RGBA{0, 123, 255, 100})
	separator.Resize(fyne.NewSize(0, 2))

	// Combine elements
	cardContent := container.NewVBox(
		header,
		separator,
		widget.NewSeparator(),
		content,
	)

	// Add padding and background
	paddedContent := container.NewPadded(cardContent)

	return container.NewStack(bg, paddedContent)
}

func main() {
	myApp := app.New()
	myApp.Settings().SetTheme(&myTheme{})
	myApp.SetIcon(theme.ComputerIcon())

	myWindow := myApp.NewWindow("🔍 Network Scanner Pro")
	myWindow.Resize(fyne.NewSize(1000, 750))
	myWindow.CenterOnScreen()

	scanner := NewScanner()

	// Enhanced input fields with better styling
	hostEntry := widget.NewEntry()
	hostEntry.SetPlaceHolder("🌐 Enter host or IP address (e.g., google.com, 192.168.1.1)")
	hostEntry.Resize(fyne.NewSize(400, 35))

	// Network entry with preset buttons
	networkEntry := widget.NewEntry()
	networkEntry.SetPlaceHolder("🌍 Enter network CIDR (e.g., 192.168.1.0/24)")
	networkEntry.Resize(fyne.NewSize(300, 35))

	// Network preset buttons for common ranges
	preset192 := widget.NewButtonWithIcon("192.168.1.0/24", theme.HomeIcon(), func() {
		networkEntry.SetText("192.168.1.0/24")
	})
	preset192.Resize(fyne.NewSize(120, 35))

	preset10 := widget.NewButtonWithIcon("10.0.0.0/24", theme.HomeIcon(), func() {
		networkEntry.SetText("10.0.0.0/24")
	})
	preset10.Resize(fyne.NewSize(120, 35))

	preset172 := widget.NewButtonWithIcon("172.16.0.0/24", theme.HomeIcon(), func() {
		networkEntry.SetText("172.16.0.0/24")
	})
	preset172.Resize(fyne.NewSize(120, 35))

	// Custom range entry
	customRangeEntry := widget.NewEntry()
	customRangeEntry.SetPlaceHolder("🎯 Custom range (e.g., 192.168.1.1-192.168.1.50)")
	customRangeEntry.Resize(fyne.NewSize(300, 35))

	startPortEntry := widget.NewEntry()
	startPortEntry.SetPlaceHolder("Start")
	startPortEntry.SetText("1")
	startPortEntry.Resize(fyne.NewSize(100, 35))

	endPortEntry := widget.NewEntry()
	endPortEntry.SetPlaceHolder("End")
	endPortEntry.SetText("1000")
	endPortEntry.Resize(fyne.NewSize(100, 35))

	// Port preset buttons for common ports
	commonPortsBtn := widget.NewButtonWithIcon("Common", theme.ListIcon(), func() {
		startPortEntry.SetText("1")
		endPortEntry.SetText("1024")
	})

	webPortsBtn := widget.NewButtonWithIcon("Web", theme.ComputerIcon(), func() {
		startPortEntry.SetText("80")
		endPortEntry.SetText("443")
	})

	allPortsBtn := widget.NewButtonWithIcon("All", theme.ViewFullScreenIcon(), func() {
		startPortEntry.SetText("1")
		endPortEntry.SetText("65535")
	})

	// Enhanced buttons with better styling
	var portScanBtn, networkScanBtn, pingRangeBtn *widget.Button

	portScanBtn = widget.NewButtonWithIcon("🔍 Port Scan", theme.SearchIcon(), func() {
		if scanner.isScanning {
			scanner.stopScanning <- true
			return
		}

		host := strings.TrimSpace(hostEntry.Text)
		if host == "" {
			scanner.addResult("❌ Error: Please enter a host", "error")
			return
		}

		startPort, err1 := strconv.Atoi(startPortEntry.Text)
		endPort, err2 := strconv.Atoi(endPortEntry.Text)

		if err1 != nil || err2 != nil {
			scanner.addResult("❌ Error: Invalid port range", "error")
			return
		}

		if startPort > endPort || startPort < 1 || endPort > 65535 {
			scanner.addResult("❌ Error: Invalid port range (1-65535)", "error")
			return
		}

		scanner.scanningBtn = portScanBtn
		go scanner.scanPorts(host, startPort, endPort)
	})
	portScanBtn.Importance = widget.MediumImportance

	networkScanBtn = widget.NewButtonWithIcon("🌐 Network Discovery", theme.ViewRefreshIcon(), func() {
		if scanner.isScanning {
			scanner.stopScanning <- true
			return
		}

		network := strings.TrimSpace(networkEntry.Text)
		if network == "" {
			scanner.addResult("❌ Error: Please enter a network", "error")
			return
		}

		scanner.scanningBtn = networkScanBtn
		go scanner.scanNetwork(network)
	})
	networkScanBtn.Importance = widget.MediumImportance

	// New ping range button
	pingRangeBtn = widget.NewButtonWithIcon("🌍 Ping Range", theme.RadioButtonIcon(), func() {
		if scanner.isScanning {
			scanner.stopScanning <- true
			return
		}

		customRange := strings.TrimSpace(customRangeEntry.Text)
		network := strings.TrimSpace(networkEntry.Text)

		if customRange != "" {
			scanner.scanningBtn = pingRangeBtn
			go scanner.pingRange(customRange)
		} else if network != "" {
			scanner.scanningBtn = pingRangeBtn
			go scanner.pingNetwork(network)
		} else {
			scanner.addResult("❌ Error: Please enter a network or custom range", "error")
		}
	})
	pingRangeBtn.Importance = widget.MediumImportance

	pingBtn := widget.NewButtonWithIcon("🏓 Quick Ping", theme.MailSendIcon(), func() {
		host := strings.TrimSpace(hostEntry.Text)
		if host == "" {
			scanner.addResult("❌ Error: Please enter a host", "error")
			return
		}

		go func() {
			scanner.updateStatus("🏓 Pinging host...")
			if scanner.pingHost(host) {
				scanner.addResult(fmt.Sprintf("✅ Host %s: ALIVE", host), "success")
			} else {
				scanner.addResult(fmt.Sprintf("❌ Host %s: NOT REACHABLE", host), "error")
			}
			scanner.updateStatus("🚀 Ready to scan networks")
		}()
	})
	pingBtn.Importance = widget.LowImportance

	clearBtn := widget.NewButtonWithIcon("🧹 Clear Results", theme.DeleteIcon(), func() {
		scanner.clearResults()
		scanner.updateStatus("✨ Results cleared - Ready to scan")
		scanner.updateProgress(0)
	})
	clearBtn.Importance = widget.LowImportance

	// Create styled cards
	targetCard := createStyledCard("🎯 Target Configuration", theme.ComputerIcon(), container.NewVBox(
		widget.NewLabelWithStyle("Host/IP Address:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		hostEntry,
		widget.NewSeparator(),
		widget.NewLabelWithStyle("Network (CIDR):", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewHBox(networkEntry, preset192, preset10, preset172),
		widget.NewSeparator(),
		widget.NewLabelWithStyle("Custom IP Range:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		customRangeEntry,
		widget.NewLabelWithStyle("💡 Format: 192.168.1.1-192.168.1.50", fyne.TextAlignLeading, fyne.TextStyle{Italic: true}),
	))

	portCard := createStyledCard("🔌 Port Configuration", theme.SettingsIcon(), container.NewVBox(
		container.NewHBox(
			widget.NewLabelWithStyle("From:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			startPortEntry,
			widget.NewLabelWithStyle("To:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			endPortEntry,
			widget.NewLabelWithStyle("Presets:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			commonPortsBtn,
			webPortsBtn,
			allPortsBtn,
		),
		widget.NewSeparator(),
		widget.NewLabelWithStyle("💡 Common ports: 21(FTP), 22(SSH), 23(Telnet), 25(SMTP), 53(DNS), 80(HTTP), 110(POP3), 443(HTTPS), 993(IMAPS), 995(POP3S)", fyne.TextAlignLeading, fyne.TextStyle{Italic: true}),
	))

	buttonCard := createStyledCard("🚀 Scan Operations", theme.MediaPlayIcon(), container.NewGridWithColumns(3,
		portScanBtn,
		networkScanBtn,
		pingRangeBtn,
		pingBtn,
		clearBtn,
		widget.NewLabel(""), // Empty space
	))

	statusCard := createStyledCard("📊 Status & Progress", theme.InfoIcon(), container.NewVBox(
		scanner.status,
		widget.NewSeparator(),
		widget.NewLabelWithStyle("Progress:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		scanner.progress,
	))

	resultsCard := createStyledCard("📋 Scan Results", theme.DocumentIcon(), scanner.results)

	// Create beautiful tabs
	inputTab := container.NewVBox(
		targetCard,
		portCard,
		buttonCard,
		statusCard,
	)

	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("⚙️ Scanner", theme.SettingsIcon(), inputTab),
		container.NewTabItemWithIcon("📊 Results", theme.DocumentIcon(), resultsCard),
	)

	// Main layout with beautiful header
	headerBg := canvas.NewLinearGradient(
		color.RGBA{0, 123, 255, 255},
		color.RGBA{108, 117, 125, 255},
		0, // horizontal gradient
	)

	headerTitle := widget.NewLabelWithStyle("🔍 Network Scanner Pro", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	headerTitle.Alignment = fyne.TextAlignCenter

	headerSubtitle := widget.NewLabelWithStyle("Advanced Network Discovery & Port Scanning Tool", fyne.TextAlignCenter, fyne.TextStyle{Italic: true})
	headerSubtitle.Alignment = fyne.TextAlignCenter

	headerContent := container.NewVBox(
		headerTitle,
		headerSubtitle,
	)

	header := container.NewStack(headerBg, container.NewPadded(headerContent))
	header.Resize(fyne.NewSize(0, 80))

	// Footer
	footer := widget.NewLabelWithStyle("💡 Tip: Always ensure you have permission to scan target networks • Use responsibly", fyne.TextAlignCenter, fyne.TextStyle{Italic: true})

	content := container.NewBorder(
		header,
		footer,
		nil,
		nil,
		tabs,
	)

	myWindow.SetContent(content)

	// Add welcome messages
	scanner.addResult("🎉 Welcome to Network Scanner Pro!", "info")
	scanner.addResult("💡 Choose your target and scan type to begin network discovery", "info")
	scanner.addResult("🔒 Remember: Only scan networks you own or have permission to test", "warning")

	myWindow.ShowAndRun()
}

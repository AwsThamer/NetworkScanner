# 🔍 Network Scanner Pro

[![Build Status](https://github.com/AwsThamer/NetworkScanner/workflows/Build%20Network%20Scanner/badge.svg)](https://github.com/AwsThamer/NetworkScanner/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/AwsThamer/NetworkScanner)](https://goreportcard.com/report/github.com/AwsThamer/NetworkScanner)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A powerful, cross-platform network scanning tool built with Go and featuring a beautiful modern GUI.

![Network Scanner Pro Screenshot](https://via.placeholder.com/800x600/007BFF/ffffff?text=Network+Scanner+Pro+GUI)

## ✨ Features

### 🎯 Enhanced Ping Range Functionality
- **Network Presets**: One-click buttons for common networks (192.168.1.0/24, 10.0.0.0/24, 172.16.0.0/24)
- **Custom IP Ranges**: Enter ranges like 192.168.1.1-192.168.1.50 for targeted scanning
- **Color-coded Results**: Green for responsive hosts, red for non-responsive
- **Real-time Progress**: Live response counts and progress tracking

### 🔍 Comprehensive Scanning
- **Port Scanning**: Scan port ranges on target hosts with preset options
- **Network Discovery**: Full CIDR network scanning capabilities
- **Host Ping**: Quick single-host connectivity testing
- **Cross-Platform**: Works seamlessly on Linux and Windows

### 🎨 Beautiful GUI Interface
- **Bootstrap-Inspired Design**: Professional color theme with gradients
- **Card-Based Layout**: Modern, organized interface sections
- **Rich Visual Feedback**: Timestamped, color-coded results
- **Interactive Elements**: Dynamic buttons and real-time progress
- **Professional Styling**: Gradient headers and modern typography

## 🚀 Quick Start

### Download Pre-built Binaries

Check the [Releases](https://github.com/AwsThamer/NetworkScanner/releases) page for pre-built binaries.

### GUI Version (Linux)
```bash
chmod +x network-scanner-gui
./network-scanner-gui
```

### CLI Version (Cross-platform)
```bash
# Ping a host
./network-scanner-cli ping google.com

# Scan ports
./network-scanner-cli portscan 192.168.1.1 80 443

# Scan network
./network-scanner-cli netscan 192.168.1.0/24
```

## 🔧 Building from Source

### Prerequisites
- Go 1.21 or later
- Git
- Linux: X11 development libraries

### Install Dependencies (Linux)
```bash
sudo apt update
sudo apt install -y libx11-dev libxrandr-dev libxinerama-dev libxcursor-dev libxi-dev libgl1-mesa-dev libxxf86vm-dev
```

### Build Instructions
```bash
# Clone the repository
git clone https://github.com/AwsThamer/NetworkScanner.git
cd NetworkScanner

# Install Go dependencies
go mod tidy

# Build all versions
./build-all.sh
```

This creates:
- `bin/network-scanner-gui` - Linux GUI version
- `bin/network-scanner-cli` - Linux CLI version  
- `bin/network-scanner-cli-windows.exe` - Windows CLI version

## 📖 Usage Guide

### GUI Features

#### 🎯 Target Configuration
- **Host/IP Input**: Enter single hosts or IP addresses
- **Network Presets**: Click preset buttons for common networks
- **Custom Ranges**: Enter IP ranges like 192.168.1.1-192.168.1.50

#### 🔌 Port Configuration  
- **Port Range**: Set start and end ports
- **Presets**: Common (1-1024), Web (80,443), All (1-65535)

#### 🚀 Scan Operations
- **Port Scan**: Comprehensive port scanning
- **Network Discovery**: Find all hosts in network
- **Ping Range**: Fast ping sweep functionality
- **Quick Ping**: Single host connectivity test

### CLI Commands

```bash
# Ping operations
./network-scanner-cli ping <host>
./network-scanner-cli ping google.com

# Port scanning  
./network-scanner-cli portscan <host> <start_port> <end_port>
./network-scanner-cli portscan 192.168.1.1 1 1000

# Network scanning
./network-scanner-cli netscan <network_cidr>
./network-scanner-cli netscan 192.168.1.0/24
```

## 🛡️ Security & Ethics

⚠️ **Important**: Only scan networks you own or have explicit permission to test.

- Uses unprivileged ping mode for compatibility
- Implements connection timeouts and rate limiting
- Respects network resources with controlled concurrency
- Designed for legitimate network administration purposes

## 🎨 GUI Screenshots

### Main Interface
The beautiful, card-based interface with gradient backgrounds and professional styling.

### Ping Range Results  
Color-coded results showing responsive (green) and non-responsive (red) hosts with timestamps.

### Port Scanning
Real-time port scan results with progress tracking and open port detection.

## 🔧 Technical Details

### Architecture
- **Language**: Go 1.21+
- **GUI Framework**: Fyne v2
- **Networking**: Native Go net package + go-ping
- **Concurrency**: Goroutines with semaphore limiting

### Scanning Methods
- **Port Scanning**: TCP connection attempts (1s timeout)
- **Host Discovery**: ICMP ping (1s timeout)  
- **Network Discovery**: CIDR range iteration
- **Concurrent Processing**: Controlled with semaphores

### Supported Platforms
- ✅ Linux (amd64) - Full GUI + CLI
- ✅ Windows (amd64) - CLI version
- ✅ macOS (amd64) - CLI version (with build flags)

## 📋 Roadmap

- [ ] Windows GUI support
- [ ] macOS GUI support  
- [ ] Service detection on open ports
- [ ] Export results to JSON/CSV
- [ ] Network topology mapping
- [ ] Custom scan profiles
- [ ] Plugin system for extensions

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### Development Setup
1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Fyne](https://fyne.io/) - Cross-platform GUI framework
- [go-ping](https://github.com/go-ping/ping) - ICMP ping implementation
- Bootstrap team for color inspiration

## 📞 Support

If you encounter any issues or have questions:
- 🐛 [Open an issue](https://github.com/AwsThamer/NetworkScanner/issues)
- 💬 [Start a discussion](https://github.com/AwsThamer/NetworkScanner/discussions)

---

⭐ If you find this project useful, please consider giving it a star!

**Made with ❤️ by [AwsThamer](https://github.com/AwsThamer)**

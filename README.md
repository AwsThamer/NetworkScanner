# Network Scanner

A cross-platform network scanning tool built with Go and Fyne GUI framework.

## Features

- **Enhanced Ping Range Functionality**: 
  - 🎯 **Network Presets**: One-click buttons for common networks (192.168.1.0/24, 10.0.0.0/24, 172.16.0.0/24)
  - 🌍 **Custom IP Ranges**: Enter ranges like 192.168.1.1-192.168.1.50 for targeted scanning
  - 🟢 **Color-coded Results**: Green for responsive hosts, red for non-responsive
  - 📊 **Real-time Progress**: Live response counts and progress tracking
- **Port Scanning**: Scan a range of ports on a target host to identify open services
- **Network Discovery**: Scan an entire network range to find active hosts  
- **Host Ping**: Simple ping functionality to check if a host is reachable
- **Cross-Platform**: Works on both Linux and Windows
- **Enhanced GUI Interface**: 
  - 🎨 **Beautiful Design**: Custom Bootstrap-inspired color theme
  - 🔲 **Card Layout**: Modern card-based interface with gradient backgrounds  
  - 📊 **Rich Results**: Timestamped, color-coded results with visual indicators
  - 🎪 **Professional Header**: Gradient header with branding
  - 📈 **Enhanced Progress**: Progress bars with percentage display
  - 🎮 **Interactive Elements**: Dynamic buttons that change during scanning
  - 💡 **User Guidance**: Helpful tips and visual feedback throughout
- **Real-time Progress**: Progress bar and status updates during scans
- **Concurrent Scanning**: Fast scanning with controlled concurrency

## Build Status

✅ **Linux CLI Version**: Working  
✅ **Linux GUI Version**: Working  
✅ **Windows CLI Version**: Working (cross-compiled)  
⚠️ **Windows GUI Version**: Limited (requires Windows build environment for GUI dependencies)

## Quick Start

The project includes both GUI and CLI versions:

### GUI Version (Linux)
```bash
./bin/network-scanner-gui
```

### CLI Version (Cross-platform)
```bash
# Ping a host
./bin/network-scanner-cli ping google.com

# Scan ports
./bin/network-scanner-cli portscan 192.168.1.1 80 443

# Scan network
./bin/network-scanner-cli netscan 192.168.1.0/24
```

### Prerequisites
- Go 1.21 or later
- Git

### Building from Source

1. Clone or download this repository
2. Navigate to the project directory
3. Install dependencies:
   ```bash
   go mod tidy
   ```

### Building for Your Platform

**For Linux:**
```bash
go build -o network-scanner main.go
```

**For Windows:**
```bash
GOOS=windows GOARCH=amd64 go build -o network-scanner.exe main.go
```

### Cross-Platform Build

Use the provided build scripts:

**On Linux/macOS:**
```bash
./build.sh
```

**On Windows:**
```batch
build.bat
```

This will create binaries for both Linux and Windows in the `bin/` directory.

## Usage

### Running the Application

**Linux:**
```bash
./network-scanner
```

**Windows:**
```batch
network-scanner.exe
```

### GUI Features (Enhanced)

The new enhanced GUI provides a professional and user-friendly experience:

#### 🎨 **Visual Design**
- **Card-based Layout**: Clean, organized sections for different functions
- **Modern Icons**: Intuitive icons throughout the interface
- **Professional Styling**: Enhanced typography and spacing
- **Responsive Design**: Optimized window sizing and layout

#### 🎯 **User Experience**
- **Tabbed Interface**: Separate tabs for Scanner controls and Results view
- **Smart Input Fields**: Helpful placeholders with examples
- **Real-time Feedback**: Live progress updates with emoji indicators
- **Error Handling**: Clear, friendly error messages

#### 🚀 **Advanced Features**
- **Scan Control**: Start/Stop buttons with dynamic states
- **Progress Tracking**: Detailed progress bars and status updates
- **Result Categorization**: Color-coded results (success/warning/error)
- **Auto-scrolling Results**: Latest results automatically visible

#### 📋 **Interface Sections**
1. **🎯 Target Configuration**: Host/IP and network input fields
2. **🔌 Port Range**: Port range selection with common port suggestions
3. **🚀 Scan Operations**: Port scan, network scan, ping, and clear buttons
4. **📊 Status**: Real-time scanning status and progress display
5. **📋 Results**: Detailed scan results with categorized display

### Examples

- **Port Scan**: Host: `192.168.1.1`, Ports: `1` to `1000`
- **Network Scan**: Network: `192.168.1.0/24`
- **Ping Test**: Host: `google.com`

## Technical Details

### Dependencies
- **Fyne v2**: Cross-platform GUI framework
- **go-ping**: ICMP ping implementation

### Scanning Methods
- **Port Scanning**: TCP connection attempts with 1-second timeout
- **Host Discovery**: ICMP ping with 1-second timeout
- **Concurrent Processing**: Limited concurrent operations to prevent resource exhaustion

### Supported Platforms
- Linux (amd64)
- Windows (amd64)
- macOS (with appropriate build flags)

## Security Considerations

- The application uses unprivileged ping mode for better compatibility
- Port scanning should only be performed on networks you own or have permission to test
- Some firewalls may block or detect scanning activities
- Use responsibly and in accordance with local laws and regulations

## Troubleshooting

### Linux Issues
- If ping doesn't work, try running with sudo (though the app is designed to work without it)
- Ensure firewall allows the application

### Windows Issues
- Windows Defender might flag the executable as suspicious (false positive)
- Add exception in Windows Defender if needed
- Run as administrator if experiencing permission issues

## Contributing

Feel free to submit issues, feature requests, or pull requests to improve the application.

## License

This project is open source. Use at your own risk and responsibility.

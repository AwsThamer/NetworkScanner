#!/bin/bash

echo "🎨 Building Beautiful Network Scanner..."

# Create bin directory if it doesn't exist
mkdir -p bin

# Build Beautiful GUI version for Linux
echo "🎨 Building Beautiful GUI version for Linux..."
go build -o bin/network-scanner-gui main.go

# Build CLI version for Linux
echo "⚡ Building CLI version for Linux..."
go build -o bin/network-scanner-cli cli.go

# Try to build for Windows (may fail in cross-compilation environment)
echo "🪟 Building CLI version for Windows..."
GOOS=windows GOARCH=amd64 go build -o bin/network-scanner-cli-windows.exe cli.go

if [ $? -eq 0 ]; then
    echo "✅ Windows CLI build successful!"
else
    echo "⚠️ Windows CLI build failed (this is expected in some environments)"
fi

echo ""
echo "🎉 Build complete!"
echo "Available binaries:"
echo "  - bin/network-scanner-gui (🎨 Beautiful Linux GUI with custom theme)"
echo "  - bin/network-scanner-cli (⚡ Linux CLI version)"
if [ -f "bin/network-scanner-cli-windows.exe" ]; then
    echo "  - bin/network-scanner-cli-windows.exe (🪟 Windows CLI version)"
fi

echo ""
echo "🚀 Usage examples:"
echo "  Beautiful GUI: ./bin/network-scanner-gui"
echo "  Demo showcase: ./demo_beautiful.sh"
echo "  CLI: ./bin/network-scanner-cli ping google.com"
echo "       ./bin/network-scanner-cli portscan 192.168.1.1 1 1000"
echo "       ./bin/network-scanner-cli netscan 192.168.1.0/24"
echo ""
echo "✨ The Beautiful GUI now features:"
echo "  🎨 Custom color theme with Bootstrap-inspired colors"
echo "  🔲 Card-based layout with gradient backgrounds"
echo "  📊 Timestamped results with color-coded backgrounds"
echo "  🎪 Beautiful gradient header and professional styling"
echo "  📈 Enhanced progress bars with percentage display"
echo "  🎮 Interactive buttons with dynamic state changes"
echo "  💡 Better user guidance and visual feedback"

@echo off
echo Building Network Scanner for Windows...

echo Building for Windows (amd64)...
set GOOS=windows
set GOARCH=amd64
go build -o bin\network-scanner-windows.exe main.go

echo Build complete!
echo Windows binary: bin\network-scanner-windows.exe
pause

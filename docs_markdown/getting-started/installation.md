# Installing Konfigo

Get Konfigo installed and ready to use on your system in just a few minutes.

## Choose Your Installation Method

### 📦 **Pre-built Binaries** (Recommended)

The fastest way to get started. Download the latest release for your platform:

**[📥 Download from GitHub Releases](https://github.com/ebogdum/konfigo/releases)**

#### macOS
```bash
# Download and install (Intel)
curl -L https://github.com/ebogdum/konfigo/releases/latest/download/konfigo-darwin-amd64 -o konfigo
chmod +x konfigo
sudo mv konfigo /usr/local/bin/

# Download and install (Apple Silicon)
curl -L https://github.com/ebogdum/konfigo/releases/latest/download/konfigo-darwin-arm64 -o konfigo
chmod +x konfigo
sudo mv konfigo /usr/local/bin/
```

#### Linux
```bash
# Download and install (x86_64)
curl -L https://github.com/ebogdum/konfigo/releases/latest/download/konfigo-linux-amd64 -o konfigo
chmod +x konfigo
sudo mv konfigo /usr/local/bin/

# Download and install (ARM64)
curl -L https://github.com/ebogdum/konfigo/releases/latest/download/konfigo-linux-arm64 -o konfigo
chmod +x konfigo
sudo mv konfigo /usr/local/bin/
```

#### Windows
```powershell
# Download using PowerShell
Invoke-WebRequest -Uri "https://github.com/ebogdum/konfigo/releases/latest/download/konfigo-windows-amd64.exe" -OutFile "konfigo.exe"

# Move to a directory in your PATH (e.g., C:\Windows\System32)
Move-Item konfigo.exe C:\Windows\System32\
```

### 🏗️ **Build from Source**

If you have Go installed and want the latest development version:

```bash
# Install directly from source
go install github.com/ebogdum/konfigo/cmd/konfigo@latest

# Or clone and build
git clone https://github.com/ebogdum/konfigo.git
cd konfigo
go build -o konfigo cmd/konfigo/main.go
```

### 🍺 **Package Managers** (Coming Soon)

We're working on adding Konfigo to popular package managers:
- Homebrew (macOS/Linux)
- Chocolatey (Windows)
- APT/YUM repositories (Linux)

## Verify Installation

Test that Konfigo is installed correctly:

```bash
konfigo --version
```

You should see output like:
```
konfigo version v1.2.3
```

### "Hello World" Test

Create a simple test to ensure everything works:

```bash
# Create a test file
echo '{"test": "hello world"}' > test.json

# Test Konfigo
konfigo -s test.json

# You should see the JSON output
# Clean up
rm test.json
```

**✅ Success!** If you see the JSON output, Konfigo is working correctly.

## Troubleshooting Installation

### Command Not Found

If you get `command not found: konfigo`:

1. **Check your PATH**: Make sure the directory containing `konfigo` is in your PATH
2. **Verify location**: Run `which konfigo` (macOS/Linux) or `where konfigo` (Windows)
3. **Restart terminal**: Some PATH changes require a new terminal session

### Permission Denied

If you get permission errors:

```bash
# macOS/Linux: Make the binary executable
chmod +x konfigo

# If moving to system directories, use sudo
sudo mv konfigo /usr/local/bin/
```

### Download Issues

If downloads fail:
1. **Check your internet connection**
2. **Try downloading manually** from the [releases page](https://github.com/ebogdum/konfigo/releases)
3. **Check firewall/proxy settings** that might block GitHub

### Version Compatibility

Konfigo requires:
- **No runtime dependencies** - it's a single binary
- **Any modern OS** - Windows 10+, macOS 10.14+, Linux with kernel 2.6+

## Platform-Specific Notes

### macOS
- **Apple Silicon (M1/M2)**: Use the `arm64` version for better performance
- **Intel Macs**: Use the `amd64` version
- **Gatekeeper**: You may need to allow the binary in System Preferences > Security & Privacy

### Linux
- **Package managers**: Standard installation works with all major distributions
- **ARM systems**: ARM64 binaries available for Raspberry Pi and similar devices
- **Alpine Linux**: Static binaries work without additional dependencies

### Windows
- **Windows Defender**: May quarantine the binary - add an exception if needed
- **PowerShell vs CMD**: Commands work in both PowerShell and Command Prompt
- **WSL**: Linux binaries work great in Windows Subsystem for Linux

## Next Steps

Now that Konfigo is installed:

1. **[Quick Start](./quick-start.md)** - Get your first merge working in 5 minutes
2. **[Basic Concepts](./concepts.md)** - Understand how Konfigo works
3. **[User Guide](../guide/)** - Learn common tasks and workflows

## Getting Help

- **Installation issues**: Check our [Troubleshooting Guide](../reference/troubleshooting.md)
- **General questions**: Browse the [FAQ](../reference/faq.md)
- **Bug reports**: Open an issue on [GitHub](https://github.com/ebogdum/konfigo/issues)

Ready to get started? Head to the **[Quick Start](./quick-start.md)** guide!

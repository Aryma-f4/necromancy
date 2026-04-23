# 🧙‍♂️ Necromancy - Advanced Post-Exploitation Shell Manager

<p align="center">
  <img src="logo.png" alt="Necromancy Banner" style="max-width: 100%; height: auto;" width="350" />
</p>

[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Latest Release](https://img.shields.io/github/v/release/Aryma-f4/necromancy?display_name=tag&label=release)](https://github.com/Aryma-f4/necromancy/releases/latest)

**Necromancy** is a powerful post-exploitation shell manager written in Go, designed for penetration testers and red team operators. It provides comprehensive reverse shell management, advanced post-exploitation modules, and an intuitive terminal-based interface.

## 🌟 Key Features

### 🔧 Core Capabilities
- **🖥️ Multi-Platform Support** - Native binaries for Linux, macOS, and Windows (AMD64 & ARM64)
- **🔗 Session Management** - Handle multiple reverse shells simultaneously with ease
- **🎯 Interactive Terminal** - Raw mode interaction with F12 detach functionality
- **🚀 Auto PTY Upgrade** - Automatic shell upgrade to full PTY for enhanced functionality
- **💾 Session Persistence** - Maintain connections across network interruptions
- **📐 Window Resizing** - Dynamic terminal size synchronization

### 🌐 Network & File Operations
- **📡 Multi-Listener Support** - Listen on multiple ports/interfaces simultaneously
- **🔗 Bind Shell Support** - Connect to listening targets
- **🌐 HTTP File Server** - Built-in file serving capability for quick transfers
- **📤 Upload/Download** - Secure file transfer with simple commands
- **💾 In-Memory Execution** - Run scripts without touching disk

### 🎯 Post-Exploitation Arsenal
- **🔍 PEASS Suite** - LinPEAS and WinPEAS integration for comprehensive enumeration
- **⚡ Linux Exploit Suggester** - Automated exploit recommendations
- **📋 LSE (Linux Smart Enumeration)** - Advanced Linux enumeration techniques
- **🥔 Potato Exploits** - Windows privilege escalation methods
- **🚇 Tunneling Tools** - Chisel, Ligolo, Ngrok integration for pivoting

### 🎨 User Experience
- **📺 Tview Dashboard** - Modern terminal-based UI with intuitive navigation
- **📚 Module Browser** - Easy access to all post-exploitation modules
- **📊 Session List** - Visual session management with detailed information
- **🚀 Payload Generator** - Built-in reverse shell payloads for quick deployment

## 📊 Use Cases

### Basic Reverse Shell Workflow
```mermaid
graph TD
    A[Start Necromancy] --> B[Configure Listener]
    B --> C[Generate Payloads]
    C --> D[Execute on Target]
    D --> E[Receive Connection]
    E --> F[Interactive Session]
    F --> G[Execute Commands]
    G --> H[File Operations]
    H --> I[Module Execution]
```

### Session Management Flow
```mermaid
graph LR
    A[Multiple Sessions] --> B[Session Browser]
    B --> C[Select Session]
    C --> D[Interactive Shell]
    D --> E[File Manager]
    D --> F[Module Execution]
```

### File Manager Operations
```mermaid
graph TD
    A[File Manager] --> B[Navigate Directories]
    B --> C[View Files]
    C --> D[Download Files]
    C --> E[Upload Files]
    C --> F[Delete Files]
    C --> G[Execute Files]
```

### Network Information Flow
```mermaid
graph TD
    A[Network Info] --> B[Detect Local IP]
    B --> C[Check Public IP]
    C --> D[Get Location Data]
    D --> E[Display Results]
    E --> F[Update Payloads]
```

## 🚀 Quick Start Guide

### 📥 Installation Options

#### Option 1: Download Pre-built Binaries
Download the latest release for your platform from [GitHub Releases](https://github.com/Aryma-f4/necromancy/releases):

```bash
# Linux AMD64
wget https://github.com/Aryma-f4/necromancy/releases/latest/download/necromancy-linux-amd64
chmod +x necromancy-linux-amd64

# macOS (Intel)
wget https://github.com/Aryma-f4/necromancy/releases/latest/download/necromancy-macos-amd64
chmod +x necromancy-macos-amd64

# macOS (Apple Silicon)
wget https://github.com/Aryma-f4/necromancy/releases/latest/download/necromancy-macos-arm64
chmod +x necromancy-macos-arm64

# Windows
# Download necromancy-windows-amd64.exe from releases
```

#### Option 2: Install with Go
```bash
# Install directly from source
go install github.com/Aryma-f4/necromancy@latest

# Or install specific version
go install github.com/Aryma-f4/necromancy@v1.2.0
```

#### Option 3: Build from Source
```bash
# Clone the repository
git clone https://github.com/Aryma-f4/necromancy.git
cd necromancy

# Build for current platform
go build -o necromancy .

# Or build for all platforms
./build-multi-platform.sh
```

## 🎮 Interactive Commands

### 🏠 Main Menu
- `s` - View active sessions
- `p` - Show reverse shell payloads
- `m` - Browse available post-exploitation modules
- `i` - List network interfaces
- `n` - Show network information
- `q` - Exit application

### 🔗 Session Management
- `interact <ID>` - Connect to specific session
- `f` - Open file manager for selected session (btop-like UI)
- `kill <ID>` - Terminate specific session
- `kill *` - Terminate all sessions
- `upload <local> <remote>` - Upload file to target
- `download <remote> <local>` - Download file from target

### 📋 Advanced Features
- **F12** - Detach from current session (return to main menu)
- **Ctrl+C** - Send interrupt to remote shell
- **Ctrl+D** - Send EOF to remote shell

## ⚙️ Command Line Options

```
Usage: ./necromancy [options]

Options:
  -p, --ports string     Port(s) to listen on (default "4444")
  -s, --serve string     Directory to serve via HTTP file server
  -w, --web-port int     HTTP server port (default 8000)
  -i, --interface string Local interface to bind (default "0.0.0.0")
  -c, --connect string   Connect to bind shell host
  -m, --maintain int     Keep N sessions per target
  -L, --no-log          Disable session log files
  -U, --no-upgrade      Disable shell auto-upgrade
  -h, --help            Show this help message
```

## 🎯 Basic Usage Examples

```bash
# Start listener on default port (4444)
./necromancy

# Start listener on custom port
./necromancy -p 8080

# Start with HTTP file server on port 8000
./necromancy -p 4444 -s /path/to/files -w 8000

# Connect to bind shell
./necromancy -c target.com -p 4444

# Multi-listener setup
./necromancy -p 4444,4445,4446
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## ⚠️ Legal Notice

**IMPORTANT**: This tool is intended for authorized penetration testing and security research purposes only. Users are responsible for complying with all applicable laws and regulations. The authors assume no liability for misuse of this software.

- **Educational Use**: Designed for learning and professional development
- **Authorized Testing**: Only use on systems you own or have explicit permission to test
- **Responsible Disclosure**: Report security vulnerabilities responsibly
- **Compliance**: Follow applicable laws and regulations

---

**Version**: 1.2.0  
**Repository**: https://github.com/Aryma-f4/necromancy  
**License**: MIT
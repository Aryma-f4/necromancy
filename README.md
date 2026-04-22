# 🧙‍♂️ Necromancy - Advanced Shell Manager

[![Release](https://github.com/Aryma-f4/necromancy/workflows/Release/badge.svg)](https://github.com/Aryma-f4/necromancy/actions)
[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Version](https://img.shields.io/badge/version-1.0.0-brightgreen.svg)](https://github.com/Aryma-f4/necromancy/releases)

**Necromancy** is a powerful post-exploitation shell manager written in Go, designed for penetration testers and red team operators. It provides comprehensive reverse shell management, advanced post-exploitation modules, and an intuitive terminal-based interface.

## 🌟 Key Features

### 🔧 Core Capabilities
- **🖥️ Multi-Platform Support** - Native binaries for Linux, macOS, and Windows (AMD64 & ARM64)
- **🔗 Session Management** - Handle multiple reverse shells simultaneously with ease
- **🎯 Interactive Terminal** - Raw mode interaction with F12 detach functionality
- **🚀 Auto PTY Upgrade** - Automatic shell upgrade to full PTY for enhanced functionality
- **💾 Session Persistence** - Maintain connections across network interruptions
- **📐 Window Resizing** - Dynamic terminal size synchronization
- **🛡️ OSCP-Safe Mode** - Compliance mode for certification exams

### 🌐 Network & File Operations
- **📡 Multi-Listener Support** - Listen on multiple ports/interfaces simultaneously
- **🔗 Bind Shell Support** - Connect to listening targets
- **🌐 HTTP File Server** - Built-in file serving capability for quick transfers
- **📤 Upload/Download** - Secure base64-based file transfer
- **💾 In-Memory Execution** - Run scripts without touching disk
- **🔄 Port Forwarding** - Local port forwarding capabilities

### 🎯 Post-Exploitation Arsenal
- **🔍 PEASS Suite** - LinPEAS and WinPEAS integration for comprehensive enumeration
- **⚡ Linux Exploit Suggester** - Automated exploit recommendations
- **📋 LSE (Linux Smart Enumeration)** - Advanced Linux enumeration techniques
- **🥔 Potato Exploits** - Windows privilege escalation methods
- **🚇 Tunneling Tools** - Chisel, Ligolo, Ngrok integration for pivoting
- **🎯 Meterpreter Integration** - Upgrade to Metasploit sessions
- **🧹 Cleanup Module** - Remove tracks and artifacts from targets
- **🦹 Traitor** - Automated Linux privilege escalation
- **🔓 UAC Bypass** - Windows UAC bypass techniques
- **🕰️ Panix** - Linux persistence via systemd
- **🧠 Process Memory Dump** - Linux memory analysis capabilities

### 🎨 User Experience
- **📺 Tview Dashboard** - Modern terminal-based UI with intuitive navigation
- **🎨 ASCII Banner** - Customizable colored banner support (BBCode format)
- **📚 Module Browser** - Easy access to all post-exploitation modules
- **📊 Session List** - Visual session management with detailed information
- **🚀 Payload Generator** - Built-in reverse shell payloads for quick deployment

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

#### Option 2: Build from Source
```bash
# Clone the repository
git clone https://github.com/Aryma-f4/necromancy.git
cd necromancy

# Build for current platform
go build -o necromancy .

# Or build for all platforms
./build-multi-platform.sh
```

### 🎯 Basic Usage Examples

```bash
# Start listener on default port (4444)
./necromancy

# Start listener on custom port
./necromancy -p 8080

# Start with HTTP file server on port 8000
./necromancy -p 4444 -s /path/to/files -w 8000

# Connect to bind shell
./necromancy -c target.com -p 4444

# OSCP-safe mode (no advanced features)
./necromancy -O

# Multi-listener setup
./necromancy -p 4444,4445,4446
```

## 🎮 Interactive Commands

### 🏠 Main Menu
- `s` - View active sessions
- `p` - Show reverse shell payloads
- `m` - Browse available post-exploitation modules
- `i` - List network interfaces
- `q` - Exit application

### 🔗 Session Management
- `interact <ID>` - Connect to specific session
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
  -O, --oscp-safe       Enable OSCP-safe mode
  -h, --help            Show this help message
```

## 🎨 Custom ASCII Banner

Necromancy supports custom ASCII banners with BBCode color formatting. Create an `ascii.txt` file in the same directory:

```bbcode
[color=#FF0000]┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓[/color]
[color=#00FF00]┃            Necromancy - Advanced Shell Manager       ┃[/color]
[color=#0000FF]┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛[/color]
[color=#808080]    Version 1.0.0 | https://github.com/Aryma-f4/necromancy[/color]
```

Colors are automatically converted to terminal ANSI colors for maximum compatibility.

## 🔧 Development

### 📁 Project Structure
```
necromancy/
├── core/           # Core functionality (sessions, networking)
├── modules/        # Post-exploitation modules
├── ui/            # Terminal user interface
├── utils/         # Utility functions (formatting, colors)
├── server/        # HTTP file server
├── pty/           # PTY upgrade functionality
├── ascii.txt      # Custom ASCII banner (optional)
├── main.go        # Application entry point
└── banner_color.go # ASCII banner color processing
```

### 🛠️ Building Custom Modules

Create a new module in `modules/` directory:

```go
package modules

import "github.com/Aryma-f4/necromancy/core"

type MyModule struct{}

func (m *MyModule) Name() string {
    return "my_module"
}

func (m *MyModule) Description() string {
    return "Description of my custom module"
}

func (m *MyModule) Execute(s *core.Session) error {
    // Your module implementation
    script := `echo "Running custom module"`
    _, err := s.Write([]byte(script + "\n"))
    return err
}
```

Register in `modules/module.go`:
```go
mm.Register(&MyModule{})
```

### 🧪 Testing
```bash
# Run basic functionality test
./necromancy --help

# Test with sample listener
./necromancy -p 9999 &
# Connect with netcat: nc localhost 9999
```

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### 📝 Quick Contribution Steps
1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## ⚠️ Legal Notice

**IMPORTANT**: This tool is intended for authorized penetration testing and security research purposes only. Users are responsible for complying with all applicable laws and regulations. The authors assume no liability for misuse of this software.

- **Educational Use**: Designed for learning and professional development
- **Authorized Testing**: Only use on systems you own or have explicit permission to test
- **Responsible Disclosure**: Report security vulnerabilities responsibly
- **Compliance**: Follow OSCP and other certification guidelines when applicable

## 🙏 Acknowledgments

- **Original Python Implementation** - Concepts and inspiration from the Python version
- **Go Community** - Excellent libraries and frameworks
- **Security Researchers** - Continuous contributions to the field
- **Penetration Testers** - Real-world feedback and improvements

---

<div align="center">

**⭐ If you find Necromancy useful, please star the repository! ⭐**

[⭐ Star on GitHub](https://github.com/Aryma-f4/necromancy) • [🐛 Report Bug](https://github.com/Aryma-f4/necromancy/issues) • [💡 Request Feature](https://github.com/Aryma-f4/necromancy/issues)

</div>
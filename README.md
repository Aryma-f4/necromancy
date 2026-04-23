# 🧙‍♂️ Necromancy - Advanced Post-Exploitation Shell Manager

<p align="center">
  <img src="logo.png" alt="Necromancy Banner" style="max-width: 100%; height: auto;" width="350" />
</p>

[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Latest Release](https://img.shields.io/github/v/release/Aryma-f4/necromancy?display_name=tag&label=release)](https://github.com/Aryma-f4/necromancy/releases/latest)

**Necromancy** is a powerful post-exploitation shell manager written in Go, designed for penetration testers and red team operators. It provides comprehensive reverse shell management, advanced post-exploitation modules, and an intuitive terminal-based interface.

## 🚀 Quick Navigation

<div align="center">

### 📖 [Installation Guide](#quick-start) | 🎯 [Usage Examples](#usage-examples) | 📚 [Documentation](#documentation)

### 🔧 [Features](#features) | ⚡ [Quick Start](#quick-start) | 🛠️ [Command Line Options](#command-line-options)

</div>

> **💡 Tip**: Click any link above to jump directly to that section!

---

## 📋 Table of Contents

- [🌟 Key Features](#features)
- [📊 Use Cases](#use-cases)
- [🚀 Quick Start](#quick-start)
- [🎮 Interactive Commands](#interactive-commands)
- [⚙️ Command Line Options](#command-line-options)
- [🎯 Basic Usage Examples](#usage-examples)
- [📚 Documentation](#documentation)
- [📄 License](#license)
- [⚠️ Legal Notice](#legal-notice)

---

<a id="features"></a>
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
- **🔓 UAC Bypass** - Windows UAC bypass techniques
- **🕰️ Panix** - Linux persistence via systemd
- **🧠 Process Memory Dump** - Linux memory analysis capabilities

### 🎨 User Experience
- **📺 Tview Dashboard** - Modern terminal-based UI with intuitive navigation
- **📚 Module Browser** - Easy access to all post-exploitation modules
- **📊 Session List** - Visual session management with detailed information
- **🚀 Payload Generator** - Built-in reverse shell payloads with automatic IP replacement
- **📁 File Manager** - Btop-like file management interface with full CRUD operations
- **🌐 Network Info** - Automatic IP detection and location services with payload updates

### 🔧 Advanced Features
- **🎯 Payload Updates**: Automatically replaces `YOUR_IP` with actual IP addresses in generated payloads
- **🌍 Multi-IP Support**: Uses public IP when available, falls back to local IP
- **📡 Port Configuration**: Dynamically updates payloads based on configured listening ports
- **🔄 Real-time Updates**: Payloads refresh automatically when network information changes

## 📊 Use Cases

### Basic Reverse Shell Workflow
```mermaid
graph TD
    subgraph "Penetration Tester"
        A[👤 Attacker] -->|"Uses"| B[🧙‍♂️ Necromancy]
    end
    
    subgraph "Target System"
        C[🖥️ Target] -->|"Connects to"| D[🌐 Listener]
        E[💻 Shell] -->|"Provides"| F[🔓 Access]
    end
    
    B -->|"Configures"| G[⚙️ Listener Setup]
    G -->|"Generates"| H[📄 Payloads]
    H -->|"Executed on"| C
    C -->|"Establishes"| I[🔗 Reverse Connection]
    I -->|"Creates"| J[🎯 Interactive Session]
    J -->|"Enables"| K[🛠️ Post-Exploitation]
    K -->|"Includes"| L[📁 File Operations]
    K -->|"Includes"| M[🔍 Module Execution]
    
    style A fill:#4CAF50,stroke:#2E7D32,color:#fff
    style B fill:#9C27B0,stroke:#6A1B9A,color:#fff
    style C fill:#FF5722,stroke:#D84315,color:#fff
    style J fill:#2196F3,stroke:#1565C0,color:#fff
    style K fill:#FFC107,stroke:#F57C00,color:#000
```

### Session Management Flow
```mermaid
graph LR
    subgraph "Session Management"
        A[📊 Session Browser] --> B[🖱️ Select Session]
        B --> C[🔗 Interactive Shell]
        C --> D[🛠️ Post-Exploitation Tools]
    end
    
    subgraph "Available Tools"
        D --> E[📁 File Manager]
        D --> F[🔍 Module Browser]
        D --> G[⚙️ System Tools]
        D --> H[📡 Network Tools]
    end
    
    subgraph "Session Operations"
        I[📝 Session Logs] --> A
        J[🔒 Session Security] --> A
        K[⏰ Session Timing] --> A
    end
    
    style A fill:#2196F3,stroke:#1565C0,color:#fff
    style C fill:#4CAF50,stroke:#2E7D32,color:#fff
    style D fill:#FF9800,stroke:#F57C00,color:#fff
    style E fill:#9C27B0,stroke:#6A1B9A,color:#fff
    style F fill:#E91E63,stroke:#C2185B,color:#fff
```

### File Manager Operations
```mermaid
graph TD
    subgraph "User Interface"
        A[👤 User] -->|"Opens"| B[📁 File Manager]
    end
    
    subgraph "Navigation Operations"
        B -->|"Browse"| C[📂 Directory Navigation]
        B -->|"View"| D[📋 File Listing]
    end
    
    subgraph "File Operations"
        D -->|"Select"| E[📄 File Selection]
        E -->|"Actions"| F[⚙️ File Actions]
        
        F -->|"Download"| G[⬇️ Download File]
        F -->|"Upload"| H[⬆️ Upload File]
        F -->|"Delete"| I[🗑️ Delete File]
        F -->|"Execute"| J[▶️ Execute File]
        F -->|"Edit"| K[✏️ Edit File]
        F -->|"Copy"| L[📋 Copy File]
    end
    
    subgraph "System Operations"
        M[🆕 Create New] --> B
        N[📊 Properties] --> E
        O[🔍 Search] --> D
    end
    
    style A fill:#4CAF50,stroke:#2E7D32,color:#fff
    style B fill:#2196F3,stroke:#1565C0,color:#fff
    style F fill:#FF9800,stroke:#F57C00,color:#000
    style G fill:#8BC34A,stroke:#558B2F,color:#fff
    style H fill:#8BC34A,stroke:#558B2F,color:#fff
    style I fill:#F44336,stroke:#C62828,color:#fff
    style J fill:#9C27B0,stroke:#6A1B9A,color:#fff
```

### Network Information Flow & Payload Updates
```mermaid
graph TD
    subgraph "Network Detection Process"
        A[🌐 Network Info System] -->|"Detects"| B[🏠 Local IP Detection]
        B -->|"Queries"| C[🌍 Public IP Services]
        C -->|"Retrieves"| D[📍 Location Data]
        D -->|"Provides"| E[🎯 IP Information]
    end
    
    subgraph "External Services"
        F[☁️ IPify API] --> C
        G[☁️ AWS CheckIP] --> C
        H[☁️ IfConfig.me] --> C
        I[☁️ ICanHazIP] --> C
    end
    
    subgraph "Information Display"
        J[📊 Network Interfaces] --> A
        K[🗺️ Geographic Location] --> D
        L[🏢 ISP Information] --> D
    end
    
    subgraph "Payload Generation Process"
        E -->|"Updates"| M[📝 Payload Templates]
        M -->|"Replaces YOUR_IP with"| N[🎯 Actual IP Address]
        N -->|"Generates"| O[📄 Updated Payloads]
        
        M --> P[⚙️ Bash Payload]
        M --> Q[🐍 Python Payload]
        M --> R[🕸️ Netcat Payload]
        M --> S[💎 PowerShell Payload]
    end
    
    subgraph "Payload Types (Auto-Updated)"
        O --> T[🔗 Reverse Shell Commands]
        T --> U[🖥️ Linux/Unix Payloads]
        T --> V[🪟 Windows Payloads]
        T --> W[🐧 Cross-Platform Payloads]
    end
    
    style A fill:#2196F3,stroke:#1565C0,color:#fff
    style B fill:#4CAF50,stroke:#2E7D32,color:#fff
    style C fill:#FF9800,stroke:#F57C00,color:#000
    style D fill:#9C27B0,stroke:#6A1B9A,color:#fff
    style E fill:#00BCD4,stroke:#00838F,color:#fff
    style M fill:#795548,stroke:#4E342E,color:#fff
    style N fill:#607D8B,stroke:#37474F,color:#fff
    style O fill:#E91E63,stroke:#C2185B,color:#fff
```

<a id="quick-start"></a>
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

<a id="interactive-commands"></a>
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

<a id="command-line-options"></a>
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

<a id="usage-examples"></a>
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

## 🚀 Quick Start

### 1. Download & Install
```bash
# Download latest release
wget https://github.com/Aryma-f4/necromancy/releases/latest/download/necromancy-linux-amd64
chmod +x necromancy-linux-amd64

# Or install with Go
go install github.com/Aryma-f4/necromancy@latest
```

### 2. Start Listener
```bash
# Start on default port
./necromancy

# Or specify custom port
./necromancy -p 8080
```

### 3. Generate Payloads
```bash
# Show available payloads
> p

# Copy payload for your target
# Execute on target system
```

### 4. Interact with Sessions
```bash
# List active sessions
> s

# Connect to session
> interact 1

# Open file manager
> f
```

<a id="documentation"></a>
## 📚 Documentation

For comprehensive documentation, please refer to:

- **[📖 Full Documentation](Documentation.md)** - Complete feature guide and usage examples
- **[⚡ Quick Reference](QUICK_REFERENCE.md)** - Essential commands and workflows
- **[⚙️ Configuration](CONFIGURATION.md)** - Configuration examples and advanced settings
- **[🤖 AI Agent Guide](AGENTS.md)** - Technical documentation for AI assistants
- **[🤝 Contributing](CONTRIBUTING.md)** - How to contribute to the project
- **[🔒 Security Policy](SECURITY.md)** - Security reporting and best practices

<a id="license"></a>
## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

<a id="legal-notice"></a>
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

<div align="right">

### [⬆️ Back to Top](#-quick-navigation)

</div>
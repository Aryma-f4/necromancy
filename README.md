# Necromancy - Advanced Shell Manager

![Release](https://github.com/Aryma-f4/necromancy/workflows/Release/badge.svg)
![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)
![License](https://img.shields.io/badge/license-MIT-green.svg)

**Necromancy** is an advanced post-exploitation shell manager written in Go, inspired by the original Python implementation. It provides a comprehensive solution for managing reverse shells, post-exploitation modules, and penetration testing workflows.

## 🌟 Features

### Core Functionality
- ✅ **Multi-Platform Support** - Linux, macOS, Windows (AMD64 & ARM64)
- ✅ **Session Management** - Handle multiple reverse shells simultaneously
- ✅ **Interactive Terminal** - Raw mode interaction with F12 detach functionality
- ✅ **Auto PTY Upgrade** - Automatic shell upgrade to full PTY
- ✅ **Session Persistence** - Maintain connections across network interruptions
- ✅ **Window Resizing** - Dynamic terminal size synchronization
- ✅ **OSCP-Safe Mode** - Compliance mode for certification exams

### Network & File Transfer
- ✅ **Multi-Listener Support** - Listen on multiple ports/interfaces
- ✅ **Bind Shell Support** - Connect to listening targets
- ✅ **HTTP File Server** - Built-in file serving capability
- ✅ **Upload/Download** - Base64-based file transfer
- ✅ **In-Memory Execution** - Run scripts without touching disk
- ✅ **Port Forwarding** - Local port forwarding capabilities

### Post-Exploitation Modules
- ✅ **PEASS Suite** - LinPEAS and WinPEAS integration
- ✅ **Linux Exploit Suggester** - Automated exploit recommendations
- ✅ **LSE (Linux Smart Enumeration)** - Comprehensive Linux enumeration
- ✅ **Potato Exploits** - Windows privilege escalation
- ✅ **Tunneling Tools** - Chisel, Ligolo, Ngrok integration
- ✅ **Meterpreter Integration** - Upgrade to Metasploit sessions
- ✅ **Cleanup Module** - Remove tracks and artifacts
- ✅ **Traitor** - Linux privilege escalation
- ✅ **UAC Bypass** - Windows UAC bypass techniques
- ✅ **Panix** - Linux persistence via systemd
- ✅ **Process Memory Dump** - Linux memory analysis

### User Interface
- ✅ **Tview Dashboard** - Modern terminal-based UI
- ✅ **ASCII Banner** - Customizable colored banner support
- ✅ **Module Browser** - Easy access to all modules
- ✅ **Session List** - Visual session management
- ✅ **Payload Generator** - Built-in reverse shell payloads

## 🚀 Quick Start

### Download Pre-built Binaries

Download the latest release for your platform:

- **Linux AMD64**: `necromancy-linux-amd64`
- **Linux ARM64**: `necromancy-linux-arm64` 
- **macOS AMD64**: `necromancy-macos-amd64`
- **macOS ARM64 (M1/M2)**: `necromancy-macos-arm64`
- **Windows AMD64**: `necromancy-windows-amd64.exe`

### Build from Source

```bash
# Clone the repository
git clone https://github.com/Aryma-f4/necromancy.git
cd necromancy

# Build for current platform
go build -o necromancy .

# Or use the multi-platform build script
./build-multi-platform.sh
```

### Basic Usage

```bash
# Start listener on default port (4444)
./necromancy

# Start listener on custom port
./necromancy -p 8080

# Start with HTTP file server
./necromancy -p 4444 -s /path/to/files

# Connect to bind shell
./necromancy -c target.com -p 4444

# OSCP-safe mode
./necromancy -O
```

## 📖 Command Reference

### Main Menu Commands
- `s` - View active sessions
- `p` - Show reverse shell payloads
- `m` - Browse available modules
- `i` - List network interfaces
- `q` - Exit application

### Session Commands
- `interact <ID>` - Connect to session
- `kill <ID>` - Terminate specific session
- `kill *` - Terminate all sessions
- `upload <local> <remote>` - Upload file
- `download <remote> <local>` - Download file

### CLI Options
- `-p, --ports` - Port to listen on (default: 4444)
- `-s, --serve` - Serve directory via HTTP
- `-i, --interface` - Local interface to bind (default: 0.0.0.0)
- `-c, --connect` - Connect to bind shell host
- `-m, --maintain` - Keep N sessions per target
- `-L, --no-log` - Disable session logging
- `-U, --no-upgrade` - Disable auto PTY upgrade
- `-O, --oscp-safe` - Enable OSCP-safe mode

## 🎨 ASCII Banner

Necromancy supports custom ASCII banners with BBCode color formatting. Place your banner in `ascii.txt` with color tags:

```
[color=#FF0000]Your[/color][color=#00FF00]Colored[/color][color=#0000FF]Banner[/color]
```

## 🔧 Development

### Project Structure
```
necromancy/
├── core/          # Core functionality (sessions, networking)
├── modules/       # Post-exploitation modules
├── ui/           # Terminal user interface
├── utils/        # Utility functions (formatting, colors)
├── server/       # HTTP file server
├── pty/          # PTY upgrade functionality
├── ascii.txt     # Custom ASCII banner (optional)
└── main.go       # Application entry point
```

### Adding New Modules

Create a new module in `modules/` directory:

```go
type MyModule struct{}

func (m *MyModule) Name() string {
    return "my_module"
}

func (m *MyModule) Description() string {
    return "Description of my module"
}

func (m *MyModule) Execute(s *core.Session) error {
    // Module implementation
    return nil
}
```

Register in `modules/module.go`:
```go
mm.Register(&MyModule{})
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## ⚠️ Disclaimer

This tool is intended for authorized penetration testing and research purposes only. Users are responsible for complying with all applicable laws and regulations. The authors assume no liability for misuse of this software.

## 🙏 Acknowledgments

- Original Python implementation concepts
- Go community for excellent libraries
- Security researchers and penetration testers

---

**⭐ Star this repository if you find it useful!**
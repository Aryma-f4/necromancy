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

### Banner Functions
- `getBannerFromFile()`: Reads and parses ascii.txt
- `parseBBCodeForTview()`: Converts BBCode to tview format
- `printColoredBanner()`: Displays banner with colors

## 🔧 Module Development

### Module Interface
```go
type Module interface {
    Name() string
    Description() string
    Execute(s *core.Session) error
}
```

### Module Registration
```go
mm.Register(&MyModule{})
```

### Example Module Structure
```go
type MyModule struct{}

func (m *MyModule) Name() string {
    return "my_module"
}

func (m *MyModule) Description() string {
    return "Module description"
}

func (m *MyModule) Execute(s *core.Session) error {
    script := `echo "Running module"`
    _, err := s.Write([]byte(script + "\n"))
    return err
}
```

## 🌍 Multi-Platform Support

### Platform Matrix
- **Linux**: amd64, arm64
- **macOS**: amd64 (Intel), arm64 (Apple Silicon)
- **Windows**: amd64

### Build Constraints
- `signal_unix.go`: Unix-specific signal handling
- `signal_windows.go`: Windows compatibility layer
- `build-multi-platform.sh`: Cross-platform build script

### CI/CD Pipeline
- GitHub Actions workflow
- Automatic binary releases
- SHA256 checksums
- Multi-architecture support

## 🚀 Command Line Interface

### Installation via Go
```bash
# Install latest version
go install github.com/Aryma-f4/necromancy@latest

# Install specific version
go install github.com/Aryma-f4/necromancy@v1.1.0

# Build with version info
go build -ldflags="-s -w -X main.Version=1.1.0 -X main.BuildDate=$(date -u +%Y-%m-%d)" -o necromancy .
```

### Basic Usage
```bash
./necromancy -p 4444                    # Single port
./necromancy -p 4444,4445,4446           # Multiple ports
./necromancy -c target.com -p 4444       # Connect to bind shell
./necromancy -s /path/to/files -w 8000   # HTTP file server
```

### Interactive Commands
- `s`: View sessions
- `p`: Show payloads
- `m`: Browse modules
- `i`: List interfaces
- `interact <ID>`: Connect to session
- `f`: Open file manager for selected session (btop-like UI)
- `kill <ID>`: Terminate session

## 🛠️ Utility Functions (`utils/`)

### Size Formatting
```go
size := NewSize(1024*1024) // 1MB
fmt.Println(size.String()) // "1.0 MB"
```

### Table Formatting
```go
table := NewTable([]string{"Name", "Size", "Type"})
table.AddRow([]string{"file.txt", "1.2 MB", "text"})
fmt.Println(table.String())
```

### Progress Bar
```go
pb := NewProgressBar(100, "Downloading", 50)
pb.Update(50)
fmt.Println(pb.String()) // "Downloading: [=====     ] 50.0%"
```

### Color Utilities
```go
fmt.Println(Red("Error message"))
fmt.Println(Green("Success message"))
fmt.Println(Blue("Info message"))
```

## 🔒 Security Considerations

### OSCP-Safe Mode
- Disables advanced features
- Compliant with certification requirements
- Maintains core functionality

### Session Security
- No hardcoded credentials
- Secure file transfer methods
- Proper session isolation

### Module Safety
- Input validation
- Safe command execution
- Error handling

## 🐛 Common Issues & Solutions

### Build Issues
- **Windows SIGWINCH**: Use build constraints
- **Missing dependencies**: Check go.mod
- **Cross-compilation**: Use proper GOOS/GOARCH

### Runtime Issues
- **Banner not displaying**: Check ascii.txt permissions
- **Sessions not connecting**: Verify network configuration
- **Modules not loading**: Check module registration

### UI Issues
- **Terminal size**: Implement SIGWINCH handling
- **Color support**: Use 256-color palette
- **Keyboard shortcuts**: Verify tcell compatibility

## 📊 Testing Guidelines

### Unit Testing
- Test individual modules
- Verify session management
- Check UI components

### Integration Testing
- Test full workflow
- Verify cross-platform compatibility
- Check module execution

### Manual Testing
- Test with real shells
- Verify file transfers
- Check OSCP-safe mode

## 📝 Code Style Guidelines

### Go Conventions
- Follow standard Go formatting
- Use meaningful variable names
- Add appropriate comments
- Handle errors properly

### Module Development
- Use descriptive names
- Provide clear descriptions
- Include error handling
- Test thoroughly

### UI Development
- Maintain consistency
- Use keyboard shortcuts
- Provide visual feedback
- Handle edge cases

## 🔄 Version Control

### Current Version
- **Version**: 1.1.0
- **Repository**: https://github.com/Aryma-f4/necromancy
- **License**: MIT

### Release Process
1. Update version numbers
2. Test all platforms
3. Create GitHub release
4. Upload binaries
5. Update documentation

## 📚 Additional Resources

### Documentation
- [README.md](README.md) - User documentation
- [ENHANCEMENT_SUMMARY.md](ENHANCEMENT_SUMMARY.md) - Development summary

### External Links
- [GitHub Repository](https://github.com/Aryma-f4/necromancy)
- [Go Documentation](https://golang.org/doc/)
- [tview Documentation](https://github.com/rivo/tview)

---

**Last Updated**: 2026-04-23  
**Version**: 1.1.0  
**Repository**: https://github.com/Aryma-f4/necromancy
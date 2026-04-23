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
    A[Start Necromancy] --> B[Configure Listener]
    B --> C[Generate Payloads]
    C --> D[Execute on Target]
    D --> E[Receive Connection]
    E --> F[Interactive Session]
    F --> G[Execute Commands]
    G --> H[File Operations]
    H --> I[Module Execution]
```

### File Manager Operations
```mermaid
graph LR
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
# 🤝 Contributing to Necromancy

Thank you for your interest in contributing to Necromancy! This document provides guidelines and information for contributors.

## 📋 Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Contribution Guidelines](#contribution-guidelines)
- [Pull Request Process](#pull-request-process)
- [Issue Reporting](#issue-reporting)
- [Development Workflow](#development-workflow)

## 🎯 Code of Conduct

By participating in this project, you agree to abide by our Code of Conduct:

- **Be respectful** to all community members
- **Be inclusive** and welcoming to newcomers
- **Be constructive** in your feedback and discussions
- **Be professional** in all interactions

## 🚀 Getting Started

### Prerequisites

- **Go 1.21+** - [Download Go](https://golang.org/dl/)
- **Git** - [Download Git](https://git-scm.com/downloads)
- **Terminal/Command Line** - Basic familiarity with command line tools

### Quick Start

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/necromancy.git
   cd necromancy
   ```
3. **Build the project**:
   ```bash
   go build -o necromancy .
   ```
4. **Test your build**:
   ```bash
   ./necromancy --help
   ```

## 🔧 Development Setup

### Environment Setup

1. **Set up Go workspace** (if not already done):
   ```bash
   export GOPATH=$HOME/go
   export PATH=$PATH:$GOPATH/bin
   ```

2. **Install development dependencies**:
   ```bash
   go mod download
   go mod tidy
   ```

3. **Run tests** (if available):
   ```bash
   go test ./...
   ```

### Code Quality Tools

We recommend using these tools for code quality:

```bash
# Install Go tools
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/mvdan/gofumpt@latest

# Format code
gofumpt -w .

# Import organization
goimports -w .
```

## 📝 Contribution Guidelines

### What We're Looking For

- **Bug fixes** - Help improve stability
- **Feature enhancements** - Add new capabilities
- **Documentation improvements** - Better guides and examples
- **Performance optimizations** - Make things faster
- **Security improvements** - Enhance security measures
- **Cross-platform support** - Better compatibility

### What We're NOT Looking For

- **Breaking changes** without discussion
- **Large refactors** without prior approval
- **Security vulnerabilities** (report these privately)
- **Malicious code** or backdoors

### Code Style Guidelines

#### Go Conventions

- Follow **standard Go formatting** (`gofmt`)
- Use **meaningful variable names**
- Add **appropriate comments** for complex logic
- Handle **errors properly** - never ignore them
- Keep **functions small** and focused

#### Example Code Style

```go
// Good: Clear function name and proper error handling
func ConnectToSession(sessionID string) (*Session, error) {
    if sessionID == "" {
        return nil, fmt.Errorf("session ID cannot be empty")
    }
    
    session, err := sessionManager.Get(sessionID)
    if err != nil {
        return nil, fmt.Errorf("failed to get session: %w", err)
    }
    
    return session, nil
}

// Bad: Ignoring errors and unclear naming
func getSess(id string) *Session {
    s, _ := sm.Get(id)
    return s
}
```

#### UI Development (tview)

- Maintain **consistent styling** across components
- Use **keyboard shortcuts** for accessibility
- Provide **visual feedback** for user actions
- Handle **edge cases** gracefully
- Keep **color schemes** consistent

### Module Development

When creating new modules:

1. **Follow the interface**:
   ```go
   type Module interface {
       Name() string
       Description() string
       Execute(s *core.Session) error
   }
   ```

2. **Use descriptive names**:
   ```go
   // Good
   func (m *NetworkScannerModule) Name() string {
       return "network_scanner"
   }
   
   // Bad
   func (m *MyModule) Name() string {
       return "m1"
   }
   ```

3. **Provide clear descriptions**:
   ```go
   func (m *NetworkScannerModule) Description() string {
       return "Performs network discovery and port scanning on target systems"
   }
   ```

4. **Handle errors gracefully**:
   ```go
   func (m *NetworkScannerModule) Execute(s *core.Session) error {
       script := `nmap -sS target_network`
       _, err := s.Write([]byte(script + "\n"))
       if err != nil {
           return fmt.Errorf("failed to execute network scan: %w", err)
       }
       return nil
   }
   ```

## 🔄 Pull Request Process

### Before You Start

1. **Check existing issues** to avoid duplicates
2. **Discuss major changes** in an issue first
3. **Ensure your fork is up-to-date** with the main repository

### Creating a Pull Request

1. **Create a feature branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes** following our guidelines

3. **Test your changes**:
   ```bash
   # Build the project
   go build -o necromancy .
   
   # Test basic functionality
   ./necromancy --help
   ```

4. **Commit your changes** with clear messages:
   ```bash
   git commit -m "feat: add network scanner module"
   ```

5. **Push to your fork**:
   ```bash
   git push origin feature/your-feature-name
   ```

6. **Create a Pull Request** on GitHub

### Pull Request Template

When creating a PR, please include:

- **Clear title** describing the change
- **Description** of what was changed and why
- **Testing** performed to verify the changes
- **Screenshots** (if UI changes)
- **Breaking changes** (if any)

### Commit Message Format

We follow conventional commits:

- `feat:` - New features
- `fix:` - Bug fixes
- `docs:` - Documentation changes
- `style:` - Code style changes
- `refactor:` - Code refactoring
- `test:` - Test additions/changes
- `chore:` - Maintenance tasks

Examples:
```
feat: add network scanner module
fix: resolve session timeout issue
docs: update installation guide
```

## 🐛 Issue Reporting

### Bug Reports

When reporting bugs, please include:

- **Environment details** (OS, Go version, etc.)
- **Steps to reproduce** the issue
- **Expected behavior** vs actual behavior
- **Error messages** or logs
- **Screenshots** (if applicable)

### Feature Requests

For feature requests, please provide:

- **Clear description** of the feature
- **Use case** - why is this needed?
- **Proposed implementation** (if you have ideas)
- **Alternatives considered**

## 🏗️ Development Workflow

### For Bug Fixes

1. **Identify the issue** and understand the root cause
2. **Create a test case** that reproduces the bug
3. **Fix the bug** with minimal changes
4. **Verify the fix** works correctly
5. **Ensure no regressions** were introduced

### For New Features

1. **Design the feature** and consider edge cases
2. **Implement incrementally** with small commits
3. **Add tests** for the new functionality
4. **Update documentation** if needed
5. **Get feedback** from maintainers

### For Documentation

1. **Identify gaps** in current documentation
2. **Write clear, concise** explanations
3. **Include examples** where helpful
4. **Test code examples** to ensure they work
5. **Proofread** for grammar and clarity

## 🧪 Testing Guidelines

### Manual Testing

Test these scenarios when making changes:

- **Basic functionality** - Does it start and show help?
- **Session management** - Can you connect and interact?
- **File operations** - Upload/download works?
- **UI navigation** - All menus accessible?
- **Cross-platform** - Works on your target OS?

### Test Commands

```bash
# Basic functionality
./necromancy --help
./necromancy --version

# Network operations
./necromancy -p 4444
./necromancy -l

# File operations (if applicable)
./necromancy -s /tmp -w 8000
```

## 📋 Pull Request Checklist

Before submitting your PR:

- [ ] **Code builds** without errors
- [ ] **Basic functionality** tested
- [ ] **No breaking changes** (or clearly documented)
- [ ] **Code follows** our style guidelines
- [ ] **Commits are clean** and well-described
- [ ] **PR description** is complete
- [ ] **Documentation updated** (if needed)

## 🆘 Getting Help

If you need help:

1. **Check existing issues** and documentation
2. **Ask in discussions** on GitHub
3. **Join our community** (if available)
4. **Create an issue** with your question

## 🙏 Recognition

Contributors will be:

- **Listed in our README** (with permission)
- **Mentioned in release notes**
- **Credited in commit history**

Thank you for contributing to Necromancy! 🎉

---

**Questions?** Feel free to open an issue or start a discussion on GitHub.
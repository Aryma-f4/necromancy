package modules

import (
	"github.com/Aryma-f4/necromancy/core"
)

// ProcessMonitorModule provides process monitoring and background process checking
type ProcessMonitorModule struct{}

func (m *ProcessMonitorModule) Name() string {
	return "process_monitor"
}

func (m *ProcessMonitorModule) Description() string {
	return "Monitor and check background processes, including necromancy instances"
}

func (m *ProcessMonitorModule) Execute(s *core.Session) error {
	script := `echo "[=== Process Monitor ===]"
echo "[*] Checking for running necromancy processes..."
echo ""

# Check for necromancy processes
echo "[+] Necromancy Process Detection:"
if command -v ps >/dev/null 2>&1; then
    NECROMANCY_PROCS=$(ps aux 2>/dev/null | grep -i "necromancy" | grep -v grep | wc -l)
    echo "    Total necromancy processes: $NECROMANCY_PROCS"
    
    if [ "$NECROMANCY_PROCS" -gt 0 ]; then
        echo "    Active necromancy processes:"
        ps aux 2>/dev/null | grep -i "necromancy" | grep -v grep | awk '{print "    - PID: " $2 " | User: " $1 " | CPU: " $3 "% | Memory: " $4 "% | Command: " $11}' | head -5
    fi
elif command -v tasklist >/dev/null 2>&1; then
    NECROMANCY_PROCS=$(tasklist 2>/dev/null | grep -i "necromancy" | wc -l)
    echo "    Total necromancy processes: $NECROMANCY_PROCS"
    
    if [ "$NECROMANCY_PROCS" -gt 0 ]; then
        echo "    Active necromancy processes:"
        tasklist 2>/dev/null | grep -i "necromancy" | head -5 | sed 's/^/    - /'
    fi
else
    echo "    [!] Unable to detect processes - ps/tasklist not available"
fi

echo ""
echo "[+] Background Process Analysis:"
# Check for common background processes
if command -v ps >/dev/null 2>&1; then
    echo "    Background shells: $(ps aux 2>/dev/null | grep -E "(bash|sh|zsh|fish)" | grep -v grep | wc -l)"
    echo "    Network listeners: $(ps aux 2>/dev/null | grep -E "(nc|netcat|socat)" | grep -v grep | wc -l)"
    echo "    SSH processes: $(ps aux 2>/dev/null | grep "ssh" | grep -v grep | wc -l)"
    echo "    Web servers: $(ps aux 2>/dev/null | grep -E "(apache|nginx|httpd)" | grep -v grep | wc -l)"
elif command -v tasklist >/dev/null 2>&1; then
    echo "    Total processes: $(tasklist 2>/dev/null | wc -l)"
    echo "    PowerShell processes: $(tasklist 2>/dev/null | grep -i "powershell" | wc -l)"
    echo "    CMD processes: $(tasklist 2>/dev/null | grep -i "cmd" | wc -l)"
fi

echo ""
echo "[+] Process Tree Analysis:"
if command -v pstree >/dev/null 2>&1; then
    echo "    Process tree (limited view):"
    pstree 2>/dev/null | head -10 | sed 's/^/    /'
elif command -v ps >/dev/null 2>&1; then
    echo "    Parent-child relationships:"
    ps -ef 2>/dev/null | head -5 | awk '{print "    PID: " $2 " | PPID: " $3 " | CMD: " $8}'
fi

echo ""
echo "[+] Network Process Detection:"
# Check for processes listening on common ports
for port in 4444 8080 8000 3000 5000; do
    if command -v lsof >/dev/null 2>&1; then
        LISTENING=$(lsof -i :$port 2>/dev/null | wc -l)
        if [ "$LISTENING" -gt 0 ]; then
            echo "    Port $port: LISTENING ($(lsof -i :$port 2>/dev/null | grep LISTEN | wc -l) processes)"
        fi
    elif command -v netstat >/dev/null 2>&1; then
        LISTENING=$(netstat -tulpn 2>/dev/null | grep :$port | wc -l)
        if [ "$LISTENING" -gt 0 ]; then
            echo "    Port $port: LISTENING ($LISTENING processes)"
        fi
    fi
done

echo ""
echo "[+] Process Status Recommendations:"
echo "    ✓ Use 'ps aux | grep necromancy' to check necromancy processes"
echo "    ✓ Use 'netstat -tulpn' to check network listeners"
echo "    ✓ Use 'kill PID' to terminate specific processes"
echo "    ✓ Use 'pkill necromancy' to terminate all necromancy processes"
if command -v systemctl >/dev/null 2>&1; then
    echo "    ✓ Use 'systemctl status' to check service status"
elif command -v service >/dev/null 2>&1; then
    echo "    ✓ Use 'service --status-all' to check service status"
elif command -v sc >/dev/null 2>&1; then
    echo "    ✓ Use 'sc query' to check Windows services"
fi

echo ""
echo "[=== Process Monitor Complete ===]"
`
	_, err := s.Write([]byte(script + "\n"))
	return err
}

// BackgroundCheckerModule provides specific background process checking
type BackgroundCheckerModule struct{}

func (m *BackgroundCheckerModule) Name() string {
	return "background_checker"
}

func (m *BackgroundCheckerModule) Description() string {
	return "Check if necromancy background processes are still running"
}

func (m *BackgroundCheckerModule) Execute(s *core.Session) error {
	script := `echo "[=== Background Process Checker ===]"
echo "[*] Checking necromancy background processes..."
echo ""

# Quick check for necromancy processes
echo "[+] Quick Status Check:"
if command -v pgrep >/dev/null 2>&1; then
    NECROMANCY_PIDS=$(pgrep -f "necromancy" 2>/dev/null)
    if [ -n "$NECROMANCY_PIDS" ]; then
        echo "    ✅ Necromancy processes found:"
        for pid in $NECROMANCY_PIDS; do
            if [ -f "/proc/$pid/cmdline" ]; then
                CMD=$(cat /proc/$pid/cmdline 2>/dev/null | tr '\0' ' ')
                echo "    - PID: $pid | Command: $CMD"
            else
                echo "    - PID: $pid (details unavailable)"
            fi
        done
    else
        echo "    ❌ No necromancy processes found"
    fi
elif command -v ps >/dev/null 2>&1; then
    NECROMANCY_COUNT=$(ps aux 2>/dev/null | grep -i "necromancy" | grep -v grep | wc -l)
    if [ "$NECROMANCY_COUNT" -gt 0 ]; then
        echo "    ✅ Found $NECROMANCY_COUNT necromancy process(es):"
        ps aux 2>/dev/null | grep -i "necromancy" | grep -v grep | awk '{print "    - PID: " $2 " | User: " $1 " | CPU: " $3 "% | Memory: " $4 "%"}'
    else
        echo "    ❌ No necromancy processes found"
    fi
elif command -v tasklist >/dev/null 2>&1; then
    NECROMANCY_COUNT=$(tasklist 2>/dev/null | grep -i "necromancy" | wc -l)
    if [ "$NECROMANCY_COUNT" -gt 0 ]; then
        echo "    ✅ Found $NECROMANCY_COUNT necromancy process(es):"
        tasklist 2>/dev/null | grep -i "necromancy" | head -3 | sed 's/^/    - /'
    else
        echo "    ❌ No necromancy processes found"
    fi
else
    echo "    [!] Unable to check processes - no suitable tools available"
fi

echo ""
echo "[+] Network Listener Check:"
# Check if necromancy is listening on common ports
for port in 4444 8080 8000 3000; do
    if command -v netstat >/dev/null 2>&1; then
        LISTENERS=$(netstat -tulpn 2>/dev/null | grep :$port | grep LISTEN | wc -l)
        if [ "$LISTENERS" -gt 0 ]; then
            echo "    🔍 Port $port: LISTENING ($LISTENERS process(es))"
            netstat -tulpn 2>/dev/null | grep :$port | grep LISTEN | sed 's/^/      /'
        fi
    elif command -v ss >/dev/null 2>&1; then
        LISTENERS=$(ss -tulpn 2>/dev/null | grep :$port | wc -l)
        if [ "$LISTENERS" -gt 0 ]; then
            echo "    🔍 Port $port: LISTENING ($LISTENERS process(es))"
            ss -tulpn 2>/dev/null | grep :$port | sed 's/^/      /'
        fi
    fi
done

echo ""
echo "[+] Process Health Check:"
if command -v ps >/dev/null 2>&1; then
    # Check if processes are in healthy state (not zombie/defunct)
    ZOMBIE_COUNT=$(ps aux 2>/dev/null | grep -i "necromancy" | grep -E "(Z|defunct)" | wc -l)
    if [ "$ZOMBIE_COUNT" -gt 0 ]; then
        echo "    ⚠️  Found $ZOMBIE_COUNT zombie/defunct necromancy process(es)"
    else
        echo "    ✅ All necromancy processes appear healthy"
    fi
    
    # Check CPU and memory usage
    if [ "$NECROMANCY_COUNT" -gt 0 ] 2>/dev/null; then
        TOTAL_CPU=$(ps aux 2>/dev/null | grep -i "necromancy" | grep -v grep | awk '{sum+=$3} END {print sum}')
        TOTAL_MEM=$(ps aux 2>/dev/null | grep -i "necromancy" | grep -v grep | awk '{sum+=$4} END {print sum}')
        echo "    📊 Total CPU usage: ${TOTAL_CPU:-0}%"
        echo "    📊 Total memory usage: ${TOTAL_MEM:-0}%"
    fi
fi

echo ""
echo "[+] Quick Actions:"
echo "    🔄 To restart necromancy: killall necromancy; ./necromancy -p 4444 &"
echo "    🛑 To stop all necromancy: pkill -f necromancy"
echo "    📋 To view all processes: ps aux | grep necromancy"
echo "    🔍 To check network listeners: netstat -tulpn | grep necromancy"

echo ""
echo "[=== Background Process Check Complete ===]"
`
	_, err := s.Write([]byte(script + "\n"))
	return err
}

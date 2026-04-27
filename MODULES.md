# Module Catalog

This page lists all built-in modules and where they fit in an operator workflow.

## How Modules Run in TUI

1. Open `Sessions` and select a session.
2. Open `Shell Tabs` (optional but recommended).
3. Open `Run Module` and choose the module.
4. Choose target tab (`Tab 1`, `Tab 2`, and so on).
5. Module command is dispatched from that selected tab.

This behavior applies to CLI-style modules such as `peass_auto`, `linpeas`, `winpeas`, and others in the module menu.

## Module List

| Module Key | Description | Platform | Category |
|---|---|---|---|
| `peass_auto` | Auto-select LinPEAS or WinPEAS based on detected OS | Linux/Windows | Enumeration |
| `linpeas` | Linux privilege escalation enumeration | Linux | Enumeration |
| `winpeas` | Windows privilege escalation enumeration | Windows | Enumeration |
| `lse` | Linux Smart Enumeration | Linux | Enumeration |
| `potato` | Placeholder launcher for Windows potato-family techniques | Windows | Privilege Escalation |
| `traitor` | Linux privilege escalation helper | Linux | Privilege Escalation |
| `uac` | Windows UAC bypass helper | Windows | Privilege Escalation |
| `chisel` | Chisel tunneling workflow helper | Multi | Tunneling |
| `ligolo` | Ligolo-ng tunneling workflow helper | Multi | Tunneling |
| `ngrok` | Ngrok tunneling workflow helper | Multi | Tunneling |
| `meterpreter` | Meterpreter upgrade helper | Multi | Session Upgrade |
| `cleanup` | Cleanup traces and shell history | Multi | Cleanup |
| `panix` | Linux persistence helper | Linux | Persistence |
| `linux_procmemdump` | Process memory dump helper | Linux | Forensics |
| `filemanager` | File manager workflow helper | Multi | Operations |
| `redsun` | Windows vulnerability enumeration helper | Windows | Enumeration |
| `bluehammer` | Windows exploitation helper | Windows | Exploitation |
| `payload_obfuscator` | Obfuscated payload generation workflow | Multi | Evasion |
| `enhanced_preflight_recon` | Recon and security-surface checks | Multi | Recon |
| `process_monitor` | Background process monitoring checks | Multi | Monitoring |
| `background_checker` | Check running Necromancy-related background processes | Multi | Monitoring |

## Recommended Usage Pattern

- Use one tab for enumeration modules.
- Use one tab for exploitation/escalation modules.
- Use one tab for cleanup/persistence modules.
- Keep module output context grouped by tab for easier review.

## Notes

- Some modules are placeholders/helpers and may print guidance commands.
- Exact behavior depends on target OS, shell capability, and available binaries on target.

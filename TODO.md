# Necromancy Todo List

## Core Features
- [x] Basic TCP Listener for reverse shells
- [x] Basic Session Multiplexing (handling multiple connections)
- [x] Tview UI Dashboard (Main Menu & Session List)
- [x] Interactive Terminal View (Reading/Writing to shell)
- [x] Raw Mode Terminal Interaction (Suspending UI to drop into shell)
- [x] F12 Detach functionality to return to UI
- [x] Basic PTY Auto-Upgrade injection (`python -c 'import pty...'`)
- [x] Window Resizing Support (SIGWINCH handling & sending `stty rows cols`)
- [x] Bind Shell Support (connecting to listening targets)
- [x] Session Logging (recording shell activity to files with timestamps)
- [x] Session Persistence (Maintain X amount of shells per target)
- [x] OSCP-Safe Mode enforcement
- [x] Multi-Listener Support (listening on multiple ports/interfaces simultaneously)

## File Transfer & Networking
- [x] Basic HTTP File Server (`-s` switch)
- [x] Upload command (local -> remote via base64 echo)
- [x] Download command (remote -> local via base64 cat)
- [x] In-memory script execution (uploading and running without touching disk)
- [x] Local Port Forwarding

## Main Menu Commands
- [x] `sessions` (View active sessions)
- [x] `interact <ID>` (Connect to session)
- [x] `payloads` (Show reverse shell payloads)
- [x] `kill <ID>` / `kill *` (Terminate sessions)
- [x] `listeners` (Manage active listeners)
- [x] `upload` / `download`
- [x] `portfwd`
- [x] `interfaces` (List local network interfaces)

## Post-Exploitation Modules
- [x] Module Framework Architecture
- [x] `peass_ng` (linpeas / winpeas injection)
- [x] `linuxexploitsuggester`
- [x] `lse` (Linux Smart Enumeration)
- [x] `potato` (Windows privilege escalation)
- [x] `chisel` / `ligolo` / `ngrok` (Tunneling/Pivoting modules)
- [x] `meterpreter` (Upgrade to MSF session)
- [x] `cleanup` (Remove tracks on target)

## CLI Options & Arguments
- [x] `-p, --ports`
- [x] `-s, --serve`
- [x] `-i, --interface`
- [x] `-c, --connect`
- [x] `-m, --maintain`
- [x] `-L, --no-log`
- [x] `-U, --no-upgrade`
- [x] `-O, --oscp-safe`

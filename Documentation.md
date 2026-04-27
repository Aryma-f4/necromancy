# Necromancy Documentation

## Overview

Necromancy is a Go-based reverse-shell manager with a terminal UI, multi-session handling, payload generation, logging, and a built-in module registry for post-session workflows. This document is the primary product manual for the repository.

This guide focuses on what the current codebase actually implements. Where behavior is partial, helper-based, or placeholder-driven, that is called out explicitly.

## Table Of Contents

1. [What Necromancy Does](#what-necromancy-does)
2. [Architecture Summary](#architecture-summary)
3. [Installation](#installation)
4. [Runtime Modes](#runtime-modes)
5. [Command-Line Reference](#command-line-reference)
6. [Terminal UI Guide](#terminal-ui-guide)
7. [Sessions And Shell Interaction](#sessions-and-shell-interaction)
8. [Payloads And Network Awareness](#payloads-and-network-awareness)
9. [File Transfer And File Manager](#file-transfer-and-file-manager)
10. [Module System](#module-system)
11. [Logging And Runtime Files](#logging-and-runtime-files)
12. [Configuration Notes](#configuration-notes)
13. [Troubleshooting](#troubleshooting)
14. [Security And Safe Use](#security-and-safe-use)
15. [Related Documents](#related-documents)

## What Necromancy Does

Necromancy is built around a simple operator workflow:

1. Start one or more listeners, or connect to a bind shell.
2. Wait for a shell to land.
3. Manage sessions from the TUI.
4. Interact with the target through raw shell mode or tabbed shell pages.
5. Use built-in helpers for upload, in-memory execution, pivoting guidance, and modules.

### Current capabilities

- Listen on one or many TCP ports.
- Connect to bind-shell targets instead of listening.
- Track multiple sessions at the same time.
- Attempt PTY upgrade automatically for non-Windows sessions.
- Log session output to per-session files.
- Show payload templates with automatic IP replacement.
- Serve files through a built-in HTTP server.
- Run the app with or without the TUI.
- Check for updates and download the latest release asset.

### Important product notes

- Necromancy is TUI-first. It is not structured as a typed command shell with built-in commands such as `interact 1` or `kill 1`; those actions are selected from menu pages.
- Module behavior varies. Some modules dispatch real command sequences, while others are scaffolds, placeholders, or operator helpers.
- The file manager UI is partially implemented. Navigation and some actions work, but several actions still show placeholder status messages.

## Architecture Summary

```text
necromancy/
├── core/       sessions, connection handling, config, logging, network helpers
├── modules/    built-in post-session modules and file-manager logic
├── pty/        PTY auto-upgrade behavior
├── server/     HTTP file server
├── ui/         tview dashboard, pages, and shell tabs
├── updater/    release checks and self-update logic
├── utils/      formatting and network utilities
└── main.go     entry point and flag parsing
```

### Main components

| Path | Purpose |
| --- | --- |
| `main.go` | Startup, flag parsing, listener setup, bind-shell mode, headless mode |
| `core/` | Session manager, OS detection, logs, listeners, config |
| `ui/` | Menu pages, session actions, payload page, shell tabs |
| `modules/` | Module registry, file-manager logic, file-transfer helpers, monitoring helpers |
| `server/` | Built-in HTTP file server for quick staging |
| `updater/` | GitHub release check and binary self-update |

## Installation

### Download a release

Use the latest asset from GitHub Releases.

```bash
# Linux amd64
curl -LO https://github.com/Aryma-f4/necromancy/releases/latest/download/necromancy-linux-amd64
chmod +x necromancy-linux-amd64
./necromancy-linux-amd64
```

### Install with Go

```bash
go install github.com/Aryma-f4/necromancy@latest
```

### Build locally

```bash
git clone https://github.com/Aryma-f4/necromancy.git
cd necromancy
go build -o necromancy .
./necromancy
```

### Multi-platform build script

The repository also includes helper scripts:

- `build.sh`
- `build-multi-platform.sh`
- `Makefile`

## Runtime Modes

### Reverse-shell listener mode

Default mode starts one or more listeners.

```bash
./necromancy
./necromancy -p 4444,4445,4446
./necromancy -p 8080 -i 0.0.0.0
```

### Bind-shell client mode

Provide `-c` to connect outward instead of listening.

```bash
./necromancy -c target.example -p 4444
```

If multiple ports are provided, the application attempts to connect to the same host on each configured port.

### HTTP file-server mode

Provide a directory with `-s`.

```bash
./necromancy -p 4444 -s ./payloads -w 8000
./necromancy -p 4444 -s ./payloads -w 8000 --prefix /static
```

### Headless mode

Headless mode skips the TUI and keeps the listeners running in non-interactive environments.

```bash
./necromancy --headless
```

Useful for:

- VPS environments
- `tmux` or `screen`
- service wrappers
- `nohup` usage
- CI or automation labs

### Update mode

```bash
./necromancy --check-update
./necromancy --update
```

`--update` checks GitHub Releases, downloads the platform-matching binary, and replaces the current executable.

## Command-Line Reference

### Flag table

| Flag | Type | Default | Description | Notes |
| --- | --- | --- | --- | --- |
| `-p` | string | `4444` | Comma-separated ports | Core listener option |
| `-s` | string | empty | Serve a directory over HTTP | Enables file server |
| `-i` | string | `0.0.0.0` | Interface or local IP to bind | Also affects listener address display |
| `-c` | string | empty | Bind-shell host | Switches app into connect mode |
| `-m` | int | `0` | Maintain at least `N` sessions | Uses the first active session to attempt respawn logic |
| `-L` | bool | `false` | Disable session log files | Does not disable `necromancy-go.log` |
| `-U` | bool | `false` | Disable PTY auto-upgrade | Mainly affects shell usability on Unix-like targets |
| `-O` | bool | `false` | Enable OSCP-safe mode | Flag exists; current code does not add extra behavior yet |
| `-w` | int | `8000` | HTTP server port | Used with `-s` |
| `-S` | bool | `false` | Accept only the first created session | Later incoming connections are rejected |
| `-C` | bool | `false` | Do not auto-attach on new sessions | Flag exists; current code path is not wired to additional behavior yet |
| `--prefix` | string | empty | URL prefix for file serving | Example: `/static` |
| `-a`, `--payloads` | bool | `false` | Print payloads and exit | Uses first configured port |
| `-l`, `--interfaces` | bool | `false` | Print interfaces and exit | Convenience output |
| `-v`, `--version` | bool | `false` | Print version and exit | Does not start listeners |
| `--headless` | bool | `false` | Disable TUI | Best for non-interactive shells |
| `--check-update` | bool | `false` | Check for a newer release | Read-only action |
| `--update` | bool | `false` | Download and replace current binary | Modifies current executable |

### Example commands

```bash
# Default listener
./necromancy

# Multiple listeners
./necromancy -p 4444,5555,6666

# File server with URL prefix
./necromancy -s ./loot -w 8080 --prefix /share

# Single-session mode
./necromancy -p 4444 -S

# Show payloads only
./necromancy --payloads -p 9001
```

## Terminal UI Guide

The TUI is built with `tview` and centers the workflow around pages rather than typed subcommands.

### Main dashboard

The root menu contains:

- `Sessions`: active reverse shells and per-session actions
- `Payloads`: payload preview, copy, and refresh
- `Modules`: browse the module catalog without selecting a session
- `Network Info`: current local/public IP information and configured ports
- `Interfaces`: local interface overview
- `Exit`: stop the UI

### Navigation basics

- Arrow keys move through lists.
- `Enter` opens the selected page or action.
- `Esc` returns to the previous page.
- Mouse support is enabled in the TUI, except on the payload page where mouse selection is intentionally freed for terminal copy behavior.

## Sessions And Shell Interaction

### Session list

When a listener receives a connection, the session is added to the session list. Each session tracks:

- a numeric ID,
- remote address,
- session type,
- output history,
- detected operating system,
- optional log file.

### Session action page

Selecting a session opens the action menu:

- `Shell Tabs`
- `Raw Interact`
- `File Manager`
- `Cancel Commands`
- `Run Module`
- `Upload File`
- `In-Memory Exec`
- `Port Forwarding`
- `Kill`

### Raw interact

`Raw Interact` suspends the TUI and hands control to a direct terminal session. If the target is not already marked as PTY-capable and the OS is not Windows, Necromancy attempts an automatic PTY upgrade before interaction.

### Tabbed shell workspace

The shell-tab workspace creates multiple local tabs that all mirror the same remote stream. This is useful for organizing operator notes or keeping separate logical contexts while still driving one shell.

#### Shell-tab shortcuts

| Key | Action |
| --- | --- |
| `Ctrl+N` | Create a new tab |
| `Ctrl+W` | Close the current tab |
| `Tab` | Next tab |
| `Shift+Tab` | Previous tab |
| `Up` | Older command from local tab history |
| `Down` | Newer command from local tab history |
| `Esc` | Return to session actions |

### Canceling running commands

The `Cancel Commands` action tries to interrupt foreground work and then sends OS-specific follow-up commands:

- Unix-like targets: attempts `kill -INT` and `kill -KILL` against shell jobs.
- Windows targets: attempts to stop PowerShell jobs.

## Payloads And Network Awareness

### Payload page behavior

The payload page renders several reverse-shell templates and substitutes `YOUR_IP` using runtime network detection.

Current payload types shown in the UI:

- Bash
- Python
- Netcat
- PowerShell
- PHP
- Ruby
- Perl

### IP selection behavior

Necromancy gathers:

- local IP information,
- public IP information,
- configured listening ports.

The UI prefers a non-private public IP when one is available; otherwise it falls back to the local IP.

### Payload-page controls

| Key | Action |
| --- | --- |
| `Enter` | Copy the selected payload to clipboard |
| `c` | Copy the selected payload to clipboard |
| `r` | Refresh network information |
| `Esc` | Return to the dashboard |

### CLI payload output

Use the CLI-only payload mode when you just need a quick payload printout:

```bash
./necromancy --payloads
./necromancy --payloads -p 8080
```

## File Transfer And File Manager

### Upload and in-memory execution

From the session action page, you can:

- upload a local file to a remote destination,
- execute a local script in-memory on the target.

These are launched through forms in the TUI.

### Built-in HTTP file server

If you start Necromancy with `-s`, the application serves the chosen directory over HTTP. This is useful for quick staging and download workflows.

```bash
./necromancy -p 4444 -s ./payloads -w 8000
```

### File manager status

The file manager exists and is wired into the session workflow, but it is not fully feature-complete.

#### Implemented or substantially wired

- directory listing through the enhanced session commands,
- path navigation,
- enter directory,
- go to parent directory,
- refresh,
- execute selected executable files,
- delete file or directory with confirmation,
- help view,
- exit back to session actions.

#### Present in UI but still partial or placeholder-driven

- direct download action,
- upload action from within the file-manager page,
- create new file,
- create new directory,
- copy and paste,
- inline edit behavior.

### File-manager keys

| Key | Action |
| --- | --- |
| `Enter` | Enter directory |
| `Backspace` | Go to parent directory |
| `r` | Refresh listing |
| `d` | Download selected file |
| `u` | Upload file |
| `x` | Delete selected file or directory |
| `n` | Create new file |
| `m` | Create new directory |
| `e` | Execute selected executable |
| `c` | Copy selected file |
| `v` | Paste file |
| `a` | Select all |
| `h` | Show help |
| `q` | Quit file manager |
| `Esc` | Return to session actions |

## Module System

The module system is available both as a catalog page and as an execution path after selecting a session.

### How module execution works

1. Open `Sessions`.
2. Select a session.
3. Open `Run Module`.
4. Pick a module.
5. Pick a shell tab to dispatch the module from.

Module output then appears through the selected session stream.

### Built-in module catalog

| Module Key | Purpose | Status |
| --- | --- | --- |
| `peass_auto` | Auto-select LinPEAS or WinPEAS by detected OS | implemented |
| `linpeas` | Linux privilege-escalation enumeration | implemented |
| `winpeas` | Windows privilege-escalation enumeration | implemented |
| `lse` | Linux Smart Enumeration helper | implemented |
| `potato` | Windows potato-family helper | placeholder/helper |
| `traitor` | Linux escalation helper | implemented |
| `uac` | Windows UAC bypass helper | implemented |
| `chisel` | Chisel tunneling helper | placeholder/helper |
| `ligolo` | Ligolo-ng tunneling helper | placeholder/helper |
| `ngrok` | Ngrok helper | placeholder/helper |
| `meterpreter` | Meterpreter upgrade helper | placeholder/helper |
| `cleanup` | Cleanup traces and history | implemented |
| `panix` | Linux persistence helper | implemented |
| `linux_procmemdump` | Linux process-memory helper | implemented/helper |
| `filemanager` | File-manager workflow helper | workflow helper |
| `redsun` | Windows enumeration helper | helper/bootstrap |
| `bluehammer` | Windows exploitation helper | helper/bootstrap |
| `payload_obfuscator` | Obfuscated payload generation helper | implemented/helper |
| `enhanced_preflight_recon` | Recon and security-surface checks | implemented/helper |
| `process_monitor` | Background process monitoring helper | implemented/helper |
| `background_checker` | Check running background processes | implemented/helper |

For the full module page and shorter descriptions, see [MODULES.md](MODULES.md).

## Logging And Runtime Files

### Runtime files created by the app

| Path | Purpose |
| --- | --- |
| `necromancy-go.log` | Main application log |
| `logs/session_<id>.log` | Per-session output logs |

### Logging behavior

- Session logs are enabled by default.
- `-L` disables per-session logs.
- Main application logging still goes to `necromancy-go.log`.

## Configuration Notes

Necromancy currently relies mainly on command-line flags. The repository includes `CONFIGURATION.md` with examples, but the live code path is flag-centric.

### Practical knobs

- Listener ports: `-p`
- Interface/IP binding: `-i`
- HTTP serving directory: `-s`
- HTTP server port: `-w`
- URL prefix: `--prefix`
- Single-session mode: `-S`
- PTY auto-upgrade disable: `-U`
- Headless mode: `--headless`

### Notes on config examples in the repo

Some older documentation examples mention environment-variable-driven behavior or extended configuration surfaces. Treat those as guidance documents, not guaranteed runtime features, unless you have verified the corresponding implementation in code.

## Troubleshooting

### No sessions appear

- Confirm the listener started on the expected port.
- Verify firewalls and network path between target and operator box.
- Check whether another process already owns the port.
- Use `./necromancy --interfaces` to confirm local network details.

### Payload IP looks wrong

- Refresh the payload page with `r`.
- If the public IP is unavailable or private, Necromancy falls back to the detected local IP.
- Bind explicitly with `-i` if you want a specific listener address.

### TUI does not start

- Use `--headless` in non-interactive shells.
- Ensure `TERM` is not `dumb`.
- Ensure stdin and stdout are real terminals when launching the TUI.

### PTY upgrade does not work

- PTY upgrade is skipped for Windows sessions.
- Some minimal shells do not have the tooling or shell behavior needed for a successful upgrade.
- Use raw interaction and verify the target has a compatible shell environment.

### Update command fails

- Confirm outbound access to GitHub Releases and the GitHub API.
- Ensure the running binary can be replaced on disk.
- On Unix-like systems, verify file permissions on the executable path.

### File-manager action looks incomplete

That may be expected. Several file-manager actions are surfaced in the UI before their full implementation is finished.

## Security And Safe Use

- Only use Necromancy in authorized environments.
- Review every payload and helper action before using it.
- Validate module behavior on lab targets first.
- Treat placeholder or helper modules as starting points, not one-click guarantees.
- Keep logs and staged files under proper operational control.

### Reporting a security issue

See [SECURITY.md](SECURITY.md) for disclosure instructions.

## Related Documents

- [README.md](README.md): project overview and quick start
- [MODULES.md](MODULES.md): module catalog
- [QUICK_REFERENCE.md](QUICK_REFERENCE.md): short commands and workflows
- [CONFIGURATION.md](CONFIGURATION.md): configuration examples
- [AGENTS.md](AGENTS.md): architecture notes for AI agents
- [CONTRIBUTING.md](CONTRIBUTING.md): contribution guide
- [SECURITY.md](SECURITY.md): security policy

---

**Version:** `v1.5.1`
**Last reviewed:** `2026-04-27`
**Repository:** <https://github.com/Aryma-f4/necromancy>

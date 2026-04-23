package core

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type Session struct {
	ID         int
	Conn       net.Conn
	RemoteAddr string
	Type       string // Raw, PTY
	History    *bytes.Buffer
	mu         sync.Mutex
	Active     bool
	OS         string
	LiveOutput chan []byte
	IsAttached bool
	LogFile    *os.File
}

func NewSession(id int, conn net.Conn) *Session {
	s := &Session{
		ID:         id,
		Conn:       conn,
		RemoteAddr: conn.RemoteAddr().String(),
		Type:       "Raw",
		History:    new(bytes.Buffer),
		Active:     true,
		LiveOutput: make(chan []byte, 100),
	}

	if GlobalConfig != nil && !GlobalConfig.NoLog {
		_ = os.MkdirAll("logs", 0755)
		f, err := os.OpenFile(fmt.Sprintf("logs/session_%d.log", id), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err == nil {
			s.LogFile = f
			f.WriteString(fmt.Sprintf("=== Session %d started at %s ===\n", id, time.Now().Format(time.RFC3339)))
		}
	}

	go s.readLoop()
	return s
}

func guessOSFromText(text string) string {
	lower := strings.ToLower(text)

	windowsIndicators := []string{
		"microsoft windows",
		"windows version",
		"c:\\",
		"\\users\\",
		"powershell",
		"cmd.exe",
	}
	for _, indicator := range windowsIndicators {
		if strings.Contains(lower, indicator) {
			return "windows"
		}
	}

	linuxIndicators := []string{
		"/bin/",
		"/usr/",
		"/tmp/",
		"uid=",
		"gid=",
		"linux",
		"bash",
		"sh:",
	}
	for _, indicator := range linuxIndicators {
		if strings.Contains(lower, indicator) {
			return "linux"
		}
	}

	return ""
}

func (s *Session) DetectedOS() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.OS == "" {
		s.OS = guessOSFromText(s.History.String())
	}

	if s.OS == "" {
		return "unknown"
	}

	return s.OS
}

func (s *Session) readLoop() {
	buf := make([]byte, 4096)
	for {
		n, err := s.Conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Printf("Session %d closed: %v\n", s.ID, err)
			}
			s.mu.Lock()
			s.Active = false
			if s.LogFile != nil {
				s.LogFile.WriteString(fmt.Sprintf("=== Session closed at %s ===\n", time.Now().Format(time.RFC3339)))
				s.LogFile.Close()
			}
			if s.IsAttached {
				s.LiveOutput <- nil // Signal EOF
			}
			s.mu.Unlock()
			return
		}

		data := make([]byte, n)
		copy(data, buf[:n])

		s.mu.Lock()
		s.History.Write(data)
		if s.OS == "" {
			s.OS = guessOSFromText(string(data))
		}
		if s.LogFile != nil {
			s.LogFile.Write(data)
		}
		if s.IsAttached {
			s.LiveOutput <- data
		}
		s.mu.Unlock()
	}
}

func (s *Session) Write(data []byte) (int, error) {
	return s.Conn.Write(data)
}

func (s *Session) CancelRunningCommands() error {
	// Always try to interrupt the current foreground process first.
	if _, err := s.Write([]byte{0x03, 0x03, '\n'}); err != nil {
		return err
	}

	switch strings.ToLower(s.DetectedOS()) {
	case "windows":
		_, err := s.Write([]byte(
			"\r\npowershell -NoProfile -ExecutionPolicy Bypass -Command \"$jobs = Get-Job -ErrorAction SilentlyContinue; if ($jobs) { $jobs | Stop-Job -Force -ErrorAction SilentlyContinue }; Write-Output '[*] Cancel signal sent to current Windows jobs'\"\r\n",
		))
		return err
	default:
		_, err := s.Write([]byte(
			"\nprintf '[*] Cancelling running commands...\\n'; jobs -p 2>/dev/null | xargs -r kill -INT 2>/dev/null; jobs -p 2>/dev/null | xargs -r kill -KILL 2>/dev/null; printf '[*] Cancel signal sent.\\n'\n",
		))
		return err
	}
}

func (s *Session) Attach() {
	s.mu.Lock()
	s.IsAttached = true
	// Drain existing buffer contents if needed, usually we just print History
	s.mu.Unlock()
}

func (s *Session) Detach() {
	s.mu.Lock()
	s.IsAttached = false
	s.mu.Unlock()
}

type SessionManager struct {
	Sessions map[int]*Session
	NextID   int
	mu       sync.Mutex
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		Sessions: make(map[int]*Session),
		NextID:   1,
	}
}

func (sm *SessionManager) Add(conn net.Conn) *Session {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	s := NewSession(sm.NextID, conn)
	sm.Sessions[sm.NextID] = s
	sm.NextID++
	return s
}

func (sm *SessionManager) Get(id int) (*Session, bool) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	s, ok := sm.Sessions[id]
	return s, ok
}

func (sm *SessionManager) Remove(id int) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if s, ok := sm.Sessions[id]; ok {
		s.Conn.Close()
		s.Active = false
		delete(sm.Sessions, id)
	}
}

func (sm *SessionManager) GetAll() []*Session {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	var all []*Session
	for _, s := range sm.Sessions {
		if s.Active {
			all = append(all, s)
		}
	}
	return all
}

func ConnectBind(target string, sm *SessionManager, onNewSession func(*Session)) {
	log.Printf("Attempting to connect to Bind Shell at %s...\n", target)
	conn, err := net.DialTimeout("tcp", target, 5*time.Second)
	if err != nil {
		log.Printf("Failed to connect to bind shell %s: %v\n", target, err)
		return
	}

	log.Printf("Successfully connected to bind shell at %s\n", target)
	s := sm.Add(conn)
	if onNewSession != nil {
		onNewSession(s)
	}
}
func StartListener(addr string, sm *SessionManager, onNewSession func(*Session)) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Error listening on %s: %v", addr, err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}

		if GlobalConfig != nil && GlobalConfig.SingleSession && len(sm.GetAll()) >= 1 {
			log.Printf("Rejecting connection from %s because single-session mode is enabled\n", conn.RemoteAddr())
			conn.Close()
			continue
		}

		s := sm.Add(conn)
		if onNewSession != nil {
			onNewSession(s)
		}
	}
}

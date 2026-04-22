package core

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"
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
func StartListener(port string, sm *SessionManager, onNewSession func(*Session)) {
	addr := fmt.Sprintf("0.0.0.0:%s", port)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Error listening: %v", err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		s := sm.Add(conn)
		if onNewSession != nil {
			onNewSession(s)
		}
	}
}

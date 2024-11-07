package core

import (
	"sync"
)

// SessionManager controla o processamento por sessionId
type SessionManager struct {
	mu       sync.Mutex
	sessions map[string]bool
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: make(map[string]bool),
	}
}

func (sm *SessionManager) StartSession(sessionId string) bool {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.sessions[sessionId] {
		return false
	}

	sm.sessions[sessionId] = true
	return true
}

func (sm *SessionManager) EndSession(sessionId string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.sessions, sessionId)
}

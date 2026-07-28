package sse

import (
	"sync"
	"time"
)

// Session represents a user session with metadata
type Session struct {
	ID        string
	UserID    string
	CreatedAt time.Time
	LastPing  time.Time
	Metadata  map[string]any
}

// SessionManager manages user sessions
type SessionManager struct {
	sessions map[string]*Session
	mu       sync.RWMutex
	hub      *Hub
}

// NewSessionManager creates a new SessionManager
func NewSessionManager(hub *Hub) *SessionManager {
	return &SessionManager{
		sessions: make(map[string]*Session),
		hub:      hub,
	}
}

// CreateSession creates a new session
func (sm *SessionManager) CreateSession(sessionID, userID string) *Session {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	session := &Session{
		ID:        sessionID,
		UserID:    userID,
		CreatedAt: time.Now(),
		LastPing:  time.Now(),
		Metadata:  make(map[string]any),
	}
	sm.sessions[sessionID] = session
	return session
}

// GetSession retrieves a session by ID
func (sm *SessionManager) GetSession(sessionID string) (*Session, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	session, ok := sm.sessions[sessionID]
	return session, ok
}

// DeleteSession removes a session
func (sm *SessionManager) DeleteSession(sessionID string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.sessions, sessionID)
}

// UpdateLastPing updates the session's last ping time
func (sm *SessionManager) UpdateLastPing(sessionID string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if session, ok := sm.sessions[sessionID]; ok {
		session.LastPing = time.Now()
	}
}

// SetMetadata sets metadata for a session
func (sm *SessionManager) SetMetadata(sessionID, key string, value any) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if session, ok := sm.sessions[sessionID]; ok {
		session.Metadata[key] = value
	}
}

// GetMetadata gets metadata from a session
func (sm *SessionManager) GetMetadata(sessionID, key string) (any, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	if session, ok := sm.sessions[sessionID]; ok {
		val, exists := session.Metadata[key]
		return val, exists
	}
	return nil, false
}

// GetSessionsByUserID returns all sessions for a user
func (sm *SessionManager) GetSessionsByUserID(userID string) []*Session {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	var sessions []*Session
	for _, session := range sm.sessions {
		if session.UserID == userID {
			sessions = append(sessions, session)
		}
	}
	return sessions
}

// SendToUser sends an event to all sessions of a user
func (sm *SessionManager) SendToUser(userID string, event EventData) int {
	sessions := sm.GetSessionsByUserID(userID)
	sent := 0
	for _, session := range sessions {
		if sm.hub.SendToClient(session.ID, event) {
			sent++
		}
	}
	return sent
}

// SendToSession sends an event to a specific session
func (sm *SessionManager) SendToSession(sessionID string, event EventData) bool {
	return sm.hub.SendToClient(sessionID, event)
}

// Broadcast sends an event to all sessions
func (sm *SessionManager) Broadcast(event EventData) {
	sm.hub.Broadcast(event)
}

// GetAllSessions returns all active sessions
func (sm *SessionManager) GetAllSessions() []*Session {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	sessions := make([]*Session, 0, len(sm.sessions))
	for _, session := range sm.sessions {
		sessions = append(sessions, session)
	}
	return sessions
}

// GetSessionCount returns the number of active sessions
func (sm *SessionManager) GetSessionCount() int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return len(sm.sessions)
}

// CleanupStaleSessions removes sessions that haven't pinged within the timeout
func (sm *SessionManager) CleanupStaleSessions(timeout time.Duration) int {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	now := time.Now()
	removed := 0
	for id, session := range sm.sessions {
		if now.Sub(session.LastPing) > timeout {
			delete(sm.sessions, id)
			removed++
		}
	}
	return removed
}

// StartCleanupRoutine starts a goroutine that periodically cleans up stale sessions
func (sm *SessionManager) StartCleanupRoutine(interval, timeout time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			sm.CleanupStaleSessions(timeout)
		}
	}()
}

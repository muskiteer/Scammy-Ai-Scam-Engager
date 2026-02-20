package internal

import (
	"sync"
	"time"
)

// SessionData holds all the state for a single session
type SessionData struct {
	SessionID      string
	Context        SessionContext
	MessageHistory []string
	Keywords       []string // Suspicious keywords from ScamDetection
	LastUpdated    time.Time
	StartTime      time.Time // Track when conversation started for engagement duration
}

// SessionStore manages all active sessions
type SessionStore struct {
	sessions map[string]*SessionData
	mu       sync.RWMutex
}

var globalStore = &SessionStore{
	sessions: make(map[string]*SessionData),
}

// GetStore returns the global session store
func GetStore() *SessionStore {
	return globalStore
}

// Get retrieves a session by ID, creates one if it doesn't exist
func (s *SessionStore) Get(sessionID string) *SessionData {
	s.mu.Lock()
	defer s.mu.Unlock()

	if session, exists := s.sessions[sessionID]; exists {
		return session
	}

	// Create new session
	newSession := &SessionData{
		SessionID: sessionID,
		Context: SessionContext{
			ScamDetected: false,
			TurnCount:    0,
			Intel: Intel{
				UPI:           []string{},
				Phone:         []string{},
				Link:          []string{},
				Bank:          []string{},
				Email:         []string{},
				CaseIDs:       []string{},
				PolicyNumbers: []string{},
				OrderNumbers:  []string{},
				CardNumbers:   []string{},
				IFSCCodes:     []string{},
			},
			CurrentState:            StateInit,
			QuestionsAsked:          0,
			InvestigativeQuestions:  0,
			RedFlagsIdentified:      []string{},
			InformationElicitations: 0,
		},
		MessageHistory: []string{},
		Keywords:       []string{},
		LastUpdated:    time.Now(),
		StartTime:      time.Now(),
	}

	s.sessions[sessionID] = newSession
	return newSession
}

// Update saves the session data
func (s *SessionStore) Update(session *SessionData) {
	s.mu.Lock()
	defer s.mu.Unlock()

	session.LastUpdated = time.Now()
	s.sessions[session.SessionID] = session
}

// Delete removes a session
func (s *SessionStore) Delete(sessionID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.sessions, sessionID)
}

// AddMessage appends a message to the session history
func (session *SessionData) AddMessage(text string) {
	session.MessageHistory = append(session.MessageHistory, text)
}

// AddKeyword adds a suspicious keyword (avoiding duplicates)
func (session *SessionData) AddKeyword(keyword string) {
	for _, k := range session.Keywords {
		if k == keyword {
			return
		}
	}
	session.Keywords = append(session.Keywords, keyword)
}

package session

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"
	"sync"
	"time"
)

type Session struct {
	ID     string
	Values map[string]interface{}
	timer  *time.Timer
}

// NewSession initializes a new Session.
func NewSession() *Session {
	return &Session{
		Values: make(map[string]interface{}),
	}
}

type Store interface {
	Load(r *http.Request) (*Session, error)
	Save(w http.ResponseWriter, r *http.Request, session *Session) error
	Delete(w http.ResponseWriter, r *http.Request) error
	Clear()
}

// Store is an in-memory session store.
type storeImpl struct {
	lock          sync.RWMutex
	name          string
	sessions      map[string]*Session
	sessionMaxAge time.Duration
}

// NewStore initializes a new Store.
func NewStore(name string, maxAge time.Duration) Store {
	return &storeImpl{
		name:          name,
		sessions:      make(map[string]*Session),
		sessionMaxAge: maxAge,
	}
}

// Load parses a Session from r and returns it from the store.
func (s *storeImpl) Load(r *http.Request) (*Session, error) {
	cookie, err := r.Cookie(s.name)
	if err != nil {
		return nil, err
	}

	s.lock.RLock()
	session := s.sessions[cookie.Value]
	s.lock.RUnlock()

	if session == nil {
		return nil, errors.New("session not found")
	}

	return session, nil
}

func (s *storeImpl) Save(w http.ResponseWriter, r *http.Request, session *Session) error {

	// generate session ID
	b := make([]byte, 33)
	if _, err := rand.Read(b); err != nil {
		return err
	}
	sessionId := base64.StdEncoding.EncodeToString(b)

	// set a timer to delete session on expiry
	t := time.AfterFunc(s.sessionMaxAge, func() {
		s.lock.Lock()
		delete(s.sessions, sessionId)
		s.lock.Unlock()
	})
	session.timer = t

	// save session
	s.lock.Lock()
	s.sessions[sessionId] = session
	s.lock.Unlock()

	// create cookie from signed session ID
	// todo Domain, Signing
	cookie := &http.Cookie{
		Name:     s.name,
		Value:    sessionId,
		MaxAge:   int(s.sessionMaxAge.Seconds()),
		Expires:  time.Now().Add(s.sessionMaxAge),
		HttpOnly: true,
		Secure:   r.URL.Scheme == "https",
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, cookie)

	return nil
}

func (s *storeImpl) Delete(w http.ResponseWriter, r *http.Request) error {

	// read cookie
	cookie, err := r.Cookie(s.name)
	if err != nil {
		return err
	}

	// delete session
	s.lock.Lock()
	if session, ok := s.sessions[cookie.Value]; ok {
		session.timer.Stop()
		delete(s.sessions, cookie.Value)
	}
	s.lock.Unlock()

	// create expired cookie
	expiredCookie := &http.Cookie{
		Name:     s.name,
		Value:    "",
		MaxAge:   -1,
		Expires:  time.Unix(1, 0),
		HttpOnly: true,
		Secure:   r.URL.Scheme == "https",
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, expiredCookie)

	return nil
}

// Clear removes all sessions from the store.
func (s *storeImpl) Clear() {
	s.lock.Lock()
	for k, v := range s.sessions {
		v.timer.Stop()
		delete(s.sessions, k)
	}
	s.lock.Unlock()
}

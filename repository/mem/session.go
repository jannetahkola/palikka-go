package mem

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	palikka "palikka-go"
	"sync"
	"time"
)

type storeEntry struct {
	Session *palikka.Session
	Timer   *time.Timer
}

type sessionStoreImpl struct {
	lock          sync.RWMutex
	delChan       chan string
	sessions      map[string]*storeEntry
	sessionMaxAge time.Duration
}

// NewSessionStore initializes a new palikka.SessionStore.
func NewSessionStore(maxAge time.Duration) palikka.SessionStore {
	return &sessionStoreImpl{
		delChan:       make(chan string),
		sessions:      make(map[string]*storeEntry),
		sessionMaxAge: maxAge,
	}
}

func (s *sessionStoreImpl) Load(ctx context.Context, id string) (*palikka.Session, error) {
	s.lock.RLock()
	session := s.sessions[id]
	s.lock.RUnlock()

	if session == nil {
		return nil, errors.New("session not found")
	}

	return session.Session, nil
}

func (s *sessionStoreImpl) Save(ctx context.Context, session *palikka.Session) error {

	// generate session ID
	b := make([]byte, 33)
	if _, err := rand.Read(b); err != nil {
		return err
	}

	// set session attributes
	session.ID = base64.StdEncoding.EncodeToString(b)
	session.MaxAge = s.sessionMaxAge
	session.Expires = time.Now().Add(s.sessionMaxAge).UTC()

	// set a timer to delete session on expiry
	t := time.AfterFunc(s.sessionMaxAge, func() {
		s.lock.Lock()
		session = s.sessions[session.ID].Session
		s.delete(session.ID)
		s.lock.Unlock()
	})

	entry := &storeEntry{
		Session: session,
		Timer:   t,
	}

	// save session
	s.lock.Lock()
	s.sessions[session.ID] = entry
	s.lock.Unlock()

	return nil
}

func (s *sessionStoreImpl) Delete(ctx context.Context, id string) error {
	s.lock.Lock()
	if session, ok := s.sessions[id]; ok {
		session.Timer.Stop()
		s.delete(id)
	}
	s.lock.Unlock()
	return nil
}

// Clear removes all sessions from the store.
func (s *sessionStoreImpl) Clear(ctx context.Context) {
	s.lock.Lock()
	for k, v := range s.sessions {
		v.Timer.Stop()
		s.delete(k)
	}
	s.lock.Unlock()
}

func (s *sessionStoreImpl) OnDelete() chan string {
	return s.delChan
}

// delete deletes a session matching the id and sends
// the id to delChan (non-blocking). Caller must handle
// operations on lock.
func (s *sessionStoreImpl) delete(id string) {
	delete(s.sessions, id)
	select {
	case s.delChan <- id:
	default:
	}
}

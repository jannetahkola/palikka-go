package mem_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	palikka "palikka-go"
	"palikka-go/repository/mem"
	"testing"
	"time"
)

func TestSessionStore_Load(t *testing.T) {
	t.Run("loads session", func(t *testing.T) {
		var (
			err     error
			store   palikka.SessionStore
			session *palikka.Session
		)

		store = mem.NewSessionStore(time.Minute)

		session = &palikka.Session{}
		err = store.Save(context.Background(), session)
		assert.NoError(t, err)

		session, err = store.Load(context.Background(), session.ID)
		assert.NoError(t, err)
		assert.NotNil(t, session)
		assert.NotEmpty(t, session.ID)
	})

	t.Run("returns error if session not found", func(t *testing.T) {
		store := mem.NewSessionStore(time.Minute)
		_, err := store.Load(context.Background(), "mock-session-id")
		assert.Error(t, err)
		assert.Equal(t, "session not found", err.Error())
	})
}

func TestSessionStore_Save(t *testing.T) {
	var (
		err     error
		store   palikka.SessionStore
		session *palikka.Session
	)

	store = mem.NewSessionStore(time.Minute)
	session = &palikka.Session{Values: make(map[string]interface{})}

	session.ID = "mock-session-id"
	session.MaxAge = time.Hour
	session.Expires = time.Now().Add(time.Hour)
	session.Values["foo"] = "bar"

	err = store.Save(context.Background(), session)
	assert.NoError(t, err)
	assert.Equal(t, "bar", session.Values["foo"])

	// these should be overwritten by the store
	assert.NotEmpty(t, session.ID)
	assert.NotEqual(t, "mock-session-id", session.ID)
	assert.Equal(t, time.Minute, session.MaxAge)
	assert.True(t, session.Expires.Before(time.Now().Add(time.Minute*10)))
}

func TestSessionStore_Delete(t *testing.T) {
	t.Run("deletes session", func(t *testing.T) {
		var (
			err     error
			store   palikka.SessionStore
			session *palikka.Session
		)

		store = mem.NewSessionStore(time.Minute)
		session = &palikka.Session{Values: make(map[string]interface{})}

		_ = store.Save(context.Background(), session)

		err = store.Delete(context.Background(), session.ID)
		assert.NoError(t, err)

		_, err = store.Load(context.Background(), session.ID)
		assert.Equal(t, "session not found", err.Error())
	})

	t.Run("returns nil if session not found", func(t *testing.T) {
		store := mem.NewSessionStore(time.Minute)
		err := store.Delete(context.Background(), "mock-session-id")
		assert.NoError(t, err)
	})
}

//func (s *SessionsTestSuite) Test_Save() {
//	req, _ := http.NewRequest("GET", "/", nil)
//	rec := httptest.NewRecorder()
//
//	maxAge := time.Hour * 24 * 7
//	maxAgeSkew := time.Minute
//
//	store := NewStore("session", maxAge)
//
//	session := NewSession()
//	session.ID = "session-id"
//	session.Values["foo"] = "bar"
//
//	err := store.Save(rec, req, session)
//	assert.NoError(s.T(), err)
//
//	// check session
//	assert.NotEmpty(s.T(), session.ID)
//	assert.NotEmpty(s.T(), "session-id", session.ID) // should be overwritten
//	assert.NotNil(s.T(), session.timer)
//	assert.Contains(s.T(), session.Values, "foo")
//	assert.Equal(s.T(), "bar", session.Values["foo"])
//
//	// check set-cookie header
//	setCookie := rec.Header().Get("Set-Cookie")
//	assert.NotEmpty(s.T(), setCookie)
//
//	cookie, err := http.ParseSetCookie(setCookie)
//	assert.NoError(s.T(), err)
//
//	// check cookie
//	assert.Equal(s.T(), cookie.Name, "session")
//	assert.NotEmpty(s.T(), cookie.Value)
//	assert.NotEqual(s.T(), "session-id", cookie.Value)
//	assert.True(s.T(), cookie.HttpOnly)
//	assert.Equal(s.T(), int(store.(*storeImpl).sessionMaxAge.Seconds()), cookie.MaxAge)
//
//	expireBefore := time.Now().Add(maxAge + maxAgeSkew).UTC()
//	assert.True(s.T(), cookie.Expires.Before(expireBefore))
//
//	expireAfter := time.Now().Add(maxAge - maxAgeSkew).UTC()
//	assert.True(s.T(), cookie.Expires.After(expireAfter))
//}
//
//func (s *SessionsTestSuite) Test_Delete() {
//	req, _ := http.NewRequest("GET", "/", nil)
//	rec := httptest.NewRecorder()
//
//	cookie, _ := s.SaveSession(NewSession())
//	req.AddCookie(cookie)
//
//	// delete session
//	err := s.store.Delete(rec, req)
//	assert.NoError(s.T(), err)
//	assert.Empty(s.T(), s.store.(*storeImpl).sessions)
//
//	// check set-cookie header
//	setCookie := rec.Header().Get("Set-Cookie")
//	assert.NotEmpty(s.T(), setCookie)
//
//	// check cookie
//	expiredCookie, err := http.ParseSetCookie(setCookie)
//	assert.NoError(s.T(), err)
//	assert.NotNil(s.T(), expiredCookie)
//	assert.Equal(s.T(), "session", expiredCookie.Name)
//	assert.Empty(s.T(), expiredCookie.Value)
//	assert.Equal(s.T(), -1, expiredCookie.MaxAge)
//	assert.True(s.T(), time.Unix(1, 0).Equal(expiredCookie.Expires))
//}
//
//func (s *SessionsTestSuite) Test_Save_Benchmark() {
//	var wg sync.WaitGroup
//
//	// create 10 million sessions
//	for i := 0; i < 20; i++ {
//		go func() {
//			defer wg.Done()
//			for j := 0; j < 500000; j++ {
//				req, _ := http.NewRequest("GET", "/", nil)
//				rec := httptest.NewRecorder()
//
//				session := NewSession()
//				session.Values["foo"] = "bar"
//
//				err := s.store.Save(rec, req, session)
//				assert.NoError(s.T(), err)
//				assert.NotEmpty(s.T(), rec.Header().Get("Set-Cookie"))
//			}
//		}()
//	}
//
//	wg.Wait()
//}

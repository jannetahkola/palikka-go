package session

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func (s *SessionsTestSuite) Test_Load() {
	req, _ := http.NewRequest("GET", "/", nil)

	cookie, _ := s.SaveSession(NewSession())
	req.AddCookie(cookie)

	session, err := s.store.Load(req)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), session)
}

func (s *SessionsTestSuite) Test_Save() {
	req, _ := http.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	maxAge := time.Hour * 24 * 7
	maxAgeSkew := time.Minute

	store := NewStore("session", maxAge)

	session := NewSession()
	session.ID = "session-id"
	session.Values["foo"] = "bar"

	err := store.Save(rec, req, session)
	assert.NoError(s.T(), err)

	// check session
	assert.NotEmpty(s.T(), session.ID)
	assert.NotEmpty(s.T(), "session-id", session.ID) // should be overwritten
	assert.NotNil(s.T(), session.timer)
	assert.Contains(s.T(), session.Values, "foo")
	assert.Equal(s.T(), "bar", session.Values["foo"])

	// check set-cookie header
	setCookie := rec.Header().Get("Set-Cookie")
	assert.NotEmpty(s.T(), setCookie)

	cookie, err := http.ParseSetCookie(setCookie)
	assert.NoError(s.T(), err)

	// check cookie
	assert.Equal(s.T(), cookie.Name, "session")
	assert.NotEmpty(s.T(), cookie.Value)
	assert.NotEqual(s.T(), "session-id", cookie.Value)
	assert.True(s.T(), cookie.HttpOnly)
	assert.Equal(s.T(), int(store.(*storeImpl).sessionMaxAge.Seconds()), cookie.MaxAge)

	expireBefore := time.Now().Add(maxAge + maxAgeSkew).UTC()
	assert.True(s.T(), cookie.Expires.Before(expireBefore))

	expireAfter := time.Now().Add(maxAge - maxAgeSkew).UTC()
	assert.True(s.T(), cookie.Expires.After(expireAfter))
}

func (s *SessionsTestSuite) Test_Delete() {
	req, _ := http.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	cookie, _ := s.SaveSession(NewSession())
	req.AddCookie(cookie)

	// delete session
	err := s.store.Delete(rec, req)
	assert.NoError(s.T(), err)
	assert.Empty(s.T(), s.store.(*storeImpl).sessions)

	// check set-cookie header
	setCookie := rec.Header().Get("Set-Cookie")
	assert.NotEmpty(s.T(), setCookie)

	// check cookie
	expiredCookie, err := http.ParseSetCookie(setCookie)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), expiredCookie)
	assert.Equal(s.T(), "session", expiredCookie.Name)
	assert.Empty(s.T(), expiredCookie.Value)
	assert.Equal(s.T(), -1, expiredCookie.MaxAge)
	assert.True(s.T(), time.Unix(1, 0).Equal(expiredCookie.Expires))
}

func (s *SessionsTestSuite) Test_Save_Benchmark() {
	var wg sync.WaitGroup

	// create 10 million sessions
	for i := 0; i < 20; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 500000; j++ {
				req, _ := http.NewRequest("GET", "/", nil)
				rec := httptest.NewRecorder()

				session := NewSession()
				session.Values["foo"] = "bar"

				err := s.store.Save(rec, req, session)
				assert.NoError(s.T(), err)
				assert.NotEmpty(s.T(), rec.Header().Get("Set-Cookie"))
			}
		}()
	}

	wg.Wait()
}

type SessionsTestSuite struct {
	suite.Suite

	maxAge time.Duration
	store  Store
}

func (s *SessionsTestSuite) SetupSuite() {
	s.maxAge = time.Minute * 5
	s.store = NewStore("session", s.maxAge)
}

func (s *SessionsTestSuite) SetupTest() {
	// clear session store between tests
	s.store.Clear()
}

func (s *SessionsTestSuite) SaveSession(session *Session) (*http.Cookie, error) {
	req, _ := http.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	_ = s.store.Save(rec, req, session)

	return http.ParseSetCookie(rec.Header().Get("Set-Cookie"))
}

func TestSessionsTestSuite(t *testing.T) {
	suite.Run(t, new(SessionsTestSuite))
}

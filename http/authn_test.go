package http_test

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	palikka "palikka-go"
	palikkahttp "palikka-go/http"
	palikkacookie "palikka-go/http/cookie"
	"testing"
	"time"
)

func (s *AuthNTestSuite) TestAuth_Login() {
	loginRequest := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		Username: "username",
		Password: "password",
	}
	loginRequestJson, _ := json.Marshal(loginRequest)

	hash, err := bcrypt.GenerateFromPassword([]byte(loginRequest.Password), bcrypt.DefaultCost)
	assert.NoError(s.T(), err)
	passwordHash := string(hash)

	doLoginRequest := func(payload []byte) (*http.Request, *http.Response) {
		req, _ := http.NewRequest(http.MethodPost, s.URL()+"/auth/login", bytes.NewBuffer(payload))
		res, _ := http.DefaultClient.Do(req)
		return req, res
	}

	s.Run("returns error when user not found", func() {
		s.UserService.FindUserByUsernameFn = func(ctx context.Context, username string) (*palikka.User, error) {
			return nil, errors.New("user not found")
		}

		_, res := doLoginRequest(loginRequestJson)
		MustRespondWithError(s.T(), res, http.StatusBadRequest, palikkahttp.ErrMsgInvalidRequest)
		assert.Empty(s.T(), res.Header.Values("Set-Cookie"))
	})

	s.Run("returns error when password is invalid", func() {
		s.UserService.FindUserByUsernameFn = func(ctx context.Context, username string) (*palikka.User, error) {
			pw := "doesnotmatch"
			return &palikka.User{Username: username, Password: &pw}, nil
		}

		_, res := doLoginRequest(loginRequestJson)
		MustRespondWithError(s.T(), res, http.StatusBadRequest, palikkahttp.ErrMsgInvalidRequest)
		assert.Empty(s.T(), res.Header.Values("Set-Cookie"))
	})

	s.Run("validates input", func() {
		s.Run("empty username", func() {
			r := &palikkahttp.LoginRequest{}
			rBytes, _ := json.Marshal(r)
			_, res := doLoginRequest(rBytes)
			MustRespondWithError(s.T(), res, http.StatusBadRequest, palikkahttp.ErrMsgInvalidRequest)
		})

		s.Run("empty password", func() {
			username := "username"
			r := &palikkahttp.LoginRequest{Username: &username}
			rBytes, _ := json.Marshal(r)
			_, res := doLoginRequest(rBytes)
			MustRespondWithError(s.T(), res, http.StatusBadRequest, palikkahttp.ErrMsgInvalidRequest)
		})

		s.Run("invalid username", func() {
			usernames := []string{
				"user name",
				"@username",
			}
			password := "password"
			for _, username := range usernames {
				r := &palikkahttp.LoginRequest{Username: &username, Password: &password}
				rBytes, _ := json.Marshal(r)
				_, res := doLoginRequest(rBytes)
				MustRespondWithError(s.T(), res, http.StatusBadRequest, palikkahttp.ErrMsgInvalidRequest)
			}
		})

		// todo invalid password
	})

	s.Run("sanitises input", func() {
		s.UserService.FindUserByUsernameFn = func(ctx context.Context, username string) (*palikka.User, error) {
			// password is not sanitized
			hash, _ := bcrypt.GenerateFromPassword([]byte(" password\t "), bcrypt.DefaultCost)
			hashString := string(hash)
			return &palikka.User{ID: 1, Username: username, Password: &hashString}, nil
		}
		s.SessionStore.SaveFn = func(ctx context.Context, session *palikka.Session) error {
			return nil
		}

		r := struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}{
			Username: " username\n ",
			Password: " password\t ",
		}
		rBytes, _ := json.Marshal(r)

		_, res := doLoginRequest(rBytes)
		assert.Equal(s.T(), http.StatusOK, res.StatusCode)
	})

	s.Run("logs user in", func() {
		var createdSession *palikka.Session

		s.UserService.FindUserByUsernameFn = func(ctx context.Context, username string) (*palikka.User, error) {
			return &palikka.User{ID: 1, Username: username, Password: &passwordHash}, nil
		}
		s.SessionStore.SaveFn = func(ctx context.Context, session *palikka.Session) error {
			// set some values like the real session store would
			session.ID = "mock-session-id"
			session.MaxAge = time.Minute
			session.Expires = time.Now().Add(time.Minute)

			createdSession = session

			return nil
		}

		_, res := doLoginRequest(loginRequestJson)
		assert.Equal(s.T(), http.StatusOK, res.StatusCode)

		// check session
		assert.Equal(s.T(), "Go-http-client/1.1", createdSession.Values["User-Agent"].(string))
		assert.Contains(s.T(), createdSession.Values["Host"], "localhost")

		// check session cookie
		setCookie := res.Header.Get("Set-Cookie")
		assert.NotNil(s.T(), setCookie)

		cookie, err := http.ParseSetCookie(setCookie)
		assert.NoError(s.T(), err)
		assert.NotNil(s.T(), cookie)
		assert.Equal(s.T(), s.CookieService.CookieName(), cookie.Name)
		assert.True(s.T(), len(cookie.Value) > sha256.Size*2)
		assert.Contains(s.T(), cookie.Value, "mock-session-id")
		assert.True(s.T(), cookie.HttpOnly)
		assert.False(s.T(), cookie.Secure)
		assert.Equal(s.T(), int(time.Minute.Seconds()), cookie.MaxAge)
		assert.True(s.T(), cookie.Expires.After(time.Now()))
		assert.True(s.T(), cookie.Expires.Before(time.Now().Add(time.Minute*2)))
	})

	s.Run("test context", func() {
		// todo test context with httptest?
	})
}

func (s *AuthNTestSuite) TestAuth_Logout() {
	s.Run("returns error without cookie", func() {
		req, _ := http.NewRequest(http.MethodPost, s.URL()+"/auth/logout", nil)
		res, _ := http.DefaultClient.Do(req)
		MustRespondWithError(s.T(), res, http.StatusUnauthorized, palikkahttp.ErrMsgUnauthorized)
	})

	s.Run("returns error with invalid cookie value", func() {
		cookie := &http.Cookie{
			Name:  s.CookieService.CookieName(),
			Value: "session-id",
		}
		req, _ := http.NewRequest(http.MethodPost, s.URL()+"/auth/logout", nil)
		req.AddCookie(cookie)

		res, _ := http.DefaultClient.Do(req)
		MustRespondWithError(s.T(), res, http.StatusUnauthorized, palikkahttp.ErrMsgUnauthorized)
	})

	s.Run("returns error with invalid cookie signature", func() {
		cookieService := palikkacookie.NewService(
			&palikkacookie.Options{
				CookieName: s.CookieService.CookieName(),
				Secret:     []byte("invalid-secret"),
			},
		)

		session := &palikka.Session{
			ID:      "mock-session-id",
			UserID:  0,
			MaxAge:  time.Minute,
			Expires: time.Now().Add(time.Minute),
			Values:  nil,
		}

		cookie := cookieService.New(session)

		req, _ := http.NewRequest(http.MethodPost, s.URL()+"/auth/logout", nil)
		req.AddCookie(cookie)

		res, _ := http.DefaultClient.Do(req)
		MustRespondWithError(s.T(), res, http.StatusUnauthorized, palikkahttp.ErrMsgUnauthorized)
	})

	s.Run("does not return error if session deletion fails", func() {
		// todo
	})

	s.Run("logs user out", func() {
		sessionDeleted := false

		session := &palikka.Session{
			ID:      "mock-session-id",
			UserID:  1,
			MaxAge:  time.Minute,
			Expires: time.Now().Add(time.Minute),
		}

		hash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		hashString := string(hash)
		user := &palikka.User{ID: 1, Username: "username", Password: &hashString}

		s.UserService.FindUserByIDFn = func(ctx context.Context, id int64) (*palikka.User, error) {
			return user, nil
		}
		s.UserService.FindUserByUsernameFn = func(ctx context.Context, username string) (*palikka.User, error) {
			return user, nil
		}
		s.SessionStore.SaveFn = func(ctx context.Context, session *palikka.Session) error {
			session.ID = "mock-session-id"
			session.MaxAge = time.Minute
			session.Expires = time.Now().Add(time.Minute)
			return nil
		}
		s.SessionStore.LoadFn = func(ctx context.Context, id string) (*palikka.Session, error) {
			return session, nil
		}
		s.SessionStore.DeleteFn = func(ctx context.Context, sessionKey string) error {
			sessionDeleted = true
			return nil
		}

		username := "username"
		password := "password"
		loginPayload := palikkahttp.LoginRequest{Username: &username, Password: &password}
		loginPayloadJson, _ := json.Marshal(loginPayload)

		loginReq, _ := http.NewRequest(http.MethodPost, s.URL()+"/auth/login", bytes.NewBuffer(loginPayloadJson))
		loginRes, _ := http.DefaultClient.Do(loginReq)
		sessionCookie, _ := http.ParseSetCookie(loginRes.Header.Get("Set-Cookie"))

		req, _ := http.NewRequest(http.MethodPost, s.URL()+"/auth/logout", nil)
		req.AddCookie(sessionCookie)

		res, _ := http.DefaultClient.Do(req)
		assert.Equal(s.T(), http.StatusOK, res.StatusCode)
		assert.True(s.T(), sessionDeleted)

		setCookie := res.Header.Get("Set-Cookie")
		assert.NotEmpty(s.T(), setCookie)

		cookie, _ := http.ParseSetCookie(setCookie)
		_, err := s.CookieService.Verify(cookie)
		assert.Error(s.T(), err)
	})
}

type AuthNTestSuite struct {
	suite.Suite
	*Server
}

func (s *AuthNTestSuite) SetupSuite() {
	s.Server = MustOpenServer(s.T())
}

func (s *AuthNTestSuite) TearDownSuite() {
	MustCloseServer(s.T(), s.Server)
}

func (s *AuthNTestSuite) TearDownTest() {
	s.Server.MustResetMocks(s.T())
}

func (s *AuthNTestSuite) TearDownSubTest() {
	s.Server.MustResetMocks(s.T())
}

func TestAuthNTestSuite(t *testing.T) {
	suite.Run(t, new(AuthNTestSuite))
}

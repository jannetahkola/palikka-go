package http_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	palikka "palikka-go"
	palikkahttp "palikka-go/http"
	palikkacookie "palikka-go/http/cookie"
	"palikka-go/mock"
	"reflect"
	"strings"
	"testing"
	"time"
)

type Server struct {
	*palikkahttp.Server

	CookieService     palikkacookie.Service
	SessionStore      mock.SessionStore
	PermissionService mock.PermissionService
	RoleService       mock.RoleService
	UserService       mock.UserService
}

func MustOpenServer(t *testing.T) *Server {
	t.Helper()

	s := &Server{
		Server: palikkahttp.NewServer(&palikkahttp.ServerOpts{
			WebAppURI: "mock-web-app-uri",
			Cors:      &palikkahttp.CorsOpts{AllowedOrigins: []string{"*"}},
		}),
	}

	// cookie service is not mocked
	s.CookieService = palikkacookie.NewService(
		&palikkacookie.Options{
			CookieName: "palikka_session",
			Secret:     []byte("test-secret"),
		},
	)

	// set dependencies
	s.Server.CookieService = s.CookieService
	s.Server.SessionStore = &s.SessionStore
	s.Server.PermissionService = &s.PermissionService
	s.Server.RoleService = &s.RoleService
	s.Server.UserService = &s.UserService

	if err := s.Server.Start(); err != nil {
		t.Fatal(err)
	}

	return s
}

func MustCloseServer(t *testing.T, s *Server) {
	t.Helper()
	err := s.Server.Stop()
	if err != nil {
		t.Fatal(err)
	}
}

func MustRespondWithError(t *testing.T, res *http.Response, status int, msgSubstring string) {
	t.Helper()

	if res.StatusCode != status {
		t.Fatalf("res.StatusCode = %d; want %d", res.StatusCode, status)
	}

	if res.Header.Get("Content-Type") != "application/json" {
		t.Fatalf("Content-Type is %s, expected application/json, status=%d", res.Header.Get("Content-Type"), res.StatusCode)
	}

	body := struct {
		Error string `json:"error"`
	}{}
	_ = json.NewDecoder(res.Body).Decode(&body)

	if !strings.Contains(body.Error, msgSubstring) {
		t.Fatalf("res.respondError() = %q; must contain %q", body.Error, msgSubstring)
	}
}

func (s *Server) MustResetMocks(t *testing.T) {
	t.Helper()

	resetCount := 0
	resetMockNames := ""

	p := reflect.ValueOf(s).Elem()
	for i := 0; i < p.NumField(); i++ {
		m := p.Field(i).Addr().MethodByName("ResetMocks")
		if m.IsValid() {
			m.Call([]reflect.Value{})
			if resetMockNames != "" {
				resetMockNames = resetMockNames + ", "
			}
			resetMockNames = resetMockNames + p.Field(i).Type().String()
			resetCount++
		}
	}

	fmt.Println(fmt.Sprintf("[test] reset %d mock(s): [%s]", resetCount, resetMockNames))
}

func (s *Server) MustAuthenticate(t *testing.T, r *http.Request, user *palikka.User) {
	t.Helper()

	session := &palikka.Session{
		ID:      "mock-session-id",
		UserID:  user.ID,
		MaxAge:  time.Minute,
		Expires: time.Now().Add(time.Minute),
	}

	s.SessionStore.LoadFn = func(ctx context.Context, id string) (*palikka.Session, error) {
		if id == session.ID {
			return session, nil
		}
		return nil, errors.New("session not found")
	}

	s.UserService.FindUserByIDFn = func(ctx context.Context, id int64) (*palikka.User, error) {
		if id == user.ID {
			return user, nil
		}
		return nil, errors.New(palikka.ErrNotFound)
	}

	cookie := s.CookieService.New(session)

	r.AddCookie(cookie)
}

func (s *Server) MustAuthorize(t *testing.T, r *http.Request, user *palikka.User, roleID int, authPerms ...string) {
	t.Helper()

	s.MustAuthenticate(t, r, user)

	// create perms
	var perms []*palikka.Permission
	for i, perm := range authPerms {
		perms = append(perms, &palikka.Permission{
			ID:   int64(i),
			Name: perm,
		})
	}

	// link perms to role and role to user
	roleID64 := int64(roleID)

	role := &palikka.Role{
		ID:          roleID64,
		Name:        fmt.Sprintf("mock-role-%d", roleID64),
		Permissions: perms,
	}
	user.RoleID = &roleID64
	user.Role = role

	s.PermissionService.FindPermissionsByRoleIDFn = func(ctx context.Context, roleID int64) ([]*palikka.Permission, error) {
		if user.RoleID != nil && *user.RoleID == roleID {
			return user.Role.Permissions, nil
		}
		return nil, errors.New(palikka.ErrNotFound)
	}
}

func TestServer_ReturnsUnauthorized_WhenNoAuthN_AndRouteNotFound(t *testing.T) {
	s := MustOpenServer(t)
	defer MustCloseServer(t, s)

	req, _ := http.NewRequest("GET", s.URL()+"/notfound", nil)
	res, _ := http.DefaultClient.Do(req)

	MustRespondWithError(t, res, http.StatusUnauthorized, palikkahttp.ErrMsgInvalidRequest)
}

//func TestServer_Sets_Request_Context(t *testing.T) {
//	s := MustOpenServer(t)
//	defer MustCloseServer(t, s)
//
//	var loggerFromContext *slog.Logger
//	s.SessionStore.LoadFn = func(ctx context.Context, id string) (*palikka.Session, error) {
//		loggerFromContext = palikka.LoggerFromContext(ctx)
//		return nil, errors.New("")
//	}
//
//	req, _ := http.NewRequest("GET", s.URL()+"/users", nil)
//	cookie := &http.Cookie{Name: palikkahttp.SessionCookieName, Value: "mock-session-id"}
//	req.AddCookie(cookie)
//
//	http.DefaultClient.Do(req)
//	assert.NotNil(t, loggerFromContext)
//}
//
//func TestServer_RequireAuth(t *testing.T) {
//	s := MustOpenServer(t)
//	defer MustCloseServer(t, s)
//
//	t.Run("returns unauthorized without cookie", func(t *testing.T) {
//		req, _ := http.NewRequest("GET", s.URL()+"/users", nil)
//		res, _ := http.DefaultClient.Do(req)
//		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
//	})
//
//	t.Run("returns unauthorized without session", func(t *testing.T) {
//		s.SessionStore.LoadFn = func(ctx context.Context, id string) (*palikka.Session, error) {
//			return nil, errors.New("session not found")
//		}
//
//		req, _ := http.NewRequest("GET", s.URL()+"/users", nil)
//		cookie := &http.Cookie{Name: palikkahttp.SessionCookieName, Value: "mock-session-id"}
//		req.AddCookie(cookie)
//
//		res, err := http.DefaultClient.Do(req)
//		assert.NoError(t, err)
//		assert.NotEmpty(t, res)
//		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
//		assert.Equal(t, "application/json", res.Header.Get("Content-Type"))
//	})
//
//	// todo should probably be in http_test.go?
//
//	t.Run("returns json errors", func(t *testing.T) {
//		req, _ := http.NewRequest("GET", s.URL()+"/users", nil)
//		res, _ := http.DefaultClient.Do(req)
//		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
//		assert.Equal(t, "application/json", res.Header.Get("Content-Type"))
//
//		var errorResponse palikkahttp.errorResponse
//		json.NewDecoder(res.Body).Decode(&errorResponse)
//		assert.Equal(t, "not authorized", errorResponse.Error)
//		fmt.Println(errorResponse.Error)
//	})
//
//	t.Run("sets context and returns ok", func(t *testing.T) {
//		// mock return values
//		var (
//			sessionID = "mock-session-id"
//			session   = &palikka.Session{
//				ID:     sessionID,
//				UserID: 1,
//				Values: nil,
//			}
//			userID = int64(1)
//			user   = &palikka.User{
//				ID:        userID,
//				Username:  "mock-username",
//				Password:  "mock-password",
//				Hash:      "mock-hash",
//				CreatedAt: "",
//				Active:    true,
//				RoleID:    nil,
//				Role:      nil,
//			}
//		)
//
//		// context values
//		var (
//			loggerFromContext  *slog.Logger
//			sessionFromContext *palikka.Session
//			userFromContext    *palikka.User
//		)
//
//		// mock functions
//		s.SessionStore.LoadFn = func(ctx context.Context, id string) (*palikka.Session, error) {
//			if id == "mock-session-id" {
//				return session, nil
//			}
//			return nil, errors.New("session not found")
//		}
//		s.UserService.FindUserByIDFn = func(ctx context.Context, id int64) (*palikka.User, error) {
//			if id == userID {
//				return user, nil
//			}
//			return nil, errors.New("user not found")
//		}
//
//		// todo context can be checked from request instead
//
//		s.UserService.FindUsersFn = func(ctx context.Context) ([]*palikka.User, error) {
//			loggerFromContext = palikka.LoggerFromContext(ctx)
//			sessionFromContext = palikka.SessionFromContext(ctx)
//			userFromContext = palikka.UserFromContext(ctx)
//			return nil, nil
//		}
//
//		// do request
//		req, _ := http.NewRequest("GET", s.URL()+"/users", nil)
//		cookie := &http.Cookie{Name: palikkahttp.SessionCookieName, Value: "mock-session-id"}
//		req.AddCookie(cookie)
//
//		res, err := http.DefaultClient.Do(req)
//		assert.NoError(t, err)
//		assert.NotEmpty(t, res)
//		assert.Equal(t, http.StatusOK, res.StatusCode)
//
//		// check context
//		assert.NotNil(t, loggerFromContext)
//		assert.NotNil(t, sessionFromContext)
//		assert.NotNil(t, userFromContext)
//
//		assert.Equal(t, sessionID, sessionFromContext.ID)
//		assert.Equal(t, userID, userFromContext.ID)
//	})
//}

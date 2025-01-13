package middleware

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"palikka-go/internal/logs"
	"palikka-go/internal/session"
	"testing"
)

func (s *MiddlewareTestSuite) TestLoggerMiddleware_EnsureAuth_Sets_Request_Attrs_To_Context_And_Logs() {
	logHandlerMock := new(logs.ContextualHandlerMock)
	logHandlerMock.
		On("Enabled", mock.Anything, slog.LevelInfo).
		Return(true).
		Times(2)
	logHandlerMock.
		On("Handle", mock.Anything, mock.Anything).
		Return(nil).
		Times(2)

	slog.SetDefault(slog.New(logHandlerMock))

	var method, path string
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method, _ = logs.ContextKeyRequestMethod().Get(r.Context())
		path, _ = logs.ContextKeyRequestPath().Get(r.Context())
	})

	req := httptest.NewRequest("GET", "/", nil)

	mw := NewLoggerMiddleware()
	mw.LogRequests(h).ServeHTTP(httptest.NewRecorder(), req)

	assert.Equal(s.T(), "GET", method)
	assert.Equal(s.T(), "/", path)

	// this is true only when "Enabled" returned true
	logHandlerMock.AssertNumberOfCalls(s.T(), "Handle", 2)
}

func (s *MiddlewareTestSuite) TestSessionMiddleware_EnsureAuth_Sets_Session_Attrs_To_Context_And_Logs() {
	logHandlerMock := new(logs.ContextualHandlerMock)
	logHandlerMock.
		On("Enabled", mock.Anything, slog.LevelInfo).
		Return(true).
		Times(1)
	logHandlerMock.
		On("Handle", mock.Anything, mock.Anything).
		Return(nil).
		Times(1)

	slog.SetDefault(slog.New(logHandlerMock))

	var sessionID string
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionID, _ = logs.ContextKeySessionID().Get(r.Context())
	})

	session := new(session.Session)
	session.ID = "mock-session-id"

	storeMock := new(session.StoreMock)
	storeMock.On("Load", mock.Anything).Return(session, nil)

	mw := NewSessionMiddleware(storeMock)
	mw.EnsureAuth(h).ServeHTTP(
		httptest.NewRecorder(),
		httptest.NewRequest("GET", "/", nil),
	)

	assert.Equal(s.T(), "mock-session-id", sessionID)

	logHandlerMock.AssertNumberOfCalls(s.T(), "Handle", 1)
	storeMock.AssertNumberOfCalls(s.T(), "Load", 1)
}

func (s *MiddlewareTestSuite) TestSessionMiddleware_EnsureAuth_Returns_Unauthorized_If_Session_Not_Loaded() {
	storeMock := new(session.StoreMock)
	storeMock.On("Load", mock.Anything).Return(nil, errors.New(""))

	rec := httptest.NewRecorder()
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	mw := NewSessionMiddleware(storeMock)
	mw.EnsureAuth(h).ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))

	assert.Equal(s.T(), http.StatusUnauthorized, rec.Code)
}

type MiddlewareTestSuite struct {
	suite.Suite
}

func TestMiddleWareTestSuite(t *testing.T) {
	suite.Run(t, new(MiddlewareTestSuite))
}

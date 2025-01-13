package http_test

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	palikka "palikka-go"
	palikkahttp "palikka-go/http"
	"strings"
	"testing"
)

func (s *SessionTestSuite) TestSession_GetCurrent() {
	s.Run("requires auth", func() {
		req, _ := http.NewRequest("GET", s.URL()+"/sessions/session", nil)
		res, err := http.DefaultClient.Do(req)
		assert.NoError(s.T(), err)
		MustRespondWithError(s.T(), res, http.StatusUnauthorized, palikkahttp.ErrMsgUnauthorized)
	})

	s.Run("does not require perms", func() {
		req, _ := http.NewRequest("GET", s.URL()+"/sessions/session", nil)
		s.MustAuthenticate(s.T(), req, &palikka.User{ID: 1, Username: "mock-user"})

		res, err := http.DefaultClient.Do(req)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), http.StatusOK, res.StatusCode)
	})

	s.Run("returns current session", func() {
		req, _ := http.NewRequest("GET", s.URL()+"/sessions/session", nil)
		s.MustAuthorize(s.T(), req, &palikka.User{ID: 1, Username: "mock-user"}, 1, "users_read", "users_create")

		res, err := http.DefaultClient.Do(req)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), http.StatusOK, res.StatusCode)
		assert.Equal(s.T(), palikkahttp.MediaTypeJSON, res.Header.Get("Content-Type"))

		// use a map to confirm JSON field naming
		var session map[string]interface{}
		err = json.NewDecoder(res.Body).Decode(&session)
		assert.NoError(s.T(), err)
		assert.Empty(s.T(), session["id"])
		assert.True(s.T(), strings.HasSuffix(session["expires"].(string), "Z"))

		user, ok := session["user"].(map[string]interface{})
		assert.True(s.T(), ok)
		assert.Equal(s.T(), 1, int(user["id"].(float64)))
		assert.Equal(s.T(), "mock-user", user["username"])
		assert.NotNil(s.T(), user["permissions"])

		// use a struct here to get permissions as []string easily
		userStr, err := json.Marshal(session["user"])
		assert.NoError(s.T(), err)

		var sessionUser = struct {
			Permissions []string `json:"permissions"`
		}{}

		err = json.Unmarshal(userStr, &sessionUser)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), []string{"users_read", "users_create"}, sessionUser.Permissions)
	})
}

type SessionTestSuite struct {
	suite.Suite
	*Server
}

func (s *SessionTestSuite) SetupSuite() {
	s.Server = MustOpenServer(s.T())
}

func (s *SessionTestSuite) TearDownSuite() {
	MustCloseServer(s.T(), s.Server)
}

func (s *SessionTestSuite) TearDownTest() {
	s.Server.MustResetMocks(s.T())
}

func (s *SessionTestSuite) TearDownSubTest() {
	s.Server.MustResetMocks(s.T())
}

func TestSessionTestSuite(t *testing.T) {
	suite.Run(t, new(SessionTestSuite))
}

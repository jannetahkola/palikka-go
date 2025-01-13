package http_test

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	palikka "palikka-go"
	palikkahttp "palikka-go/http"
	"testing"
	"time"
)

func (s *AuthZTestSuite) TestAuthZ_GetAuthorizations() {
	s.Run("returns authorizations", func() {
		session := &palikka.Session{
			ID:      "mock-session-id",
			UserID:  1,
			MaxAge:  time.Minute,
			Expires: time.Now().Add(time.Minute),
		}
		hash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		hashString := string(hash)
		user := &palikka.User{ID: 1, Username: "username", Password: &hashString}

		s.Run("for user without roles", func() {
			user.Role = nil

			s.SessionStore.LoadFn = func(ctx context.Context, id string) (*palikka.Session, error) {
				return session, nil
			}
			s.UserService.FindUserByIDFn = func(ctx context.Context, id int64) (*palikka.User, error) {
				return user, nil
			}

			req, _ := http.NewRequest("GET", s.URL()+"/authorizations", nil)
			s.MustAuthenticate(s.T(), req, user)
			res, err := http.DefaultClient.Do(req)
			assert.NoError(s.T(), err)
			assert.Equal(s.T(), http.StatusOK, res.StatusCode)

			var authorizations palikkahttp.AuthorizationsResponse
			err = json.NewDecoder(res.Body).Decode(&authorizations)
			assert.NoError(s.T(), err)
			assert.False(s.T(), authorizations.Features.Users.Read)
			assert.False(s.T(), authorizations.Features.Roles.Read)
		})

		s.Run("for user with roles", func() {

			user.Role = &palikka.Role{
				ID:   1,
				Name: "mock-role",
				Rank: int64(0),
				Permissions: []*palikka.Permission{
					{
						ID:   int64(1),
						Name: "users_read",
					},
				},
			}

			s.SessionStore.LoadFn = func(ctx context.Context, id string) (*palikka.Session, error) {
				return session, nil
			}
			s.UserService.FindUserByIDFn = func(ctx context.Context, id int64) (*palikka.User, error) {
				return user, nil
			}

			req, _ := http.NewRequest("GET", s.URL()+"/authorizations", nil)
			s.MustAuthenticate(s.T(), req, user)
			res, err := http.DefaultClient.Do(req)
			assert.NoError(s.T(), err)
			assert.Equal(s.T(), http.StatusOK, res.StatusCode)

			var authorizations palikkahttp.AuthorizationsResponse
			err = json.NewDecoder(res.Body).Decode(&authorizations)
			assert.NoError(s.T(), err)
			assert.True(s.T(), authorizations.Features.Users.Read)
			assert.False(s.T(), authorizations.Features.Roles.Read)
		})
	})
}

type AuthZTestSuite struct {
	suite.Suite
	*Server
}

func (s *AuthZTestSuite) SetupSuite() {
	s.Server = MustOpenServer(s.T())
}

func (s *AuthZTestSuite) TearDownSuite() {
	MustCloseServer(s.T(), s.Server)
}

func (s *AuthZTestSuite) TearDownTest() {
	s.Server.MustResetMocks(s.T())
}

func (s *AuthZTestSuite) TearDownSubTest() {
	s.Server.MustResetMocks(s.T())
}

func TestAuthZTestSuite(t *testing.T) {
	suite.Run(t, new(AuthZTestSuite))
}

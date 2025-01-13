package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	palikka "palikka-go"
	palikkahttp "palikka-go/http"
	"testing"
)

func (s *UserTestSuite) TestUser_GetUsers() {
	s.Run("requires auth", func() {
		req, _ := http.NewRequest("GET", s.URL()+"/users", nil)
		res, err := http.DefaultClient.Do(req)
		assert.NoError(s.T(), err)
		MustRespondWithError(s.T(), res, http.StatusUnauthorized, palikkahttp.ErrMsgUnauthorized)
	})

	s.Run("requires perms", func() {
		req, _ := http.NewRequest("GET", s.URL()+"/users", nil)
		s.MustAuthorize(s.T(), req, &palikka.User{ID: 1}, 1, "invalid")

		res, err := http.DefaultClient.Do(req)
		assert.NoError(s.T(), err)
		MustRespondWithError(s.T(), res, http.StatusForbidden, palikkahttp.ErrMsgForbidden)
	})

	s.Run("returns list of users", func() {
		userRoleID := int64(1)

		user := &palikka.User{ID: int64(1), RoleID: &userRoleID}
		user2 := &palikka.User{ID: int64(2)}

		s.UserService.FindUsersFn = func(ctx context.Context) ([]*palikka.User, error) {
			return []*palikka.User{user, user2}, nil
		}
		s.RoleService.FindRoleByIDFn = func(ctx context.Context, roleID int64) (*palikka.Role, error) {
			return &palikka.Role{ID: roleID}, nil
		}

		req, _ := http.NewRequest("GET", s.URL()+"/users", nil)
		s.MustAuthorize(s.T(), req, &palikka.User{ID: 1}, 1, palikka.PermUsersRead)

		res, err := http.DefaultClient.Do(req)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), 200, res.StatusCode)

		var usersResult []map[string]interface{}
		err = json.NewDecoder(res.Body).Decode(&usersResult)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), 2, len(usersResult))

		userResult := usersResult[0]
		assert.Equal(s.T(), user.ID, int64(userResult["id"].(float64)))
		assert.Equal(s.T(), *user.RoleID, int64(userResult["role_id"].(float64)))
		assert.NotNil(s.T(), userResult["role"])

		userResult = usersResult[1]
		assert.Equal(s.T(), user2.ID, int64(userResult["id"].(float64)))
		assert.Nil(s.T(), userResult["role_id"])
		assert.Nil(s.T(), userResult["role"])
	})
}

func (s *UserTestSuite) TestUser_GetUserByID() {
	s.Run("requires auth", func() {
		req, _ := http.NewRequest("GET", s.URL()+"/users/1", nil)
		res, err := http.DefaultClient.Do(req)
		assert.NoError(s.T(), err)
		MustRespondWithError(s.T(), res, http.StatusUnauthorized, palikkahttp.ErrMsgUnauthorized)
	})

	s.Run("requires perms", func() {
		req, _ := http.NewRequest("GET", s.URL()+"/users/1", nil)
		s.MustAuthorize(s.T(), req, &palikka.User{ID: 1}, 1, "invalid")

		res, err := http.DefaultClient.Do(req)
		assert.NoError(s.T(), err)
		MustRespondWithError(s.T(), res, http.StatusForbidden, palikkahttp.ErrMsgForbidden)
	})

	s.Run("returns error if path value is invalid", func() {
		req, _ := http.NewRequest("GET", s.URL()+"/users/1_", nil)
		s.MustAuthorize(s.T(), req, &palikka.User{ID: 1}, 1, palikka.PermUsersRead)

		res, err := http.DefaultClient.Do(req)
		assert.NoError(s.T(), err)
		MustRespondWithError(s.T(), res, http.StatusBadRequest, "invalid user id, must be an integer")
	})

	s.Run("returns error if user not found", func() {
		req, _ := http.NewRequest("GET", s.URL()+"/users/2", nil)
		s.MustAuthorize(s.T(), req, &palikka.User{ID: 1}, 1, palikka.PermUsersRead)

		res, err := http.DefaultClient.Do(req)
		assert.NoError(s.T(), err)
		MustRespondWithError(s.T(), res, http.StatusNotFound, "user not found")
	})

	s.Run("returns user", func() {
		req, _ := http.NewRequest("GET", s.URL()+"/users/1", nil) // same user as authenticated

		authUserPassword := "mock-pass"
		authUser := &palikka.User{ID: 1, Username: "mock-user", Password: &authUserPassword}
		s.MustAuthorize(s.T(), req, authUser, 1, palikka.PermUsersRead)

		res, err := http.DefaultClient.Do(req)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), 200, res.StatusCode)

		// check JSON
		var user map[string]interface{}
		err = json.NewDecoder(res.Body).Decode(&user)
		assert.NoError(s.T(), err)

		assert.Equal(s.T(), "mock-user", user["username"])
		assert.NotContains(s.T(), user, "password")
		assert.Contains(s.T(), user, "role_id")
		assert.Contains(s.T(), user, "role")
	})
}

func (s *UserTestSuite) TestUser_CreateUser() {
	s.Run("requires auth", func() {
		req, _ := http.NewRequest("POST", s.URL()+"/users", nil)
		res, err := http.DefaultClient.Do(req)
		assert.NoError(s.T(), err)
		MustRespondWithError(s.T(), res, http.StatusUnauthorized, palikkahttp.ErrMsgUnauthorized)
	})

	s.Run("requires perms", func() {
		req, _ := http.NewRequest("POST", s.URL()+"/users", nil)
		s.MustAuthorize(s.T(), req, &palikka.User{ID: 1}, 1, "invalid")
		res, err := http.DefaultClient.Do(req)
		assert.NoError(s.T(), err)
		MustRespondWithError(s.T(), res, http.StatusForbidden, palikkahttp.ErrMsgForbidden)
	})

	// todo check invalid args, no need to check everything because we have input_test.go

	s.Run("creates user", func() {
		var (
			username = "mock-username"
			password = " mock-password\t " // test that spaces do not get stripped
		)

		s.UserService.InsertUserFn = func(ctx context.Context, args *palikka.UserCreate) (int64, error) {

			// check password is hashed
			assert.NotEqual(s.T(), password, *args.Password)
			err := bcrypt.CompareHashAndPassword([]byte(*args.Password), []byte(password))
			assert.NoError(s.T(), err)

			return 2, nil
		}

		requestBody := &palikka.UserCreate{Username: &username, Password: &password}
		requestBodyJson, err := json.Marshal(requestBody)
		assert.NoError(s.T(), err)

		req, _ := http.NewRequest("POST", s.URL()+"/users", bytes.NewBuffer(requestBodyJson))
		s.Server.MustAuthorize(s.T(), req, &palikka.User{ID: 1}, 1, palikka.PermUsersCreate)

		res, err := http.DefaultClient.Do(req)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), http.StatusCreated, res.StatusCode)
		assert.Contains(s.T(), "/users/2", res.Header.Get(palikkahttp.HeaderLocation))
	})
}

type UserTestSuite struct {
	suite.Suite
	*Server
}

func (s *UserTestSuite) SetupSuite() {
	s.Server = MustOpenServer(s.T())
}

func (s *UserTestSuite) TearDownSuite() {
	MustCloseServer(s.T(), s.Server)
}

func (s *UserTestSuite) TearDownTest() {
	s.Server.MustResetMocks(s.T())
}

func (s *UserTestSuite) TearDownSubTest() {
	s.Server.MustResetMocks(s.T())
}

func TestUserTestSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}

package sqlite_test

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	palikka "palikka-go"
	"palikka-go/repository/sqlite"
	"testing"
)

func (s *UserServiceTestSuite) TestUserService_FindUsers() {
	s.Run("returns empty list on empty result set", func() {
		users, err := s.userService.FindUsers(context.Background())
		assert.NoError(s.T(), err)
		assert.NotNil(s.T(), users)
		assert.Empty(s.T(), users)
	})

	s.Run("returns list of users", func() {

		// todo continue with more users

		username := "mock-user"
		active := true
		userCreate := &palikka.UserCreate{
			Username: &username,
			Active:   &active,
		}
		id, err := s.userService.InsertUser(context.Background(), userCreate)
		assert.NoError(s.T(), err)

		users, err := s.userService.FindUsers(context.Background())
		assert.NoError(s.T(), err)
		assert.Len(s.T(), users, 1)

		assert.Equal(s.T(), id, users[0].ID)
		assert.Equal(s.T(), *userCreate.Username, users[0].Username)
		assert.Equal(s.T(), *userCreate.Active, users[0].Active)
	})
}

func (s *UserServiceTestSuite) TestUserService_FindUserByID() {
	s.Run("returns error if user not found", func() {
		user, err := s.userService.FindUserByID(context.Background(), int64(999))
		assert.Error(s.T(), err)
		assert.ErrorContains(s.T(), err, "sqlite: query find user by id")
		assert.ErrorContains(s.T(), err, sql.ErrNoRows.Error())
		assert.Empty(s.T(), user)
	})

	s.Run("returns user", func() {
		args := palikka.UserCreate{}
		id, _ := s.userService.InsertUser(context.Background(), &args)
		user, err := s.userService.FindUserByID(context.Background(), id)
		assert.NoError(s.T(), err)
		assert.NotNil(s.T(), user)
		assert.Equal(s.T(), id, user.ID)
	})
}

func (s *UserServiceTestSuite) TestUserService_FindUserByUsername() {
	s.Run("returns error if user not found", func() {
		user, err := s.userService.FindUserByUsername(context.Background(), "mock-user")
		assert.Error(s.T(), err)
		assert.ErrorContains(s.T(), err, "sqlite: query find user by username")
		assert.ErrorContains(s.T(), err, sql.ErrNoRows.Error())
		assert.Empty(s.T(), user)
	})

	s.Run("returns user", func() {
		username, username2 := "mock-user", "mock-user2"
		args := palikka.UserCreate{Username: &username}
		args2 := palikka.UserCreate{Username: &username2}

		id, _ := s.userService.InsertUser(context.Background(), &args)
		_, _ = s.userService.InsertUser(context.Background(), &args2)

		user, err := s.userService.FindUserByID(context.Background(), id)
		assert.NoError(s.T(), err)
		assert.NotNil(s.T(), user)
		assert.Equal(s.T(), username, user.Username)
	})
}

func (s *UserServiceTestSuite) TestUserService_InsertUser() {
	s.Run("converts nil bool to false", func() {
		args := palikka.UserCreate{}
		id, _ := s.userService.InsertUser(context.Background(), &args)
		user, _ := s.userService.FindUserByID(context.Background(), id)
		assert.False(s.T(), user.Active)
	})

	s.Run("converts false bool to false", func() {
		active := false
		args := palikka.UserCreate{Active: &active}
		id, _ := s.userService.InsertUser(context.Background(), &args)
		user, _ := s.userService.FindUserByID(context.Background(), id)
		assert.False(s.T(), user.Active)
	})

	s.Run("converts true bool to true", func() {
		active := true
		args := palikka.UserCreate{Active: &active}
		id, _ := s.userService.InsertUser(context.Background(), &args)
		user, _ := s.userService.FindUserByID(context.Background(), id)
		assert.True(s.T(), user.Active)
	})

	s.Run("returns error on duplicate username", func() {
		username := "mock-user"

		args := palikka.UserCreate{Username: &username}
		_, _ = s.userService.InsertUser(context.Background(), &args)

		args2 := palikka.UserCreate{Username: &username}
		_, err := s.userService.InsertUser(context.Background(), &args2)
		assert.Error(s.T(), err)
		assert.ErrorContains(s.T(), err, "sqlite: query insert user")
		assert.Contains(s.T(), err.Error(), "UNIQUE constraint failed: user.username")
	})

	s.Run("inserts user", func() {
		username, password := "mock-user", "mock-password"
		active := true
		roleID := int64(1) // todo this should fail or else need to check if role exists before insert
		args := palikka.UserCreate{
			Username: &username,
			Password: &password,
			Active:   &active,
			RoleID:   &roleID,
		}

		id, _ := s.userService.InsertUser(context.Background(), &args)
		user, _ := s.userService.FindUserByID(context.Background(), id)

		assert.Equal(s.T(), username, user.Username)
		assert.Equal(s.T(), password, user.Password)
		assert.True(s.T(), user.Active)
		assert.Equal(s.T(), roleID, *user.RoleID)
	})
}

func (s *UserServiceTestSuite) TestUserService_UpdateUser() {
	s.Run("returns error on duplicate username", func() {
		var (
			username  = "mock-user"
			username2 = "mock-user2"
		)

		args := palikka.UserCreate{Username: &username}
		id, _ := s.userService.InsertUser(context.Background(), &args)

		args2 := palikka.UserCreate{Username: &username2}
		_, _ = s.userService.InsertUser(context.Background(), &args2)

		updateArgs := palikka.UserUpdate{ID: &id, Username: &username2}
		err := s.userService.UpdateUser(context.Background(), &updateArgs)
		assert.NotNil(s.T(), err)
		assert.ErrorContains(s.T(), err, "sqlite: query update user")
		assert.Contains(s.T(), err.Error(), "UNIQUE constraint failed: user.username")
	})

	s.Run("updates user", func() {
		var (
			username      = "mock-user"
			active        = true
			activeUpdated = false
		)

		args := palikka.UserCreate{Username: &username, Active: &active}
		id, _ := s.userService.InsertUser(context.Background(), &args)

		updateArgs := palikka.UserUpdate{ID: &id, Active: &activeUpdated}
		err := s.userService.UpdateUser(context.Background(), &updateArgs)
		assert.NoError(s.T(), err)

		user, err := s.userService.FindUserByID(context.Background(), id)
		assert.Empty(s.T(), user.Username)
		assert.Equal(s.T(), activeUpdated, user.Active)
	})
}

type UserServiceTestSuite struct {
	suite.Suite

	db          *sqlite.DB
	userService palikka.UserService
}

func (s *UserServiceTestSuite) SetupSuite() {
	s.db = MustOpenDB(s.T())
	s.userService = sqlite.NewUserService(s.db)
}

func (s *UserServiceTestSuite) TearDownSubTest() {
	// reset database for each subtest
	MustCloseDB(s.T(), s.db)
	s.SetupSuite()
}

func (s *UserServiceTestSuite) TearDownSuite() {
	MustCloseDB(s.T(), s.db)
}

func TestUserServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}

package sqlite

//
//import (
//	"context"
//	"database/sql"
//	_ "embed"
//	"github.com/stretchr/testify/suite"
//	sqlc "palikka-go/internal/data/.sqlc"
//	"testing"
//)
//
////go:embed schema_reset.sql
//var ddlReset string
//
////go:embed schema.sql
//var ddlCreate string
//
////go:embed schema_seed.sql
//var seed string
//
//// todo
////
////func (s *RepositoryTestSuite) TestRepository_FindUsers() {
////	s.Run("returns empty list when result set empty", func() {
////		args := sqlc.FindUsersParams{Limit: 10, Offset: 0}
////		users, err := s.querier.FindUsers(context.Background(), args)
////		assert.NoError(s.T(), err)
////		assert.Empty(s.T(), users)
////	})
////
////	s.Run("returns list of users when result set not empty", func() {
////		s.SeedDatabase()
////
////		args := sqlc.FindUsersParams{Limit: 10, Offset: 0}
////		users, err := s.querier.FindUsers(context.Background(), args)
////		assert.NoError(s.T(), err)
////		assert.NotEmpty(s.T(), users)
////		assert.True(s.T(), len(users) > 1) // limit = 10
////
////		for _, user := range users {
////			assert.True(s.T(), user.ID > int64(0))
////			assert.NotEmpty(s.T(), user.Username)
////			assert.NotEmpty(s.T(), user.Password)
////			assert.NotEmpty(s.T(), user.Hash)
////			assert.NotEmpty(s.T(), user.CreatedAt)
////			assert.True(s.T(), user.Active)
////		}
////	})
////
////	s.Run("returns next list of users with offset", func() {
////		s.SeedDatabase()
////
////		args := sqlc.FindUsersParams{Limit: 1, Offset: 0}
////		users, err := s.querier.FindUsers(args)
////		assert.NoError(s.T(), err)
////		assert.Equal(s.T(), args.Limit, int64(len(users)))
////		assert.Equal(s.T(), int64(1), users[0].ID)
////
////		args = sqlc.FindUsersParams{Limit: 1, Offset: 1}
////		users, err = s.querier.FindUsers(args)
////		assert.NoError(s.T(), err)
////		assert.Equal(s.T(), args.Limit, int64(len(users)))
////		assert.Equal(s.T(), int64(2), users[0].ID)
////	})
////}
////
////func (s *RepositoryTestSuite) TestRepository_FindUserByID() {
////	s.Run("returns error when result empty", func() {
////		user, err := s.querier.FindUserByID(1)
////		assert.Nil(s.T(), user)
////		assert.Error(s.T(), err)
////		assert.Equal(s.T(), sql.ErrNoRows, err)
////	})
////
////	s.Run("returns user when result not empty", func() {
////		s.SeedDatabase()
////		user, err := s.querier.FindUserByID(1)
////
////		assert.NoError(s.T(), err)
////		assert.NotEmpty(s.T(), user)
////		assert.Equal(s.T(), int64(1), user.ID)
////		assert.NotEmpty(s.T(), user.Username)
////		assert.NotEmpty(s.T(), user.Password)
////		assert.NotEmpty(s.T(), user.Hash)
////		assert.NotEmpty(s.T(), user.CreatedAt)
////		assert.True(s.T(), user.Active)
////		assert.Equal(s.T(), int64(1), *user.RoleID)
////	})
////}
////
////func (s *RepositoryTestSuite) TestRepository_FindUserByUsername() {
////	s.Run("returns error when result empty", func() {
////		user, err := s.querier.FindUserByUsername("mock-user-1")
////		assert.Nil(s.T(), user)
////		assert.Error(s.T(), err)
////		assert.Equal(s.T(), sql.ErrNoRows, err)
////	})
////
////	s.Run("returns user when result not empty", func() {
////		s.SeedDatabase()
////		user, err := s.querier.FindUserByUsername("mock-user-1")
////
////		assert.NoError(s.T(), err)
////		assert.NotEmpty(s.T(), user)
////		assert.Equal(s.T(), int64(1), user.ID)
////		assert.Equal(s.T(), "mock-user-1", user.Username)
////		assert.NotEmpty(s.T(), user.Password)
////		assert.NotEmpty(s.T(), user.Hash)
////		assert.NotEmpty(s.T(), user.CreatedAt)
////		assert.True(s.T(), user.Active)
////		assert.Equal(s.T(), int64(1), *user.RoleID)
////	})
////}
////
////func (s *RepositoryTestSuite) TestRepository_UpdateUser() {
////	s.Run("updates user", func() {
////		s.SeedDatabase()
////
////		args := sqlc.UpdateUserParams{
////			Username: "new-username",
////			Password: "new-password",
////			Hash:     "new-hash",
////			Active:   types.BoolFrom(false),
////			RoleID:   nil,
////			ID:       1,
////		}
////		err := s.querier.UpdateUser(args)
////		assert.NoError(s.T(), err)
////
////		user, err := s.querier.FindUserByID(1)
////		assert.NoError(s.T(), err)
////		assert.NotEmpty(s.T(), user)
////		assert.Equal(s.T(), "new-username", user.Username)
////		assert.Equal(s.T(), "new-password", user.Password)
////		assert.Equal(s.T(), "new-hash", user.Hash)
////		assert.False(s.T(), user.Active)
////		assert.Nil(s.T(), user.RoleID)
////	})
////}
////
////func (s *RepositoryTestSuite) TestRepository_FindRoles() {
////	s.Run("returns empty list when result set empty", func() {
////		args := sqlc.FindRolesParams{Limit: 10, Offset: 0}
////		roles, err := s.querier.FindRoles(args)
////		assert.NoError(s.T(), err)
////		assert.Empty(s.T(), roles)
////	})
////
////	s.Run("returns list of roles when result set not empty", func() {
////		s.SeedDatabase()
////
////		args := sqlc.FindRolesParams{Limit: 10, Offset: 0}
////		roles, err := s.querier.FindRoles(args)
////		assert.NoError(s.T(), err)
////		assert.NotEmpty(s.T(), roles)
////		assert.True(s.T(), len(roles) > 1) // limit = 10
////
////		for _, role := range roles {
////			assert.True(s.T(), role.ID > int64(0))
////			assert.NotEmpty(s.T(), role.Name)
////			assert.True(s.T(), role.Rank >= int64(0))
////		}
////	})
////
////	s.Run("returns next list of roles with offset", func() {
////		s.SeedDatabase()
////
////		args := sqlc.FindRolesParams{Limit: 1, Offset: 0}
////		roles, err := s.querier.FindRoles(args)
////		assert.NoError(s.T(), err)
////		assert.Equal(s.T(), args.Limit, int64(len(roles)))
////		assert.Equal(s.T(), int64(1), roles[0].ID)
////
////		args = sqlc.FindRolesParams{Limit: 1, Offset: 1}
////		roles, err = s.querier.FindRoles(args)
////		assert.NoError(s.T(), err)
////		assert.Equal(s.T(), args.Limit, int64(len(roles)))
////		assert.Equal(s.T(), int64(2), roles[0].ID)
////	})
////}
////
////func (s *RepositoryTestSuite) TestRepository_FindRoleByID() {
////	s.Run("returns error when result empty", func() {
////		role, err := s.querier.FindRoleByID(1)
////		assert.Nil(s.T(), role)
////		assert.Error(s.T(), err)
////		assert.Equal(s.T(), sql.ErrNoRows, err)
////	})
////
////	s.Run("returns user when result not empty", func() {
////		s.SeedDatabase()
////		role, err := s.querier.FindRoleByID(1)
////
////		assert.NoError(s.T(), err)
////		assert.NotEmpty(s.T(), role)
////		assert.Equal(s.T(), int64(1), role.ID)
////		assert.NotEmpty(s.T(), role.Name)
////		assert.True(s.T(), role.Rank >= int64(0))
////		assert.Nil(s.T(), role.Permissions)
////	})
////}
////
////func (s *RepositoryTestSuite) TestRepository_FindPermissions() {
////	s.Run("returns empty list when result set empty", func() {
////		args := sqlc.FindPermissionsParams{Limit: 10, Offset: 0}
////		perms, err := s.querier.FindPermissions(args)
////		assert.NoError(s.T(), err)
////		assert.Empty(s.T(), perms)
////	})
////
////	s.Run("returns list of permissions when result set not empty", func() {
////		s.SeedDatabase()
////
////		args := sqlc.FindPermissionsParams{Limit: 10, Offset: 0}
////		perms, err := s.querier.FindPermissions(args)
////		assert.NoError(s.T(), err)
////		assert.NotEmpty(s.T(), perms)
////		assert.True(s.T(), len(perms) > 1) // limit = 10
////
////		for _, perm := range perms {
////			assert.True(s.T(), perm.ID > int64(0))
////			assert.NotEmpty(s.T(), perm.Name)
////		}
////	})
////
////	s.Run("returns next list of permissions with offset", func() {
////		s.SeedDatabase()
////
////		args := sqlc.FindPermissionsParams{Limit: 1, Offset: 0}
////		perms, err := s.querier.FindPermissions(args)
////		assert.NoError(s.T(), err)
////		assert.Equal(s.T(), args.Limit, int64(len(perms)))
////		assert.Equal(s.T(), int64(1), perms[0].ID)
////
////		args = sqlc.FindPermissionsParams{Limit: 1, Offset: 1}
////		perms, err = s.querier.FindPermissions(args)
////		assert.NoError(s.T(), err)
////		assert.Equal(s.T(), args.Limit, int64(len(perms)))
////		assert.Equal(s.T(), int64(2), perms[0].ID)
////	})
////}
////
////func (s *RepositoryTestSuite) TestRepository_FindPermissionsByRoleID() {
////	s.Run("returns empty list when result set empty", func() {
////		perms, err := s.querier.FindPermissionsByRoleID(1)
////		assert.NoError(s.T(), err)
////		assert.Empty(s.T(), perms)
////	})
////
////	s.Run("returns list of permissions when result set not empty", func() {
////		s.SeedDatabase()
////		perms, err := s.querier.FindPermissionsByRoleID(1)
////
////		assert.NoError(s.T(), err)
////		assert.NotEmpty(s.T(), perms)
////		assert.Equal(s.T(), 1, len(perms))
////
////		perm := perms[0]
////		assert.Equal(s.T(), int64(1), perm.ID)
////		assert.NotEmpty(s.T(), perm.Name)
////		assert.Equal(s.T(), int64(1), *perm.RoleID)
////	})
////
////	s.Run("returns list of permissions based on role rank", func() {
////		s.SeedDatabase()
////		perms, err := s.querier.FindPermissionsByRoleID(2)
////
////		assert.NoError(s.T(), err)
////		assert.NotEmpty(s.T(), perms)
////		assert.Equal(s.T(), 3, len(perms))
////
////		containsPermsAssignedToOtherRoles := false
////		for _, perm := range perms {
////			if *perm.RoleID != int64(2) {
////				containsPermsAssignedToOtherRoles = true
////			}
////		}
////
////		assert.True(s.T(), containsPermsAssignedToOtherRoles)
////	})
////}
//
//type RepositoryTestSuite struct {
//	suite.Suite
//	*sql.DB
//	sqlc.querier
//}
//
//func (s *RepositoryTestSuite) SetupSuite() {
//	db, err := sql.Open("sqlite", ":memory:") // todo foreign keys etc.
//	if err != nil {
//		s.FailNow(err.Error())
//	}
//
//	s.DB = db
//	s.querier = NewQuerier(sqlc.New(db))
//
//	s.ResetDatabase()
//}
//
//func (s *RepositoryTestSuite) TearDownTest() {
//	// testify has no way to "SetupSubTest" so instead we use
//	// "TearDownSubTest" to prepare the database for the next
//	// one. Match that for normal tests here.
//	s.ResetDatabase()
//}
//
//func (s *RepositoryTestSuite) TearDownSubTest() {
//	s.ResetDatabase()
//}
//
//func (s *RepositoryTestSuite) ResetDatabase() {
//	ctx := context.Background()
//
//	if _, err := s.DB.ExecContext(ctx, ddlReset); err != nil {
//		s.FailNow(err.Error())
//	}
//	if _, err := s.DB.ExecContext(ctx, ddlCreate); err != nil {
//		s.FailNow(err.Error())
//	}
//}
//
//func (s *RepositoryTestSuite) SeedDatabase() {
//	if _, err := s.DB.ExecContext(context.Background(), seed); err != nil {
//		s.FailNow(err.Error())
//	}
//}
//
//func TestRepositoryTestSuite(t *testing.T) {
//	suite.Run(t, new(RepositoryTestSuite))
//}

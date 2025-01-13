package sqlite_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	palikka "palikka-go"
	"palikka-go/repository/sqlite"
	"testing"
)

func (s *PermissionTestSuite) TestPermissionService_FindPermissions() {
	s.Run("returns empty list on empty result set", func() {
		_ = s.permissionService.DeletePermissions(context.Background())
		perms, err := s.permissionService.FindPermissions(context.Background())
		assert.NoError(s.T(), err)
		assert.NotNil(s.T(), perms)
		assert.Empty(s.T(), perms)
	})

	s.Run("returns list of permissions", func() {
		perms, err := s.permissionService.FindPermissions(context.Background())
		assert.NoError(s.T(), err)
		assert.Len(s.T(), perms, 1)
		for _, perm := range perms {
			assert.NotEmpty(s.T(), perm.Name)
		}
	})
}

func (s *PermissionTestSuite) TestPermissionService_FindPermissionsByRoleID() {
	s.Run("returns empty list on empty result set", func() {
		perms, err := s.permissionService.FindPermissionsByRoleID(context.Background(), int64(999))
		assert.NoError(s.T(), err)
		assert.NotNil(s.T(), perms)
		assert.Empty(s.T(), perms)
	})

	s.Run("returns list of permissions", func() {
		var (
			roleName = "mock-role"
			roleRank = int64(0)
		)
		args := palikka.RoleCreate{
			Name: &roleName,
			Rank: &roleRank,
		}
		roleID, err := s.roleService.InsertRole(context.Background(), &args)
		assert.NoError(s.T(), err)

		permsInDb, err := s.permissionService.FindPermissions(context.Background())
		assert.NoError(s.T(), err)

		updatePermissionArgs := &palikka.PermissionUpdate{ID: &permsInDb[0].ID, RoleID: &roleID}
		err = s.permissionService.UpdatePermission(context.Background(), updatePermissionArgs)
		assert.NoError(s.T(), err)

		perms, err := s.permissionService.FindPermissionsByRoleID(context.Background(), roleID)
		assert.NoError(s.T(), err)
		assert.Len(s.T(), perms, 1)

		perm := perms[0]
		assert.NotEmpty(s.T(), perm.Name)
		assert.Equal(s.T(), roleID, *perm.RoleID)
	})
}

type PermissionTestSuite struct {
	suite.Suite

	db                *sqlite.DB
	roleService       palikka.RoleService
	permissionService palikka.PermissionService
}

func (s *PermissionTestSuite) SetupSuite() {
	s.db = MustOpenDB(s.T())
	s.roleService = sqlite.NewRoleService(s.db)
	s.permissionService = sqlite.NewPermissionService(s.db)
}

func (s *PermissionTestSuite) TearDownSubTest() {
	MustCloseDB(s.T(), s.db)
	s.SetupSuite()
}

func (s *PermissionTestSuite) TearDownSuite() {
	MustCloseDB(s.T(), s.db)
}

func TestPermissionTestSuite(t *testing.T) {
	suite.Run(t, new(PermissionTestSuite))
}

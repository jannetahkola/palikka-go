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

func (s *RoleServiceTestSuite) TestRoleService_FindRoles() {
	s.Run("returns empty list on empty result set", func() {
		roles, err := s.roleService.FindRoles(context.Background())
		assert.NoError(s.T(), err)
		assert.NotNil(s.T(), roles)
		assert.Empty(s.T(), roles)
	})

	s.Run("returns list of roles", func() {

		// create roles

		roleName := "mock-role"
		roleRank := int64(1)
		createRoleArgs := &palikka.RoleCreate{Name: &roleName, Rank: &roleRank}
		roleID, err := s.roleService.InsertRole(context.Background(), createRoleArgs)
		assert.NoError(s.T(), err)

		roleName2 := "mock-role-2"
		roleRank2 := int64(2)
		createRoleArgs2 := &palikka.RoleCreate{Name: &roleName2, Rank: &roleRank2}
		roleID2, err := s.roleService.InsertRole(context.Background(), createRoleArgs2)
		assert.NoError(s.T(), err)

		roleName3 := "mock-role-3"
		roleRank3 := int64(0) // lowest rank here, no permission inheritance
		createRoleArgs3 := &palikka.RoleCreate{Name: &roleName3, Rank: &roleRank3}
		_, err = s.roleService.InsertRole(context.Background(), createRoleArgs3)
		assert.NoError(s.T(), err)

		// create permissions for roles

		permissionName := "mock-permission"
		createPermissionArgs := &palikka.PermissionCreate{Name: &permissionName, RoleID: &roleID}
		_, err = s.permissionService.InsertPermission(context.Background(), createPermissionArgs)
		assert.NoError(s.T(), err)

		permissionName2 := "mock-permission-2"
		createPermissionArgs2 := &palikka.PermissionCreate{Name: &permissionName2, RoleID: &roleID2}
		_, err = s.permissionService.InsertPermission(context.Background(), createPermissionArgs2)
		assert.NoError(s.T(), err)

		// find roles

		roles, err := s.roleService.FindRoles(context.Background())
		assert.NoError(s.T(), err)
		assert.Len(s.T(), roles, 3)

		role := roles[0]
		assert.Equal(s.T(), int64(1), role.ID)
		assert.Equal(s.T(), roleName, role.Name)
		assert.Equal(s.T(), roleRank, role.Rank)
		assert.NotEmpty(s.T(), role.Permissions)
		assert.Equal(s.T(), 1, len(role.Permissions))
		assert.Equal(s.T(), permissionName, role.Permissions[0].Name)

		role = roles[1]
		assert.Equal(s.T(), int64(2), role.ID)
		assert.Equal(s.T(), roleName2, role.Name)
		assert.Equal(s.T(), roleRank2, role.Rank)
		assert.NotEmpty(s.T(), role.Permissions)
		assert.Equal(s.T(), 2, len(role.Permissions))
		assert.Equal(s.T(), permissionName, role.Permissions[0].Name)
		assert.Equal(s.T(), permissionName2, role.Permissions[1].Name)

		role = roles[2]
		assert.Equal(s.T(), int64(3), role.ID)
		assert.Equal(s.T(), roleName3, role.Name)
		assert.Equal(s.T(), roleRank3, role.Rank)
		assert.Empty(s.T(), role.Permissions)
	})
}

func (s *RoleServiceTestSuite) TestRoleService_FindRoleByID() {
	s.Run("returns error if role not found", func() {
		role, err := s.roleService.FindRoleByID(context.Background(), int64(1))
		assert.Error(s.T(), err)
		assert.ErrorContains(s.T(), err, "sqlite: query find role by id")
		assert.ErrorContains(s.T(), err, sql.ErrNoRows.Error())
		assert.Nil(s.T(), role)
	})

	s.Run("returns role", func() {
		// todo continue with perms?

		name, name2 := "mock-role", "mock-role2"
		rank, rank2 := int64(1), int64(2)
		args := palikka.RoleCreate{Name: &name, Rank: &rank}
		args2 := palikka.RoleCreate{Name: &name2, Rank: &rank2}

		id, err := s.roleService.InsertRole(context.Background(), &args)
		assert.NoError(s.T(), err)
		_, err = s.roleService.InsertRole(context.Background(), &args2)
		assert.NoError(s.T(), err)

		role, err := s.roleService.FindRoleByID(context.Background(), id)
		assert.NoError(s.T(), err)
		assert.NotNil(s.T(), role)
		assert.Equal(s.T(), name, role.Name)
	})
}

func (s *RoleServiceTestSuite) TestRoleService_InsertRole() {
	s.Run("returns error on duplicate name", func() {
		name, name2 := "mock-role", "mock-role"
		rank, rank2 := int64(1), int64(2)
		args := palikka.RoleCreate{Name: &name, Rank: &rank}
		args2 := palikka.RoleCreate{Name: &name2, Rank: &rank2}

		_, _ = s.roleService.InsertRole(context.Background(), &args)
		_, err := s.roleService.InsertRole(context.Background(), &args2)
		assert.Error(s.T(), err)
		assert.ErrorContains(s.T(), err, "sqlite: query insert role")
		assert.Contains(s.T(), err.Error(), "UNIQUE constraint failed: role.name")
	})

	s.Run("returns error on duplicate rank", func() {
		name, name2 := "mock-role", "mock-role-2"
		rank, rank2 := int64(1), int64(1)
		args := palikka.RoleCreate{Name: &name, Rank: &rank}
		args2 := palikka.RoleCreate{Name: &name2, Rank: &rank2}

		_, _ = s.roleService.InsertRole(context.Background(), &args)
		_, err := s.roleService.InsertRole(context.Background(), &args2)
		assert.Error(s.T(), err)
		assert.ErrorContains(s.T(), err, "sqlite: query insert role")
		assert.Contains(s.T(), err.Error(), "UNIQUE constraint failed: role.rank")
	})

	s.Run("inserts role", func() {
		name := "mock-role"
		rank := int64(1)
		args := palikka.RoleCreate{Name: &name, Rank: &rank}

		id, err := s.roleService.InsertRole(context.Background(), &args)
		assert.NoError(s.T(), err)

		role, _ := s.roleService.FindRoleByID(context.Background(), id)
		assert.NotNil(s.T(), role)
		assert.Equal(s.T(), name, role.Name)
		assert.Equal(s.T(), rank, role.Rank)
	})
}

func (s *RoleServiceTestSuite) TestRoleService_UpdateRole() {
	s.Run("returns error on duplicate name", func() {
		name, name2 := "mock-role", "mock-role-2"
		rank, rank2 := int64(1), int64(2)
		args := palikka.RoleCreate{Name: &name, Rank: &rank}
		args2 := palikka.RoleCreate{Name: &name2, Rank: &rank2}

		id, err := s.roleService.InsertRole(context.Background(), &args)
		assert.NoError(s.T(), err)
		_, err = s.roleService.InsertRole(context.Background(), &args2)
		assert.NoError(s.T(), err)

		args3 := palikka.RoleUpdate{ID: &id, Name: &name2, Rank: &rank}
		err = s.roleService.UpdateRole(context.Background(), &args3)
		assert.Error(s.T(), err)
		assert.ErrorContains(s.T(), err, "sqlite: query update role")
		assert.Contains(s.T(), err.Error(), "UNIQUE constraint failed: role.name")
	})

	s.Run("returns error on duplicate rank", func() {
		name, name2 := "mock-role", "mock-role-2"
		rank, rank2 := int64(1), int64(2)
		args := palikka.RoleCreate{Name: &name, Rank: &rank}
		args2 := palikka.RoleCreate{Name: &name2, Rank: &rank2}

		id, err := s.roleService.InsertRole(context.Background(), &args)
		assert.NoError(s.T(), err)
		_, err = s.roleService.InsertRole(context.Background(), &args2)
		assert.NoError(s.T(), err)

		args3 := palikka.RoleUpdate{ID: &id, Name: &name, Rank: &rank2}
		err = s.roleService.UpdateRole(context.Background(), &args3)
		assert.Error(s.T(), err)
		assert.ErrorContains(s.T(), err, "sqlite: query update role")
		assert.Contains(s.T(), err.Error(), "UNIQUE constraint failed: role.rank")
	})

	s.Run("updates role", func() {
		var (
			name        = "mock-role"
			rank        = int64(1)
			rankUpdated = int64(2)
		)

		args := palikka.RoleCreate{Name: &name, Rank: &rank}
		id, err := s.roleService.InsertRole(context.Background(), &args)
		assert.NoError(s.T(), err)

		args2 := palikka.RoleUpdate{ID: &id, Rank: &rankUpdated}
		err = s.roleService.UpdateRole(context.Background(), &args2)
		assert.NoError(s.T(), err)

		role, _ := s.roleService.FindRoleByID(context.Background(), id)
		assert.NotNil(s.T(), role)
		assert.Empty(s.T(), role.Name)
		assert.Equal(s.T(), rankUpdated, role.Rank)
	})
}

type RoleServiceTestSuite struct {
	suite.Suite

	db                *sqlite.DB
	permissionService palikka.PermissionService
	roleService       palikka.RoleService
}

func (s *RoleServiceTestSuite) SetupSuite() {
	s.db = MustOpenDB(s.T())
	s.permissionService = sqlite.NewPermissionService(s.db)
	s.roleService = sqlite.NewRoleService(s.db)
}

func (s *RoleServiceTestSuite) TearDownSubTest() {
	// reset database for each subtest
	MustCloseDB(s.T(), s.db)
	s.SetupSuite()
}

func (s *RoleServiceTestSuite) TearDownSuite() {
	MustCloseDB(s.T(), s.db)
}

func TestRoleServiceTestSuite(t *testing.T) {
	suite.Run(t, new(RoleServiceTestSuite))
}

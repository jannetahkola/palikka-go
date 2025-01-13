package sqlite_test

//import (
//	"context"
//	"database/sql"
//	. "github.com/ovechkin-dm/mockio/mock"
//	"github.com/ovechkin-dm/mockio/tests/common"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/suite"
//	_ "modernc.org/sqlite"
//	palikka "palikka-go"
//	sqlc "palikka-go/repository/sqlite/.sqlc"
//	"palikka-go/test"
//	"testing"
//	"time"
//)
//
//func (s *RoleServiceTestSuite) Test_FindRoles() {
//	s.Run("returns error on context timeout", func() {
//		r := common.NewMockReporter(s.T())
//		SetUp(r)
//
//		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
//		defer cancel()
//
//		m := Mock[RoleService]()
//		WhenDouble(m.FindRoles(AnyContext())).
//			ThenAnswer(func(args []any) ([]*palikka.Role, error) {
//				time.Sleep(time.Millisecond * 200)
//				return s.RoleService.FindRoles(ctx)
//			})
//
//		roles, err := m.FindRoles(ctx)
//		assert.Error(s.T(), err)
//		assert.Equal(s.T(), "context deadline exceeded", err.Error())
//		assert.Nil(s.T(), roles)
//	})
//
//	s.Run("returns empty list when no rows", func() {
//		result, err := s.RoleService.FindRoles(context.Background())
//		assert.NoError(s.T(), err)
//		assert.NotNil(s.T(), result)
//		assert.Empty(s.T(), result)
//	})
//
//	s.Run("returns roles with and without permissions", func() {
//		_ = test.SeedDatabase(s.DB)
//
//		result, err := s.RoleService.FindRoles(context.Background())
//		assert.NoError(s.T(), err)
//		assert.NotNil(s.T(), result)
//		assert.NotEmpty(s.T(), result)
//
//		containsPerms := false
//		containsNoPerms := false
//		for _, row := range result {
//			if len(row.Permissions) > 0 {
//				containsPerms = true
//			}
//			if len(row.Permissions) == 0 {
//				containsNoPerms = true
//			}
//			assert.NotEmpty(s.T(), row.ID)
//			assert.NotEmpty(s.T(), row.Name)
//		}
//		assert.True(s.T(), containsPerms)
//		assert.True(s.T(), containsNoPerms)
//	})
//}
//
//func (s *RoleServiceTestSuite) Test_FindRoleByID() {
//	s.Run("returns error when no rows", func() {
//		_, err := s.RoleService.FindRoleByID(context.Background(), 999)
//		assert.Error(s.T(), err)
//		assert.Equal(s.T(), sql.ErrNoRows, err)
//	})
//
//	s.Run("returns role without permissions", func() {
//		_ = test.SeedDatabase(s.DB)
//
//		role, err := s.RoleService.FindRoleByID(context.Background(), 3)
//		assert.NoError(s.T(), err)
//		assert.NotNil(s.T(), role)
//
//		assert.Equal(s.T(), 3, int(role.ID))
//		assert.Equal(s.T(), "mock-role-3", role.Name)
//		assert.Equal(s.T(), 0, int(role.Rank))
//		assert.Empty(s.T(), role.Permissions)
//	})
//
//	s.Run("returns role with permissions", func() {
//		_ = test.SeedDatabase(s.DB)
//
//		role, err := s.RoleService.FindRoleByID(context.Background(), 2)
//		assert.NoError(s.T(), err)
//		assert.NotNil(s.T(), role)
//
//		assert.Equal(s.T(), 2, int(role.ID))
//		assert.Equal(s.T(), "mock-role-2", role.Name)
//		assert.Equal(s.T(), 2, int(role.Rank))
//		assert.Equal(s.T(), 3, len(role.Permissions))
//	})
//}
//
//func (s *RoleServiceTestSuite) Test_InsertRole() {
//	s.Run("inserts role", func() {
//
//		name := "test_role"
//		rank := int64(0)
//		args := palikka.RoleCreate{
//			Name: &name,
//			Rank: &rank,
//		}
//
//		id, err := s.RoleService.InsertRole(context.Background(), args)
//		assert.NoError(s.T(), err)
//		assert.NotNil(s.T(), id)
//
//		row, err := s.RoleService.FindRoleByID(context.Background(), id)
//		assert.NoError(s.T(), err)
//		assert.NotNil(s.T(), row)
//		assert.Equal(s.T(), id, row.ID)
//		assert.Equal(s.T(), "test_role", row.Name)
//	})
//
//	s.Run("returns error when duplicate role name", func() {
//		ctx := context.Background()
//
//		args := sqlc.InsertRoleParams{
//			Name: "test_role",
//			Rank: 0,
//		}
//		_, err := s.RoleService.InsertRole(ctx, args)
//		assert.NoError(s.T(), err)
//
//		args2 := sqlc.InsertRoleParams{
//			Name: "test_role",
//			Rank: 1,
//		}
//
//		_, err = s.RoleService.InsertRole(ctx, args2)
//		assert.Error(s.T(), err)
//		assert.Contains(s.T(), err.Error(), "UNIQUE constraint failed: role.name")
//	})
//
//	s.Run("returns error when duplicate rank", func() {
//		ctx := context.Background()
//
//		args := sqlc.InsertRoleParams{
//			Name: "test_role",
//			Rank: 0,
//		}
//		_, err := s.RoleService.InsertRole(ctx, args)
//		assert.NoError(s.T(), err)
//
//		args2 := sqlc.InsertRoleParams{
//			Name: "test_role_2",
//			Rank: 0,
//		}
//		_, err = s.RoleService.InsertRole(ctx, args2)
//		assert.Error(s.T(), err)
//		assert.Contains(s.T(), err.Error(), "UNIQUE constraint failed: role.rank")
//	})
//}
//
//func (s *RoleServiceTestSuite) Test_UpdateRole() {
//	s.Run("updates role", func() {
//		ctx := context.Background()
//
//		insertArgs := sqlc.InsertRoleParams{
//			Name: "test_role",
//			Rank: 0,
//		}
//		id, err := s.RoleService.InsertRole(ctx, insertArgs)
//		assert.NoError(s.T(), err)
//
//		args := sqlc.UpdateRoleParams{
//			ID:   id,
//			Name: "updated",
//			Rank: 1,
//		}
//		err = s.RoleService.UpdateRole(ctx, args)
//		assert.NoError(s.T(), err)
//
//		row, err := s.RoleService.FindRoleByID(ctx, id)
//		assert.NoError(s.T(), err)
//		assert.Equal(s.T(), "updated", row.Name)
//		assert.Equal(s.T(), 1, int(row.Rank))
//	})
//
//	s.Run("returns error when duplicate role name", func() {
//		ctx := context.Background()
//		_, err := test.InsertRolesTx(ctx, s.DB, s.Queries,
//			sqlc.InsertRoleParams{
//				Name: "test_role",
//				Rank: 0,
//			},
//			sqlc.InsertRoleParams{
//				Name: "test_role_2",
//				Rank: 1,
//			},
//		)
//		assert.NoError(s.T(), err)
//
//		role, err := s.RoleService.FindRoles(ctx)
//
//		id := int64(0)
//		for _, r := range role {
//			if r.Name != "test_role" {
//				id = r.ID
//				break
//			}
//		}
//
//		updateArgs := sqlc.UpdateRoleParams{
//			ID:   id,
//			Name: "test_role",
//			Rank: 2,
//		}
//		err = s.RoleService.UpdateRole(ctx, updateArgs)
//		assert.Error(s.T(), err)
//		assert.Error(s.T(), err)
//		assert.Contains(s.T(), err.Error(), "UNIQUE constraint failed: role.name")
//	})
//
//	s.Run("returns error when duplicate rank", func() {
//		ctx := context.Background()
//		_, err := test.InsertRolesTx(ctx, s.DB, s.Queries,
//			sqlc.InsertRoleParams{
//				Name: "test_role",
//				Rank: 0,
//			},
//			sqlc.InsertRoleParams{
//				Name: "test_role_2",
//				Rank: 1,
//			},
//		)
//		assert.NoError(s.T(), err)
//
//		role, err := s.RoleService.FindRoles(ctx)
//
//		id := int64(0)
//		for _, r := range role {
//			if r.Name != "test_role" {
//				id = r.ID
//				break
//			}
//		}
//
//		updateArgs := sqlc.UpdateRoleParams{
//			ID:   id,
//			Name: "test_role_3",
//			Rank: 0,
//		}
//		err = s.RoleService.UpdateRole(ctx, updateArgs)
//		assert.Error(s.T(), err)
//		assert.Error(s.T(), err)
//		assert.Contains(s.T(), err.Error(), "UNIQUE constraint failed: role.rank")
//	})
//}
//
//type RoleServiceTestSuite struct {
//	suite.Suite
//
//	DB          *sql.DB
//	Queries     *sqlc.Queries
//	RoleService RoleService
//}
//
//func (s *RoleServiceTestSuite) SetupSuite() {
//	db, err := test.InitDatabase()
//	assert.NoError(s.T(), err)
//
//	err = test.ResetDatabase(db)
//	assert.NoError(s.T(), err)
//
//	s.DB = db
//	s.Queries = sqlc.New(db)
//	s.RoleService = NewRoleService(s.Queries)
//}
//
//func (s *RoleServiceTestSuite) TearDownTest() {
//	assert.NoError(s.T(), test.ResetDatabase(s.DB))
//}
//
//func (s *RoleServiceTestSuite) TearDownSubTest() {
//	assert.NoError(s.T(), test.ResetDatabase(s.DB))
//}
//
//func TestRoleServiceTestSuite(t *testing.T) {
//	suite.Run(t, new(RoleServiceTestSuite))
//}

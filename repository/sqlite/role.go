package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jinzhu/copier"
	palikka "palikka-go"
	sqlc "palikka-go/repository/sqlite/.sqlc"
)

func NewRole() *palikka.Role {
	return &palikka.Role{
		Permissions: make([]*palikka.Permission, 0),
	}
}

type roleServiceImpl struct {
	DB *DB
}

func NewRoleService(db *DB) palikka.RoleService {
	return &roleServiceImpl{db}
}

func (s *roleServiceImpl) FindRoles(ctx context.Context) ([]*palikka.Role, error) {
	fmt.Print(sqlc.FindRoles)

	rows, err := s.DB.querier.FindRoles(ctx)
	if err != nil {
		return nil, fmt.Errorf("sqlite: query find roles: %w", err)
	}

	// role id - role
	roles := map[int64]*palikka.Role{}
	// role id - permission id - permission
	perms := map[int64]map[int64]*palikka.Permission{}

	for _, row := range rows {

		// check if role already in the map
		if _, ok := roles[row.ID]; !ok {
			role := NewRole()

			if err = copier.Copy(&role, &row); err != nil {
				return nil, fmt.Errorf("sqlite: copy find roles: %w", err)
			}

			roles[row.ID] = role
			perms[row.ID] = make(map[int64]*palikka.Permission)
		}

		// check if row has a permission
		if row.PermissionID != nil {
			// add permission for corresponding role, replacing duplicates
			perms[row.ID][*row.PermissionID] = &palikka.Permission{
				ID:     *row.PermissionID,
				Name:   *row.PermissionName,
				RoleID: row.PermissionRoleID,
			}
		}
	}

	results := make([]*palikka.Role, 0)
	for _, role := range roles {
		if rolePerms, ok := perms[role.ID]; ok {
			for _, perm := range rolePerms {
				role.Permissions = append(role.Permissions, perm)
			}
		}
		results = append(results, role)
	}

	return results, nil
}

func (s *roleServiceImpl) FindRoleByID(ctx context.Context, id int64) (*palikka.Role, error) {
	fmt.Print(sqlc.FindRoleByID)

	rows, err := s.DB.querier.FindRoleByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("sqlite: query find role by id: %w", err)
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("sqlite: query find role by id: %w", sql.ErrNoRows)
	}

	role := NewRole()
	if err = copier.Copy(&role, &rows[0]); err != nil {
		return nil, fmt.Errorf("sqlite: copy find role by id: %w", err)
	}

	for _, row := range rows {
		if row.PermissionID == nil {
			continue
		}
		perm := &palikka.Permission{
			ID:   *row.PermissionID,
			Name: *row.PermissionName,
		}
		role.Permissions = append(role.Permissions, perm)
	}

	return role, nil
}

func (s *roleServiceImpl) InsertRole(ctx context.Context, args *palikka.RoleCreate) (int64, error) {
	fmt.Print(sqlc.InsertRole)

	var params sqlc.InsertRoleParams
	if err := copier.Copy(&params, args); err != nil {
		return 0, fmt.Errorf("sqlite: copy insert role params: %w", err)
	}

	id, err := s.DB.querier.InsertRole(ctx, params)

	if err != nil {
		return id, fmt.Errorf("sqlite: query insert role: %w", err)
	}

	return id, nil
}

func (s *roleServiceImpl) UpdateRole(ctx context.Context, args *palikka.RoleUpdate) error {
	fmt.Print(sqlc.UpdateRole)

	var params sqlc.UpdateRoleParams

	if err := copier.Copy(&params, &args); err != nil {
		return fmt.Errorf("sqlite: copy update role params: %w", err)
	}

	if err := s.DB.querier.UpdateRole(ctx, params); err != nil {
		return fmt.Errorf("sqlite: query update role: %w", err)
	}

	return nil
}

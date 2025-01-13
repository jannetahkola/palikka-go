package sqlite

import (
	"context"
	"fmt"
	"github.com/jinzhu/copier"
	palikka "palikka-go"
	sqlc "palikka-go/repository/sqlite/.sqlc"
)

type permissionServiceImpl struct {
	DB *DB
}

func NewPermissionService(db *DB) palikka.PermissionService {
	return &permissionServiceImpl{
		DB: db,
	}
}

func (s *permissionServiceImpl) FindPermissions(ctx context.Context) ([]*palikka.Permission, error) {
	fmt.Print(sqlc.FindPermissions)

	// todo take params in
	params := sqlc.FindPermissionsParams{
		Offset: 0,
		Limit:  10,
	}

	rows, err := s.DB.querier.FindPermissions(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("sqlite: query find permissions: %w", err)
	}

	perms := make([]*palikka.Permission, 0)
	if err := copier.Copy(&perms, rows); err != nil {
		return nil, fmt.Errorf("sqlite: copy find permissions: %w", err)
	}

	return perms, nil
}

func (s *permissionServiceImpl) FindPermissionsByRoleID(ctx context.Context, roleID int64) ([]*palikka.Permission, error) {
	fmt.Print(sqlc.FindPermissionsByRoleID)

	rows, err := s.DB.querier.FindPermissionsByRoleID(ctx, roleID)
	if err != nil {
		// todo not found error
		return nil, fmt.Errorf("sqlite: query find permissions by role ID: %w", err)
	}

	// todo handle duplicates?
	perms := make([]*palikka.Permission, 0)
	if err := copier.Copy(&perms, rows); err != nil {
		return nil, fmt.Errorf("sqlite: copy find permissions by role ID: %w", err)
	}

	return perms, nil
}

func (s *permissionServiceImpl) InsertPermission(ctx context.Context, args *palikka.PermissionCreate) (int64, error) {
	fmt.Print(sqlc.InsertPermission)

	var params sqlc.InsertPermissionParams
	if err := copier.Copy(&params, args); err != nil {
		return 0, fmt.Errorf("sqlite: copy insert permission params: %w", err)
	}

	id, err := s.DB.querier.InsertPermission(ctx, params)
	if err != nil {
		return 0, fmt.Errorf("sqlite: query insert permission: %w", err)
	}

	return id, nil
}

func (s *permissionServiceImpl) UpdatePermission(ctx context.Context, args *palikka.PermissionUpdate) error {
	fmt.Print(sqlc.UpdatePermission)

	var params sqlc.UpdatePermissionParams
	if err := copier.Copy(&params, args); err != nil {
		return fmt.Errorf("sqlite: copy update permission params: %w", err)
	}

	if err := s.DB.querier.UpdatePermission(ctx, params); err != nil {
		return fmt.Errorf("sqlite: query insert permission: %w", err)
	}

	return nil
}

func (s *permissionServiceImpl) DeletePermissions(ctx context.Context) error {
	fmt.Print(sqlc.DeletePermissions)

	if err := s.DB.querier.DeletePermissions(ctx); err != nil {
		return fmt.Errorf("sqlite: query delete permissions: %w", err)
	}

	return nil
}

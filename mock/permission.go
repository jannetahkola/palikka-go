package mock

import (
	"context"
	palikka "palikka-go"
	"reflect"
)

type PermissionService struct {
	FindPermissionsFn         func(ctx context.Context) ([]*palikka.Permission, error)
	FindPermissionsByRoleIDFn func(ctx context.Context, roleID int64) ([]*palikka.Permission, error)
	InsertPermissionFn        func(ctx context.Context, args *palikka.PermissionCreate) (int64, error)
	UpdatePermissionFn        func(ctx context.Context, args *palikka.PermissionUpdate) error
	DeletePermissionsFn       func(ctx context.Context) error
}

func (s *PermissionService) FindPermissions(ctx context.Context) ([]*palikka.Permission, error) {
	return s.FindPermissionsFn(ctx)
}

func (s *PermissionService) FindPermissionsByRoleID(ctx context.Context, roleID int64) ([]*palikka.Permission, error) {
	return s.FindPermissionsByRoleIDFn(ctx, roleID)
}

func (s *PermissionService) InsertPermission(ctx context.Context, args *palikka.PermissionCreate) (int64, error) {
	return s.InsertPermissionFn(ctx, args)
}

func (s *PermissionService) UpdatePermission(ctx context.Context, args *palikka.PermissionUpdate) error {
	return s.UpdatePermissionFn(ctx, args)
}

func (s *PermissionService) DeletePermissions(ctx context.Context) error {
	return s.DeletePermissionsFn(ctx)
}

func (s *PermissionService) ResetMocks() {
	p := reflect.ValueOf(s).Elem()
	p.Set(reflect.Zero(p.Type()))
}

package mock

import (
	"context"
	palikka "palikka-go"
	"reflect"
)

//var _ palikka.RoleService = (*RoleService)(nil)

type RoleService struct {
	FindRolesFn    func(ctx context.Context) ([]*palikka.Role, error)
	FindRoleByIDFn func(ctx context.Context, id int64) (*palikka.Role, error)
	InsertRoleFn   func(ctx context.Context, args *palikka.RoleCreate) (int64, error)
	UpdateRoleFn   func(ctx context.Context, args *palikka.RoleUpdate) error
}

func (s *RoleService) FindRoles(ctx context.Context) ([]*palikka.Role, error) {
	return s.FindRolesFn(ctx)
}

func (s *RoleService) FindRoleByID(ctx context.Context, id int64) (*palikka.Role, error) {
	return s.FindRoleByIDFn(ctx, id)
}

func (s *RoleService) InsertRole(ctx context.Context, args *palikka.RoleCreate) (int64, error) {
	return s.InsertRoleFn(ctx, args)
}

func (s *RoleService) UpdateRole(ctx context.Context, args *palikka.RoleUpdate) error {
	return s.UpdateRoleFn(ctx, args)
}

func (s *RoleService) ResetMocks() {
	p := reflect.ValueOf(s).Elem()
	p.Set(reflect.Zero(p.Type()))
}

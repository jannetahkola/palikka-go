package mock

import (
	"context"
	palikka "palikka-go"
	"reflect"
)

type UserService struct {
	FindUsersFn          func(ctx context.Context) ([]*palikka.User, error)
	FindUserByIDFn       func(ctx context.Context, id int64) (*palikka.User, error)
	FindUserByUsernameFn func(ctx context.Context, username string) (*palikka.User, error)
	InsertUserFn         func(ctx context.Context, args *palikka.UserCreate) (int64, error)
	UpdateUserFn         func(ctx context.Context, args *palikka.UserUpdate) error
}

func (s *UserService) FindUsers(ctx context.Context) ([]*palikka.User, error) {
	return s.FindUsersFn(ctx)
}

func (s *UserService) FindUserByID(ctx context.Context, id int64) (*palikka.User, error) {
	return s.FindUserByIDFn(ctx, id)
}

func (s *UserService) FindUserByUsername(ctx context.Context, username string) (*palikka.User, error) {
	return s.FindUserByUsernameFn(ctx, username)
}

func (s *UserService) InsertUser(ctx context.Context, args *palikka.UserCreate) (int64, error) {
	return s.InsertUserFn(ctx, args)
}

func (s *UserService) UpdateUser(ctx context.Context, args *palikka.UserUpdate) error {
	return s.UpdateUserFn(ctx, args)
}

func (s *UserService) ResetMocks() {
	p := reflect.ValueOf(s).Elem()
	p.Set(reflect.Zero(p.Type()))
}

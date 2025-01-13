package mock

import (
	"context"
	palikka "palikka-go"
	"reflect"
)

type SessionStore struct {
	LoadFn   func(ctx context.Context, id string) (*palikka.Session, error)
	SaveFn   func(ctx context.Context, session *palikka.Session) error
	DeleteFn func(ctx context.Context, id string) error
	ClearFn  func(ctx context.Context)
}

func (s *SessionStore) Load(ctx context.Context, id string) (*palikka.Session, error) {
	return s.LoadFn(ctx, id)
}

func (s *SessionStore) Save(ctx context.Context, session *palikka.Session) error {
	return s.SaveFn(ctx, session)
}

func (s *SessionStore) Delete(ctx context.Context, id string) error {
	return s.DeleteFn(ctx, id)
}
func (s *SessionStore) Clear(ctx context.Context) {
	s.ClearFn(ctx)
}

func (s *SessionStore) ResetMocks() {
	p := reflect.ValueOf(s).Elem()
	p.Set(reflect.Zero(p.Type()))
}

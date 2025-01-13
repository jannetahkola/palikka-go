package session

import (
	"github.com/stretchr/testify/mock"
	"net/http"
)

type StoreMock struct {
	mock.Mock
}

func (m *StoreMock) Load(r *http.Request) (*Session, error) {
	args := m.Called(r)
	s, _ := args.Get(0).(*Session)
	err := args.Error(1)
	return s, err
}

func (m *StoreMock) Save(w http.ResponseWriter, r *http.Request, session *Session) error {
	args := m.Called(w, r, session)
	return args.Error(0)
}

func (m *StoreMock) Delete(w http.ResponseWriter, r *http.Request) error {
	args := m.Called(w, r)
	return args.Error(0)
}

func (m *StoreMock) Clear() {
	m.Called()
}

package palikka

import (
	"context"
	"time"
)

const (
	PermUsersCreate  = "users_create"
	PermUsersRead    = "users_read"
	PermsUsersUpdate = "users_update"
	PermsUsersDelete = "users_delete"
	PermRolesCreate  = "roles_create"
	PermRolesRead    = "roles_read"
	PermRolesUpdate  = "roles_update"
	PermRolesDelete  = "roles_delete"
)

type Permission struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	RoleID *int64 `json:"role_id"`
}

type Role struct {
	ID          int64         `json:"id"`
	Name        string        `json:"name"`
	Rank        int64         `json:"rank"`
	Permissions []*Permission `json:"permissions"`
}

type User struct {
	ID        int64   `json:"id"`
	Username  string  `json:"username"`
	Password  *string `json:"-"`
	CreatedAt string  `json:"created_at"`
	Active    bool    `json:"active"`
	RoleID    *int64  `json:"role_id,omitempty"`
	Role      *Role   `json:"role,omitempty"` // todo test omitempty
}

// PermissionCreate is a struct for creating a Permission.
type PermissionCreate struct {
	Name   *string `json:"name"`
	RoleID *int64  `json:"role_id"`
}

// PermissionUpdate is a struct for updating a Permission.
type PermissionUpdate struct {
	ID     *int64 `json:"-"`
	RoleID *int64 `json:"role_id"`
}

// RoleCreate is a struct for creating a Role.
type RoleCreate struct {
	Name *string `json:"name"`
	Rank *int64  `json:"rank"`
}

// RoleUpdate is a struct for updating a Role.
type RoleUpdate struct {
	ID   *int64  `json:"-"`
	Name *string `json:"name"`
	Rank *int64  `json:"rank"`
}

// UserCreate is a struct for creating a User.
type UserCreate struct {
	Username *string `json:"username"`
	Password *string `json:"password"`
	Active   *bool   `json:"active"`
	RoleID   *int64  `json:"role_id"`
}

// UserUpdate is a struct for updating a User.
type UserUpdate struct {
	ID       *int64  `json:"-"`
	Username *string `json:"username"`
	Password *string `json:"password"`
	Active   *bool   `json:"active"`
	RoleID   *int64  `json:"role_id"`
}

type PermissionService interface {
	FindPermissions(ctx context.Context) ([]*Permission, error)
	FindPermissionsByRoleID(ctx context.Context, roleID int64) ([]*Permission, error)
	InsertPermission(ctx context.Context, args *PermissionCreate) (int64, error)
	UpdatePermission(ctx context.Context, args *PermissionUpdate) error
	DeletePermissions(ctx context.Context) error
}

type RoleService interface {
	FindRoles(ctx context.Context) ([]*Role, error)
	FindRoleByID(ctx context.Context, id int64) (*Role, error)
	InsertRole(ctx context.Context, args *RoleCreate) (int64, error)
	UpdateRole(ctx context.Context, args *RoleUpdate) error
}

type UserService interface {
	FindUsers(ctx context.Context) ([]*User, error)
	FindUserByID(ctx context.Context, id int64) (*User, error)
	FindUserByUsername(ctx context.Context, username string) (*User, error)
	InsertUser(ctx context.Context, args *UserCreate) (int64, error)
	UpdateUser(ctx context.Context, args *UserUpdate) error
}

type Subscription struct {
	Messages chan []byte
}

type Session struct {
	ID      string
	UserID  int64
	MaxAge  time.Duration
	Expires time.Time // todo should be a pointer?
	// Values holds additional values for the session. May be nil.
	Values       map[string]interface{}
	Subscription *Subscription
}

// SessionStore is an in-memory session store.
type SessionStore interface {
	Load(ctx context.Context, id string) (*Session, error)
	Save(ctx context.Context, session *Session) error
	Delete(ctx context.Context, id string) error
	Clear(ctx context.Context)
	OnDelete() chan string
}

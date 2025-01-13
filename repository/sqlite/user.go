package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
	palikka "palikka-go"
	sqlc "palikka-go/repository/sqlite/.sqlc"
	"palikka-go/repository/sqlite/types"
)

var userCopyOptions = copier.Option{
	Converters: []copier.TypeConverter{
		{
			SrcType: types.Bool{},
			DstType: copier.Bool,
			Fn: func(src interface{}) (dst interface{}, err error) {
				// todo check below error
				// Cannot assign bool to err (type error) in multiple assignment
				// Type does not implement error as some methods are missing: Error() string
				b, ok := src.(types.Bool)
				if !ok {
					return false, err
				}
				v, _ := b.Value()
				return v, nil
			},
		},
		{
			SrcType: types.BoolPtr,
			DstType: types.Bool{},
			Fn: func(src interface{}) (dst interface{}, err error) {
				b, ok := src.(*bool)
				if !ok {
					return types.BoolFrom(false), err
				}

				var v types.Bool
				if b == nil {
					v = types.BoolFrom(false)
				} else {
					v = types.BoolFrom(*b)
				}

				return v, nil
			},
		},
	},
}

func NewUser() *palikka.User {
	return &palikka.User{}
}

type userServiceImpl struct {
	DB *DB
}

func NewUserService(db *DB) palikka.UserService {
	return &userServiceImpl{db}
}

func (s *userServiceImpl) FindUsers(ctx context.Context) ([]*palikka.User, error) {
	fmt.Print(sqlc.FindUsers)

	// todo take params in
	params := sqlc.FindUsersParams{
		Offset: 0,
		Limit:  10,
	}

	rows, err := s.DB.querier.FindUsers(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("sqlite: query find users: %w", err)
	}

	//todo rest of mappings, can copy list directly without looping here
	users := make([]*palikka.User, 0)
	for _, row := range rows {
		var user palikka.User
		if err := copier.CopyWithOption(&user, &row, userCopyOptions); err != nil {
			return nil, fmt.Errorf("sqlite: copy find users: %w", err)
		}
		users = append(users, &user)
	}

	return users, nil
}

func (s *userServiceImpl) FindUserByID(ctx context.Context, id int64) (*palikka.User, error) {
	fmt.Print(sqlc.FindUserByID)

	row, err := s.DB.querier.FindUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("sqlite: query find user by id: %s", palikka.ErrNotFound)
		}
		return nil, fmt.Errorf("sqlite: query find user by id: %w", err)
	}

	user := NewUser()
	if err = copier.CopyWithOption(&user, &row, userCopyOptions); err != nil {
		return nil, fmt.Errorf("sqlite: copy find user by id: %w", err)
	}

	return user, nil
}

func (s *userServiceImpl) FindUserByUsername(ctx context.Context, username string) (*palikka.User, error) {
	fmt.Print(sqlc.FindUserByUsername)

	row, err := s.DB.querier.FindUserByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("sqlite: query find user by username: %w", err)
	}

	user := NewUser()
	if err = copier.CopyWithOption(&user, &row, userCopyOptions); err != nil {
		return nil, fmt.Errorf("sqlite: copy find user by username: %w", err)
	}

	return user, nil
}

func (s *userServiceImpl) InsertUser(ctx context.Context, args *palikka.UserCreate) (int64, error) {
	fmt.Print(sqlc.InsertUser)

	var params sqlc.InsertUserParams
	if err := copier.CopyWithOption(&params, args, userCopyOptions); err != nil {
		return 0, fmt.Errorf("sqlite: copy insert user params: %w", err)
	}

	id, err := s.DB.querier.InsertUser(ctx, params)
	if err != nil {
		return 0, fmt.Errorf("sqlite: query insert user: %w", err)
	}

	return id, err
}

func (s *userServiceImpl) UpdateUser(ctx context.Context, args *palikka.UserUpdate) error {
	fmt.Print(sqlc.UpdateUser)

	var params sqlc.UpdateUserParams
	if err := copier.CopyWithOption(&params, args, userCopyOptions); err != nil {
		return fmt.Errorf("sqlite: copy update user params: %w", err)
	}

	if err := s.DB.querier.UpdateUser(ctx, params); err != nil {
		return fmt.Errorf("sqlite: query update user: %w", err)
	}

	return nil
}

// todo format errors for other services in sqlite package

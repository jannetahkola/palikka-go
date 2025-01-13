package sqlite

//
//import (
//	"context"
//	"fmt"
//	"github.com/jinzhu/copier"
//	_ "modernc.org/sqlite"
//	sqlc "palikka-go/internal/data/.sqlc"
//	"palikka-go/internal/data/types"
//)
//
//var userCopyOptions = copier.Option{
//	Converters: []copier.TypeConverter{
//		{
//			SrcType: types.Bool{},
//			DstType: copier.Bool,
//			Fn: func(src interface{}) (dst interface{}, err error) {
//				// todo check below error
//				// Cannot assign bool to err (type error) in multiple assignment
//				// Type does not implement error as some methods are missing: Error() string
//				b, _ := src.(types.Bool)
//				v, _ := b.Value()
//				return v, nil
//			},
//		},
//	},
//}
//
//type querierImpl struct {
//	querier *sqlc.Queries
//}
//
//func NewQuerier(querier *sqlc.Queries) sqlc.querier {
//	return &querierImpl{querier: querier}
//}
//
//// todo interface + factory?
//
//// **********
//// Users
//// **********
//
//func (r *querierImpl) FindUsers(ctx context.Context, args sqlc.FindUsersParams) ([]sqlc.User, error) {
//	fmt.Print(sqlc.FindUsers)
//
//	results, err := r.querier.FindUsers(ctx, args)
//	if err != nil {
//		return nil, err
//	}
//
//	return results, nil
//}
//
//func (r *querierImpl) FindUserByID(ctx context.Context, id int64) (sqlc.User, error) {
//	fmt.Print(sqlc.FindUserByID)
//
//	result := sqlc.User{}
//	result, err := r.querier.FindUserByID(ctx, id)
//	if err != nil {
//		return result, err
//	}
//
//	return result, nil
//}
//
//func (r *querierImpl) FindUserByUsername(ctx context.Context, username string) (sqlc.User, error) {
//	fmt.Print(sqlc.FindUserByUsername)
//
//	result := sqlc.User{}
//	result, err := r.querier.FindUserByUsername(ctx, username)
//	if err != nil {
//		return result, err
//	}
//
//	return result, nil
//}
//
//func (r *querierImpl) UpdateUser(ctx context.Context, args sqlc.UpdateUserParams) error {
//	fmt.Print(sqlc.UpdateUser)
//
//	if err := r.querier.UpdateUser(ctx, args); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//// **********
//// Roles
//// **********
//
//func (r *querierImpl) FindRoles(ctx context.Context, args sqlc.FindRolesParams) ([]sqlc.Role, error) {
//	fmt.Print(sqlc.FindRoles)
//
//	results, err := r.querier.FindRoles(ctx, args)
//	if err != nil {
//		return nil, err
//	}
//
//	return results, nil
//}
//
//func (r *querierImpl) FindRoleByID(ctx context.Context, id int64) (sqlc.Role, error) {
//	fmt.Print(sqlc.FindRoleByID)
//
//	result := sqlc.Role{}
//	result, err := r.querier.FindRoleByID(ctx, id)
//	if err != nil {
//		return result, err
//	}
//
//	return result, nil
//}
//
//// **********
//// Permissions
//// **********
//
//func (r *querierImpl) FindPermissions(ctx context.Context, args sqlc.FindPermissionsParams) ([]sqlc.Permission, error) {
//	fmt.Print(sqlc.FindPermissions)
//
//	results, err := r.querier.FindPermissions(ctx, args)
//	if err != nil {
//		return nil, err
//	}
//
//	return results, nil
//}
//
//func (r *querierImpl) FindPermissionsByRoleID(ctx context.Context, roleID int64) ([]sqlc.Permission, error) {
//	fmt.Print(sqlc.FindPermissionsByRoleID)
//
//	results, err := r.querier.FindPermissionsByRoleID(ctx, roleID)
//	if err != nil {
//		return nil, err
//	}
//
//	return results, nil
//}

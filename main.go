package palikka

import (
	_ "embed"
	_ "modernc.org/sqlite"
)

////go:embed repository/data/schema.sql
//var CreateSchemaFS embed.FS
//
////go:embed repository/data/schema_seed.sql
//var seed embed.FS
//
////go:embed repository/data/schema_reset.sql
//var reset embed.FS

func runSqlc() {
	const (
		driverName        = "sqlite"
		connectionString2 = "file:palikka.sqlite3"
		connectionString  = "file:palikka.sqlite3?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)&_pragma=synchronous(2)"
	)

	//
	//// https://pkg.go.dev/modernc.org/sqlite?utm_source=godoc#Driver.Open
	//db, err := sql.Open(driverName, connectionString2)
	//if err != nil {
	//	panic(err)
	//}
	//
	//// reset
	//if _, err = db.ExecContext(context.Background(), reset); err != nil {
	//	panic(err)
	//}
	//
	//// create tables
	//if _, err = db.ExecContext(context.Background(), ddl); err != nil {
	//	panic(err)
	//}
	//
	//// seed data
	//if _, err = db.ExecContext(context.Background(), seed); err != nil {
	//	panic(err)
	//}

	//queries := sqlc.New(db)
	//
	//repo := data.NewQuerier(queries)
	//users, _ := repo.FindUsers(0)
	//
	//jsonText, _ := json.MarshalIndent(users, "", "\t")
	//fmt.Println(string(jsonText))

	// -----

	//roles, err := queries.
	//	FindRoleByIdWithPermissions(context.Background(), 2)
	//if err != nil {
	//	panic(err)
	//}
	//
	//jsonText, _ := json.MarshalIndent(roles, "", "\t")
	//fmt.Println(string(jsonText))
	//
	//var perms []internal.Permission
	//for _, r := range roles {
	//	var perm internal.Permission
	//	if err = copier.Copy(&perm, &r.Permission); err != nil {
	//		panic(err)
	//	}
	//	perms = append(perms, perm)
	//}
	//
	//var role internal.Role
	//if err = copier.Copy(&role, &roles[0]); err != nil {
	//	panic(err)
	//}
	//
	//role.Permissions = perms
	//
	//jsonText, _ = json.MarshalIndent(role, "", "\t")
	//fmt.Println(string(jsonText))

	// -----

	//fmt.Print(data.FindUserByUsernameWithCredentials)
	//users, err := queries.FindUserByUsername(context.Background(), "mock-user")
	//if err != nil {
	//	panic(err)
	//}
	//
	//jsonText, _ := json.MarshalIndent(users, "", "\t")
	//fmt.Println(string(jsonText))
	//
	//// todo could ditch the custom Bool except can't because don't know which ints represent bool
	//copierOpts := copier.Option{
	//	Converters: []copier.TypeConverter{
	//		{
	//			SrcType: types.Bool{},
	//			DstType: copier.Bool,
	//			Fn: func(src interface{}) (dst interface{}, err error) {
	//				// todo Cannot assign bool to err (type error) in multiple assignment Type does not implement error as some methods are missing: Error() string
	//				b, _ := src.(types.Bool)
	//				v, _ := b.Value()
	//				return v, nil
	//			},
	//		},
	//	},
	//}
	//
	//var u internal.UserModel
	//if err = copier.CopyWithOption(&u, &users, copierOpts); err != nil {
	//	panic(err)
	//}
	//
	//jsonText, _ = json.MarshalIndent(u, "", "\t")
	//fmt.Println(string(jsonText))
}

func main() {
	runSqlc()
}

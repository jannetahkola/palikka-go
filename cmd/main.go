package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"palikka-go/http"
	"palikka-go/http/cookie"
	"palikka-go/repository/mem"
	"palikka-go/repository/sqlite"
	"time"
)

func main() {
	// open database connection
	db := sqlite.NewDB("file:palikka.sqlite3")
	if err := db.Open(); err != nil {
		panic(err)
	}

	// create session store
	sessionStore := mem.NewSessionStore(time.Hour)

	// create services
	cookieServiceOpts := &cookie.Options{
		CookieName: "palikka_session",
		Secret:     []byte("changeme"),
		Secure:     false,
	}
	cookieService := cookie.NewService(cookieServiceOpts)
	permissionService := sqlite.NewPermissionService(db)
	roleService := sqlite.NewRoleService(db)
	userService := sqlite.NewUserService(db)

	// create server
	serverOpts := &http.ServerOpts{
		WebAppURI: "http://localhost:4200",
		Cors: &http.CorsOpts{
			AllowedOrigins: []string{"*"},
		},
	}
	s := http.NewServer(serverOpts)

	// load server config
	s.Addr = ":8080"

	// attach dependencies to server
	s.CookieService = cookieService
	s.SessionStore = sessionStore
	s.PermissionService = permissionService
	s.RoleService = roleService
	s.UserService = userService

	// start server
	if err := s.Start(); err != nil {
		_ = s.Stop()
		fmt.Println(err)
		os.Exit(1)
	}

	// wait for ctrl+c
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		cancel()
	}()
	<-ctx.Done()

	// stop the server
	if err := s.Stop(); err != nil {
		fmt.Println(err)
	}

	fmt.Println("server close")
}

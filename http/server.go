package http

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	palikka "palikka-go"
	"palikka-go/http/assets"
	palikkacookie "palikka-go/http/cookie"
	"slices"
	"strings"
	"time"
)

type CorsOpts struct {
	AllowedOrigins []string
}

type ServerOpts struct {
	Cors      *CorsOpts
	WebAppURI string
}

func (o *ServerOpts) validate() error {
	if o.WebAppURI == "" {
		return errors.New("web app uri is required")
	}
	if o.Cors == nil {
		return errors.New("cors is required")
	}
	if o.Cors.AllowedOrigins == nil {
		return errors.New("cors allowed origins is required")
	}
	return nil
}

type Server struct {
	ln     net.Listener
	server *http.Server
	router *http.ServeMux

	opts *ServerOpts

	Domain string
	Addr   string

	EventService      *eventService
	CookieService     palikkacookie.Service
	SessionStore      palikka.SessionStore
	PermissionService palikka.PermissionService
	RoleService       palikka.RoleService
	UserService       palikka.UserService
}

func NewServer(opts *ServerOpts) *Server {
	s := &Server{
		server: &http.Server{},
		router: http.NewServeMux(),
		opts:   opts,
	}

	s.server.Handler = http.HandlerFunc(s.serveHTTP)

	// html
	// todo tests
	loginFS, _ := fs.Sub(assets.LoginFS, "login")
	s.router.Handle("/login/", s.serveLogin(http.StripPrefix("/login", http.FileServer(http.FS(loginFS)))))

	// authN
	s.router.HandleFunc("POST /auth/login", s.handleLogin)
	s.router.HandleFunc("POST /auth/logout", s.requireAuth(s.handleLogout, authorizeAll()))

	// authZ
	s.router.HandleFunc("GET /authorizations", s.requireAuth(s.handleGetAuthorizations, authorizeAll()))

	// session
	s.router.HandleFunc("GET /sessions/session", s.requireAuth(s.handleGetCurrentSession, authorizeAll()))

	// users
	s.router.HandleFunc("GET /users", s.requireAuth(s.handleGetUsers, authorizeAllOf(palikka.PermUsersRead)))
	s.router.HandleFunc("GET /users/{id}", s.requireAuth(s.handleGetUserByID, authorizeAllOf(palikka.PermUsersRead)))
	s.router.HandleFunc("POST /users", s.requireAuth(s.handleCreateUser, authorizeAllOf(palikka.PermUsersCreate)))

	// roles
	// todo in frontend this is needed for both users and roles currently; figure it out
	s.router.HandleFunc("GET /roles", s.requireAuth(s.handleGetRoles, authorizeAllOf(palikka.PermRolesRead)))

	// websocket
	s.router.HandleFunc("GET /ws", s.handleUpgradeConn)

	// not found
	s.router.HandleFunc("/", s.handleNotFound())

	return s
}

func (s *Server) Scheme() string {
	return "http"
}

func (s *Server) Port() int {
	return s.ln.Addr().(*net.TCPAddr).Port
}

func (s *Server) URL() string {

	domain := "localhost"
	if s.Domain != "" {
		domain = s.Domain
	}

	return fmt.Sprintf("%s://%s:%d", s.Scheme(), domain, s.Port())
}

func (s *Server) Start() error {
	var (
		err error
	)

	// todo check dependencies here?

	if err = s.opts.validate(); err != nil {
		return err
	}

	if s.ln, err = net.Listen("tcp", s.Addr); err != nil {
		return err
	}

	go s.server.Serve(s.ln)

	fmt.Println("[http] server listening at " + s.URL())

	return nil
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	return s.server.Shutdown(ctx)
}

func (s *Server) serveHTTP(w http.ResponseWriter, r *http.Request) {
	w = &responseWrapper{ResponseWriter: w}

	// set request context
	r = copyContextWithLogger(r)
	r = copyContextWithRequestAttributes(r)

	fmt.Println("request to:", r.URL)

	// CORS
	origin := r.Header.Get("Origin")
	if origin == "" && slices.Contains(s.opts.Cors.AllowedOrigins, "*") {
		origin = "http://localhost:8080/" // todo check can we set it dynamically?
	}
	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	s.router.ServeHTTP(w, r)
}

func (s *Server) requireAuth(next http.HandlerFunc, authorizePerms []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logInfo(r, "require auth")

		var err error
		if r, err = s.loadSession(r); err != nil {
			cause := fmt.Errorf("requireAuth: load session: %w", err)
			s.respondError(w, r, palikka.Errorf2(cause, palikka.ErrUnauthorized, ErrMsgUnauthorized))
			return
		}

		if err = s.authorize(r, authorizePerms); err != nil {
			cause := fmt.Errorf("requireAuth: authorize: %w", err)
			s.respondError(w, r, palikka.Errorf2(cause, palikka.ErrForbidden, ErrMsgForbidden))
			return
		}

		next(w, r)
	}
}

func (s *Server) handleNotFound() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("handleNotFound")

		if _, err := s.loadSession(r); err != nil {
			// return 401 like any other protected endpoint to obfuscate the API
			s.respondError(w, r, palikka.Errorf(palikka.ErrUnauthorized, ErrMsgUnauthorized))
			return
		}

		s.respondError(w, r, palikka.Errorf(palikka.ErrNotFound, ErrMsgNotFound))
		return
	}
}

func (s *Server) loadSession(r *http.Request) (*http.Request, error) {
	var (
		err       error
		cookie    *http.Cookie
		session   *palikka.Session
		sessionID string
		user      *palikka.User
	)

	// parse & verify cookie
	if cookie, err = s.CookieService.Parse(r); err != nil {
		return r, fmt.Errorf("CookieService.Parse: %w", err)
	}
	if sessionID, err = s.CookieService.Verify(cookie); err != nil {
		return r, fmt.Errorf("CookieService.Verify: %w", err)
	}

	// load session
	// todo nuke the cookie if can't load session (or user below)
	if session, err = s.SessionStore.Load(r.Context(), sessionID); err != nil {
		return r, fmt.Errorf("SessionStore.Load: %w", err)
	}
	r = r.WithContext(palikka.CopyContextWithSession(r.Context(), session))

	// todo do NOT include password in context
	// load user
	if user, err = s.UserService.FindUserByID(r.Context(), session.UserID); err != nil {
		return r, fmt.Errorf("UserService.FindUserByID: %w", err)
	}
	r = r.WithContext(palikka.CopyContextWithUser(r.Context(), user))

	logInfo(r, "authenticate user")

	return r, nil
}

func (s *Server) authorize(r *http.Request, requiredPerms []string) error {
	if len(requiredPerms) > 0 {
		u := palikka.UserFromContext(r.Context())
		if u == nil {
			return errors.New("no user in context")
		}

		if u.RoleID == nil {
			return errors.New("user has no role")
		}

		userPerms, err := s.PermissionService.FindPermissionsByRoleID(r.Context(), *u.RoleID)
		if err != nil || len(userPerms) == 0 {
			return errors.New("user has no permissions")
		}

		var userPermStrings []string
		for _, perm := range userPerms { // todo make a separate query instead of loop
			userPermStrings = append(userPermStrings, perm.Name)
		}

		// todo should we add the roles/perms to context here?

		authorize := false
		for _, requiredPerm := range requiredPerms {
			authorize = slices.Contains(userPermStrings, requiredPerm)
		}

		if !authorize {
			return errors.New(
				fmt.Sprintf("user has no required permission(s), user=[%s], required=[%s]",
					strings.Join(userPermStrings, ", "), strings.Join(requiredPerms, ", ")))
		}
	}

	logInfo(r, "authorize user")

	return nil
}

func (s *Server) serveLogin(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var err error
		if r, err = s.loadSession(r); err == nil {
			http.Redirect(w, r, s.opts.WebAppURI, http.StatusFound)
			return
		}

		// disable caching
		w.Header().Add("Cache-Control", "no-cache")

		next.ServeHTTP(w, r)
	}
}

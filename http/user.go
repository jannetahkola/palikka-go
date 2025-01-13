package http

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	palikka "palikka-go"
	"palikka-go/http/input"
	"strconv"
	"strings"
)

func (s *Server) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	var (
		err       error
		users     []*palikka.User
		usersJSON []byte
	)

	// find users
	if users, err = s.UserService.FindUsers(r.Context()); err != nil {
		cause := fmt.Errorf("handleGetUsers: find users: %w", err)
		s.respondError(w, r, palikka.Errorf2(cause, palikka.ErrInternal, ErrMsgInternal))
		return
	}

	// find roles for users
	for _, user := range users {
		if user.RoleID != nil {
			role, err := s.RoleService.FindRoleByID(r.Context(), *user.RoleID)
			if err != nil {
				cause := fmt.Errorf("handleGetUsers: find role by id: %w", err)
				s.respondError(w, r, palikka.Errorf2(cause, palikka.ErrInternal, ErrMsgInternal))
				return
			}
			user.Role = role
		}
	}

	// marshal response
	if usersJSON, err = json.Marshal(users); err != nil {
		cause := fmt.Errorf("handleGetUsers: marshal json: %w", err)
		s.respondError(w, r, palikka.Errorf2(cause, palikka.ErrInternal, ErrMsgInternal))
		return
	}

	w.Header().Set(HeaderContentType, MediaTypeJSON)
	_, _ = w.Write(usersJSON)
}

func (s *Server) handleGetUserByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		cause := fmt.Errorf("handleGetUserByID: convert path value to int: %w", err)
		s.respondError(w, r, palikka.Errorf2(cause, palikka.ErrInvalid, "invalid user id, must be an integer"))
		return
	}

	u, err := s.UserService.FindUserByID(r.Context(), int64(id))
	if err != nil {
		cause := fmt.Errorf("handleGetUserByID: find user by id: %w", err)
		if strings.Contains(err.Error(), palikka.ErrNotFound) {
			s.respondError(w, r, palikka.Errorf2(cause, palikka.ErrNotFound, "user not found"))
		} else {
			s.respondError(w, r, palikka.Errorf2(cause, palikka.ErrInternal, ErrMsgInternal))
		}
		return
	}

	// todo add role?

	userJSON, err := json.Marshal(u)
	if err != nil {
		cause := fmt.Errorf("handleGetUserByID: marshal json: %w", err)
		s.respondError(w, r, palikka.Errorf2(cause, palikka.ErrInternal, ErrMsgInternal))
		return
	}

	w.Header().Set(HeaderContentType, MediaTypeJSON)
	_, _ = w.Write(userJSON)
}

func (s *Server) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var args palikka.UserCreate
	if err := json.NewDecoder(r.Body).Decode(&args); err != nil {
		cause := fmt.Errorf("handleCreateUser: decode request body: %w", err)
		s.respondError(w, r, palikka.Errorf2(cause, palikka.ErrInvalid, ErrMsgInvalidRequest))
		return
	}

	// sanitize input
	args.Username = input.Sanitize(args.Username)

	// validate input
	if err := validateUserCreate(&args); err != nil {
		cause := fmt.Errorf("handleCreateUser: validate request: %w", err)
		s.respondError(w, r, palikka.Errorf2(cause, palikka.ErrInvalid, err.Error())) // can use .Error()
		return
	}

	// hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(*args.Password), bcrypt.DefaultCost)
	if err != nil {
		cause := fmt.Errorf("handleCreateUser: generate bcrypt hash: %w", err)
		s.respondError(w, r, palikka.Errorf2(cause, palikka.ErrInternal, ErrMsgInternal))
		return
	}

	hashString := string(hash)
	args.Password = &hashString

	// insert user
	id, err := s.UserService.InsertUser(r.Context(), &args)
	if err != nil {
		cause := fmt.Errorf("handleCreateUser: insert user: %w", err)
		s.respondError(w, r, palikka.Errorf2(cause, palikka.ErrInternal, ErrMsgInternal))
		return
	}

	logInfo(r, fmt.Sprintf("user create, id=%d", id))

	w.WriteHeader(http.StatusCreated)
	w.Header().Set(HeaderLocation, fmt.Sprintf("%s/users/%d", s.URL(), id))
}

func validateUserCreate(args *palikka.UserCreate) error {

	// todo validate roleID

	if err := input.ValidateUsername(args.Username); err != nil {
		return err
	}

	if err := input.ValidatePassword(args.Password); err != nil {
		return err
	}

	return nil
}

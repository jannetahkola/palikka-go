package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	palikka "palikka-go"
	"palikka-go/http/input"
)

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	var request *LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		cause := fmt.Errorf("handleLogin: decode request body: %w", err)
		s.respondError(w, r, palikka.Errorf2(cause, palikka.ErrInvalid, "invalid request body"))
		return
	}

	// NOTE: open api so do not reveal any details in error responses!

	// sanitize input (except password)
	request.Username = input.Sanitize(request.Username)

	// validate input
	if err := validateLoginRequest(request); err != nil {
		cause := fmt.Errorf("handleLogin: validate request: %w", err)
		s.respondError(w, r, palikka.Errorf2(cause, palikka.ErrInvalid, ErrMsgInvalidRequest))
		return
	}

	// load user
	u, err := s.UserService.FindUserByUsername(r.Context(), *request.Username)
	if err != nil {
		cause := fmt.Errorf("handleLogin: find user by username: %w", err)
		s.respondError(w, r, palikka.Errorf2(cause, palikka.ErrInvalid, ErrMsgInvalidRequest))
		return
	}

	r = r.WithContext(palikka.CopyContextWithUser(r.Context(), u))

	if !u.Active {
		// todo user not active
	}

	// validate password
	if err := bcrypt.CompareHashAndPassword([]byte(*u.Password), []byte(*request.Password)); err != nil {
		cause := fmt.Errorf("handleLogin: compare password: %w", err)
		s.respondError(w, r, palikka.Errorf2(cause, palikka.ErrInvalid, ErrMsgInvalidRequest))
		return
	}

	// create session
	session := &palikka.Session{UserID: u.ID, Values: make(map[string]interface{})}
	session.Values["User-Agent"] = r.UserAgent()
	session.Values["Host"] = r.Host

	if err := s.SessionStore.Save(r.Context(), session); err != nil {
		// todo do not return failed to saved session here!
		s.respondError(w, r, palikka.Errorf2(err, palikka.ErrInternal, "failed to save session"))
		return
	}

	r = r.WithContext(palikka.CopyContextWithSession(r.Context(), session))

	// create session cookie
	cookie := s.CookieService.New(session)
	http.SetCookie(w, cookie)

	logInfo(r, "user log in")

	var response = &LoginResponse{RedirectURI: s.opts.WebAppURI}
	var responseJSON []byte

	// marshal response
	if responseJSON, err = json.Marshal(response); err != nil {
		cause := fmt.Errorf("handleLogin: marshal json: %w", err)
		s.respondError(w, r, palikka.Errorf2(cause, palikka.ErrInternal, ErrMsgInternal))
		return
	}

	w.Header().Set(HeaderContentType, MediaTypeJSON)
	_, _ = w.Write(responseJSON)
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	session := palikka.SessionFromContext(r.Context())

	// delete session
	if err := s.SessionStore.Delete(r.Context(), session.ID); err != nil {
		cause := fmt.Errorf("handleLogout: delete session: %w", err)
		logError(r, palikka.Errorf2(cause, palikka.ErrInternal, "error logging out").Error())
		// todo test this error case
	}

	// expire session cookie
	expiredCookie := s.CookieService.NewExpired()
	http.SetCookie(w, expiredCookie)

	// no redirect to let web app handle it instead of browser
}

// todo do not export
type LoginRequest struct {
	Username *string `json:"username"`
	Password *string `json:"password"`
}

type LoginResponse struct {
	RedirectURI string `json:"redirect_uri"`
}

// validateLoginRequest validates request fields and returns client-friendly errors.
func validateLoginRequest(lr *LoginRequest) error {
	if lr == nil {
		return errors.New("invalid request")
	}

	if err := input.ValidateUsername(lr.Username); err != nil {
		return err
	}

	if err := input.ValidatePassword(lr.Password); err != nil {
		return err
	}

	return nil
}

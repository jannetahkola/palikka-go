package http

import (
	"encoding/json"
	"net/http"
	palikka "palikka-go"
	"time"
)

func (s *Server) handleGetCurrentSession(w http.ResponseWriter, r *http.Request) {
	userCtx := palikka.UserFromContext(r.Context())
	sessionCtx := palikka.SessionFromContext(r.Context())

	user := struct {
		ID          int64    `json:"id"`
		Username    string   `json:"username"`
		Permissions []string `json:"permissions"`
	}{
		ID:          userCtx.ID,
		Username:    userCtx.Username,
		Permissions: []string{}, // todo not needed, we have authz and roles?
	}

	if userCtx.RoleID != nil {
		userPerms, err := s.PermissionService.FindPermissionsByRoleID(r.Context(), *userCtx.RoleID)
		if err == nil {
			for _, perm := range userPerms { // todo make a separate query instead of loop
				user.Permissions = append(user.Permissions, perm.Name)
			}
		}
	}

	session := struct {
		Expires time.Time   `json:"expires"`
		User    interface{} `json:"user"`
	}{
		Expires: sessionCtx.Expires.UTC(),
		User:    user,
	}

	responseJson, err := json.Marshal(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set(HeaderContentType, MediaTypeJSON)
	_, _ = w.Write(responseJson)
}

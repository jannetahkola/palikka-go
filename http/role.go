package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	palikka "palikka-go"
)

func (s *Server) handleGetRoles(w http.ResponseWriter, r *http.Request) {
	var (
		roles     []*palikka.Role
		rolesJSON []byte
		err       error
	)

	if roles, err = s.RoleService.FindRoles(r.Context()); err != nil {
		cause := fmt.Errorf("handleGetRoles: find roles: %w", err)
		s.respondError(w, r, palikka.Errorf2(cause, palikka.ErrInternal, ErrMsgInternal))
		return
	}

	// marshal response
	if rolesJSON, err = json.Marshal(roles); err != nil {
		cause := fmt.Errorf("handleGetRoles: marshal json: %w", err)
		s.respondError(w, r, palikka.Errorf2(cause, palikka.ErrInternal, ErrMsgInternal))
		return
	}

	w.Header().Set(HeaderContentType, MediaTypeJSON)
	_, _ = w.Write(rolesJSON)
}

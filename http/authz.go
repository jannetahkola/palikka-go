package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	palikka "palikka-go"
	"slices"
)

type CRUDs struct {
	Create bool `json:"create"`
	Read   bool `json:"read"`
	Update bool `json:"update"`
	Delete bool `json:"delete"`
}

type FeatureAuthorizations struct {
	Users CRUDs `json:"users"`
	Roles CRUDs `json:"roles"`
}

type AuthorizationsResponse struct {
	Features FeatureAuthorizations `json:"features"`
}

func (s *Server) handleGetAuthorizations(w http.ResponseWriter, r *http.Request) {
	var (
		authzJSON []byte
		err       error
	)

	u := palikka.UserFromContext(r.Context())

	var perms []string
	if u.RoleID != nil {
		p, err := s.PermissionService.FindPermissionsByRoleID(r.Context(), *u.RoleID)
		if err != nil {
			cause := fmt.Errorf("handleGetAuthorizations: find permissions by role ID: %w", err)
			s.respondError(w, r, palikka.Errorf2(cause, palikka.ErrInternal, ErrMsgInternal))
			return
		}

		fmt.Println("user perms: ", p)

		for _, perm := range p {
			perms = append(perms, perm.Name)
		}

		fmt.Println("user perms mapped: ", perms)
	}

	authz := AuthorizationsResponse{
		Features: FeatureAuthorizations{
			Users: CRUDs{
				Create: slices.Contains(perms, palikka.PermUsersCreate),
				Read:   slices.Contains(perms, palikka.PermUsersRead),
				Update: slices.Contains(perms, palikka.PermsUsersUpdate),
				Delete: slices.Contains(perms, palikka.PermsUsersDelete),
			},
			Roles: CRUDs{
				Create: slices.Contains(perms, palikka.PermRolesCreate),
				Read:   slices.Contains(perms, palikka.PermRolesRead),
				Update: slices.Contains(perms, palikka.PermRolesUpdate),
				Delete: slices.Contains(perms, palikka.PermRolesDelete),
			},
		},
	}

	fmt.Println("user authz: ", authz)

	// marshal response
	if authzJSON, err = json.Marshal(authz); err != nil {
		cause := fmt.Errorf("handleGetAuthorizations: marshal json: %w", err)
		s.respondError(w, r, palikka.Errorf2(cause, palikka.ErrInternal, ErrMsgInternal))
		return
	}

	w.Header().Set(HeaderContentType, MediaTypeJSON)
	_, _ = w.Write(authzJSON)
}

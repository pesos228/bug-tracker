package appmw

import (
	"net/http"
	"strings"
)

func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userRoles, ok := UserRolesFromContext(r.Context())
		if !ok || len(userRoles) == 0 {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		var isAdmin bool
		for _, role := range userRoles {
			if strings.EqualFold(role, "admin") {
				isAdmin = true
				break
			}
		}

		if !isAdmin {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

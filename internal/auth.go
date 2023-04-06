package internal

import (
	"net/http"
	"strings"

	"github.com/elliotchance/pie/v2"
)

func authTokenHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(appConfig.AuthTokens) > 0 {
			auth := r.Header.Get("Authorization")
			if auth == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// extract Bearer token from Authorization header
			parts := strings.Split(auth, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			token := parts[1]

			// check if token is valid
			if !pie.Contains(appConfig.AuthTokens, token) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

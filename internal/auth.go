package internal

import (
	"context"
	"net/http"
)

type authTokenKey struct{}

func getAuthToken(ctx context.Context) string {
	return ctx.Value(authTokenKey{}).(string)
}

func authTokenHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// auth := r.Header.Get("Authorization")
		// if auth == "" {
		// 	w.WriteHeader(http.StatusUnauthorized)
		// 	return
		// }

		// // extract Bearer token from Authorization header
		// parts := strings.Split(auth, " ")
		// if len(parts) != 2 || parts[0] != "Bearer" {
		// 	w.WriteHeader(http.StatusUnauthorized)
		// 	return
		// }

		// token := parts[1]

		// ctx := context.WithValue(r.Context(), authTokenKey{}, token)
		// r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

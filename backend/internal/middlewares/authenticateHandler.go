package middlewares

import "net/http"

func AuthenticateHandler(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			Authenticate(next).ServeHTTP(w, r)
		})
}
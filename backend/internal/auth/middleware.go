package auth

import (
	"net/http"

	"github.com/gorilla/sessions"
)

// IsAuthenticated is middleware that redirects unauthenticated requests to "/".
// Wrap any handler that requires a logged-in user with this middleware.
func IsAuthenticated(store sessions.Store, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, "auth-session")
		if err != nil || session.Values["profile"] == nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

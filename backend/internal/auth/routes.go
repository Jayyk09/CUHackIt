package auth

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"

	"github.com/Jayyk09/CUHackIt/config"
)

// RegisterRoutes wires up the Auth0 login/callback/logout/profile routes.
// It must be called once during application startup.
func RegisterRoutes(r *http.ServeMux, cfg *config.Config, store sessions.Store) error {

	auth, err := New(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialise Auth0 authenticator: %w", err)
	}

	h := newHandler(auth, store, cfg)

	r.HandleFunc("GET /login", h.Login)
	r.HandleFunc("GET /callback", h.Callback)
	r.HandleFunc("GET /logout", h.Logout)

	return nil
}

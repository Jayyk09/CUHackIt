package auth

import (
	"encoding/gob"
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"

	"github.com/Jayyk09/CUHackIt/config"
)

// RegisterRoutes wires up the Auth0 login/callback/logout/profile routes.
// It must be called once during application startup.
func RegisterRoutes(r *http.ServeMux, cfg *config.Config) error {
	// gorilla/sessions stores arbitrary types in the cookie via gob encoding.
	// map[string]interface{} is the type we use for the Auth0 profile claim.
	gob.Register(map[string]interface{}{})

	store := sessions.NewCookieStore([]byte(cfg.Auth0.SessionSecret))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 1 week
		HttpOnly: true,
	}

	auth, err := New(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialise Auth0 authenticator: %w", err)
	}

	h := newHandler(auth, store, cfg)

	r.HandleFunc("GET /login", h.Login)
	r.HandleFunc("GET /callback", h.Callback)
	r.HandleFunc("GET /logout", h.Logout)
	// /user is protected — only accessible once authenticated.
	r.Handle("GET /user", IsAuthenticated(store, http.HandlerFunc(h.Profile)))

	return nil
}

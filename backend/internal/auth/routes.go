package auth

import (
	"encoding/gob"
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"

	"github.com/Jayyk09/CUHackIt/config"
	"github.com/Jayyk09/CUHackIt/internal/database"
)

// RegisterRoutes wires up the Auth0 login/callback/logout/profile routes.
// It must be called once during application startup.
func RegisterRoutes(r *http.ServeMux, cfg *config.Config, db *database.DB) (sessions.Store, error) {
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
		return nil, fmt.Errorf("failed to initialise Auth0 authenticator: %w", err)
	}

	h := newHandler(auth, store, cfg, db)

	r.HandleFunc("GET /login", h.Login)
	r.HandleFunc("GET /callback", h.Callback)
	r.HandleFunc("GET /logout", h.Logout)

	return store, nil
}

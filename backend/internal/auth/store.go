package auth

import (
	"encoding/gob"

	"github.com/gorilla/sessions"

	"github.com/Jayyk09/CUHackIt/config"
)

// NewSessionStore builds the cookie store used for Auth0 sessions.
func NewSessionStore(cfg *config.Config) sessions.Store {
	// gorilla/sessions stores arbitrary types in the cookie via gob encoding.
	// map[string]interface{} is the type we use for the Auth0 profile claim.
	gob.Register(map[string]interface{}{})

	store := sessions.NewCookieStore([]byte(cfg.Auth0.SessionSecret))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 1 week
		HttpOnly: true,
	}

	return store
}

package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/gorilla/sessions"

	"github.com/Jayyk09/CUHackIt/config"
	"github.com/Jayyk09/CUHackIt/internal/database"
	"github.com/Jayyk09/CUHackIt/internal/users"
)

// Handler holds the authenticator and session store used by auth routes.
type Handler struct {
	auth  *Authenticator
	store sessions.Store
	cfg   *config.Config
	db    *database.DB
}

func newHandler(auth *Authenticator, store sessions.Store, cfg *config.Config, db *database.DB) *Handler {
	return &Handler{auth: auth, store: store, cfg: cfg, db: db}
}

// Login initiates the Auth0 authorization code flow.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	state, err := generateRandomState()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session, err := h.store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["state"] = state
	if err := session.Save(r, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, h.auth.AuthCodeURL(state), http.StatusTemporaryRedirect)
}

// Callback handles the redirect from Auth0 after the user authenticates.
func (h *Handler) Callback(w http.ResponseWriter, r *http.Request) {
	session, err := h.store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.URL.Query().Get("state") != session.Values["state"] {
		http.Error(w, "Invalid state parameter.", http.StatusBadRequest)
		return
	}

	token, err := h.auth.Exchange(r.Context(), r.URL.Query().Get("code"))
	if err != nil {
		http.Error(w, "Failed to exchange authorization code for a token.", http.StatusUnauthorized)
		return
	}

	idToken, err := h.auth.VerifyIDToken(r.Context(), token)
	if err != nil {
		http.Error(w, "Failed to verify ID Token.", http.StatusInternalServerError)
		return
	}

	var profile map[string]interface{}
	if err := idToken.Claims(&profile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["access_token"] = token.AccessToken
	session.Values["profile"] = profile
	if err := session.Save(r, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Extract the Auth0 sub (user ID), email and name from the profile.
	sub, _ := profile["sub"].(string)
	email, _ := profile["email"].(string)
	name, _ := profile["name"].(string)
	if name == "" {
		name, _ = profile["nickname"].(string)
	}

	// Check if this is a new user.
	isNew := false
	if sub != "" {
		repo := users.NewRepository(h.db.Pool)
		_, err := repo.GetByAuth0ID(r.Context(), sub)
		if errors.Is(err, users.ErrUserNotFound) {
			_, _ = repo.Create(r.Context(), users.CreateUserInput{
				Auth0ID: sub,
				Email:   email,
				Name:    name,
			})
			isNew = true
		}
	}

	// Redirect to the frontend with the Auth0 sub so the frontend can resolve the internal user ID.
	// New users go to /onboarding; existing users go to /dashboard.
	frontendURL := h.cfg.App.FrontendURL
	if isNew {
		http.Redirect(w, r, frontendURL+"/onboarding?uid="+url.QueryEscape(sub), http.StatusTemporaryRedirect)
	} else {
		http.Redirect(w, r, frontendURL+"/dashboard?uid="+url.QueryEscape(sub), http.StatusTemporaryRedirect)
	}
}

// Logout clears the session and redirects to Auth0's logout endpoint.
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	logoutURL, err := url.Parse("https://" + h.cfg.Auth0.Domain + "/v2/logout")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect back to the frontend after Auth0 logout.
	returnTo, err := url.Parse(h.cfg.App.FrontendURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	params := url.Values{}
	params.Add("returnTo", returnTo.String())
	params.Add("client_id", h.cfg.Auth0.ClientID)
	logoutURL.RawQuery = params.Encode()

	// Invalidate the session cookie.
	session, _ := h.store.Get(r, "auth-session")
	session.Options.MaxAge = -1
	_ = session.Save(r, w)

	http.Redirect(w, r, logoutURL.String(), http.StatusTemporaryRedirect)
}

// Profile returns the authenticated user's profile as JSON.
// Protected by IsAuthenticated middleware, so profile is guaranteed non-nil.
func (h *Handler) Profile(w http.ResponseWriter, r *http.Request) {
	session, err := h.store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	profile := session.Values["profile"]
	if profile == nil {
		http.Error(w, "not authenticated", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(profile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func generateRandomState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

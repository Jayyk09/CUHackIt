package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/gorilla/sessions"

	"github.com/Jayyk09/CUHackIt/config"
)

// Handler holds the authenticator and session store used by auth routes.
type Handler struct {
	auth  *Authenticator
	store sessions.Store
	cfg   *config.Config
}

func newHandler(auth *Authenticator, store sessions.Store, cfg *config.Config) *Handler {
	return &Handler{auth: auth, store: store, cfg: cfg}
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

	http.Redirect(w, r, "/user", http.StatusTemporaryRedirect)
}

// Logout clears the session and redirects to Auth0's logout endpoint.
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	logoutURL, err := url.Parse("https://" + h.cfg.Auth0.Domain + "/v2/logout")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	returnTo, err := url.Parse(scheme + "://" + r.Host)
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

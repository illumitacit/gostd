package chistd

import (
	"net/http"

	"gitea.com/go-chi/session"
	"github.com/go-chi/chi/v5"
	"github.com/ory/nosurf"
	"go.uber.org/zap"

	"github.com/fensak-io/gostd/webstd"
)

const (
	OIDCLoginPath    = "/oidc/login"
	OIDCLogoutPath   = "/oidc/logout"
	OIDCCallbackPath = "/oidc/callback"
)

type oidcHandlerContext struct {
	auth     *webstd.Authenticator
	logger   *zap.Logger
	homePath string
}

func NewOIDCHandlerContext(
	auth *webstd.Authenticator,
	logger *zap.Logger,
	homePath string,
) *oidcHandlerContext {
	return &oidcHandlerContext{auth: auth, logger: logger, homePath: homePath}
}

// AddOIDCHandlerRoutes will add a group of routes that can be used to implement OIDC client protocol to manage
// authentication into an existing go-chi based web app. Note that this depends on the following two middlewares:
// - gitea.com/go-chi/session
// - github.com/ory/nosurf
func (h oidcHandlerContext) AddOIDCHandlerRoutes(router chi.Router) {
	router.Get(OIDCLoginPath, h.oidcLoginHandler)
	router.Get(OIDCLogoutPath, h.oidcLogoutHandler)
	router.Get(OIDCCallbackPath, h.oidcCallbackHandler)
}

func (h oidcHandlerContext) oidcLoginHandler(w http.ResponseWriter, r *http.Request) {
	stateToken := nosurf.Token(r)
	http.Redirect(
		w, r,
		h.auth.AuthCodeURL(stateToken),
		http.StatusTemporaryRedirect,
	)
}

func (h oidcHandlerContext) oidcLogoutHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: there needs to be a way to invalidate the auth token, but it looks like logout is dependent on the platform.
	sess := session.GetSession(r)
	if err := sess.Destroy(w, r); err != nil {
		h.logger.Sugar().Errorf("Error clearing session on logout: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(
		w, r,
		OIDCLoginPath,
		http.StatusTemporaryRedirect,
	)
}

func (h oidcHandlerContext) oidcCallbackHandler(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.Sugar()
	ctx := r.Context()

	state := r.URL.Query().Get("state")
	if !nosurf.VerifyToken(nosurf.Token(r), state) {
		logger.Warn("Detected wrong state parameter on oidc login. Possible XSRF attack.")
		http.Redirect(
			w, r,
			OIDCLoginPath,
			http.StatusTemporaryRedirect,
		)
		return
	}

	// Exchange an authorization code for a token.
	token, err := h.auth.Exchange(ctx, r.URL.Query().Get("code"))
	if err != nil {
		logger.Errorf("Error exchaning authorization code for a token: %s", err)
		http.Redirect(
			w, r,
			OIDCLoginPath,
			http.StatusTemporaryRedirect,
		)
		return
	}

	idToken, err := h.auth.VerifyIDToken(ctx, token)
	if err != nil {
		logger.Errorf("Error validating exchanged id token: %s", err)
		http.Redirect(
			w, r,
			OIDCLoginPath,
			http.StatusTemporaryRedirect,
		)
		return
	}

	// TODO: define a concrete struct
	var profile map[string]interface{}
	if err := idToken.Claims(&profile); err != nil {
		logger.Errorf("Error parsing id token claims: %s", err)
		http.Redirect(
			w, r,
			OIDCLoginPath,
			http.StatusTemporaryRedirect,
		)
		return
	}

	sess := session.GetSession(r)
	if err := sess.Set("access_token", token.AccessToken); err != nil {
		logger.Errorf("Error setting access token on session: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := sess.Set("profile", profile); err != nil {
		logger.Errorf("Error setting profile on session: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Redirect to logged in page.
	http.Redirect(
		w, r,
		h.homePath,
		http.StatusTemporaryRedirect,
	)
}

package chistd

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/ory/nosurf"
	"go.uber.org/zap"

	"github.com/fensak-io/gostd/webstd"
)

const (
	// URL paths
	OIDCLoginPath    = "/oidc/login"
	OIDCLogoutPath   = "/oidc/logout"
	OIDCCallbackPath = "/oidc/callback"

	// Session keys
	AccessTokenSessionKey = "access_token"
	UserProfileSessionKey = "profile"
)

type OIDCHandlerContext[T any] struct {
	logger   *zap.Logger
	auth     *webstd.Authenticator
	sessMgr  *scs.SessionManager
	homePath string
}

// NewOIDCHandlerContext returns a new handler context for the OIDC pages. The generic type parameter represents the
// profile struct to marshal the ID token claims to.
func NewOIDCHandlerContext[T any](
	logger *zap.Logger,
	auth *webstd.Authenticator,
	sessMgr *scs.SessionManager,
	homePath string,
) *OIDCHandlerContext[T] {
	return &OIDCHandlerContext[T]{
		logger: logger, auth: auth, sessMgr: sessMgr, homePath: homePath,
	}
}

// AddOIDCHandlerRoutes will add a group of routes that can be used to implement OIDC client protocol to manage
// authentication into an existing go-chi based web app. Note that this depends on the following two middlewares:
// - github.com/alexedwards/scs/v2
// - github.com/ory/nosurf
func (h OIDCHandlerContext[T]) AddOIDCHandlerRoutes(router chi.Router) {
	router.Get(OIDCLoginPath, h.oidcLoginHandler)
	router.Get(OIDCLogoutPath, h.oidcLogoutHandler)
	router.Get(OIDCCallbackPath, h.oidcCallbackHandler)
}

func (h OIDCHandlerContext[T]) oidcLoginHandler(w http.ResponseWriter, r *http.Request) {
	stateToken := nosurf.Token(r)
	http.Redirect(
		w, r,
		h.auth.AuthCodeURL(stateToken),
		http.StatusTemporaryRedirect,
	)
}

func (h OIDCHandlerContext[T]) oidcLogoutHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: there needs to be a way to invalidate the auth token, but it looks like logout is dependent on the platform.
	ctx := r.Context()
	if err := h.sessMgr.Destroy(ctx); err != nil {
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

func (h OIDCHandlerContext[T]) oidcCallbackHandler(w http.ResponseWriter, r *http.Request) {
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

	var profile T
	if err := idToken.Claims(&profile); err != nil {
		logger.Errorf("Error parsing id token claims: %s", err)
		http.Redirect(
			w, r,
			OIDCLoginPath,
			http.StatusTemporaryRedirect,
		)
		return
	}

	h.sessMgr.Put(ctx, AccessTokenSessionKey, token.AccessToken)
	h.sessMgr.Put(ctx, UserProfileSessionKey, profile)

	// Redirect to logged in page.
	http.Redirect(
		w, r,
		h.homePath,
		http.StatusTemporaryRedirect,
	)
}

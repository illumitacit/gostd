package chistd

import (
	"encoding/gob"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/ory/nosurf"
	"go.uber.org/zap"
	"golang.org/x/oauth2"

	"github.com/illumitacit/gostd/webstd"
)

func init() {
	// Register time.Time to gob so it can be stored in the session
	gob.Register(time.Time{})
}

const (
	// URL paths
	OIDCRegisterPath = "/oidc/register"
	OIDCLoginPath    = "/oidc/login"
	OIDCLogoutPath   = "/oidc/logout"
	OIDCCallbackPath = "/oidc/callback"

	// Session keys
	IDTokenSessionKey           = "id_token"
	AccessTokenSessionKey       = "access_token"
	AccessTokenExpirySessionKey = "access_token_expiry"
	RefreshTokenSessionKey      = "refresh_token"
	UserProfileSessionKey       = "profile"
	ContinueToURLSessionKey     = "continue_to"
	PKCECodeVerifierSessionKey  = "pkce_code_verifier"
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
	router.Get(OIDCRegisterPath, h.oidcRegisterHandler)
	router.Get(OIDCLoginPath, h.oidcLoginHandler)
	router.Get(OIDCLogoutPath, h.oidcLogoutHandler)
	router.Get(OIDCCallbackPath, h.oidcCallbackHandler)
}

func (h OIDCHandlerContext[T]) oidcRegisterHandler(w http.ResponseWriter, r *http.Request) {
	stateToken, opts, err := h.oidcSetupLogin(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	opts = append(
		opts,
		oauth2.SetAuthURLParam("prompt", "create"),
	)

	http.Redirect(
		w, r,
		h.auth.AuthCodeURL(stateToken, opts...),
		http.StatusTemporaryRedirect,
	)
}

func (h OIDCHandlerContext[T]) oidcLoginHandler(w http.ResponseWriter, r *http.Request) {
	stateToken, opts, err := h.oidcSetupLogin(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(
		w, r,
		h.auth.AuthCodeURL(stateToken, opts...),
		http.StatusTemporaryRedirect,
	)
}

func (h OIDCHandlerContext[T]) oidcSetupLogin(r *http.Request) (string, []oauth2.AuthCodeOption, error) {
	stateToken := nosurf.Token(r)
	opts := []oauth2.AuthCodeOption{}
	if h.auth.WithPKCE {
		codeVerifier, err := h.auth.NewCodeVerifier()
		if err != nil {
			h.logger.Sugar().Errorf("%s", err)
			return "", nil, err
		}
		opts = append(
			opts,
			oauth2.SetAuthURLParam("code_challenge_method", "S256"),
			oauth2.SetAuthURLParam("code_challenge", codeVerifier.Challenge),
		)

		// Store the verifier in the session so it can be used after the login finishes
		h.sessMgr.Put(r.Context(), PKCECodeVerifierSessionKey, codeVerifier.Verifier)
	}
	return stateToken, opts, nil
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

	opts := []oauth2.AuthCodeOption{}
	if h.auth.WithPKCE {
		rawCodeVerifier := h.sessMgr.Get(ctx, PKCECodeVerifierSessionKey)
		if rawCodeVerifier == nil {
			logger.Errorf("Required PKCE, but no verifier in session.")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		codeVerifier := rawCodeVerifier.(string)
		opts = append(opts, oauth2.SetAuthURLParam("code_verifier", codeVerifier))
	}

	// Exchange an authorization code for a token.
	token, err := h.auth.Exchange(ctx, r.URL.Query().Get("code"), opts...)
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
	rawIDToken := token.Extra("id_token").(string)

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

	h.sessMgr.Put(ctx, IDTokenSessionKey, rawIDToken)
	h.sessMgr.Put(ctx, AccessTokenSessionKey, token.AccessToken)
	h.sessMgr.Put(ctx, AccessTokenExpirySessionKey, token.Expiry)
	h.sessMgr.Put(ctx, RefreshTokenSessionKey, token.RefreshToken)
	h.sessMgr.Put(ctx, UserProfileSessionKey, profile)

	// Clear the PKCE code verifier from the session now that the token is verified
	h.sessMgr.Put(r.Context(), PKCECodeVerifierSessionKey, nil)

	// If there is a continue URL recorded in the session, redirect to there.
	// Otherwise, redirect to the default home page.
	continueTo, hasContinueTo := h.sessMgr.Get(ctx, ContinueToURLSessionKey).(string)
	if !hasContinueTo {
		continueTo = h.homePath
	}
	http.Redirect(
		w, r,
		continueTo,
		http.StatusTemporaryRedirect,
	)
}

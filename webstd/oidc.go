package webstd

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

// Authenticator is used to authenticate our users.
type Authenticator struct {
	*oidc.Provider
	oauth2.Config
	WithPKCE bool
}

// NewAuthenticator instantiates the Authenticator object using the provided configuration options.
func NewAuthenticator(ctx context.Context, cfg *OIDCProvider) (*Authenticator, error) {
	discoveryURL := cfg.IssuerURL
	if cfg.SkipIssuerVerification {
		ctx = oidc.InsecureIssuerURLContext(ctx, cfg.IssuerURL)
		discoveryURL = cfg.DiscoveryURL
	}

	provider, err := oidc.NewProvider(ctx, discoveryURL)
	if err != nil {
		return nil, err
	}

	conf := oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.CallbackURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       append([]string{oidc.ScopeOpenID}, cfg.AdditionalScopes...),
	}

	return &Authenticator{
		Provider: provider,
		Config:   conf,
		WithPKCE: cfg.WithPKCE,
	}, nil
}

// PKCECodeVerifier captures the code verifier string, as well as the hashed string that can be used as the code
// challenge for the PKCE flow.
type PKCECodeVerifier struct {
	Verifier  string
	Challenge string
}

// NewCodeVerifier creates cryptographically secure code verification string for the PKCE flow.
func (a Authenticator) NewCodeVerifier() (PKCECodeVerifier, error) {
	codeVerifier, err := randomBytesInHex(32)
	if err != nil {
		return PKCECodeVerifier{}, fmt.Errorf("Could not create a code verifier: %s", err)
	}
	sha2 := sha256.New()
	if _, err := io.WriteString(sha2, codeVerifier); err != nil {
		return PKCECodeVerifier{}, fmt.Errorf("Could not write codeVerifier string to sha256: %s", err)
	}
	codeChallenge := base64.RawURLEncoding.EncodeToString(sha2.Sum(nil))
	return PKCECodeVerifier{
		Verifier:  codeVerifier,
		Challenge: codeChallenge,
	}, nil
}

// VerifyIDToken verifies that an *oauth2.Token is a valid *oidc.IDToken.
func (a Authenticator) VerifyIDToken(ctx context.Context, token *oauth2.Token) (*oidc.IDToken, error) {
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("no id_token field in oauth2 token")
	}

	oidcConfig := &oidc.Config{
		ClientID: a.ClientID,
	}

	return a.Verifier(oidcConfig).Verify(ctx, rawIDToken)
}

// RefreshIDToken obtains a new OIDC ID token using the provided refresh token.
func (a Authenticator) RefreshIDToken(ctx context.Context, refreshToken string) (string, *oidc.IDToken, error) {
	ts := a.TokenSource(ctx, &oauth2.Token{RefreshToken: refreshToken})
	token, err := ts.Token()
	if err != nil {
		return "", nil, err
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return "", nil, errors.New("no id_token field in oauth2 token")
	}
	idToken, err := a.VerifyIDToken(ctx, token)
	return rawIDToken, idToken, err
}

// LogoutURL returns the logout URL to end the session, if it exists. Note that there is no OIDC standard for RP
// initiated logout. As such, there is no guarantee that this will always return a valid logout URL. For IdPs where we
// can not determine a valid logout URL, this will return an empty string.
// NOTE: for now, we only support the `end_session_endpoint` claim, which is used by Azure AD B2C.
func (a Authenticator) LogoutURL() string {
	var claims struct {
		EndSessionEndpoint string `json:"end_session_endpoint"`
	}
	if err := a.Claims(&claims); err != nil {
		// Does not contain the end_session_endpoint in the claims, so return empty string.
		return ""
	}
	return claims.EndSessionEndpoint
}

func randomBytesInHex(count int) (string, error) {
	buf := make([]byte, count)
	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		return "", fmt.Errorf("Could not generate %d random bytes: %v", count, err)
	}
	return hex.EncodeToString(buf), nil
}

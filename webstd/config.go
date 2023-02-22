package webstd

import (
	"time"
)

// OIDCProvider represents configuration options for the OIDC Provider that handles authentication for the web app.
// This can be embedded in a viper compatible config struct.
type OIDCProvider struct {
	// IssuerURL is the full URL (including scheme and path) of the OIDC provider issuer.
	IssuerURL string `mapstructure:"issuer_url"`

	// ClientID is the oauth2 application client ID to use for the OIDC protocol.
	ClientID string `mapstructure:"clientid"`

	// ClientSecret is the oauth2 application client secret to use for the OIDC protocol.
	ClientSecret string `mapstructure:"secret"`

	// SkipIssuerVerification determines whether the issuer URL should be verified against the discovery base URL. This
	// should ONLY be set to true for OIDC providers that are off-spec, such as Azure where the discovery URL
	// (/.well-known/openid-configuration) is different from the issuer URL. When true, the discovery URL must be
	// provided under the DiscoveryURL config.
	SkipIssuerVerification bool `mapstructure:"skip_iss_verification"`

	// DiscoveryURL is the full base URL of the discovery page for OIDC. The authenticator will look for the OIDC
	// configuration under the page DISCOVERY_URL/.well-known/openid-configuration. Only used if SkipIssuerVerification is
	// true; when SkipIssuerVerification is false, the IssuerURL will be used instead.
	DiscoveryURL string `mapstructure:"discovery_url"`

	// AdditionalScopes is the list of Oauth2 scopes to request for the OIDC token. Note that the library will always
	// request the required "openid" scope.
	AdditionalScopes []string `mapstructure:"additional_scopes"`

	// CallbackURL is the full URL (including scheme) of the endpoint that handles the access token returned from the OIDC
	// protocol. This should be automatically configured by the application instead of being configured in the config
	// chain.
	CallbackURL string
}

// Session represents configuration options for the Session object and cookie.
// This can be embedded in a viper compatible config struct.
type Session struct {
	// Lifetime indicates how long a session is valid for.
	Lifetime time.Duration `mapstructure:"lifetime"`

	// CookieName is the name of the cookie to use to store the session ID on the client side.
	CookieName string `mapstructure:"cookie_name"`

	// CookieSecure determines whether the secure flag should be set on the cookie.
	CookieSecure bool `mapstructure:"cookie_secure"`

	// CookieSameSiteStr is the string representation of the samesite mode to set on the session cookie.
	CookieSameSiteStr string `mapstructure:"cookie_samesite"`
}

// CSRF represents configuration options for CSRF protection.
// This can be embedded in a viper compatible config struct.
type CSRF struct {
	MaxAge int `mapstructure:"maxage"`

	// Dev determines whether to use dev mode for CSRF validation. When true, disables the secure flag on the CSRF cookie.
	Dev bool `mapstructure:"dev"`
}

// IdP represents configuration options for interacting with the Identity Provider that handles authentication for the
// web app. This can be embedded in a viper compatible config struct.
type IdP struct {
	// Provider represents one of the supported identity providers.
	Provider IdPProvider `mapstructure:"provider"`

	// The ID of the AAD B2C Tenant. Only used if the provider is set to aadb2c.
	TenantID string `mapstructure:"tenantid"`

	// The name of the AAD B2C Tenant. Only used if the provider is set to aadb2c.
	TenantName string `mapstructure:"tenant_name"`
}

// IdPProvider is an enum describing the possible options for the IdP.Provider setting.
type IdPProvider string

const (
	IdPProviderAADB2C IdPProvider = "aadb2c"
)

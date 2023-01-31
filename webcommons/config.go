package webcommons

// OIDCProvider represents configuration options for the OIDC Provider that handles authentication for the web app.
// This can be embedded in a viper compatible config struct.
type OIDCProvider struct {
	// IssuerURL is the full URL (including scheme and path) of the OIDC provider issuer.
	IssuerURL string `mapstructure:"issuerurl"`

	// ClientID is the oauth2 application client ID to use for the OIDC protocol.
	ClientID string `mapstructure:"clientid"`

	// ClientSecret is the oauth2 application client secret to use for the OIDC protocol.
	ClientSecret string `mapstructure:"secret"`

	// SkipIssuerVerification determines whether the issuer URL should be verified against the discovery base URL. This
	// should ONLY be set to true for OIDC providers that are off-spec, such as Azure where the discovery URL
	// (/.well-known/oidc-configuration) is different from the issuer URL. When true, the discovery URL must be
	// provided under the DiscoveryURL config.
	SkipIssuerVerification bool `mapstructure:"skipissverification"`

	// DiscoveryURL is the full base URL of the discovery page for OIDC. The authenticator will look for the OIDC
	// configuration under the page DISCOVERY_URL/.well-known/oidc-configuration. Only used if SkipIssuerVerification is
	// true; when SkipIssuerVerification is false, the IssuerURL will be used instead.
	DiscoveryURL string `mapstructure:"discoveryurl"`

	// CallbackURL is the full URL (including scheme) of the endpoint that handles the access token returned from the OIDC
	// protocol.
	CallbackURL string `mapstructure:"callbackurl"`
}

// CSRF represents configuration options for CSRF protection.
// This can be embedded in a viper compatible config struct.
type CSRF struct {
	MaxAge int `mapstructure:"maxage"`

	// Dev determines whether to use dev mode for CSRF validation. When true, disables the secure flag on the CSRF cookie.
	Dev bool `mapstructure:"dev"`
}

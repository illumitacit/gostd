package webcommons

// OIDCProvider represents configuration options for the OIDC Provider that handles authentication for the web app.
// This can be embedded in a viper compatible config struct.
type OIDCProvider struct {
	// Domain is the domain name of the OIDC provider. E.g., if you are using Azure AD B2C, then this is the domain of the
	// tenant (the one that ends in onmicrosoft.com).
	Domain string `mapstructure:"domain"`

	// ClientID is the oauth2 application client ID to use for the OIDC protocol.
	ClientID string `mapstructure:"clientid"`

	// ClientSecret is the oauth2 application client secret to use for the OIDC protocol.
	ClientSecret string `mapstructure:"secret"`

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

package webcli

import (
	"time"

	"github.com/illumitacit/gostd/clistd"
	"github.com/illumitacit/gostd/webstd"
	"github.com/spf13/pflag"
)

// BindOIDCCfgFlags binds the necessary cobra CLI flags for configuring OIDC. This will also make sure to bind the CLI
// flags to viper as well so that the config is loaded.
func BindOIDCCfgFlags(flags *pflag.FlagSet, flagPrefix, cfgPrefix string) {
	flags.String(flagPrefix+"oidc-issuer", "", "The full URL (including domain and path) of the OIDC provider issuer.")
	clistd.MustBindPFlag(cfgPrefix+"oidc.issuer_url", flags.Lookup(flagPrefix+"oidc-issuer"))

	flags.String(flagPrefix+"oidc-clientid", "", "The oauth2 application client ID to use for the OIDC protocol.")
	clistd.MustBindPFlag(cfgPrefix+"oidc.clientid", flags.Lookup(flagPrefix+"oidc-clientid"))

	flags.String(flagPrefix+"oidc-secret", "", "The oauth2 application client secret to use for the OIDC protocol. Recommended to be set using an environment variable.")
	clistd.MustBindPFlag(cfgPrefix+"oidc.secret", flags.Lookup(flagPrefix+"oidc-secret"))

	flags.StringSlice(flagPrefix+"oidc-raw-token-clientids", []string{}, "The oauth2 application client ID that is approved for authenticating to the API.")
	clistd.MustBindPFlag(cfgPrefix+"oidc.raw_token_client_ids", flags.Lookup(flagPrefix+"oidc-raw-token-clientids"))

	flags.Bool(flagPrefix+"oidc-with-pkce", false, "Whether the oauth2 flow associated with OIDC should use the PKCE flow.")
	clistd.MustBindPFlag(cfgPrefix+"oidc.with_pkce", flags.Lookup(flagPrefix+"oidc-with-pkce"))

	flags.Bool(flagPrefix+"oidc-skip-issuer-verification", false, "Whether the OIDC discovery process should skip verifying the issuer URL against the discovery URL. This should only be used for off-spec providers where the discovery URL is different from the issuer URL, like Azure. When true, --oidc-discovery must be provided.")
	clistd.MustBindPFlag(cfgPrefix+"oidc.skip_iss_verification", flags.Lookup(flagPrefix+"oidc-skip-issuer-verification"))

	flags.String(flagPrefix+"oidc-discovery", "", "The full base URL (including domain and path) of the OIDC provider discovery page.")
	clistd.MustBindPFlag(cfgPrefix+"oidc.discovery_url", flags.Lookup(flagPrefix+"oidc-discovery"))

	flags.StringSlice(flagPrefix+"oidc-scopes", nil, "The list of Oauth2 scopes that should be requested for the OIDC token.")
	clistd.MustBindPFlag(cfgPrefix+"oidc.additional_scopes", flags.Lookup(flagPrefix+"oidc-scopes"))
}

// BindSessionCfgFlags binds the necessary cobra CLI flags for configuring the web session. This will also make sure to
// bind the CLI flags to viper as well so that the config is loaded.
func BindSessionCfgFlags(flags *pflag.FlagSet, cfgPrefix, cookieNameDefault string) {
	flags.Duration("session-lifetime", 336*time.Hour, "The lifetime of the session cookie.")
	clistd.MustBindPFlag(cfgPrefix+"session.lifetime", flags.Lookup("session-lifetime"))

	flags.String("session-cookie", cookieNameDefault, "The name of the cookie to use for storing the web session ID.")
	clistd.MustBindPFlag("web.session.cookie_name", flags.Lookup("session-cookie"))

	flags.Bool("session-cookie-secure", true, "Whether the secure flag should be set on the session cookie.")
	clistd.MustBindPFlag("web.session.cookie_secure", flags.Lookup("session-cookie-secure"))

	flags.String("session-cookie-samesite", "lax", "The samesite mode to be set on the session cookie.")
	clistd.MustBindPFlag("web.session.cookie_samesite", flags.Lookup("session-cookie-samesite"))
}

// BindCSRFCfgFlags binds the necessary cobra CLI flags for configuring the CSRF protection. This will also make sure to
// bind the CLI flags to viper as well so that the config is loaded.
func BindCSRFCfgFlags(flags *pflag.FlagSet, cfgPrefix string) {
	flags.Int("csrf-maxage", 0, "Maximum age of CSRF Token cookies.")
	clistd.MustBindPFlag(cfgPrefix+"csrf.maxage", flags.Lookup("csrf-maxage"))

	flags.Bool("csrf-dev", false, "When true, CSRF Tokens will not be secured, allowing for passage via http.")
	clistd.MustBindPFlag(cfgPrefix+"csrf.dev", flags.Lookup("csrf-dev"))
}

// BindIdPCfgFlags binds the necessary cobra CLI flags for configuring the IdP interaction. This will also make sure to
// bind the CLI flags to viper as well so that the config is loaded.
func BindIdPCfgFlags(flags *pflag.FlagSet, cfgPrefix string, defaultIdPProvider webstd.IdPProvider) {
	flags.String("idp-provider", string(defaultIdPProvider), "The identity provider service that manages authentication to Fensak. Must be one of: zitadel, nopidp")
	clistd.MustBindPFlag(cfgPrefix+"idp.provider", flags.Lookup("idp-provider"))

	flags.String("zitadel-instance-name", "", "The name of the Zitadel Instance used for hosting users for the app. Only used if the provider is set to zitadel.")
	clistd.MustBindPFlag(cfgPrefix+"idp.zitadel.instance_name", flags.Lookup("zitadel-instance-name"))

	flags.String("zitadel-jwt-key", "", "The base64 encoded JWT key to use for authenticating to the Zitadel Admin API. Only used if the provider is set to zitadel. Recommended to be set with environment variables.")
	clistd.MustBindPFlag(cfgPrefix+"idp.zitadel.jwt_key_base64", flags.Lookup("zitadel-jwt-key"))
}

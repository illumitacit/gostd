package webstd

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"go.uber.org/zap"
)

// SetSessionSettings configures the session manager based on the provided session configuration.
func SetSessionSettings(
	logger *zap.Logger, sessMgr *scs.SessionManager, cfg *Session,
) {
	sugar := logger.Sugar()

	sessMgr.Lifetime = cfg.Lifetime
	sessMgr.Cookie.Name = cfg.CookieName
	sessMgr.Cookie.HttpOnly = true
	sessMgr.Cookie.Secure = cfg.CookieSecure

	switch cfg.CookieSameSiteStr {
	case "lax":
		sessMgr.Cookie.SameSite = http.SameSiteLaxMode
	case "strict":
		sessMgr.Cookie.SameSite = http.SameSiteStrictMode
	case "none":
		sugar.Warn("The session cookie is set with samesite none mode. This is prone to XSRF attacks! Reconsider configuring with samesite strict or samesite lax.")
		sessMgr.Cookie.SameSite = http.SameSiteNoneMode
	default:
		sugar.Warnf(
			`Could not parse cookie samesite mode "%s". Using default mode.`,
			cfg.CookieSameSiteStr,
		)
		sessMgr.Cookie.SameSite = http.SameSiteDefaultMode
	}
}

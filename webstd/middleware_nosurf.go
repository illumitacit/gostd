package webstd

import (
	"net/http"

	"github.com/ory/nosurf"
)

// NewNosurfHandler returns a nosurf handler function that can be used as a http middleware. The nosurf handler will
// take care to ensure that a valid CSRF token is provided in every PUT, POST, DELETE request.
func NewNosurfHandler(cfg *CSRF) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		csrfH := nosurf.New(h)
		csrfH.SetBaseCookieFunc(func(w http.ResponseWriter, r *http.Request) http.Cookie {
			secure := !cfg.Dev
			cookie := http.Cookie{
				Name:     nosurf.CookieName,
				MaxAge:   cfg.MaxAge,
				HttpOnly: true,
				Path:     "/",
				Secure:   secure,

				// We have to use Lax same site mode for CSRF tokens to support usage across OAuth exchanges. Otherwise, the
				// CSRF Token is never sent to the OAuth callback handlers on return since the originating request is the OAuth
				// server.
				SameSite: http.SameSiteLaxMode,
			}

			return cookie
		})
		return csrfH
	}
}

package chistd

import (
	"github.com/go-chi/chi/v5"
	"github.com/illumitacit/gostd/webstd"
)

// AddNosurfMiddleware will add the nosurf middleware into the chi stack.
func AddNosurfMiddleware(cfg *webstd.CSRF, router chi.Router) {
	router.Use(webstd.NewNosurfHandler(cfg))
}

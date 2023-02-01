package chistd

import (
	"github.com/fensak-io/gostd/webstd"
	"github.com/go-chi/chi/v5"
)

// AddNosurfMiddleware will add the nosurf middleware into the chi stack.
func AddNosurfMiddleware(cfg *webstd.CSRF, router chi.Router) {
	router.Use(webstd.NewNosurfHandler(cfg))
}

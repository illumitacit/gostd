package chistd

import (
	"github.com/fensak-io/gostd/webcommons"
	"github.com/go-chi/chi/v5"
)

// AddNosurfMiddleware will add the nosurf middleware into the chi stack.
func AddNosurfMiddleware(cfg *webcommons.CSRF, router chi.Router) {
	router.Use(webcommons.NewNosurfHandler(cfg))
}

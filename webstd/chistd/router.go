package chistd

import (
	"net/http"

	"github.com/fensak-io/httpzaplog"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type DefaultMiddlewareOptions struct {
	Logger          *zap.Logger
	ErrorMiddleware func(http.Handler) http.Handler
}

// NewRouterWithDefaultMiddlewares returns a new go-chi router that has a set of recommended default routers
// configured.
func NewRouterWithDefaultMiddlewares(opts DefaultMiddlewareOptions) chi.Router {
	router := chi.NewRouter()
	AddDefaultMiddlewares(router, opts)
	return router
}

// AddDefaultMiddlewares adds all the chi middlewares that we want for applied to all requests.
func AddDefaultMiddlewares(router chi.Router, opts DefaultMiddlewareOptions) {
	logOpts := &httpzaplog.Options{
		Logger:          opts.Logger,
		ErrorMiddleware: opts.ErrorMiddleware,

		// TODO: make configurable
		Concise:     false,
		SkipHeaders: nil,
	}

	router.Use(middleware.RealIP)

	// NOTE: httpzaplog includes RequestID and Recoverer middlewares
	router.Use(httpzaplog.RequestLogger(logOpts))
}

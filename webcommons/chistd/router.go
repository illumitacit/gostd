package chistd

import (
	"github.com/fensak-io/httpzaplog"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

// NewRouterWithDefaultMiddlewares returns a new go-chi router that has a set of recommended default routers
// configured.
func NewRouterWithDefaultMiddlewares(logger *zap.Logger) chi.Router {
	router := chi.NewRouter()
	AddDefaultMiddlewares(router, logger)
	return router
}

// AddDefaultMiddlewares adds all the chi middlewares that we want for applied to all requests.
func AddDefaultMiddlewares(router chi.Router, logger *zap.Logger) {
	logOpts := &httpzaplog.Options{
		Logger: logger,

		// TODO: make configurable
		Concise:     false,
		SkipHeaders: nil,
	}

	router.Use(middleware.RealIP)

	// NOTE: httpzaplog includes RequestID and Recoverer middlewares
	router.Use(httpzaplog.RequestLogger(logOpts))
}

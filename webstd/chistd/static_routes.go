package chistd

import (
	"embed"
	"io/fs"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// AddStaticRoutes creates routes for static assets. This will directly serve static files under the /static route, but
// also serve some favicon related files from the root path.
func AddStaticRoutes(logger *zap.Logger, router chi.Router, staticFS embed.FS) error {
	sugar := logger.Sugar()
	sugar.Debug("Setting up static routes")

	router.Handle("/static/*", http.FileServer(http.FS(staticFS)))

	faviconPaths, err := fs.Glob(staticFS, "static/favicons/*")
	if err != nil {
		sugar.Errorf("Error loading favicons: %s", err)
		return err
	}
	for _, path := range faviconPaths {
		sugar.Debugf("Loading favicon path %s", path)

		// Bring range var in scope of for loop, so that it is bound when request function is run
		path := path

		file, err := filepath.Rel("static/favicons", path)
		if err != nil {
			sugar.Errorf("Error finding favicon (%s) relative path: %s", path, err)
			return err
		}
		reqPath := filepath.Join("/", file)

		router.Get(reqPath, func(w http.ResponseWriter, r *http.Request) {
			data, err := staticFS.ReadFile(path)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(data)
		})
	}

	return nil
}

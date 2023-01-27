package gincommons

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fensak-io/gostd/quit"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type App struct {
	Engine          *gin.Engine
	Logger          *zap.SugaredLogger
	Port            int
	ShutdownTimeout time.Duration

	// Any addiitonal close routine should be handled in the custom close function passed in here.
	CloseFn func() error
}

// RunWithSignalHandler runs the Gin app described by the App struct in the background, and implements a signal handler
// in the foreground that traps the INT and TERM signals. When the INT or TERM signal is sent to the process, this will
// start a graceful shutdown of the Gin app, waiting up to ShutdownTimeout duration for all Gin threads to stop
// processing.
func RunWithSignalHandler(app *App) (returnErr error) {
	listen := fmt.Sprintf(":%d", app.Port)
	srv := &http.Server{
		Addr:    listen,
		Handler: app.Engine,
	}

	// Start the server in the background so that we can handle shutdown signals gracefully.
	waiterDoneCh := make(chan struct{})
	waiter := quit.GetWaiter()
	waiter.Add(1)
	go func() {
		defer waiter.Done()

		app.Logger.Debug("Starting Gin server")
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			app.Logger.Debugf("Error shutting down: %s", err)
			if returnErr == nil {
				returnErr = err
			}
		}

		if app.CloseFn != nil {
			app.Logger.Debug("Handling additional shutdown tasks")
			if err := app.CloseFn(); err != nil {
				app.Logger.Debugf("Error running additional shutdown tasks: %s", err)
				if returnErr == nil {
					returnErr = err
				}
			}
		}
	}()
	go func() {
		waiter.Wait()
		close(waiterDoneCh)
		app.Logger.Debug("Shutting down server")
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a configurable timeout.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh

	// At this point, an interrupt signal was triggered (SIGINT or SIGTERM), so start to gracefully shutdown the server.
	app.Logger.Infof("Received interrupt signal. Gracefully shutting down server...")
	timeout, cancel := context.WithTimeout(context.Background(), app.ShutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(timeout); err != nil && !errors.Is(err, context.DeadlineExceeded) {
		app.Logger.Errorf("Error shutting down server: %s", err)
		return err
	}
	select {
	case <-timeout.Done():
		app.Logger.Errorf("Timed out waiting for server to shutdown.")
		return fmt.Errorf("timeout")
	case <-waiterDoneCh:
		app.Logger.Infof("All services gracefully shutdown.")
	}
	return nil
}

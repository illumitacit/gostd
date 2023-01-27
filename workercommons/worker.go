package workercommons

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fensak-io/gostd/quit"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type App[T proto.Message] struct {
	Broker          *Broker
	Logger          *zap.SugaredLogger
	TaskHandler     TaskHandler[T]
	ShutdownTimeout time.Duration

	// Any addiitonal close routine should be handled in the custom close function passed in here.
	CloseFn func() error
}

// RunWithSignalHandler runs a worker process described by the App struct in the background, and implements a signal
// handler in the foreground that traps the INT and TERM signals. When the INT or TERM signal is sent to the process,
// this will start a graceful shutdown of the worker app, waiting up to ShutdownTimeout duration for all the worker
// threads to stop processing.
func RunWithSignalHandler[T proto.Message](app *App[T]) (returnErr error) {
	subscription, err := NewSubClient[T](app.Logger, app.Broker, context.Background())
	if err != nil {
		return err
	}

	app.Logger.Infof("Reading tasks from broker")

	// Start the worker in the background so that we can handle shutdown signals gracefully.
	errCh := make(chan error, 1)
	waiterDoneCh := make(chan struct{})
	waiter := quit.GetWaiter()
	waiter.Add(1)
	go func() {
		defer waiter.Done()
		defer func() {
			if err := subscription.Close(); err != nil {
				app.Logger.Errorf("Error closing subscription: %s", err)
				if returnErr == nil {
					returnErr = err
				}
			}

			if err := app.CloseFn(); err != nil {
				app.Logger.Errorf("Error closing connections: %s", err)
				if returnErr == nil {
					returnErr = err
				}
			}
		}()

		// TODO: handle panics for graceful recovery

		// The main loop pulls messages from the pubsub broker and exectues the tasks. This uses a few techniques:
		// - To ensure we can shutdown the worker, we run the receive task with a timeout. This is necessary so that the
		//   main goroutine doesn't endlessly wait for a task even if there has been a message sent in the quit channel.
		// - Use select to watch for the shutdown message in a non-blocking fashion.
		// - Add a time.After to ensure that we don't endlessly block on waiting for the quit channel.
		for {
			// Use a timeout context to avoid blocking the thread on receive. This allows the worker to able to handle shutdown
			// messages from the main thread.
			timeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			err := receiveMsgWithTimeout(
				app.Logger, subscription, app.TaskHandler,
				timeout, cancel,
			)
			if err != nil && !errors.Is(err, context.DeadlineExceeded) {
				errCh <- err
				return
			}
			select {
			case <-quit.GetQuitChannel():
				app.Logger.Debugf("Received shutdown message. Exiting loop.")
				errCh <- nil
				return
			case <-timeout.Done():
				app.Logger.Debugf("Receive loop broker wait timeout reached")
			// NOTE: We add a tiny delay here so that the other channels get read preference. Otherwise, when they have
			// messages, there is a small chance the time.After wins the concurrency.
			case <-time.After(1 * time.Millisecond):
				app.Logger.Debugf("Receive loop broker successfully handled message.")
			}
		}
	}()
	go func() {
		waiter.Wait()
		close(waiterDoneCh)
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a configurable timeout.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh

	// At this point, an interrupt signal was triggered (SIGINT or SIGTERM), so start to gracefully shutdown the server.
	app.Logger.Infof("Received interrupt signal. Gracefully shutting down worker...")
	timeout, cancel := context.WithTimeout(context.Background(), app.ShutdownTimeout)
	defer cancel()
	quit.BroadcastShutdown()
	select {
	case <-timeout.Done():
		app.Logger.Errorf("Timed out waiting for worker to shutdown.")
		return fmt.Errorf("timeout")
	case <-waiterDoneCh:
		app.Logger.Infof("All services gracefully shutdown.")
	}

	loopErr := <-errCh
	return loopErr
}

func receiveMsgWithTimeout[T proto.Message](
	logger *zap.SugaredLogger, subscription *SubClient[T], h TaskHandler[T],
	timeout context.Context, cancel context.CancelFunc,
) (returnErr error) {
	defer cancel()

	task, msg, err := subscription.ReceiveTask(timeout)
	if err != nil {
		// TODO: differentiate fatal error from ignorable errors
		if !errors.Is(err, context.DeadlineExceeded) {
			logger.Errorf("Error receiving message from broker: %s", err)
		}
		return err
	}
	// TODO: spawn goroutine for handling the task, but with work pooling to prevent overflowing.
	if err := h.HandleTaskMsg(task, msg); err != nil {
		// NOTE: we don't halt on task errors so that the worker continues to process other messages.
		logger.Errorf("Error processing task TODO from broker: %s", err)
	}
	logger.Infof("Successfully processed task TODO")
	return nil
}

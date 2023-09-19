package workerstd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/illumitacit/gostd/quit"
	"go.uber.org/zap"
	"gocloud.dev/pubsub"
	"google.golang.org/protobuf/proto"
)

type App struct {
	Broker          *Broker
	Logger          *zap.SugaredLogger
	TaskHandler     TaskHandler
	ShutdownTimeout time.Duration

	// ReceiveTaskFn is called to receive a task from the Sub client. Ideally this is not necessary, but because
	// proto.Message is a pointer type, we can't instantiate the struct without knowing what protobuf message we want to
	// unmarshal to.
	ReceiveTaskFn func(ctx context.Context, subClt *SubClient) (proto.Message, *pubsub.Message, error)

	// CloseFn is called on close. contain Any additional close routine should be handled in the custom close function passed in here.
	CloseFn func() error
}

// RunWithSignalHandler runs a worker process described by the App struct in the background, and implements a signal
// handler in the foreground that traps the INT and TERM signals. When the INT or TERM signal is sent to the process,
// this will start a graceful shutdown of the worker app, waiting up to ShutdownTimeout duration for all the worker
// threads to stop processing.
func RunWithSignalHandler(app *App) (returnErr error) {
	subscription, err := NewSubClient(app.Logger, app.Broker, context.Background())
	if err != nil {
		return err
	}

	app.Logger.Infof("Reading tasks from broker")

	// Start the worker in the background so that we can handle shutdown signals gracefully.
	errCh := make(chan error, 1)
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
			err := receiveMsgWithTimeout(app, subscription, timeout, cancel)
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
		app.Logger.Debug("Shutting down worker")
		quit.BroadcastShutdown()
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
	case <-quit.GetQuitChannel():
		app.Logger.Infof("All services gracefully shutdown.")
	}

	loopErr := <-errCh
	return loopErr
}

func receiveMsgWithTimeout(
	app *App, subscription *SubClient,
	timeout context.Context, cancel context.CancelFunc,
) (returnErr error) {
	defer cancel()

	task, msg, err := app.ReceiveTaskFn(timeout, subscription)
	if err != nil {
		// TODO: differentiate fatal error from ignorable errors
		if !errors.Is(err, context.DeadlineExceeded) {
			app.Logger.Errorf("Error receiving message from broker: %s", err)
		}
		return err
	}
	// TODO: spawn goroutine for handling the task, but with work pooling to prevent overflowing.
	if err := app.TaskHandler.HandleTaskMsg(task, msg); err != nil {
		// NOTE: we don't halt on task errors so that the worker continues to process other messages.
		app.Logger.Errorf("Error processing task TODO from broker: %s", err)
	}
	app.Logger.Infof("Successfully processed task TODO")
	return nil
}

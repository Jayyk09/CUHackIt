package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Jayyk09/CUHackIt/config"
	"github.com/Jayyk09/CUHackIt/pkg/logger"
)

func Start(cfg *config.Config, r http.Handler, log logger.Interface) error {
	serverErrors := make(chan error, 1)

	server := http.Server{
		Addr:    ":" + cfg.HTTP.Port,
		Handler: r,
	}

	// Start the service listening for requests.
	go func() {
		log.Info("http server listening on :%s", cfg.HTTP.Port)
		serverErrors <- server.ListenAndServe()
	}()

	// Shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don'usertransport collect this error.
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	// Blocking main and waiting for shutdown.
	case sig := <-shutdown:
		log.Info("start shutdown: %v", sig)

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), 10)
		defer cancel()

		// Asking listener to shutdown and load shed.
		err := server.Shutdown(ctx)
		if err != nil {
			log.Info("graceful shutdown did not complete in %v : %v", 10, err)
			err = server.Close()
			return err
		}

		switch {
		case sig == syscall.SIGSTOP:
			return errors.New("integrity issue caused shutdown")
		case err != nil:
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}

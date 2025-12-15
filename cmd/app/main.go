package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"techno/internal/app"
)

func main() {
	ctx := context.Background()

	application, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("failed to init app: %s", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	errCh := make(chan error, 1)
	go func() {
		if err := application.Run(); err != nil {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case <-quit:
	case err := <-errCh:
		if err != nil {
			log.Printf("Application error: %s", err)
			os.Exit(1)
		}
	}

	shutdownCtx := context.Background()
	if err := application.Stop(shutdownCtx); err != nil {
		log.Fatalf("failed to stop app: %v", err)
	}
}

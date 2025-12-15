package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"techno/internal/app"
)

func main() {
	ctx := context.Background()

	application, err := app.NewWorkerApp(ctx)
	if err != nil {
		log.Fatalf("failed to init worker app: %s", err)
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
		log.Println("Received shutdown signal")
	case err := <-errCh:
		if err != nil {
			log.Printf("Worker error: %s", err)
			os.Exit(1)
		}
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := application.Stop(shutdownCtx); err != nil {
		log.Fatalf("failed to stop worker app: %v", err)
	}

	log.Println("Worker stopped gracefull")
}

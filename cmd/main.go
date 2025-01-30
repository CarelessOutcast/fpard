package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	caddy "github.com/caddyserver/caddy/v2/cmd"
	_ "github.com/caddyserver/caddy/v2/modules/standard"

	fpard "github.com/carelessoutcast/fpard/cmd/fpard"
)

func main() {
	os.Args = []string{"fpard", "run"}

	// Create context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Capture OS interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	go func() {
		// Start the backend service
		log.Println("[Main] Starting backend service (fpard)...")
		fpard.Start(ctx)
		log.Println("[Main] Backend service stopped.")
	}()

	go func() {
		// Start caddy
		log.Println("[Main] Starting Caddy...")
		caddy.Main()
		log.Println("[Main] Caddy has exited.")
		cancel() // Ensure the backend also stops when Caddy exits
	}()

	// wait for interrupt
	sig := <-sigChan
	log.Printf("[Main] Received signal: %v. Initiating shutdown...\n", sig)
	cancel()

	// wait for cleanup
	<-ctx.Done()
	log.Println("[Main] Application shutdown complete.")
}

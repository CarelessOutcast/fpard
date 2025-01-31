package fpard

import (
	"context"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/carelessoutcast/fpard/cmd/backend/handlers"
)

var addr = flag.String("address", "localhost:8081", "The listening address of the service")

func Start(ctx context.Context) {
	flag.Parse()
	log.Println("[fpard] Starting backend service...")

	err := run(ctx, *addr)
	if err != nil {
		log.Fatal(err)
		log.Fatalf("[fpard] Fatal error: %v", err)
	}

	log.Println("[fpard] Backend service shutdown complete.")
}

func run(ctx context.Context, addr string) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/upload", handlers.UploadPDF)

	mux.Handle("/",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientAddr := r.Header.Get("X-Forwarded-For")
			log.Printf("[fpard] Request received: %s -> %s -> %s", clientAddr, r.RemoteAddr, r.URL)
			_, _ = w.Write([]byte("Welcome to fpard!"))
		}),
	)

	server := &http.Server{
		Addr:        addr,
		Handler:     mux,
		IdleTimeout: time.Minute,
		ReadTimeout: 30 * time.Second,
	}
	go func() {
		log.Printf("[fpard] Listening on %s...\n", server.Addr)
		err := server.ListenAndServe()
		if err == http.ErrServerClosed {
			err = nil
		}
	}()

	<-ctx.Done()
	log.Println("[fpard] Shutdown signal received. Stopping server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("[fpard] Error during shutdown: %v", err)
	} else {
		log.Println("[fpard] Server shut down successfully.")
	}

	return nil
}

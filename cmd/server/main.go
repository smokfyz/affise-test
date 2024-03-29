package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/smokfyz/affise-test/pkg/id"
	"github.com/smokfyz/affise-test/pkg/log"
	"github.com/smokfyz/affise-test/pkg/urls"
)

func setupLogger() {
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel != "" {
		log.Init(logLevel)
	}
}

func main() {
	setupLogger()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", port),
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		<-ctx.Done()
		log.Info.Print("shutting down server")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Error.Printf("server shutdown failed: %v", err)
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		requestID := id.New()
		ctx := context.WithValue(r.Context(), urls.RequestIDKey, requestID)

		urls.IndexHandler(w, r.WithContext(ctx))
	})

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Error.Fatalf("server crashed: %v", err)
	}
}

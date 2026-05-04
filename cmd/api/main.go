package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/quangho/yoga-be/internal/adapter/firestoredb"
	"github.com/quangho/yoga-be/internal/config"
	"github.com/quangho/yoga-be/internal/transport/rest"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		slog.Error("configuration error", "error", err)
		os.Exit(1)
	}

	client, err := firestoredb.NewClient(ctx, cfg)
	if err != nil {
		slog.Error("firestore connection failed", "error", err)
		os.Exit(1)
	}
	defer client.Close()

	router := rest.NewRouter(firestoredb.NewRegistry(client))
	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		slog.Info("api listening", "port", cfg.Port, "project_id", cfg.ProjectID, "database_id", cfg.FirestoreDatabase)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("server shutdown failed", "error", err)
	}
}

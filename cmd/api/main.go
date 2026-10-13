package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	"github.com/Vitorpignataro/go-job-orchestrator/internal/config"
	"github.com/Vitorpignataro/go-job-orchestrator/internal/database"
	"github.com/go-chi/chi/v5"
)

func main(){
	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, nil),
	)

	cfg := config.Load()

	db, err := database.Open(cfg.DatabaseDSN())
	if err != nil {
		logger.Error(
			"failed to connect to database",
			"error",
			err,
		)
		os.Exit(1)
	}

	defer db.Close()

	router := chi.NewRouter()

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(http.StatusOK)

		_ = json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		})
	})

	logger.Info(
		"api started",
		"port",
		cfg.HTTPPort,
	)

	err = http.ListenAndServe(
		":"+cfg.HTTPPort,
		router,
	)

	if err != nil {
		logger.Error(
			"http server stopped",
			"error",
			err,
		)

		os.Exit(1)
	}
}
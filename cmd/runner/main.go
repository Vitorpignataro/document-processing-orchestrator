package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/Vitorpignataro/go-job-orchestrator/internal/config"
	"github.com/Vitorpignataro/go-job-orchestrator/internal/database"
)

func main(){
	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, nil),
	)

	cfg := config.Load()

	db, err := database.Open(cfg.DatabaseDSN())
	if err != nil{
		logger.Error(
			"failed to connect to database",
			"error",
			err,
		)
		os.Exit(1)
	}

	defer db.Close()
	logger.Info("runner started")
	
	for {
		logger.Info("looking for available jobs")
		time.Sleep(5 * time.Second)
	}
}
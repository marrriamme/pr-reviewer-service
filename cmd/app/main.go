package main

import (
	"log"
	"os"

	"github.com/marrria_mme/pr-reviewer-service/config"
	"github.com/marrria_mme/pr-reviewer-service/internal/app"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("config error: %v", err)
		os.Exit(1)
	}

	application, err := app.NewApp(cfg)
	if err != nil {
		log.Fatalf("failed to create app: %v", err)
	}

	application.Run()
}

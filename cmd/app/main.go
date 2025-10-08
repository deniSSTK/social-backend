package main

import (
	"log"
	"os"
	"social-backend/internal/infrastructure/http"
	"social-backend/internal/infrastructure/logger"
	"social-backend/scripts"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}

	runMode := os.Getenv("RUN_MODE")
	if runMode == "" {
		panic("Set RUN_MODE env variable")
	}

	if err := logger.Init(runMode); err != nil {
		panic(err)
	}

	scripts.Migrate()

	http.StartServer()
}

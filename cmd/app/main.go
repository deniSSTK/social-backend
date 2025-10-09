package main

import (
	"social-backend/internal/infrastructure/http"
	"social-backend/internal/infrastructure/logger"
	"social-backend/scripts"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}

	if err := logger.Init(); err != nil {
		panic(err)
	}

	scripts.Migrate()

	http.StartServer()
}

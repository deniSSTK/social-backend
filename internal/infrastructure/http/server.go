package http

import (
	"fmt"
	"log"
	"os"
	"social-backend/internal/infrastructure/db"
	"social-backend/internal/infrastructure/db/repository"
	"social-backend/internal/infrastructure/http/handler"
	"social-backend/internal/usecase"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func StartServer() {
	dbPool := db.ConnectDB()
	defer dbPool.Close()

	userRepo := repository.NewUserRepository(dbPool)

	userUC := usecase.NewUserUsecase(userRepo)

	userHandler := handler.NewUserHandler(userUC)

	r := gin.Default()

	frontendUrl := os.Getenv("FRONTEND_URL")
	if frontendUrl == "" {
		log.Fatal("Environment variable FRONTEND_URL is not set")
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{frontendUrl},
		AllowMethods:     []string{"POST", "GET", "OPTIONS", "DELETE", "PATCH", "PUT"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
	}))

	api := r.Group("/api")

	userHandler.RegisterRoutes(api)

	port := ":8080"

	if err := r.Run(port); err != nil {
		fmt.Println("Failed to start server:", err)
		panic(err)
	}
}

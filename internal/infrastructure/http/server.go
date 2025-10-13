package http

import (
	"fmt"
	"os"
	"social-backend/internal/infrastructure/auth"
	"social-backend/internal/infrastructure/db"
	"social-backend/internal/infrastructure/db/repository"
	"social-backend/internal/infrastructure/http/handler"
	"social-backend/internal/infrastructure/imgbb"
	"social-backend/internal/infrastructure/logger"
	"social-backend/internal/usecase"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func StartServer() {
	dbPool := db.ConnectDB()
	defer dbPool.Close()

	log := logger.Get().Sugar()

	//TODO create error "environment variable not set"

	jwtKey := os.Getenv("JWT_SECRET")
	if jwtKey == "" {
		log.Fatal("JWT_SECRET environment variable not set")
	}

	imgBBApiKey := os.Getenv("IMGBB_API_KEY")
	if imgBBApiKey == "" {
		log.Fatal("IMG_BAPI_KEY environment variable not set")
	}

	imgBBApiUrl := os.Getenv("IMGBB_API_URL")
	if imgBBApiUrl == "" {
		log.Fatal("IMGBB_API_URL environment variable not set")
	}

	jwtService := auth.NewJWTService(jwtKey)
	ImgBBService := imgbb.NewImgBBService(imgBBApiUrl, imgBBApiKey)

	baseRepo := repository.NewBaseRepo(dbPool)

	userRepo := repository.NewUserRepository(dbPool)
	postRepo := repository.NewPostRepository(dbPool)
	imageRepo := repository.NewImageRepository(dbPool)
	hashtagRepo := repository.NewHashtagRepository(dbPool)

	userUC := usecase.NewUserUsecase(userRepo)
	postUC := usecase.NewPostUsecase(baseRepo, postRepo, imageRepo, hashtagRepo, ImgBBService)

	userHandler := handler.NewUserHandler(userUC, jwtService)
	postHandler := handler.NewPostHandler(postUC, jwtService)

	r := gin.Default()

	frontendUrl := os.Getenv("FRONTEND_URL")
	if frontendUrl == "" {
		log.Fatal("FRONTEND_URL environment variable not set")
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{frontendUrl},
		AllowMethods:     []string{"POST", "GET", "OPTIONS", "DELETE", "PATCH", "PUT"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
	}))

	api := r.Group("/api")

	userHandler.RegisterRoutes(api)
	postHandler.RegisterRoutes(api)

	port := ":8080"

	if err := r.Run(port); err != nil {
		fmt.Println("Failed to start server:", err)
		panic(err)
	}
}

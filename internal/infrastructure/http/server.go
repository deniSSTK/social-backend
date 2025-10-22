package http

import (
	"fmt"
	"os"
	"social-backend/internal/infrastructure/auth"
	"social-backend/internal/infrastructure/db"
	"social-backend/internal/infrastructure/db/repository"
	"social-backend/internal/infrastructure/errors"
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

	jwtKey := os.Getenv("JWT_SECRET")
	if jwtKey == "" {
		log.Fatal(errors.EnvironmentVariableNotSet.Error() + "JWT_SECRET")
	}

	imgBBApiKey := os.Getenv("IMGBB_API_KEY")
	if imgBBApiKey == "" {
		log.Fatal(errors.EnvironmentVariableNotSet.Error() + "IMGBB_API_KEY")
	}

	imgBBApiUrl := os.Getenv("IMGBB_API_URL")
	if imgBBApiUrl == "" {
		log.Fatal(errors.EnvironmentVariableNotSet.Error() + "IMGBB_API_URL")
	}

	baseRepo := repository.NewBaseRepo(dbPool)

	userRepo := repository.NewUserRepository(dbPool)
	postRepo := repository.NewPostRepository(dbPool)
	imageRepo := repository.NewImageRepository(dbPool)
	followRepo := repository.NewFollowRepository(dbPool)
	hashtagRepo := repository.NewHashtagRepository(dbPool)

	jwtService := auth.NewJWTService(jwtKey)
	userService := auth.NewUserService(userRepo)
	authServie := auth.NewAuthService(jwtService, userService)
	ImgBBService := imgbb.NewImgBBService(imgBBApiKey, imgBBApiUrl)

	userUC := usecase.NewUserUsecase(userRepo)
	postUC := usecase.NewPostUsecase(baseRepo, postRepo, imageRepo, hashtagRepo, ImgBBService)
	followUC := usecase.NewFollowUsecase(baseRepo, followRepo)
	hashtagUC := usecase.NewHashtagUsecase(hashtagRepo)

	userHandler := handler.NewUserHandler(userUC, authServie)
	postHandler := handler.NewPostHandler(postUC, authServie)
	followHandler := handler.NewFollowHandler(followUC, authServie)
	hashtagHandler := handler.NewHashtagHandler(hashtagUC, authServie)

	r := gin.Default()

	frontendUrl := os.Getenv("FRONTEND_URL")
	if frontendUrl == "" {
		log.Fatal(errors.EnvironmentVariableNotSet.Error() + "FRONTEND_URL")
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
	followHandler.RegisterRoutes(api)
	hashtagHandler.RegisterRoutes(api)

	port := ":8080"

	if err := r.Run(port); err != nil {
		fmt.Println("Failed to start server:", err)
		panic(err)
	}
}

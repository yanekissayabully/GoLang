package main

import (
	"practice-7/internal/controller/http/v1"
	"practice-7/internal/middleware"
	"practice-7/internal/usecase"
	"practice-7/internal/usecase/repo"
	"practice-7/pkg/postgres"

	"github.com/gin-gonic/gin"
)

func main() {
	db := postgres.New()
	
	userRepo := repo.NewUserRepo(db)
	userUseCase := usecase.NewUserUseCase(userRepo)
	
	limiter := middleware.NewRateLimiter(10, 60)
	
	r := gin.Default()
	r.Use(middleware.RateLimitMiddleware(limiter))
	
	v1Group := r.Group("/v1")
	v1.NewUserRoutes(v1Group, userUseCase)
	
	r.Run(":8090")
}
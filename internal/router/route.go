package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/sung2708/shorten_url/internal/config"
	"github.com/sung2708/shorten_url/internal/handle"
	"github.com/sung2708/shorten_url/internal/middleware"
	"github.com/sung2708/shorten_url/internal/repository"
	"github.com/sung2708/shorten_url/internal/service"
)

func Setup(cfg *config.Config, db *gorm.DB) *gin.Engine {

	r := gin.Default()

	userRepo := repository.NewUserRepository(db)

	userService := service.NewUserService(userRepo, cfg.JWTSecret)

	userHandler := handle.NewUserHandler(userService)

	urlRepo := repository.NewURLRepository(db)

	urlService := service.NewUrlService(urlRepo)
	urlHandler := handle.NewURLHandler(urlService)

	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/register", userHandler.Register)
		authRoutes.POST("/login", userHandler.Login)
	}

	urlRoutes := r.Group("/api/v1")
	{
		urlRoutes.POST("/shorten", urlHandler.Shorten)
	}

	privateRoutes := r.Group("/api/v1/links")
	privateRoutes.Use(middleware.AuthMiddleware(cfg.JWTSecret)) // <-- Dùng Middleware
	{
		// privateRoutes.GET("/", urlHandler.GetMyLinks)
		// privateRoutes.DELETE("/:id", urlHandler.DeleteLink)
	}
	r.GET("/:code", urlHandler.Resolve)

	return r
}

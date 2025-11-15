package router

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/sung2708/shorten_url/internal/config"
	"github.com/sung2708/shorten_url/internal/handle"
	"github.com/sung2708/shorten_url/internal/middleware"
	"github.com/sung2708/shorten_url/internal/repository"
	"github.com/sung2708/shorten_url/internal/service"
	"gorm.io/gorm"
)

func RequiredAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		useridValue, _ := c.Get("user_id")
		if useridValue == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Login must require"})
			return
		}
		c.Next()
	}
}

func Setup(cfg *config.Config, db *gorm.DB, rdb *redis.Client) *gin.Engine {

	r := gin.Default()
	r.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://ui.tinyr.site", "http://localhost:5173"},
		AllowMethods:     []string{"POST", "GET", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	userRepo := repository.NewUserRepository(db)
	otpRepo := repository.NewOTPRepository(db)

	notificationService := service.NewNotificationService(
		otpRepo,
		userRepo,
		cfg.Email.SMTPHost,
		cfg.Email.SMTPUser,
		cfg.Email.SMTPPass,
		cfg.Email.SenderEmail,
		cfg.Email.SMTPPort,
	)

	userService := service.NewUserService(userRepo, otpRepo, notificationService, cfg.JWTSecret)
	userHandler := handle.NewUserHandler(userService, notificationService)

	urlRepo := repository.NewURLRepository(db, rdb)
	urlService := service.NewUrlService(urlRepo)
	urlHandler := handle.NewURLHandler(urlService)

	// API routes - must be registered before catch-all route
	apiV1 := r.Group("/api/v1")
	{
		// Auth routes
		authRoutes := apiV1.Group("/auth")
		{
			authRoutes.POST("/register", userHandler.Register)
			authRoutes.POST("/login", userHandler.Login)
			authRoutes.POST("/verify-code", userHandler.VerifyCode)
		}

		// URL routes
		apiV1.POST("/shorten", urlHandler.Shorten)

		// Private routes (require authentication)
		privateRoutes := apiV1.Group("/links")
		privateRoutes.Use(RequiredAuthMiddleware())
		{
			privateRoutes.GET("/", urlHandler.GetMyLinks)
			privateRoutes.DELETE("/:code", urlHandler.DeleteLink)
		}
	}

	// Catch-all route for short code resolution (must be last)
	r.GET("/:code", urlHandler.Resolve)
	return r
}

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
		AllowOrigins:     []string{"https://ui.tinyr.site", "http://localhost:5173", "https://tinyr.site"}, // Đã thêm tinyr.site
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

	// Health check endpoint (for debugging deployment)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"routes": []string{
				"POST /api/v1/auth/register",
				"POST /api/v1/auth/login",
				"POST /api/v1/auth/verify-code",
				"POST /api/v1/shorten",
				"GET /api/v1/links",
				"DELETE /api/v1/links/:code",
			},
		})
	})

	// --- ĐĂNG KÝ CÁC ROUTE API TRỰC TIẾP TRÊN ROUTER GỐC (r) ---

	// Auth routes
	r.POST("/api/v1/auth/register", userHandler.Register)
	r.POST("/api/v1/auth/login", userHandler.Login)
	r.POST("/api/v1/auth/verify-code", userHandler.VerifyCode)

	// URL routes
	r.POST("/api/v1/shorten", urlHandler.Shorten)

	// Private routes (require authentication)
	// Sử dụng Subrouter cho các route yêu cầu xác thực
	privateRoutes := r.Group("/api/v1/links")
	privateRoutes.Use(RequiredAuthMiddleware())
	{
		privateRoutes.GET("/", urlHandler.GetMyLinks)
		privateRoutes.DELETE("/:code", urlHandler.DeleteLink)
	}

	// API routes (Các route cũ đã xóa)

	// Catch-all route for short code resolution (must be last)
	r.GET("/:code", urlHandler.Resolve)
	return r
}

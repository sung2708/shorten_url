package router

import (
	"github.com/gin-gonic/gin"
	"github.com/sung2708/shorten_url/internal/handle"
	"github.com/sung2708/shorten_url/internal/repository"
	"github.com/sung2708/shorten_url/internal/service"
	"gorm.io/gorm"
)

func Router(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	rp := repository.NewURLRepository(db)
	svc := service.NewUrlService(rp)
	handler := handle.NewURLHandler(svc)

	v1 := r.Group("/api/v1")
	{
		v1.POST("/shorten", handler.Shorten)
	}
	r.GET("/code", handler.Resolve)
	return r
}

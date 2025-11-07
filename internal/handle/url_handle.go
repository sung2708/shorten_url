package handle

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sung2708/shorten_url/internal/service"
)

type URLHandle interface {
	Shorten(url string) (string, error)
	Resolve(url string) (string, error)
}

type URLHandleImpl struct {
	service *service.UrlServiceImpl
}

type NewURLRequest struct {
	URL string `json:"url"`
}

func NewURLHandler(service *service.UrlServiceImpl) *URLHandleImpl {
	return &URLHandleImpl{service: service}
}

func (h *URLHandleImpl) Shorten(ctx *gin.Context) {
	var req NewURLRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	url, err := h.service.Shorten(req.URL)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"url": ctx.Request.Host + "/" + url.ShortCode,
	})
}

func (h *URLHandleImpl) Resolve(ctx *gin.Context) {
	code := ctx.Param("code")
	url, err := h.service.GetById(code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Redirect(http.StatusFound, url.LongURl)
}

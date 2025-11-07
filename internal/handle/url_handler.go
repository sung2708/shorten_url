package handle

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sung2708/shorten_url/internal/service"
)

type URLHandleImpl struct {
	urlService service.UrlService
}
type NewURLRequest struct {
	URL string `json:"url"`
}

func NewURLHandler(service service.UrlService) *URLHandleImpl {
	return &URLHandleImpl{urlService: service}
}

func (handler *URLHandleImpl) Shorten(ctx *gin.Context) {
	var req NewURLRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userID *uint = nil

	useridval, exists := ctx.Get("user_id")

	if exists {
		uid := useridval.(uint)
		userID = &uid
	}
	url, err := handler.urlService.Shorten(req.URL, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"url": url})
}

func (handler *URLHandleImpl) Resolve(ctx *gin.Context) {
	code := ctx.Param("code")
	url, err := handler.urlService.GetById(code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Redirect(http.StatusFound, url.LongURl)
}

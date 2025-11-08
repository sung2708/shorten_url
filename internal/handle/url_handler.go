package handle

import (
	"net/http"
	"strings"

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
	host := strings.TrimPrefix(ctx.Request.Host, "www.")
	ctx.JSON(http.StatusOK, gin.H{"url": host + "/" + url.ShortCode})
}

func (handler *URLHandleImpl) Resolve(ctx *gin.Context) {
	code := ctx.Param("code")
	url, err := handler.urlService.GetById(code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Redirect(http.StatusFound, url.LongURL)
}

func (handler *URLHandleImpl) GetMyLinks(ctx *gin.Context) {
	useridValue, _ := ctx.Get("userID")
	userID := useridValue.(uint)

	links, err := handler.urlService.FindByUserID(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	ctx.JSON(http.StatusOK, gin.H{"links": links})
}

func (handler *URLHandleImpl) DeleteLink(ctx *gin.Context) {

	useridValue, _ := ctx.Get("userID")
	userID := useridValue.(uint)

	shortCode := ctx.Param("code")

	err := handler.urlService.DeleteLink(shortCode, userID)

	if err != nil {
		if err.Error() == "user is not own link" {
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else if err.Error() == "url not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	ctx.Status(http.StatusNoContent)
}

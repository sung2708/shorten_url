package handle

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sung2708/shorten_url/internal/model"
	"github.com/sung2708/shorten_url/internal/service"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (handler *UserHandler) Register(ctx *gin.Context) {
	var input model.User

	if err := ctx.ShouldBind(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(input.Password) < 8 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "password is too short"})
		return
	}
	createUser, err := handler.userService.Register(input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": createUser})
}

func (handler *UserHandler) Login(ctx *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := ctx.ShouldBind(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	if len(input.Password) < 8 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "password is too short"})
	}
	token, err := handler.userService.Login(input.Email, input.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}
	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

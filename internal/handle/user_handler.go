package handle

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sung2708/shorten_url/internal/model"
	"github.com/sung2708/shorten_url/internal/service"
)

type UserHandler struct {
	userService         service.UserService
	notificationService service.NotificationService
}

func NewUserHandler(userService service.UserService, notificationService service.NotificationService) *UserHandler {
	return &UserHandler{
		userService:         userService,
		notificationService: notificationService,
	}
}

func (h *UserHandler) Register(ctx *gin.Context) {
	var input struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := model.User{Name: input.Name, Email: input.Email, Password: input.Password}

	token, createdUser, err := h.userService.Register(user) // JWT issued (is_verified = false)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// OTP is already sent by userService.Register
	ctx.JSON(http.StatusCreated, gin.H{
		"user":    createdUser,
		"token":   token, // FE auto-login
		"message": "OTP sent to your email",
	})
}

func (h *UserHandler) Login(ctx *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, user, err := h.userService.Login(input.Email, input.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// If user is not verified, include message about OTP
	response := gin.H{"token": token, "user": user}
	if !user.IsActive {
		response["message"] = "OTP sent to your email. Please verify your account."
	}

	ctx.JSON(http.StatusOK, response)
}

// VerifyCode handles POST /api/auth/verify-code
// Receives 6-digit OTP code and current JWT token
// Returns new JWT token with is_verified: true
func (h *UserHandler) VerifyCode(ctx *gin.Context) {
	// Get userID from JWT token (set by auth middleware)
	userIDValue, exists := ctx.Get("user_id")
	if !exists || userIDValue == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	userIDPtr, ok := userIDValue.(*uint)
	if !ok || userIDPtr == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id"})
		return
	}
	userID := *userIDPtr

	// Get OTP code from request body (JSON)
	var input struct {
		Code string `json:"code" binding:"required,len=6"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "6-digit code is required"})
		return
	}

	// Verify OTP and get new token with full permissions
	token, user, err := h.userService.VerifyAccount(userID, input.Code)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Account verified successfully",
		"user":    user,
		"token":   token, // New JWT token with is_verified: true
	})
}

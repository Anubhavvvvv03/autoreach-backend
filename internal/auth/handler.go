package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/autoreach-backend/internal/config"
	"github.com/yourusername/autoreach-backend/internal/dto/request"
	"github.com/yourusername/autoreach-backend/internal/dto/response"
)


func SignupHandler(c *gin.Context) {
	var req request.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, false, err.Error(), nil)
		return
	}

	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		response.JSON(c, http.StatusInternalServerError, false, "Failed to hash password", nil)
		return
	}

	// Check if user already exists
	var existingUser User
	if err := config.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		response.JSON(c, http.StatusConflict, false, "Email already registered", nil)
		return
	}

	user := User{
		Email:    req.Email,
		Password: hashedPassword,
		FullName: req.FullName,
	}

	result := config.DB.Create(&user)
	if result.Error != nil {
		response.JSON(c, http.StatusInternalServerError, false, "Failed to create user", nil)
		return
	}

	token, err := GenerateJWT(user.ID, user.Email)
	if err != nil {
		response.JSON(c, http.StatusInternalServerError, false, "Failed to generate token", nil)
		return
	}

	response.JSON(c, http.StatusCreated, true, "User created successfully", response.SignupResponse{
		Token: token,
		User: response.UserResponse{
			ID:       user.ID,
			Email:    user.Email,
			FullName: user.FullName,
		},
	})
}

func LoginHandler(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, false, err.Error(), nil)
		return
	}

	var user User
	result := config.DB.Where("email = ?", req.Email).First(&user)
	if result.Error != nil {
		response.JSON(c, http.StatusUnauthorized, false, "Invalid email or password", nil)
		return
	}

	if !CheckPasswordHash(req.Password, user.Password) {
		response.JSON(c, http.StatusUnauthorized, false, "Invalid email or password", nil)
		return
	}

	token, err := GenerateJWT(user.ID, user.Email)
	if err != nil {
		response.JSON(c, http.StatusInternalServerError, false, "Failed to generate token", nil)
		return
	}

	response.JSON(c, http.StatusOK, true, "Login successful", response.LoginResponse{
		Token: token,
		User: response.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FullName:  user.FullName,
			AvatarUrl: user.AvatarUrl,
		},
	})
}

func MeHandler(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.JSON(c, http.StatusUnauthorized, false, "Unauthorized", nil)
		return
	}

	var user User
	if err := config.DB.First(&user, "id = ?", userID).Error; err != nil {
		response.JSON(c, http.StatusNotFound, false, "User not found", nil)
		return
	}

	response.JSON(c, http.StatusOK, true, "User fetched successfully", response.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FullName:  user.FullName,
		AvatarUrl: user.AvatarUrl,
	})
}

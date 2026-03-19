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
		Name:     req.Name,
	}

	result := config.DB.Create(&user)
	if result.Error != nil {
		response.JSON(c, http.StatusInternalServerError, false, "Failed to create user", nil)
		return
	}

	response.JSON(c, http.StatusCreated, true, "User created successfully", response.SignupResponse{
		Name:  user.Name,
		Email: user.Email,
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
		Email: user.Email,
	})
}

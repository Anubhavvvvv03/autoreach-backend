package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/autoreach-backend/internal/config"
	"github.com/yourusername/autoreach-backend/pkg/response"
)

type SignupRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func SignupHandler(c *gin.Context) {
	var req SignupRequest
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

	response.JSON(c, http.StatusCreated, true, "User created successfully", gin.H{
		"name":  user.Name,
		"email": user.Email,
	})
}

func LoginHandler(c *gin.Context) {
	var req LoginRequest
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

	response.JSON(c, http.StatusOK, true, "Login successful", gin.H{"token": token})
}

package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/yourusername/autoreach-backend/internal/auth"
	"github.com/yourusername/autoreach-backend/internal/config"
	"github.com/yourusername/autoreach-backend/internal/dto/response"
	"github.com/yourusername/autoreach-backend/internal/message"
	"github.com/yourusername/autoreach-backend/internal/profile"
	"github.com/yourusername/autoreach-backend/internal/resume"
)

func SetupRouter() *gin.Engine {
    r := gin.Default()

    // Add CORS middleware
    r.Use(cors.New(cors.Config{
        AllowOrigins:     config.AppConfig.AllowedOrigins,
        AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
    }))

    r.GET("/api/v1/health", func(c *gin.Context) {
        response.JSON(c, 200, true, "AutoReach backend is running", nil)
    })

    authGroup := r.Group("/api/v1/auth")
    {
        authGroup.POST("/signup", auth.SignupHandler)
        authGroup.POST("/login", auth.LoginHandler)
    }

    api := r.Group("/api/v1")
    api.Use(auth.AuthMiddleware())
    {
        api.POST("/generate-message", message.GenerateMessageHandler)
        api.GET("/profile", profile.GetProfileHandler)
        api.POST("/profile", profile.CreateProfileHandler)
        api.PUT("/profile", profile.UpdateProfileHandler)
        api.POST("/resume/upload", resume.UploadResumeHandler)
        api.GET("/resume/:id", resume.GetResumeStatusHandler)
    }

    return r
}

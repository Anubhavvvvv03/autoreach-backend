package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yourusername/autoreach-backend/internal/auth"
	"github.com/yourusername/autoreach-backend/internal/message"
	"github.com/yourusername/autoreach-backend/pkg/response"
)

func SetupRouter() *gin.Engine {
    r := gin.Default()

    r.GET("/", func(c *gin.Context) {
        response.JSON(c, 200, true, "AutoReach backend is running", nil)
    })

    authGroup := r.Group("/api/v1/auth")
    {
        authGroup.POST("/signup", auth.SignupHandler)
        authGroup.POST("/login", auth.LoginHandler)
    }

    api := r.Group("/api")
    api.Use(auth.AuthMiddleware())
    {
        api.POST("/generate-message", message.GenerateMessageHandler)
    }

    return r
}

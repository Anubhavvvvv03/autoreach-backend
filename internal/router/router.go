package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yourusername/autoreach-backend/internal/auth"
	"github.com/yourusername/autoreach-backend/internal/message"
)

func SetupRouter() *gin.Engine {
    r := gin.Default()

    r.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "AutoReach backend is running"})
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

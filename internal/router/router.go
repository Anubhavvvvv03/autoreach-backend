package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yourusername/autoreach-backend/internal/handler"
)

func SetupRouter() *gin.Engine {
    r := gin.Default()

    r.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "AutoReach backend is running"})
    })

    api := r.Group("/api")
    {
        api.POST("/generate-message", handler.GenerateMessage)
    }

    return r
}

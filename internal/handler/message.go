package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/yourusername/autoreach-backend/internal/model"
	"github.com/yourusername/autoreach-backend/internal/service"
)

func GenerateMessage(c *gin.Context) {
    var req model.MessageRequest
    if err := c.BindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    msg, err := service.GenerateMessage(req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": msg})
}

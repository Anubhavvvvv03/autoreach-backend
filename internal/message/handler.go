package message

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func GenerateMessageHandler(c *gin.Context) {
    var req MessageRequest
    if err := c.BindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    msg, err := GenerateMessageService(req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": msg})
}

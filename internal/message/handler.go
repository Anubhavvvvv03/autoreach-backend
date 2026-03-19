package message

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/yourusername/autoreach-backend/internal/dto/request"
	"github.com/yourusername/autoreach-backend/internal/dto/response"
)

func GenerateMessageHandler(c *gin.Context) {
    var req request.MessageRequest
    if err := c.BindJSON(&req); err != nil {
        response.JSON(c, http.StatusBadRequest, false, err.Error(), nil)
        return
    }

    msg, err := GenerateMessageService(req)
    if err != nil {
        response.JSON(c, http.StatusInternalServerError, false, err.Error(), nil)
        return
    }

    response.JSON(c, http.StatusOK, true, "Message generated successfully", response.MessageResponse{Message: msg})
}

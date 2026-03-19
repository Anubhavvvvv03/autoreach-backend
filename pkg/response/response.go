package response

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func JSON(c *gin.Context, code int, success bool, message string, data interface{}) {
	c.JSON(code, Response{
		Success: success,
		Message: message,
		Data:    data,
	})
}

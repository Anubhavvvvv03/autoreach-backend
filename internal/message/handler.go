package message

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/autoreach-backend/internal/config"
	"github.com/yourusername/autoreach-backend/internal/dashboard"
	"github.com/yourusername/autoreach-backend/internal/dto/request"
	"github.com/yourusername/autoreach-backend/internal/dto/response"
)

func GenerateMessageHandler(c *gin.Context) {
	var req request.MessageRequest
	if err := c.BindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, false, err.Error(), nil)
		return
	}

	userID := c.GetString("userID")
	msg, err := GenerateMessageService(req, userID)
	if err != nil {
		response.JSON(c, http.StatusInternalServerError, false, err.Error(), nil)
		return
	}

	// Fetch the last created message for this user to get its ID
	var lastMsg Message
	config.DB.Where("user_id = ?", userID).Order("created_at desc").First(&lastMsg)

	// Log Activity
	config.DB.Create(&dashboard.Activity{
		UserID: userID,
		Type:   "OUTREACH",
		Title:  "Message generated for " + req.Company,
		Status: "SUCCESS",
	})

	response.JSON(c, http.StatusOK, true, "Message generated successfully", response.MessageResponse{
		MessageID:   lastMsg.ID,
		Text:        msg,
		GeneratedAt: time.Now().Format(time.RFC3339),
	})
}

func GetMessageHistoryHandler(c *gin.Context) {
	userID := c.GetString("userID")
	messages, err := GetMessageHistory(userID)
	if err != nil {
		response.JSON(c, http.StatusInternalServerError, false, "Failed to fetch history", nil)
		return
	}

	var results []gin.H
	for _, m := range messages {
		results = append(results, gin.H{
			"messageId": m.ID,
			"company":   m.Company,
			"role":      m.Role,
			"text":      m.Text,
			"status":    m.Status,
		})
	}

	response.JSON(c, http.StatusOK, true, "History fetched successfully", gin.H{
		"messages": results,
	})
}

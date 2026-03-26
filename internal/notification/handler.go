package notification

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/yourusername/autoreach-backend/internal/config"
	"github.com/yourusername/autoreach-backend/internal/dto/response"
)

func GetNotificationsHandler(c *gin.Context) {
	userID := c.GetString("userID")

	var notifications []Notification
	config.DB.Where("user_id = ? AND is_read = ?", userID, false).Find(&notifications)

	var allNotifications []Notification
	config.DB.Where("user_id = ?", userID).Order("created_at desc").Limit(20).Find(&allNotifications)

	response.JSON(c, http.StatusOK, true, "Notifications fetched", gin.H{
		"unreadCount": len(notifications),
		"list":        allNotifications,
	})
}

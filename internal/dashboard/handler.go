package dashboard

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/autoreach-backend/internal/config"
	"github.com/yourusername/autoreach-backend/internal/dto/response"
)

func GetSummaryHandler(c *gin.Context) {
	userID := c.GetString("userID")

	var totalOutreach int64
	config.DB.Table("messages").Where("user_id = ?", userID).Count(&totalOutreach)

	// Stubbed stats until we have real campaign data, but count is real
	response.JSON(c, http.StatusOK, true, "Summary fetched", gin.H{
		"stats": gin.H{
			"totalOutreach":   totalOutreach,
			"successRate":     24.5, // To be calculated once we track replies
			"pendingTasks":    3,
			"activeCampaigns": 1,
		},
		"performanceHistory": []gin.H{
			{"date": "2024-03-20", "value": totalOutreach}, // Simple stub for now
		},
	})
}

func GetActivityHandler(c *gin.Context) {
	userID := c.GetString("userID")

	var activities []Activity
	config.DB.Where("user_id = ?", userID).Order("created_at desc").Limit(10).Find(&activities)

	response.JSON(c, http.StatusOK, true, "Activity fetched", gin.H{
		"activities": activities,
	})
}

func GetInsightsHandler(c *gin.Context) {
	userID := c.GetString("userID")

	var insights []Insight
	config.DB.Where("user_id = ?", userID).Find(&insights)

	response.JSON(c, http.StatusOK, true, "Insights fetched", gin.H{
		"insights": insights,
	})
}

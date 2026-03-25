package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/yourusername/autoreach-backend/internal/auth"
	"github.com/yourusername/autoreach-backend/internal/config"
	"github.com/yourusername/autoreach-backend/internal/dashboard"
	"github.com/yourusername/autoreach-backend/internal/dto/response"
	"github.com/yourusername/autoreach-backend/internal/message"
	"github.com/yourusername/autoreach-backend/internal/notification"
	"github.com/yourusername/autoreach-backend/internal/profile"
	"github.com/yourusername/autoreach-backend/internal/resume"
)

func SetupRouter() *gin.Engine {
    r := gin.Default()

    // Add CORS middleware
    r.Use(cors.New(cors.Config{
        AllowOrigins:     config.AppConfig.AllowedOrigins,
        AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
    }))

    r.GET("/api/v1/health", func(c *gin.Context) {
        response.JSON(c, 200, true, "AutoReach backend is running", nil)
    })

    authGroup := r.Group("/api/v1/auth")
    {
        authGroup.POST("/signup", auth.SignupHandler)
        authGroup.POST("/login", auth.LoginHandler)
    }

    api := r.Group("/api/v1")
    api.Use(auth.AuthMiddleware())
    {
        api.GET("/auth/me", auth.MeHandler)

        // Dashboard
        api.GET("/dashboard/summary", dashboard.GetSummaryHandler)
        api.GET("/dashboard/activity", dashboard.GetActivityHandler)
        api.GET("/dashboard/insights", dashboard.GetInsightsHandler)

        api.POST("/generate/message", message.GenerateMessageHandler)
        api.GET("/generate/history", message.GetMessageHistoryHandler)
        api.GET("/profile", profile.GetProfileHandler)
        api.POST("/profile", profile.CreateProfileHandler)
        api.PUT("/profile", profile.UpdateProfileHandler)
        api.GET("/profile/sync-status", profile.SyncStatusHandler)
        api.POST("/resume/upload", resume.UploadResumeHandler)
        api.GET("/resume/jobs/:jobId", resume.GetResumeStatusHandler)
        api.GET("/resume/capabilities", func(c *gin.Context) {
            capabilities := []gin.H{
                {"title": "Skill Extraction", "description": "AI-powered skill identification"},
                {"title": "Experience Summarization", "description": "Structured experience mapping"},
            }
            response.JSON(c, 200, true, "Capabilities fetched", gin.H{"capabilities": capabilities})
        })

        // Notifications
        api.GET("/notifications", notification.GetNotificationsHandler)
    }

    return r
}

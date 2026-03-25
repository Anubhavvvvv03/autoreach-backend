package profile

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/autoreach-backend/internal/dto/request"
	"github.com/yourusername/autoreach-backend/internal/dto/response"
)

func GetProfileHandler(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		response.JSON(c, http.StatusUnauthorized, false, "Unauthorized", nil)
		return
	}

	p, err := GetProfileByUserID(userID)
	if err != nil {
		response.JSON(c, http.StatusInternalServerError, false, "Failed to fetch profile", nil)
		return
	}

	if p == nil {
		response.JSON(c, http.StatusNotFound, false, "Profile not found", nil)
		return
	}

	response.JSON(c, http.StatusOK, true, "Profile fetched successfully", p)
}

func CreateProfileHandler(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		response.JSON(c, http.StatusUnauthorized, false, "Unauthorized", nil)
		return
	}

	var req request.UpsertProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, false, err.Error(), nil)
		return
	}

	p, err := CreateProfile(userID, req)
	if err != nil {
		response.JSON(c, http.StatusBadRequest, false, err.Error(), nil)
		return
	}

	response.JSON(c, http.StatusCreated, true, "Profile created successfully", p)
}

func UpdateProfileHandler(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		response.JSON(c, http.StatusUnauthorized, false, "Unauthorized", nil)
		return
	}

	var req request.UpsertProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, false, err.Error(), nil)
		return
	}

	p, err := UpdateProfile(userID, req)
	if err != nil {
		// Use 404 for not found, 400 for other business logic errors
		status := http.StatusBadRequest
		if err.Error() == "profile not found" {
			status = http.StatusNotFound
		}
		response.JSON(c, status, false, err.Error(), nil)
		return
	}

	response.JSON(c, http.StatusOK, true, "Profile updated successfully", p)
}

func SyncStatusHandler(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		response.JSON(c, http.StatusUnauthorized, false, "Unauthorized", nil)
		return
	}

	prof, _ := GetProfileByUserID(userID)

	linkedInStatus := "PENDING"
	if prof != nil && prof.SocialLinks.LinkedIn != "" {
		linkedInStatus = "SUCCESS"
	}

	items := []gin.H{
		{"id": "linkedin", "label": "LinkedIn Sync", "status": linkedInStatus},
		{"id": "resume", "label": "Resume Analysis", "status": "SUCCESS"}, // If they have a profile, resume was parsed
		{"id": "github", "label": "GitHub Sync", "status": "PENDING"},
	}

	response.JSON(c, http.StatusOK, true, "Sync status fetched successfully", gin.H{
		"items": items,
	})
}

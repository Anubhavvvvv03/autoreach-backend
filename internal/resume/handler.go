package resume

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/autoreach-backend/internal/dto/response"
	"github.com/yourusername/autoreach-backend/internal/profile"
)

type ResumeUploadRequest struct {
	ResumeRaw string `json:"resume_raw" binding:"required"`
}

func ParseResumeHandler(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		response.JSON(c, http.StatusUnauthorized, false, "Unauthorized", nil)
		return
	}

	var req ResumeUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, false, err.Error(), nil)
		return
	}

	// 1. Parse the resume (Mock)
	parsedData, err := ParseResume(req.ResumeRaw)
	if err != nil {
		response.JSON(c, http.StatusInternalServerError, false, "Failed to parse resume", nil)
		return
	}

	// 2. Check if profile exists to determine whether to Create or Update
	existingProfile, err := profile.GetProfileByUserID(userID)
	if err != nil {
		response.JSON(c, http.StatusInternalServerError, false, "Failed to check existing profile", nil)
		return
	}

	var p *response.ProfileResponse
	if existingProfile != nil {
		p, err = profile.UpdateProfile(userID, *parsedData)
	} else {
		p, err = profile.CreateProfile(userID, *parsedData)
	}

	if err != nil {
		response.JSON(c, http.StatusInternalServerError, false, "Failed to store parsed profile", nil)
		return
	}

	response.JSON(c, http.StatusOK, true, "Resume parsed and profile updated successfully", p)
}

package resume

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourusername/autoreach-backend/internal/config"
	"github.com/yourusername/autoreach-backend/internal/dto/response"
	"github.com/yourusername/autoreach-backend/pkg/logger"
	"github.com/yourusername/autoreach-backend/pkg/storage"
)

const maxUploadSize = 5 << 20 // 5 MB

// UploadResumeHandler handles PDF upload, validates, stores in S3, and triggers the parsing pipeline.
func UploadResumeHandler(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		response.JSON(c, http.StatusUnauthorized, false, "Unauthorized", nil)
		return
	}

	// Parse multipart form
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.JSON(c, http.StatusBadRequest, false, "File is required", nil)
		return
	}
	defer file.Close()

	// Validate file type
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext != ".pdf" {
		response.JSON(c, http.StatusBadRequest, false, "Only PDF files are allowed", nil)
		return
	}

	// Validate file size
	if header.Size > maxUploadSize {
		response.JSON(c, http.StatusBadRequest, false, "File size must be less than 5MB", nil)
		return
	}

	// Generate S3 key
	fileID := uuid.NewString()
	s3Key := fmt.Sprintf("resumes/%s/%s.pdf", userID, fileID)

	// Upload to S3
	logger.Info(fmt.Sprintf("Uploading resume to S3: %s", s3Key))
	_, err = storage.Client.Upload(c.Request.Context(), s3Key, file, "application/pdf")
	if err != nil {
		logger.Error("S3 upload failed", err)
		response.JSON(c, http.StatusInternalServerError, false, "Failed to upload file", nil)
		return
	}

	// Create DB record
	record, err := CreateResumeRecord(config.DB, userID, s3Key)
	if err != nil {
		logger.Error("Failed to create resume record", err)
		response.JSON(c, http.StatusInternalServerError, false, "Failed to save resume record", nil)
		return
	}

	logger.Info(fmt.Sprintf("Resume uploaded: id=%s, s3_key=%s", record.ID, s3Key))

	// Trigger parsing pipeline asynchronously
	// context.Background() is used because the request context (c.Request.Context())
	// will be cancelled as soon as this handler returns the 202 response.
	go func() {
		// Create a separate background context for the long-running task
		ctx := context.Background()
		_, err := ProcessResume(ctx, userID, record.ID, s3Key)
		if err != nil {
			logger.Error(fmt.Sprintf("Async resume processing failed for id=%s", record.ID), err)
		}
	}()

	// Returns 202 Accepted immediately
	response.JSON(c, http.StatusAccepted, true, "Resume uploaded and processing started", gin.H{
		"resume_id": record.ID,
		"status":    StatusPending,
	})
}

// GetResumeStatusHandler returns the status of a previously uploaded resume.
func GetResumeStatusHandler(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		response.JSON(c, http.StatusUnauthorized, false, "Unauthorized", nil)
		return
	}

	resumeID := c.Param("id")
	if resumeID == "" {
		response.JSON(c, http.StatusBadRequest, false, "Resume ID is required", nil)
		return
	}

	record, err := GetResumeByID(resumeID)
	if err != nil {
		response.JSON(c, http.StatusNotFound, false, "Resume not found", nil)
		return
	}

	// Verify ownership
	if record.UserID != userID {
		response.JSON(c, http.StatusForbidden, false, "Access denied", nil)
		return
	}

	var parsedData interface{}
	if record.ParsedData != "" {
		json.Unmarshal([]byte(record.ParsedData), &parsedData)
	}

	response.JSON(c, http.StatusOK, true, "Resume status fetched", gin.H{
		"resume_id":   record.ID,
		"status":      record.Status,
		"parsed_data": parsedData,
		"fail_reason": record.FailReason,
		"created_at":  record.CreatedAt,
	})
}

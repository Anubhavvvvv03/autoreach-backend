package resume

import (
	"github.com/yourusername/autoreach-backend/internal/config"
	"gorm.io/gorm"
)

func CreateResumeRecord(db *gorm.DB, userID, s3Key string) (*ResumeFile, error) {
	record := &ResumeFile{
		UserID: userID,
		S3Key:  s3Key,
		Status: StatusPending,
	}
	if err := db.Create(record).Error; err != nil {
		return nil, err
	}
	return record, nil
}

func UpdateResumeStatus(db *gorm.DB, id, status, parsedData, failReason string) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if parsedData != "" {
		updates["parsed_data"] = parsedData
	}
	if failReason != "" {
		updates["fail_reason"] = failReason
	}
	return db.Model(&ResumeFile{}).Where("id = ?", id).Updates(updates).Error
}

func GetResumeByID(id string) (*ResumeFile, error) {
	var r ResumeFile
	if err := config.DB.Where("id = ?", id).First(&r).Error; err != nil {
		return nil, err
	}
	return &r, nil
}

func GetExpiredResumes(db *gorm.DB, days int) ([]ResumeFile, error) {
	var resumes []ResumeFile
	err := db.Where("created_at < NOW() - INTERVAL '1 day' * ?", days).Find(&resumes).Error
	return resumes, err
}

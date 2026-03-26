package resume

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	StatusPending = "PENDING"
	StatusSuccess = "SUCCESS"
	StatusFailed  = "FAILED"
)

type ResumeFile struct {
	ID         string         `gorm:"type:uuid;primaryKey" json:"id"`
	UserID     string         `gorm:"type:uuid;index;not null" json:"user_id"`
	S3Key      string         `gorm:"not null" json:"s3_key"`
	Status     string         `gorm:"type:varchar(20);default:PENDING;not null" json:"status"`
	ParsedData string         `gorm:"type:text" json:"parsed_data"`
	FailReason string         `gorm:"type:text" json:"fail_reason,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (r *ResumeFile) BeforeCreate(tx *gorm.DB) (err error) {
	r.ID = uuid.NewString()
	return
}

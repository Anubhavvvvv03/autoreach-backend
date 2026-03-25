package dashboard

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Activity struct {
	ID        string         `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    string         `gorm:"type:uuid;index;not null" json:"user_id"`
	Type      string         `json:"type"` // OUTREACH, SYSTEM, RESUME
	Title     string         `json:"title"`
	Status    string         `json:"status"` // SUCCESS, FAILED, PENDING
	CreatedAt time.Time      `json:"timestamp"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (a *Activity) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ID == "" {
		a.ID = uuid.NewString()
	}
	return
}

type Insight struct {
	ID          string         `gorm:"type:uuid;primaryKey" json:"id"`
	UserID      string         `gorm:"type:uuid;index;not null" json:"user_id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Priority    string         `json:"priority"` // HIGH, LOW
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (i *Insight) BeforeCreate(tx *gorm.DB) (err error) {
	i.ID = uuid.NewString()
	return
}

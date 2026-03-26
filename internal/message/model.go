package message

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Message struct {
	ID          string         `gorm:"type:uuid;primaryKey" json:"id"`
	UserID      string         `gorm:"type:uuid;index;not null" json:"user_id"`
	Company     string         `json:"company"`
	Role        string         `json:"role"`
	Recruiter   string         `json:"recruiter"`
	Tone        string         `json:"tone"`
	Context     string         `json:"context"`
	Text        string         `gorm:"type:text" json:"text"`
	Status      string         `json:"status"` // E.g., GENERATED, SENT
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (m *Message) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.NewString()
	return
}

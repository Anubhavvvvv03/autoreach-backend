package profile

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Project struct {
	Name        string `json:"name"`
	TechStack   string `json:"tech_stack"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Meta        string `json:"meta"`
}

type Experience struct {
	CompanyName string `json:"company_name"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	WorkDone    string `json:"work_done"`
}

type Profile struct {
	ID         string         `gorm:"type:uuid;primaryKey" json:"id"`
	UserID     string         `gorm:"type:uuid;uniqueIndex;not null" json:"user_id"`
	Skills     []string       `gorm:"serializer:json" json:"skills"`
	Experience []Experience   `gorm:"serializer:json" json:"experience"`
	Projects   []Project      `gorm:"serializer:json" json:"projects"`
	ResumeRaw  string         `gorm:"type:text" json:"resume_raw"`
	Meta       string         `gorm:"type:text" json:"meta"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (p *Profile) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.NewString()
	return
}

package auth

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserSettings struct {
	NotificationsEnabled bool `json:"notificationsEnabled"`
}

type User struct {
	ID        string         `gorm:"type:uuid;primaryKey" json:"id"`
	FullName  string         `json:"fullName"`
	Email     string         `gorm:"uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"not null" json:"-"`
	AvatarUrl string         `json:"avatarUrl"`
	Settings  UserSettings   `gorm:"serializer:json" json:"settings"`
	Meta      string         `gorm:"type:text" json:"meta"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.NewString()
	return
}

package models

import (
	modelMotor "motorbike-rental-backend/internal/app/motorbike/models"
	modelUser "motorbike-rental-backend/internal/app/user-and-auth/models"
	"time"
)

// todo : buraya ek olarak başlangıç ve bitiş noktası/konumu eklemeliyiz!
type Ride struct {
	BaseModel
	UserID      uint       `gorm:"not null"`
	MotorbikeID uint       `gorm:"not null"`
	StartTime   time.Time  `gorm:"not null"`
	EndTime     *time.Time `gorm:"not null"`
	Duration    string     `gorm:"type:interval;not null"`
	Cost        float64    `gorm:"not null"`

	User      modelUser.User       `gorm:"foreignKey:UserID"`
	Motorbike modelMotor.Motorbike `gorm:"foreignKey:MotorbikeID"`
}

func (Ride) TableName() string {
	return "rides"
}

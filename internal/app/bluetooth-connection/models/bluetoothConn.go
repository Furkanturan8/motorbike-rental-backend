package models

import (
	modelMotor "motorbike-rental-backend/internal/app/motorbike/models"
	modelUser "motorbike-rental-backend/internal/app/user-and-auth/models"
	"time"
)

type BluetoothConnection struct {
	BaseModel
	UserID         uint       `gorm:"not null"`
	MotorbikeID    uint       `gorm:"not null"`
	ConnectedAt    time.Time  `gorm:"not null"` // Bağlantı zamanı
	DisconnectedAt *time.Time // Bağlantının kesildiği zaman

	User      modelUser.User       `gorm:"foreignKey:UserID"`
	Motorbike modelMotor.Motorbike `gorm:"foreignKey:MotorbikeID"`
}

func (BluetoothConnection) TableName() string {
	return "bluetooth_connection"
}

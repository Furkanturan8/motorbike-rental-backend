package models

import (
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	ID        int64          `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index;column:deleted_at" json:"deleted_at"`
}

func (BaseModel) ModelName() string {
	return "base_model"
}

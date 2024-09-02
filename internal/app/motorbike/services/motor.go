package services

import (
	"encoding/json"
	"gorm.io/gorm"
	"motorbike-rental-backend/internal/app/motorbike/models"
)

type IMotorService interface {
	CreateMotorbike(motorbike *models.Motorbike) error
}

type MotorService struct {
	DB *gorm.DB
}

func NewMotorService(db *gorm.DB) IMotorService {
	return &MotorService{DB: db}
}

func (s *MotorService) CreateMotorbike(motorbike *models.Motorbike) error {
	photoURLsJSON, err := json.Marshal(motorbike.PhotoURLs)
	if err != nil {
		return err
	}

	motorbike.PhotoURLs = []string{string(photoURLsJSON)}

	return s.DB.Create(motorbike).Error
}

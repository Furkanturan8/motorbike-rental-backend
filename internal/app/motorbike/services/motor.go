package services

import (
	"context"
	"gorm.io/gorm"
	"motorbike-rental-backend/internal/app/motorbike/models"
)

type IMotorService interface {
	CreateMotorbike(ctx context.Context, motorbike *models.Motorbike) error
	AddPhotosToMotorbike(ctx context.Context, photos []models.MotorbikePhoto) error
	GetPhotosByMotorbikeID(ctx context.Context, motorbikeID string, photos *[]models.MotorbikePhoto) error
	GetAllMotors(ctx context.Context) (*[]models.Motorbike, error)
}
type MotorService struct {
	DB *gorm.DB
}

func NewMotorService(db *gorm.DB) IMotorService {
	return &MotorService{DB: db}
}

func (s *MotorService) CreateMotorbike(ctx context.Context, motorbike *models.Motorbike) error {
	return s.DB.WithContext(ctx).Create(motorbike).Error
}

func (s *MotorService) AddPhotosToMotorbike(ctx context.Context, photos []models.MotorbikePhoto) error {
	return s.DB.WithContext(ctx).Create(&photos).Error
}

func (s *MotorService) GetAllMotors(ctx context.Context) (*[]models.Motorbike, error) {
	var motors []models.Motorbike
	if err := s.DB.WithContext(ctx).Find(&motors).Error; err != nil {
		return nil, err
	}

	return &motors, nil
}

func (s *MotorService) GetPhotosByMotorbikeID(ctx context.Context, motorbikeID string, photos *[]models.MotorbikePhoto) error {
	return s.DB.WithContext(ctx).Where("motorbike_id = ?", motorbikeID).Find(photos).Error
}

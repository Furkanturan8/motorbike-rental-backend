package services

import (
	"context"
	"gorm.io/gorm"
	"motorbike-rental-backend/internal/app/motorbike/models"
)

type IMotorService interface {
	CreateMotor(ctx context.Context, motorbike *models.Motorbike) error
	UpdateMotor(ctx context.Context, motorbike *models.Motorbike) error
	UpdatePhotosForMotor(ctx context.Context, newPhotos []models.MotorbikePhoto, motorbikeID int) error
	AddPhotosToMotor(ctx context.Context, photos []models.MotorbikePhoto) error
	GetPhotosByID(ctx context.Context, motorbikeID string, photos *[]models.MotorbikePhoto) error
	GetAllMotors(ctx context.Context) (*[]models.Motorbike, error)
	GetMotorByID(ctx context.Context, motorbikeID int) (*models.Motorbike, error)
	GetMotorsForStatus(ctx context.Context, status string) (*[]models.Motorbike, error)
}
type MotorService struct {
	DB *gorm.DB
}

func NewMotorService(db *gorm.DB) IMotorService {
	return &MotorService{DB: db}
}

func (s *MotorService) CreateMotor(ctx context.Context, motorbike *models.Motorbike) error {
	return s.DB.WithContext(ctx).Create(motorbike).Error
}

func (s *MotorService) UpdateMotor(ctx context.Context, motorbike *models.Motorbike) error {
	return s.DB.WithContext(ctx).Save(motorbike).Error
}

func (s *MotorService) UpdatePhotosForMotor(ctx context.Context, newPhotos []models.MotorbikePhoto, motorbikeID int) error {
	// Önce mevcut fotoğrafları sil
	if err := s.DB.WithContext(ctx).Where("motorbike_id = ?", motorbikeID).Delete(&models.MotorbikePhoto{}).Error; err != nil {
		return err
	}

	// Yeni fotoğrafları ekle
	for _, photo := range newPhotos {
		photo.MotorbikeID = motorbikeID
		if err := s.DB.WithContext(ctx).Create(&photo).Error; err != nil {
			return err
		}
	}

	return nil
}

func (s *MotorService) AddPhotosToMotor(ctx context.Context, photos []models.MotorbikePhoto) error {
	return s.DB.WithContext(ctx).Create(&photos).Error
}

func (s *MotorService) GetAllMotors(ctx context.Context) (*[]models.Motorbike, error) {
	var motors []models.Motorbike
	if err := s.DB.WithContext(ctx).Find(&motors).Error; err != nil {
		return nil, err
	}

	return &motors, nil
}

func (s *MotorService) GetPhotosByID(ctx context.Context, motorbikeID string, photos *[]models.MotorbikePhoto) error {
	return s.DB.WithContext(ctx).Where("motorbike_id = ?", motorbikeID).Find(photos).Error
}

func (s *MotorService) GetMotorByID(ctx context.Context, motorbikeID int) (*models.Motorbike, error) {
	var motor models.Motorbike
	if err := s.DB.WithContext(ctx).Where("id = ?", motorbikeID).First(&motor).Error; err != nil {
		return nil, err
	}

	return &motor, nil
}

func (s *MotorService) GetMotorsForStatus(ctx context.Context, status string) (*[]models.Motorbike, error) {
	var motors []models.Motorbike
	if err := s.DB.WithContext(ctx).Where("status = ?", status).Find(&motors).Error; err != nil {
		return nil, err
	}

	return &motors, nil
}

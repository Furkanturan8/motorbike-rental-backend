package services

import (
	"context"
	"gorm.io/gorm"
	"motorbike-rental-backend/internal/app/ride/models"
)

type IRideService interface {
	GetAllRides(ctx context.Context) (*[]models.Ride, error)
}

type RideService struct {
	DB *gorm.DB
}

func NewRideService(db *gorm.DB) IRideService {
	return &RideService{DB: db}
}

func (s *RideService) GetAllRides(ctx context.Context) (*[]models.Ride, error) {
	var rides []models.Ride
	if err := s.DB.WithContext(ctx).Find(&rides).Error; err != nil {
		return nil, err
	}

	return &rides, nil
}

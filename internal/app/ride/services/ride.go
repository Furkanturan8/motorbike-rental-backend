package services

import (
	"context"
	"gorm.io/gorm"
	"motorbike-rental-backend/internal/app/ride/models"
)

type IRideService interface {
	GetAllRides(ctx context.Context) (*[]models.Ride, error)
	GetRideByID(ctx context.Context, id int) (*models.Ride, error)
	CreateRide(ctx context.Context, ride *models.Ride) error
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

func (s *RideService) GetRideByID(ctx context.Context, id int) (*models.Ride, error) {
	var ride models.Ride

	if err := s.DB.WithContext(ctx).Where("id = ?", id).First(&ride).Error; err != nil {
		return nil, err
	}
	return &ride, nil
}

func (s *RideService) CreateRide(ctx context.Context, ride *models.Ride) error {
	return s.DB.WithContext(ctx).Create(ride).Error
}

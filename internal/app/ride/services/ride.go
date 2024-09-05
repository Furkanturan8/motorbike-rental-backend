package services

import (
	"context"
	"gorm.io/gorm"
	modelBike "motorbike-rental-backend/internal/app/motorbike/models"
	"motorbike-rental-backend/internal/app/ride/models"
	modelUser "motorbike-rental-backend/internal/app/user-and-auth/models"
)

type IRideService interface {
	GetAllRides(ctx context.Context) (*[]models.Ride, error)
	GetRideByID(ctx context.Context, id int) (*models.Ride, error)
	CreateRide(ctx context.Context, ride *models.Ride) error
	GetRidesByUserID(ctx context.Context, userID int) (*[]models.Ride, error)
	GetRideByUserID(ctx context.Context, userID int, rideID int) (*models.Ride, error)
	GetRidesByBikeID(ctx context.Context, bikeID int) (*[]models.Ride, error)
	UpdateRide(ctx context.Context, ride *models.Ride) error
	DeleteRide(ctx context.Context, id int) error
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

func (s *RideService) GetRidesByUserID(ctx context.Context, userID int) (*[]models.Ride, error) {
	var user modelUser.User
	var rides []models.Ride

	if err := s.DB.WithContext(ctx).Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}

	if err := s.DB.WithContext(ctx).Where("user_id = ?", userID).Find(&rides).Error; err != nil {
		return nil, err
	}
	return &rides, nil
}

func (s *RideService) GetRideByUserID(ctx context.Context, userID int, rideID int) (*models.Ride, error) {
	var user modelUser.User
	var ride models.Ride

	if err := s.DB.WithContext(ctx).Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}

	if err := s.DB.WithContext(ctx).Where("user_id = ? AND id = ?", userID, rideID).First(&ride).Error; err != nil {
		return nil, err
	}
	return &ride, nil
}

func (s *RideService) GetRidesByBikeID(ctx context.Context, bikeID int) (*[]models.Ride, error) {
	var motor modelBike.Motorbike
	var rides []models.Ride

	if err := s.DB.WithContext(ctx).Where("id = ?", bikeID).First(&motor).Error; err != nil {
		return nil, err
	}

	if err := s.DB.WithContext(ctx).Where("motorbike_id = ?", bikeID).Find(&rides).Error; err != nil {
		return nil, err
	}

	return &rides, nil
}

func (s *RideService) UpdateRide(ctx context.Context, ride *models.Ride) error {
	return s.DB.WithContext(ctx).Save(ride).Error
}

func (s *RideService) DeleteRide(ctx context.Context, id int) error {
	var ride models.Ride
	if err := s.DB.WithContext(ctx).First(&ride, id).Error; err != nil {
		return err
	}

	return s.DB.WithContext(ctx).Delete(&ride).Error
}

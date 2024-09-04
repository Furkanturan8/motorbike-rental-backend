package viewmodels

import (
	"motorbike-rental-backend/internal/app/ride/models"
	"time"
)

type RideCreateVM struct {
	UserID      uint       `json:"user_id" validate:"required,numeric"`
	MotorbikeID uint       `json:"motorbike_id" validate:"required,numeric"`
	StartTime   time.Time  `json:"start_time" validate:"required"`
	EndTime     *time.Time `json:"end_time,omitempty"` // Zorunlu değil, boş olabilir
	Duration    string     `json:"duration" validate:"required"`
	Cost        float64    `json:"cost" validate:"required"`
}

// RideCreateVM'den Ride modeline dönüştürme
func (vm *RideCreateVM) ToDBModel() models.Ride {
	now := time.Now()

	return models.Ride{
		UserID:      vm.UserID,
		MotorbikeID: vm.MotorbikeID,
		StartTime:   now,        // StartTime, o anki zaman olacak şekilde ayarlandı
		EndTime:     vm.EndTime, // EndTime isteğe bağlı olarak null olabilir
		Duration:    vm.Duration,
		Cost:        vm.Cost,
	}
}

type RideUpdateVM struct {
	UserID      uint       `json:"user_id" validate:"required,numeric"`
	MotorbikeID uint       `json:"motorbike_id" validate:"required,numeric"`
	StartTime   time.Time  `json:"start_time" validate:"required"`
	EndTime     *time.Time `json:"end_time,omitempty"` // Zorunlu değil, boş olabilir
	Duration    string     `json:"duration" validate:"required"`
	Cost        float64    `json:"cost" validate:"required"`
}

func (vm *RideUpdateVM) ToDBModel() models.Ride {

	return models.Ride{
		UserID:      vm.UserID,
		MotorbikeID: vm.MotorbikeID,
		EndTime:     vm.EndTime,
		Duration:    vm.Duration,
		Cost:        vm.Cost,
	}
}

type RideDetailVM struct {
	UserID      uint       `json:"user_id"`
	MotorbikeID uint       `json:"motorbike_id"`
	StartTime   time.Time  `json:"start_time"`
	EndTime     *time.Time `json:"end_time"`
	Duration    string     `json:"duration"`
	Cost        float64    `json:"cost"`
}

func (vm *RideDetailVM) ToDBModel(ride models.Ride) RideDetailVM {
	return RideDetailVM{
		UserID:      ride.UserID,
		MotorbikeID: ride.MotorbikeID,
		StartTime:   ride.StartTime,
		EndTime:     ride.EndTime,
		Duration:    ride.Duration,
		Cost:        ride.Cost,
	}
}

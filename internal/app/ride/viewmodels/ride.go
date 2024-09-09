package viewmodels

import (
	modelBike "motorbike-rental-backend/internal/app/motorbike/models"
	"motorbike-rental-backend/internal/app/ride/models"
	modelUser "motorbike-rental-backend/internal/app/user-and-auth/models"
	"time"
)

type RideCreateVM struct {
	UserID      uint `json:"user_id" validate:"required,numeric"`
	MotorbikeID uint `json:"motorbike_id" validate:"required,numeric"`
}

// RideCreateVM'den Ride modeline dönüştürme
func (vm *RideCreateVM) ToDBModel() models.Ride {
	now := time.Now().UTC()

	// EndTime varsayılan olarak null olabilir veya şu anki zamanı ayarlayabilirsiniz
	var endTime *time.Time

	return models.Ride{
		UserID:      vm.UserID,
		MotorbikeID: vm.MotorbikeID,
		StartTime:   now,     // StartTime, o anki zaman olacak şekilde ayarlandı
		EndTime:     endTime, // EndTime isteğe bağlı olarak null olabilir
		Duration:    "0 years 0 mons 0 days 0 hours 0 mins 0 secs",
		Cost:        0,
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

func (vm *RideUpdateVM) ToDBModel(m models.Ride) models.Ride {
	m.UserID = vm.UserID
	m.MotorbikeID = vm.MotorbikeID
	m.EndTime = vm.EndTime
	m.Duration = vm.Duration
	m.Cost = vm.Cost

	return m
}

type RideDetailVM struct {
	ID          uint                `json:"id"`
	UserID      uint                `json:"user_id"`
	MotorbikeID uint                `json:"motorbike_id"`
	StartTime   time.Time           `json:"start_time"`
	EndTime     *time.Time          `json:"end_time"`
	Duration    string              `json:"duration"`
	Cost        float64             `json:"cost"`
	User        modelUser.User      `json:"user"`
	Motorbike   modelBike.Motorbike `json:"bike"`
}

func (vm *RideDetailVM) ToViewModel(ride models.Ride) RideDetailVM {
	return RideDetailVM{
		ID:          uint(ride.ID),
		UserID:      ride.UserID,
		MotorbikeID: ride.MotorbikeID,
		StartTime:   ride.StartTime,
		EndTime:     ride.EndTime,
		Duration:    ride.Duration,
		Cost:        ride.Cost,
		User:        ride.User,
		Motorbike:   ride.Motorbike,
	}
}

package viewmodels

import (
	"motorbike-rental-backend/internal/app/motorbike/models"
	"motorbike-rental-backend/pkg/utils"
)

type BikeCreateVM struct {
	Model             string   `json:"model" validate:"required,max=100"`
	LocationLatitude  float64  `json:"location_latitude" validate:"required,numeric"`
	LocationLongitude float64  `json:"location_longitude" validate:"required,numeric"`
	Status            string   `json:"status" validate:"required,oneof=available maintenance rented"`
	PhotoURLs         []string `json:"photo_urls" validate:"dive,url"`
	LockStatus        string   `json:"lock_status" validate:"required,oneof=locked unlocked"`
}

func (vm BikeCreateVM) ToDBModel(m models.Motorbike) models.Motorbike {
	m.Model = utils.ToTitle(vm.Model)
	m.Status = models.MotorBikeStatus(vm.Status)
	m.LocationLatitude = vm.LocationLatitude
	m.LocationLongitude = vm.LocationLongitude
	m.LockStatus = models.LockStatus(vm.LockStatus)
	m.PhotoURLs = vm.PhotoURLs

	return m
}

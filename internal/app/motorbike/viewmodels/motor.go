package viewmodels

import (
	"motorbike-rental-backend/internal/app/motorbike/models"
	"time"
)

// Fotoğraflar için ayrı bir view model
type PhotoCreateVM struct {
	PhotoURL string `json:"photo_url" validate:"required,url"`
}

// Motorbike oluşturma için view model
type BikeCreateVM struct {
	Model             string          `json:"model" validate:"required,max=100"`
	LocationLatitude  float64         `json:"location_latitude" validate:"required,numeric"`
	LocationLongitude float64         `json:"location_longitude" validate:"required,numeric"`
	Status            string          `json:"status" validate:"required,oneof=available maintenance rented"`
	Photos            []PhotoCreateVM `json:"photos"`
	LockStatus        string          `json:"lock_status" validate:"required,oneof=locked unlocked"`
}

// Motorbike modeline dönüştürme
func (vm BikeCreateVM) ToDBModel() models.Motorbike {
	return models.Motorbike{
		Model:             vm.Model,
		LocationLatitude:  vm.LocationLatitude,
		LocationLongitude: vm.LocationLongitude,
		Status:            models.MotorBikeStatus(vm.Status),
		LockStatus:        models.LockStatus(vm.LockStatus),
	}
}

// Fotoğraf modellerine dönüştürme
func (vm BikeCreateVM) ToPhotoModels(motorbikeID int) []models.MotorbikePhoto {
	var photos []models.MotorbikePhoto
	for _, photoVM := range vm.Photos {
		photos = append(photos, models.MotorbikePhoto{
			MotorbikeID: motorbikeID,
			PhotoURL:    photoVM.PhotoURL,
		})
	}
	return photos
}

// Motorbike güncelleme için view model
type BikeUpdateVM struct {
	Model             string          `json:"model" validate:"required,max=100"`
	LocationLatitude  float64         `json:"location_latitude" validate:"required,numeric"`
	LocationLongitude float64         `json:"location_longitude" validate:"required,numeric"`
	Status            string          `json:"status" validate:"required,oneof=available maintenance rented"`
	Photos            []PhotoCreateVM `json:"photos"`
	LockStatus        string          `json:"lock_status" validate:"required,oneof=locked unlocked"`
}

// Güncellenmiş Motorbike modeline dönüştürme
func (vm BikeUpdateVM) ToDBModel(m models.Motorbike) models.Motorbike {
	m.Model = vm.Model
	m.LocationLatitude = vm.LocationLatitude
	m.LocationLongitude = vm.LocationLongitude
	m.Status = models.MotorBikeStatus(vm.Status)
	m.LockStatus = models.LockStatus(vm.LockStatus)
	return m
}

// Fotoğraf detayları için view model
type PhotoDetailVM struct {
	ID          int    `json:"id"`
	MotorbikeID int    `json:"motorbike_id"`
	PhotoURL    string `json:"photo_url"`
}

// Motorbike detayları için view model
type BikeDetailVM struct {
	ID                int             `json:"id"`
	Model             string          `json:"model"`
	LocationLatitude  float64         `json:"location_latitude"`
	LocationLongitude float64         `json:"location_longitude"`
	Status            string          `json:"status"`
	Photos            []PhotoDetailVM `json:"photos"`
	LockStatus        string          `json:"lock_status"`
}

// Motorbike modelini detay view modeline dönüştürme
func NewBikeDetailVM(motorbike models.Motorbike, photos []models.MotorbikePhoto) BikeDetailVM {
	var photoVMs []PhotoDetailVM
	for _, photo := range photos {
		photoVMs = append(photoVMs, PhotoDetailVM{
			ID:          int(photo.ID),
			MotorbikeID: photo.MotorbikeID,
			PhotoURL:    photo.PhotoURL,
		})
	}

	return BikeDetailVM{
		ID:                int(motorbike.ID),
		Model:             motorbike.Model,
		LocationLatitude:  motorbike.LocationLatitude,
		LocationLongitude: motorbike.LocationLongitude,
		Status:            motorbike.Status.String(),
		Photos:            photoVMs,
		LockStatus:        motorbike.LockStatus.String(),
	}
}

// Nullable time formatlama fonksiyonu
func formatNullableTime(t *time.Time) *string {
	if t == nil {
		return nil
	}
	formattedTime := t.Format("2006-01-02 15:04:05")
	return &formattedTime
}

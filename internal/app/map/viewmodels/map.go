package viewmodels

import (
	"motorbike-rental-backend/internal/app/map/models"
	"motorbike-rental-backend/internal/app/motorbike/viewmodels"
)

type MapCreateVM struct {
	MotorbikeID       uint    `json:"motorbike_id" validate:"required"`               // Zorunlu alan
	Name              string  `json:"name" validate:"required,min=3,max=255"`         // Zorunlu, min 3, max 255 karakter
	Description       string  `json:"description" validate:"max=500"`                 // Maksimum 500 karakter
	LocationLatitude  float64 `json:"latitude" validate:"required,gte=-90,lte=90"`    // Enlem -90 ile 90 arasında olmalı
	LocationLongitude float64 `json:"longitude" validate:"required,gte=-180,lte=180"` // Boylam -180 ile 180 arasında olmalı
	ZoomLevel         int     `json:"zoom_level" validate:"gte=1,lte=20"`             // Yakınlaştırma seviyesi 1-20 arası olmalı
	MapType           string  `json:"map_type" validate:"required"`                   // Harita türü belirtilmeli (sınırlı değerler)
}

func (vm *MapCreateVM) ToDBModel() models.Map {
	return models.Map{
		MotorbikeID:       vm.MotorbikeID,
		Name:              vm.Name,
		Description:       vm.Description,
		LocationLatitude:  vm.LocationLatitude,
		LocationLongitude: vm.LocationLongitude,
		ZoomLevel:         vm.ZoomLevel,
		MapType:           vm.MapType,
	}
}

type MapUpdateVM struct {
	MotorbikeID       uint    `json:"motorbike_id" validate:"required"`
	Name              string  `json:"name" validate:"required,min=3,max=255"`
	Description       string  `json:"description" validate:"max=500"`
	LocationLatitude  float64 `json:"latitude" validate:"required,gte=-90,lte=90"`
	LocationLongitude float64 `json:"longitude" validate:"required,gte=-180,lte=180"`
	ZoomLevel         int     `json:"zoom_level" validate:"gte=1,lte=20"`
	MapType           string  `json:"map_type" validate:"required"`
}

func (vm *MapUpdateVM) ToDBModel(m models.Map) models.Map {
	m.MotorbikeID = vm.MotorbikeID
	m.Name = vm.Name
	m.Description = vm.Description
	m.LocationLatitude = vm.LocationLatitude
	m.LocationLongitude = vm.LocationLongitude
	m.ZoomLevel = vm.ZoomLevel
	m.MapType = vm.MapType

	return m
}

type MapDetailVM struct {
	ID                uint                    `json:"id"`
	MotorbikeID       uint                    `json:"motorbike_id"`
	Name              string                  `json:"name"`
	Description       string                  `json:"description"`
	LocationLatitude  float64                 `json:"latitude"`  // Enlem bilgisi
	LocationLongitude float64                 `json:"longitude"` // Boylam bilgisi
	ZoomLevel         int                     `json:"zoom_level"`
	MapType           string                  `json:"map_type"`
	MotorBike         viewmodels.BikeDetailVM `json:"motorbike"`
}

func (vm *MapDetailVM) ToViewModel(m models.Map) MapDetailVM {
	return MapDetailVM{
		ID:                uint(m.ID),
		MotorbikeID:       m.MotorbikeID,
		Name:              m.Name,
		Description:       m.Description,
		LocationLatitude:  m.LocationLatitude,
		LocationLongitude: m.LocationLongitude,
		ZoomLevel:         m.ZoomLevel,
		MapType:           m.MapType,
		MotorBike:         viewmodels.NewBikeDetailVM(m.Motorbike, m.Motorbike.Photos),
	}
}

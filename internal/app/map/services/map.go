package services

import (
	"context"
	"gorm.io/gorm"
	"motorbike-rental-backend/internal/app/map/models"
)

type IMapService interface {
	GetAllMaps(ctx context.Context) (*[]models.Map, error)
	GetMapByID(ctx context.Context, id int) (*models.Map, error)
}

type MapService struct {
	DB *gorm.DB
}

func NewMapService(db *gorm.DB) IMapService {
	return &MapService{DB: db}
}

func (s *MapService) GetAllMaps(ctx context.Context) (*[]models.Map, error) {
	var maps []models.Map

	if err := s.DB.WithContext(ctx).Preload("Motorbike").Preload("Motorbike.Photos").Find(&maps).Error; err != nil {
		return nil, err
	}

	return &maps, nil
}

func (s *MapService) GetMapByID(ctx context.Context, id int) (*models.Map, error) {
	var _map models.Map

	if err := s.DB.WithContext(ctx).Preload("Motorbike").Preload("Motorbike.Photos").Where("id = ?", id).First(&_map).Error; err != nil {
		return nil, err
	}
	return &_map, nil
}

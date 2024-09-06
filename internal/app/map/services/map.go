package services

import (
	"context"
	"gorm.io/gorm"
	"motorbike-rental-backend/internal/app/map/models"
)

type IMapService interface {
	GetAllMaps(ctx context.Context) (*[]models.Map, error)
}

type MapService struct {
	DB *gorm.DB
}

func NewMapService(db *gorm.DB) IMapService {
	return &MapService{DB: db}
}

func (s *MapService) GetAllMaps(ctx context.Context) (*[]models.Map, error) {
	var maps []models.Map

	if err := s.DB.WithContext(ctx).Find(&maps).Error; err != nil {
		return nil, err
	}

	return &maps, nil
}

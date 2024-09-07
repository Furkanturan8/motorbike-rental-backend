package services

import (
	"context"
	"gorm.io/gorm"
	"motorbike-rental-backend/internal/app/bluetooth-connection/models"
)

type IConnService interface {
	GetAllConnections(ctx context.Context) (*[]models.BluetoothConnection, error)
}

type ConnService struct {
	DB *gorm.DB
}

func NewConnService(db *gorm.DB) IConnService {
	return &ConnService{DB: db}
}

func (s *ConnService) GetAllConnections(ctx context.Context) (*[]models.BluetoothConnection, error) {
	var connections []models.BluetoothConnection
	if err := s.DB.WithContext(ctx).Preload("User").Preload("Motorbike").Find(&connections).Error; err != nil {
		return nil, err
	}

	return &connections, nil
}

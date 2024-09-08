package services

import (
	"context"
	"gorm.io/gorm"
	"motorbike-rental-backend/internal/app/bluetooth-connection/models"
)

type IConnService interface {
	GetAllConnections(ctx context.Context) (*[]models.BluetoothConnection, error)
	GetConnByID(ctx context.Context, id int) (*models.BluetoothConnection, error)
	CreateConn(ctx context.Context, conn *models.BluetoothConnection) error
	DeleteConn(ctx context.Context, id int) error
	UpdateConn(ctx context.Context, connection *models.BluetoothConnection) error
}

type ConnService struct {
	DB *gorm.DB
}

func NewConnService(db *gorm.DB) IConnService {
	return &ConnService{DB: db}
}

func (s *ConnService) GetAllConnections(ctx context.Context) (*[]models.BluetoothConnection, error) {
	var connections []models.BluetoothConnection
	if err := s.DB.WithContext(ctx).Preload("User").Preload("Motorbike").Preload("Motorbike.Photos").Find(&connections).Error; err != nil {
		return nil, err
	}

	return &connections, nil
}

func (s *ConnService) GetConnByID(ctx context.Context, id int) (*models.BluetoothConnection, error) {
	var connection models.BluetoothConnection

	if err := s.DB.WithContext(ctx).Where("id=?", id).Preload("User").Preload("Motorbike").Preload("Motorbike.Photos").First(&connection).Error; err != nil {
		return nil, err
	}

	return &connection, nil
}

// Connect func from handler -> this func gonna create a new connection!
func (s *ConnService) CreateConn(ctx context.Context, conn *models.BluetoothConnection) error {
	return s.DB.WithContext(ctx).Create(conn).Error
}

func (s *ConnService) DeleteConn(ctx context.Context, id int) error {
	var connection models.BluetoothConnection

	if err := s.DB.WithContext(ctx).First(&connection, id).Error; err != nil {
		return err
	}

	return s.DB.WithContext(ctx).Delete(&connection).Error
}

// Disconnect func from handler -> this func gonna disconnect a spesific connection!
func (s *ConnService) UpdateConn(ctx context.Context, connection *models.BluetoothConnection) error {
	return s.DB.WithContext(ctx).Save(connection).Error
}

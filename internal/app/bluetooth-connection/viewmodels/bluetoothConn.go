package viewmodels

import (
	"motorbike-rental-backend/internal/app/bluetooth-connection/models"
	motorViewmodel "motorbike-rental-backend/internal/app/motorbike/viewmodels"
	userViewmodel "motorbike-rental-backend/internal/app/user-and-auth/viewmodels"
	"time"
)

// BluetoothConnectionCreateVM is the view model for creating a new Bluetooth connection
type BluetoothConnectionCreateVM struct {
	UserID      uint      `json:"user_id" validate:"required"`
	MotorbikeID uint      `json:"motorbike_id" validate:"required"`
	ConnectedAt time.Time `json:"connected_at" validate:"required"`
}

func (vm *BluetoothConnectionCreateVM) ToDBModel() models.BluetoothConnection {
	return models.BluetoothConnection{
		UserID:      vm.UserID,
		MotorbikeID: vm.MotorbikeID,
		ConnectedAt: vm.ConnectedAt,
	}
}

// BluetoothConnectionUpdateVM is the view model for updating an existing Bluetooth connection
type BluetoothConnectionUpdateVM struct {
	UserID         uint      `json:"user_id" validate:"required"`
	MotorbikeID    uint      `json:"motorbike_id" validate:"required"`
	DisconnectedAt time.Time `json:"disconnected_at" validate:"required"`
}

func (vm *BluetoothConnectionUpdateVM) ToDBModel(m models.BluetoothConnection) models.BluetoothConnection {
	m.UserID = vm.UserID
	m.MotorbikeID = vm.MotorbikeID
	m.DisconnectedAt = &vm.DisconnectedAt
	return m
}

// BluetoothConnectionDetailVM is the view model for retrieving detailed information about a Bluetooth connection
type BluetoothConnectionDetailVM struct {
	ID             uint                              `json:"id"`
	UserID         uint                              `json:"user_id"`
	MotorbikeID    uint                              `json:"motorbike_id"`
	ConnectedAt    time.Time                         `json:"connected_at"`
	DisconnectedAt *time.Time                        `json:"disconnected_at"`
	User           userViewmodel.UserDetailVMForUser `json:"user"`
	Motorbike      motorViewmodel.BikeDetailVM       `json:"motorbike"`
}

func (vm *BluetoothConnectionDetailVM) ToViewModel(m models.BluetoothConnection) BluetoothConnectionDetailVM {
	return BluetoothConnectionDetailVM{
		ID:             uint(m.ID),
		UserID:         m.UserID,
		MotorbikeID:    m.MotorbikeID,
		ConnectedAt:    m.ConnectedAt,
		DisconnectedAt: m.DisconnectedAt,
		User:           userViewmodel.UserToUserDetailVMForUser(m.User),
		Motorbike:      motorViewmodel.NewBikeDetailVM(m.Motorbike, m.Motorbike.Photos),
	}
}

package models

type MotorBikeStatus string

const (
	BikeAvailable     MotorBikeStatus = "available"
	BikeInMaintenance MotorBikeStatus = "maintenance"
	BikeRented        MotorBikeStatus = "rented"
)

type LockStatus string

const (
	Locked   LockStatus = "locked"
	Unlocked LockStatus = "unlocked"
)

// Motorbike modeli
type Motorbike struct {
	BaseModel
	Model             string          `gorm:"not null"`
	LocationLatitude  float64         `gorm:"not null"`
	LocationLongitude float64         `gorm:"not null"`
	Status            MotorBikeStatus `gorm:"type:varchar(20);not null"` // ENUM gibi çalışacak şekilde varchar tanımlandı
	PhotoURLs         []string        `gorm:"type:json"`
	LockStatus        LockStatus      `gorm:"type:varchar(10);not null"` // ENUM gibi çalışacak şekilde varchar tanımlandı
}

func (Motorbike) TableName() string {
	return "motorbike"
}

func (r MotorBikeStatus) String() string {
	switch r {
	case BikeAvailable:
		return "available"
	case BikeInMaintenance:
		return "maintenance"
	case BikeRented:
		return "rented"
	default:
		return "unknown"
	}
}

func (r LockStatus) String() string {
	switch r {
	case Locked:
		return "locked"
	case Unlocked:
		return "unlocked"
	default:
		return "unknown"
	}
}

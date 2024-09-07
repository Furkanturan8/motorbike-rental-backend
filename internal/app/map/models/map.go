package models

import "motorbike-rental-backend/internal/app/motorbike/models"

type Map struct {
	BaseModel
	MotorbikeID       uint    `gorm:"not null" json:"motorbike_id"`
	Name              string  `gorm:"type:varchar(255)" json:"name"`
	Description       string  `gorm:"type:text" json:"description"`
	LocationLatitude  float64 `gorm:"not null" json:"latitude"`         // Enlem bilgisi
	LocationLongitude float64 `gorm:"not null" json:"longitude"`        // Boylam bilgisi
	ZoomLevel         int     `gorm:"default:12" json:"zoom_level"`     // Varsayılan yakınlaştırma seviyesi
	MapType           string  `gorm:"type:varchar(50)" json:"map_type"` // Harita türü

	Motorbike models.Motorbike `gorm:"foreignKey:MotorbikeID" json:"motorbike"`
}

func (Map) TableName() string {
	return "maps"
}

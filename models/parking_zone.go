package models
 
import (
	"time"
)
 
type ParkingZone struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Name          string         `gorm:"column:name;not null" json:"name"`
	Type          string         `gorm:"column:type;not null" json:"type"` // 'general', 'ev_charging', 'covered'
	TotalCapacity int            `gorm:"column:total_capacity;not null" json:"total_capacity"`
	PricePerHour  float64        `gorm:"column:price_per_hour;not null" json:"price_per_hour"`
	CreatedAt     time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
 
	// Relationships
	Reservations []Reservation `gorm:"foreignKey:ZoneID;references:ID" json:"-"`
}

 
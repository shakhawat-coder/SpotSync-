package models
 
import (
	"time"
 
)
 
type Reservation struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	UserID       uint           `gorm:"column:user_id;not null" json:"user_id"`
	ZoneID       uint           `gorm:"column:zone_id;not null" json:"zone_id"`
	LicensePlate string         `gorm:"column:license_plate;not null" json:"license_plate"`
	Status       string         `gorm:"column:status;default:'active';not null" json:"status"` // 'active', 'completed', 'cancelled'
	CreatedAt    time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
 
	// Relationships
	User *User        `gorm:"foreignKey:UserID;references:ID" json:"-"`
	Zone *ParkingZone `gorm:"foreignKey:ZoneID;references:ID" json:"-"`
}
 

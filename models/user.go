package models
 
import (
	"time"
 
)
 
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"column:name;not null" json:"name"`
	Email     string         `gorm:"column:email;uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"column:password;not null" json:"-"` // Never expose password
	Role      string         `gorm:"column:role;default:'driver';not null" json:"role"` // 'driver' or 'admin'
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
 
	// Relationships
	Reservations []Reservation `gorm:"foreignKey:UserID;references:ID" json:"-"`
}

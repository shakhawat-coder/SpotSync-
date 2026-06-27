package migrations

import (
	"fmt"
	"spotsync/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&models.User{},
		&models.ParkingZone{},
		&models.Reservation{},
	)
	if err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}
	fmt.Println("Migration completed successfully")
	return nil
}

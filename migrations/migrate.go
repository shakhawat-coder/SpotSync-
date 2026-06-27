package migrations

import (
	"fmt"
	"spotsync/models"

	"gorm.io/gorm"
)
 
func Migrate(db *gorm.DB) error {
	// Auto-migrate all models
	err :=	 db.AutoMigrate(
		&models.User{},
		&models.ParkingZone{},
		&models.Reservation{},
	)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Migration completed successfully")
	return nil
}
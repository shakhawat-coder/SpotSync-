package repository
 
import (
	"spotsync/models"
 
	"gorm.io/gorm"
)
 
type ParkingZoneRepository interface {
	Create(zone *models.ParkingZone) error
	GetByID(id uint) (*models.ParkingZone, error)
	GetAll() ([]models.ParkingZone, error)
	Update(zone *models.ParkingZone) error
	Delete(id uint) error
	GetActiveReservationCount(zoneID uint) (int64, error)
}
 
type parkingZoneRepository struct {
	db *gorm.DB
}
 
func NewParkingZoneRepository(db *gorm.DB) ParkingZoneRepository {
	return &parkingZoneRepository{db: db}
}
 
func (r *parkingZoneRepository) Create(zone *models.ParkingZone) error {
	if err := r.db.Create(zone).Error; err != nil {
		return err
	}
	return nil
}
 
func (r *parkingZoneRepository) GetByID(id uint) (*models.ParkingZone, error) {
	var zone models.ParkingZone
	if err := r.db.First(&zone, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &zone, nil
}
 
func (r *parkingZoneRepository) GetAll() ([]models.ParkingZone, error) {
	var zones []models.ParkingZone
	if err := r.db.Find(&zones).Error; err != nil {
		return nil, err
	}
	return zones, nil
}
 
func (r *parkingZoneRepository) Update(zone *models.ParkingZone) error {
	if err := r.db.Save(zone).Error; err != nil {
		return err
	}
	return nil
}
 
func (r *parkingZoneRepository) Delete(id uint) error {
	if err := r.db.Delete(&models.ParkingZone{}, id).Error; err != nil {
		return err
	}
	return nil
}
 
// Count active (non-cancelled) reservations for a zone
func (r *parkingZoneRepository) GetActiveReservationCount(zoneID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&models.Reservation{}).
		Where("zone_id = ? AND status = ?", zoneID, "active").
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
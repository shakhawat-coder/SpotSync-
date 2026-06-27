package repository

import (
	"spotsync/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ReservationRepository interface {
	Create(reservation *models.Reservation) error
	GetByID(id uint) (*models.Reservation, error)
	GetByUserID(userID uint) ([]models.Reservation, error)
	GetAll() ([]models.Reservation, error)
	Update(reservation *models.Reservation) error
	Cancel(id uint) error
	// CRITICAL: Atomic create with row-level locking to prevent over-capacity
	CreateWithLocking(reservation *models.Reservation, zoneID uint) error
}

type reservationRepository struct {
	db *gorm.DB
}

func NewReservationRepository(db *gorm.DB) ReservationRepository {
	return &reservationRepository{db: db}
}

func (r *reservationRepository) Create(reservation *models.Reservation) error {
	if err := r.db.Create(reservation).Error; err != nil {
		return err
	}
	return nil
}

func (r *reservationRepository) GetByID(id uint) (*models.Reservation, error) {
	var reservation models.Reservation
	if err := r.db.
		Preload("User").
		Preload("Zone").
		First(&reservation, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &reservation, nil
}

func (r *reservationRepository) GetByUserID(userID uint) ([]models.Reservation, error) {
	var reservations []models.Reservation
	if err := r.db.
		Preload("Zone").
		Where("user_id = ?", userID).
		Find(&reservations).Error; err != nil {
		return nil, err
	}
	return reservations, nil
}

func (r *reservationRepository) GetAll() ([]models.Reservation, error) {
	var reservations []models.Reservation
	if err := r.db.
		Preload("User").
		Preload("Zone").
		Find(&reservations).Error; err != nil {
		return nil, err
	}
	return reservations, nil
}

func (r *reservationRepository) Update(reservation *models.Reservation) error {
	if err := r.db.Save(reservation).Error; err != nil {
		return err
	}
	return nil
}

func (r *reservationRepository) Cancel(id uint) error {
	if err := r.db.Model(&models.Reservation{}).Where("id = ?", id).Update("status", "cancelled").Error; err != nil {
		return err
	}
	return nil
}

// ⚠️ CRITICAL: CreateWithLocking prevents race condition on EV spot reservation
// Uses database transaction + row-level locking (FOR UPDATE) to ensure atomicity
//
// Problem: Two drivers might both read "19 active spots" and both reserve,
// resulting in 21 cars in a 20-spot zone.
//
// Solution: Lock the parking_zone row, count active reservations, check capacity,
// and create reservation all in one atomic transaction.
func (r *reservationRepository) CreateWithLocking(reservation *models.Reservation, zoneID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var zone models.ParkingZone
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&zone, zoneID).Error; err != nil {
			return err
		}

		var duplicateCount int64
		if err := tx.Model(&models.Reservation{}).
			Where("license_plate = ? AND status = ?", reservation.LicensePlate, "active").
			Count(&duplicateCount).Error; err != nil {
			return err
		}
		if duplicateCount > 0 {
			return ErrDuplicateActiveLicense
		}

		var activeCount int64
		if err := tx.Model(&models.Reservation{}).
			Where("zone_id = ? AND status = ?", zoneID, "active").
			Count(&activeCount).Error; err != nil {
			return err
		}

		if int(activeCount) >= zone.TotalCapacity {
			return ErrZoneAtCapacity
		}

		if err := tx.Create(reservation).Error; err != nil {
			return err
		}

		return nil
	})
}

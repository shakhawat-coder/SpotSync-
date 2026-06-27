package service

import (
	"spotsync/dto"
	"spotsync/errors"
	"spotsync/models"
	"spotsync/repository"
)

type ParkingZoneService interface {
	CreateZone(req *dto.CreateZoneRequest) (*dto.ParkingZoneDetail, error)
	GetAllZones() ([]dto.ParkingZoneListItem, error)
	GetZoneByID(id uint) (*dto.ParkingZoneListItem, error)
}

type parkingZoneService struct {
	zoneRepo repository.ParkingZoneRepository
}

func NewParkingZoneService(zoneRepo repository.ParkingZoneRepository) ParkingZoneService {
	return &parkingZoneService{
		zoneRepo: zoneRepo,
	}
}

func (s *parkingZoneService) CreateZone(req *dto.CreateZoneRequest) (*dto.ParkingZoneDetail, error) {
	zone := &models.ParkingZone{
		Name:          req.Name,
		Type:          req.Type,
		TotalCapacity: req.TotalCapacity,
		PricePerHour:  req.PricePerHour,
	}

	if err := s.zoneRepo.Create(zone); err != nil {
		return nil, errors.ErrDatabaseError
	}

	return s.mapZoneToDetail(zone), nil
}

func (s *parkingZoneService) GetAllZones() ([]dto.ParkingZoneListItem, error) {
	zones, err := s.zoneRepo.GetAll()
	if err != nil {
		return nil, errors.ErrDatabaseError
	}

	var response []dto.ParkingZoneListItem
	for _, zone := range zones {
		activeCount, err := s.zoneRepo.GetActiveReservationCount(zone.ID)
		if err != nil {
			return nil, errors.ErrDatabaseError
		}

		response = append(response, *s.mapZoneToListItem(&zone, int(activeCount)))
	}

	return response, nil
}

func (s *parkingZoneService) GetZoneByID(id uint) (*dto.ParkingZoneListItem, error) {
	zone, err := s.zoneRepo.GetByID(id)
	if err != nil {
		return nil, errors.ErrDatabaseError
	}
	if zone == nil {
		return nil, errors.ErrZoneNotFound
	}

	activeCount, err := s.zoneRepo.GetActiveReservationCount(id)
	if err != nil {
		return nil, errors.ErrDatabaseError
	}

	return s.mapZoneToListItem(zone, int(activeCount)), nil
}

func (s *parkingZoneService) mapZoneToListItem(zone *models.ParkingZone, activeReservations int) *dto.ParkingZoneListItem {
	availableSpots := zone.TotalCapacity - activeReservations
	if availableSpots < 0 {
		availableSpots = 0
	}

	return &dto.ParkingZoneListItem{
		ID:             zone.ID,
		Name:           zone.Name,
		Type:           zone.Type,
		TotalCapacity:  zone.TotalCapacity,
		AvailableSpots: availableSpots,
		PricePerHour:   zone.PricePerHour,
		CreatedAt:      zone.CreatedAt,
	}
}

func (s *parkingZoneService) mapZoneToDetail(zone *models.ParkingZone) *dto.ParkingZoneDetail {
	return &dto.ParkingZoneDetail{
		ID:            zone.ID,
		Name:          zone.Name,
		Type:          zone.Type,
		TotalCapacity: zone.TotalCapacity,
		PricePerHour:  zone.PricePerHour,
		CreatedAt:     zone.CreatedAt,
		UpdatedAt:     zone.UpdatedAt,
	}
}

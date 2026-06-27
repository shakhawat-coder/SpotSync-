package service
 
import (
	"spotsync/dto"
	"spotsync/errors"
	"spotsync/models"
	"spotsync/repository"
)
 
type ParkingZoneService interface {
	CreateZone(req *dto.CreateZoneRequest) (*dto.ParkingZoneResponse, error)
	GetAllZones() ([]dto.ParkingZoneResponse, error)
	GetZoneByID(id uint) (*dto.ParkingZoneResponse, error)
	UpdateZone(id uint, req *dto.CreateZoneRequest) (*dto.ParkingZoneResponse, error)
	DeleteZone(id uint) error
}
 
type parkingZoneService struct {
	zoneRepo repository.ParkingZoneRepository
}
 
func NewParkingZoneService(zoneRepo repository.ParkingZoneRepository) ParkingZoneService {
	return &parkingZoneService{
		zoneRepo: zoneRepo,
	}
}
 
func (s *parkingZoneService) CreateZone(req *dto.CreateZoneRequest) (*dto.ParkingZoneResponse, error) {
	zone := &models.ParkingZone{
		Name:          req.Name,
		Type:          req.Type,
		TotalCapacity: req.TotalCapacity,
		PricePerHour:  req.PricePerHour,
	}
 
	if err := s.zoneRepo.Create(zone); err != nil {
		return nil, errors.ErrDatabaseError
	}
 
	return s.mapZoneToResponse(zone, 0), nil
}
 
func (s *parkingZoneService) GetAllZones() ([]dto.ParkingZoneResponse, error) {
	zones, err := s.zoneRepo.GetAll()
	if err != nil {
		return nil, errors.ErrDatabaseError
	}
 
	var response []dto.ParkingZoneResponse
	for _, zone := range zones {
		// Calculate available spots for each zone
		activeCount, err := s.zoneRepo.GetActiveReservationCount(zone.ID)
		if err != nil {
			return nil, errors.ErrDatabaseError
		}
 
		response = append(response, *s.mapZoneToResponse(&zone, int(activeCount)))
	}
 
	return response, nil
}
 
func (s *parkingZoneService) GetZoneByID(id uint) (*dto.ParkingZoneResponse, error) {
	zone, err := s.zoneRepo.GetByID(id)
	if err != nil {
		return nil, errors.ErrDatabaseError
	}
	if zone == nil {
		return nil, errors.ErrZoneNotFound
	}
 
	// Count active reservations for this zone
	activeCount, err := s.zoneRepo.GetActiveReservationCount(id)
	if err != nil {
		return nil, errors.ErrDatabaseError
	}
 
	return s.mapZoneToResponse(zone, int(activeCount)), nil
}
 
func (s *parkingZoneService) UpdateZone(id uint, req *dto.CreateZoneRequest) (*dto.ParkingZoneResponse, error) {
	zone, err := s.zoneRepo.GetByID(id)
	if err != nil {
		return nil, errors.ErrDatabaseError
	}
	if zone == nil {
		return nil, errors.ErrZoneNotFound
	}
 
	zone.Name = req.Name
	zone.Type = req.Type
	zone.TotalCapacity = req.TotalCapacity
	zone.PricePerHour = req.PricePerHour
 
	if err := s.zoneRepo.Update(zone); err != nil {
		return nil, errors.ErrDatabaseError
	}
 
	return s.mapZoneToResponse(zone, 0), nil
}
 
func (s *parkingZoneService) DeleteZone(id uint) error {
	zone, err := s.zoneRepo.GetByID(id)
	if err != nil {
		return errors.ErrDatabaseError
	}
	if zone == nil {
		return errors.ErrZoneNotFound
	}
 
	if err := s.zoneRepo.Delete(id); err != nil {
		return errors.ErrDatabaseError
	}
 
	return nil
}
 
func (s *parkingZoneService) mapZoneToResponse(zone *models.ParkingZone, activeReservations int) *dto.ParkingZoneResponse {
	availableSpots := zone.TotalCapacity - activeReservations
	if availableSpots < 0 {
		availableSpots = 0
	}
 
	return &dto.ParkingZoneResponse{
		ID:             zone.ID,
		Name:           zone.Name,
		Type:           zone.Type,
		TotalCapacity:  zone.TotalCapacity,
		AvailableSpots: availableSpots,
		PricePerHour:   zone.PricePerHour,
		CreatedAt:      zone.CreatedAt,
	}
}
 
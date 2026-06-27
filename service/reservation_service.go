package service

import (
	"errors"

	"spotsync/dto"
	appErrors "spotsync/errors"
	"spotsync/models"
	"spotsync/repository"
)

type ReservationService interface {
	CreateReservation(req *dto.CreateReservationRequest, userID uint) (*dto.ReservationCreateResponse, error)
	GetMyReservations(userID uint) ([]dto.MyReservationResponse, error)
	GetAllReservations() ([]dto.AdminReservationResponse, error)
	CancelReservation(reservationID, userID uint) error
}

type reservationService struct {
	reservationRepo repository.ReservationRepository
	zoneRepo        repository.ParkingZoneRepository
}

func NewReservationService(
	reservationRepo repository.ReservationRepository,
	zoneRepo repository.ParkingZoneRepository,
) ReservationService {
	return &reservationService{
		reservationRepo: reservationRepo,
		zoneRepo:        zoneRepo,
	}
}

func (s *reservationService) CreateReservation(
	req *dto.CreateReservationRequest,
	userID uint,
) (*dto.ReservationCreateResponse, error) {
	zone, err := s.zoneRepo.GetByID(req.ZoneID)
	if err != nil {
		return nil, appErrors.ErrDatabaseError
	}
	if zone == nil {
		return nil, appErrors.ErrZoneNotFound
	}

	reservation := &models.Reservation{
		UserID:       userID,
		ZoneID:       req.ZoneID,
		LicensePlate: req.LicensePlate,
		Status:       "active",
	}

	if err := s.reservationRepo.CreateWithLocking(reservation, req.ZoneID); err != nil {
		if errors.Is(err, repository.ErrZoneAtCapacity) {
			return nil, appErrors.ErrZoneFull
		}
		if errors.Is(err, repository.ErrDuplicateActiveLicense) {
			return nil, appErrors.ErrDuplicateLicensePlate
		}
		return nil, appErrors.ErrDatabaseError
	}

	return &dto.ReservationCreateResponse{
		ID:           reservation.ID,
		UserID:       reservation.UserID,
		ZoneID:       reservation.ZoneID,
		LicensePlate: reservation.LicensePlate,
		Status:       reservation.Status,
		CreatedAt:    reservation.CreatedAt,
		UpdatedAt:    reservation.UpdatedAt,
	}, nil
}

func (s *reservationService) GetMyReservations(userID uint) ([]dto.MyReservationResponse, error) {
	reservations, err := s.reservationRepo.GetByUserID(userID)
	if err != nil {
		return nil, appErrors.ErrDatabaseError
	}

	var response []dto.MyReservationResponse
	for _, res := range reservations {
		response = append(response, s.mapMyReservation(&res))
	}

	return response, nil
}

func (s *reservationService) GetAllReservations() ([]dto.AdminReservationResponse, error) {
	reservations, err := s.reservationRepo.GetAll()
	if err != nil {
		return nil, appErrors.ErrDatabaseError
	}

	var response []dto.AdminReservationResponse
	for _, res := range reservations {
		response = append(response, s.mapAdminReservation(&res))
	}

	return response, nil
}

func (s *reservationService) CancelReservation(reservationID, userID uint) error {
	reservation, err := s.reservationRepo.GetByID(reservationID)
	if err != nil {
		return appErrors.ErrDatabaseError
	}
	if reservation == nil {
		return appErrors.ErrReservationNotFound
	}

	if reservation.UserID != userID {
		return appErrors.ErrForbidden
	}

	if err := s.reservationRepo.Cancel(reservationID); err != nil {
		return appErrors.ErrDatabaseError
	}

	return nil
}

func (s *reservationService) mapMyReservation(res *models.Reservation) dto.MyReservationResponse {
	zoneDTO := dto.ZoneDTO{}
	if res.Zone != nil {
		zoneDTO = dto.ZoneDTO{
			ID:   res.Zone.ID,
			Name: res.Zone.Name,
			Type: res.Zone.Type,
		}
	}

	return dto.MyReservationResponse{
		ID:           res.ID,
		LicensePlate: res.LicensePlate,
		Status:       res.Status,
		Zone:         zoneDTO,
		CreatedAt:    res.CreatedAt,
	}
}

func (s *reservationService) mapAdminReservation(res *models.Reservation) dto.AdminReservationResponse {
	zoneDTO := dto.ZoneDTO{}
	if res.Zone != nil {
		zoneDTO = dto.ZoneDTO{
			ID:   res.Zone.ID,
			Name: res.Zone.Name,
			Type: res.Zone.Type,
		}
	}

	userDTO := dto.UserDTO{}
	if res.User != nil {
		userDTO = dto.UserDTO{
			ID:    res.User.ID,
			Name:  res.User.Name,
			Email: res.User.Email,
			Role:  res.User.Role,
		}
	}

	return dto.AdminReservationResponse{
		ID:           res.ID,
		UserID:       res.UserID,
		ZoneID:       res.ZoneID,
		LicensePlate: res.LicensePlate,
		Status:       res.Status,
		User:         userDTO,
		Zone:         zoneDTO,
		CreatedAt:    res.CreatedAt,
		UpdatedAt:    res.UpdatedAt,
	}
}

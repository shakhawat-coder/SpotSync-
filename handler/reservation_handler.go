package handler

import (
	"net/http"
	"strconv"
	"spotsync/dto"
	"spotsync/errors"
	"spotsync/middleware"
	"spotsync/service"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

type ReservationHandler struct {
	reservationService service.ReservationService
	validator          *validator.Validate
}

func NewReservationHandler(reservationService service.ReservationService, validator *validator.Validate) *ReservationHandler {
	return &ReservationHandler{
		reservationService: reservationService,
		validator:          validator,
	}
}

func (h *ReservationHandler) CreateReservation(c *echo.Context) error {
	userID := middleware.GetUserIDFromContext(c)
	if userID == 0 {
		return respondError(c, http.StatusUnauthorized, errors.ErrUnauthorized.Message, "User not authenticated")
	}

	req := new(dto.CreateReservationRequest)
	if err := c.Bind(req); err != nil {
		return respondError(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := h.validator.Struct(req); err != nil {
		return respondError(c, http.StatusBadRequest, "Validation failed", err.Error())
	}

	reservation, err := h.reservationService.CreateReservation(req, userID)
	if err != nil {
		return respondAppError(c, err, "Failed to create reservation")
	}

	return respondSuccess(c, http.StatusCreated, "Reservation confirmed successfully", reservation)
}

func (h *ReservationHandler) GetMyReservations(c *echo.Context) error {
	userID := middleware.GetUserIDFromContext(c)
	if userID == 0 {
		return respondError(c, http.StatusUnauthorized, errors.ErrUnauthorized.Message, "User not authenticated")
	}

	reservations, err := h.reservationService.GetMyReservations(userID)
	if err != nil {
		return respondAppError(c, err, "Failed to fetch reservations")
	}

	return respondSuccess(c, http.StatusOK, "My reservations retrieved successfully", reservations)
}

func (h *ReservationHandler) GetAllReservations(c *echo.Context) error {
	if err := requireAdmin(c); err != nil {
		return err
	}

	reservations, err := h.reservationService.GetAllReservations()
	if err != nil {
		return respondAppError(c, err, "Failed to fetch reservations")
	}

	return respondSuccess(c, http.StatusOK, "All reservations retrieved successfully", reservations)
}

func (h *ReservationHandler) CancelReservation(c *echo.Context) error {
	userID := middleware.GetUserIDFromContext(c)
	if userID == 0 {
		return respondError(c, http.StatusUnauthorized, errors.ErrUnauthorized.Message, "User not authenticated")
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return respondError(c, http.StatusBadRequest, "Invalid reservation ID", err.Error())
	}

	if err := h.reservationService.CancelReservation(uint(id), userID); err != nil {
		return respondAppError(c, err, "Failed to cancel reservation")
	}

	return respondSuccess(c, http.StatusOK, "Reservation cancelled successfully", nil)
}

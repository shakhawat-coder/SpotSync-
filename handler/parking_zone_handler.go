package handler

import (
	"net/http"
	"spotsync/dto"
	"spotsync/service"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

type ParkingZoneHandler struct {
	zoneService service.ParkingZoneService
	validator   *validator.Validate
}

func NewParkingZoneHandler(zoneService service.ParkingZoneService, validator *validator.Validate) *ParkingZoneHandler {
	return &ParkingZoneHandler{
		zoneService: zoneService,
		validator:   validator,
	}
}

func (h *ParkingZoneHandler) CreateZone(c *echo.Context) error {
	if err := requireAdmin(c); err != nil {
		return err
	}

	req := new(dto.CreateZoneRequest)
	if err := c.Bind(req); err != nil {
		return respondError(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := h.validator.Struct(req); err != nil {
		return respondError(c, http.StatusBadRequest, "Validation failed", err.Error())
	}

	zone, err := h.zoneService.CreateZone(req)
	if err != nil {
		return respondAppError(c, err, "Failed to create zone")
	}

	return respondSuccess(c, http.StatusCreated, "Parking zone created successfully", zone)
}

func (h *ParkingZoneHandler) GetAllZones(c *echo.Context) error {
	zones, err := h.zoneService.GetAllZones()
	if err != nil {
		return respondAppError(c, err, "Failed to fetch zones")
	}

	return respondSuccess(c, http.StatusOK, "Parking zones retrieved successfully", zones)
}

func (h *ParkingZoneHandler) GetZoneByID(c *echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return respondError(c, http.StatusBadRequest, "Invalid zone ID", err.Error())
	}

	zone, err := h.zoneService.GetZoneByID(uint(id))
	if err != nil {
		return respondAppError(c, err, "Failed to fetch zone")
	}

	return respondSuccess(c, http.StatusOK, "Parking zone retrieved successfully", zone)
}


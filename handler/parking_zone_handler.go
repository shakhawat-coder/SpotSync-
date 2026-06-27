package handler

import (
	"net/http"
	"spotsync/dto"
	"spotsync/errors"
	"spotsync/middleware"
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

// POST /api/v1/zones (Admin only)
func (h *ParkingZoneHandler) CreateZone(c *echo.Context) error {
	// Check admin role
	role := middleware.GetRoleFromContext(c)
	if role != "admin" {
		return c.JSON(http.StatusForbidden, map[string]interface{}{
			"success": false,
			"message": errors.ErrForbidden.Message,
			"errors":  "Admin access required",
		})
	}

	req := new(dto.CreateZoneRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid request body",
			"errors":  err.Error(),
		})
	}

	// Validate
	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Validation failed",
			"errors":  err.Error(),
		})
	}

	zone, err := h.zoneService.CreateZone(req)
	if err != nil {
		appErr, ok := err.(*errors.AppError)
		if ok {
			return c.JSON(appErr.Status, map[string]interface{}{
				"success": false,
				"message": appErr.Message,
				"errors":  appErr.Code,
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to create zone",
			"errors":  err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "Parking zone created successfully",
		"data":    zone,
	})
}

// GET /api/v1/zones (Public)
func (h *ParkingZoneHandler) GetAllZones(c *echo.Context) error {
	zones, err := h.zoneService.GetAllZones()
	if err != nil {
		appErr, ok := err.(*errors.AppError)
		if ok {
			return c.JSON(appErr.Status, map[string]interface{}{
				"success": false,
				"message": appErr.Message,
				"errors":  appErr.Code,
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to fetch zones",
			"errors":  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Parking zones retrieved successfully",
		"data":    zones,
	})
}

// GET /api/v1/zones/:id (Public)
func (h *ParkingZoneHandler) GetZoneByID(c *echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid zone ID",
			"errors":  err.Error(),
		})
	}

	zone, err := h.zoneService.GetZoneByID(uint(id))
	if err != nil {
		appErr, ok := err.(*errors.AppError)
		if ok {
			return c.JSON(appErr.Status, map[string]interface{}{
				"success": false,
				"message": appErr.Message,
				"errors":  appErr.Code,
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to fetch zone",
			"errors":  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Parking zone retrieved successfully",
		"data":    zone,
	})
}

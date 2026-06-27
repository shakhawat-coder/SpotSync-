package handler
 
import (
	"net/http"
	"spotsync/dto"
	"spotsync/errors"
	"spotsync/service"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

type AuthHandler struct {
	authService service.AuthService
	validator   *validator.Validate
}
 
func NewAuthHandler(authService service.AuthService, validator *validator.Validate) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validator:   validator,
	}
}
 
// POST /api/v1/auth/register
func (h *AuthHandler) Register(c echo.Context) error {
	req := new(dto.RegisterRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid request body",
			"errors":  err.Error(),
		})
	}
 
	// Validate request
	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Validation failed",
			"errors":  err.Error(),
		})
	}
 
	// Call service
	user, err := h.authService.Register(req)
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
			"message": "Registration failed",
			"errors":  err.Error(),
		})
	}
 
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "User registered successfully",
		"data":    user,
	})
}
 
// POST /api/v1/auth/login
func (h *AuthHandler) Login(c echo.Context) error {
	req := new(dto.LoginRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid request body",
			"errors":  err.Error(),
		})
	}
 
	// Validate request
	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Validation failed",
			"errors":  err.Error(),
		})
	}
 
	// Call service
	response, err := h.authService.Login(req)
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
			"message": "Login failed",
			"errors":  err.Error(),
		})
	}
 
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Login successful",
		"data":    response,
	})
}
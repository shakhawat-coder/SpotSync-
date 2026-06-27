package handler

import (
	"net/http"
	"spotsync/dto"
	"spotsync/errors"
	"spotsync/middleware"
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

func (h *AuthHandler) Register(c *echo.Context) error {
	req := new(dto.RegisterRequest)
	if err := c.Bind(req); err != nil {
		return respondError(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := h.validator.Struct(req); err != nil {
		return respondError(c, http.StatusBadRequest, "Validation failed", err.Error())
	}

	user, err := h.authService.Register(req)
	if err != nil {
		return respondAppError(c, err, "Registration failed")
	}

	return respondSuccess(c, http.StatusCreated, "User registered successfully", user)
}

func (h *AuthHandler) Login(c *echo.Context) error {
	req := new(dto.LoginRequest)
	if err := c.Bind(req); err != nil {
		return respondError(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := h.validator.Struct(req); err != nil {
		return respondError(c, http.StatusBadRequest, "Validation failed", err.Error())
	}

	response, err := h.authService.Login(req)
	if err != nil {
		return respondAppError(c, err, "Login failed")
	}

	return respondSuccess(c, http.StatusOK, "Login successful", response)
}

func requireAdmin(c *echo.Context) error {
	if middleware.GetRoleFromContext(c) != "admin" {
		return respondError(c, http.StatusForbidden, errors.ErrForbidden.Message, "Admin access required")
	}
	return nil
}

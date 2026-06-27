package handler

import (
	"net/http"

	"spotsync/errors"

	"github.com/labstack/echo/v5"
)

type JSONResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

func respondError(c *echo.Context, status int, message, details string) error {
	return c.JSON(status, JSONResponse{
		Success: false,
		Message: message,
		Errors:  details,
	})
}

func respondAppError(c *echo.Context, err error, fallbackMessage string) error {
	if appErr, ok := err.(*errors.AppError); ok {
		return respondError(c, appErr.Status, appErr.Message, appErr.Code)
	}
	return respondError(c, http.StatusInternalServerError, fallbackMessage, err.Error())
}

func respondSuccess(c *echo.Context, status int, message string, data interface{}) error {
	return c.JSON(status, JSONResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

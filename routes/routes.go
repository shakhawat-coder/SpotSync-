package routes

import (
	"spotsync/config"
	"spotsync/handler"
	"spotsync/middleware"

	"github.com/labstack/echo/v5"
)

type RouterConfig struct {
	Echo           *echo.Echo
	AuthHandler    *handler.AuthHandler
	ZoneHandler    *handler.ParkingZoneHandler
	ReserveHandler *handler.ReservationHandler
	Config         *config.Config
}

func SetupRoutes(cfg RouterConfig) {
	e := cfg.Echo
	jwt := middleware.JWTMiddleware(cfg.Config.JWTSecret)

	e.GET("/", func(c *echo.Context) error {
		return c.JSON(200, map[string]interface{}{
			"message": "Welcome to SpotSync API",
			"status":  "Server is running",
			"version": "1.0.0",
		})
	})

	api := e.Group("/api/v1")

	auth := api.Group("/auth")
	{
		auth.POST("/register", cfg.AuthHandler.Register)
		auth.POST("/login", cfg.AuthHandler.Login)
	}

	zones := api.Group("/zones")
	{
		zones.GET("", cfg.ZoneHandler.GetAllZones)
		zones.GET("/:id", cfg.ZoneHandler.GetZoneByID)

		adminZones := zones.Group("")
		adminZones.Use(jwt)
		adminZones.POST("", cfg.ZoneHandler.CreateZone)
	}

	reservations := api.Group("/reservations")
	reservations.Use(jwt)
	{
		reservations.POST("", cfg.ReserveHandler.CreateReservation)
		reservations.GET("/my-reservations", cfg.ReserveHandler.GetMyReservations)
		reservations.GET("", cfg.ReserveHandler.GetAllReservations)
		reservations.DELETE("/:id", cfg.ReserveHandler.CancelReservation)
	}
}

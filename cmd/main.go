package main

import (
	"log"
	"os"

	"spotsync/config"
	"spotsync/handler"
	"spotsync/migrations"
	"spotsync/repository"
	"spotsync/routes"
	"spotsync/service"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	echoMiddleware "github.com/labstack/echo/v5/middleware"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}
	cfg := config.LoadConfig()

	if cfg.DatabaseURL == "" || cfg.JWTSecret == "" {
		log.Fatal("Missing required environment variables: DATABASE_URL, JWT_SECRET")
	}
	db, err := cfg.GetDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := migrations.Migrate(db); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	validate := validator.New()

	// ═══════════════════════════════════════════════════════════════════════
	// 🔌 DEPENDENCY INJECTION WIRING (Clean Architecture)
	// ═══════════════════════════════════════════════════════════════════════

	// Repositories: Data Access Layer
	userRepo := repository.NewUserRepository(db)
	parkingZoneRepo := repository.NewParkingZoneRepository(db)
	reservationRepo := repository.NewReservationRepository(db)

	// Services: Business Logic Layer
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	parkingZoneService := service.NewParkingZoneService(parkingZoneRepo)
	reservationService := service.NewReservationService(reservationRepo, parkingZoneRepo)

	// Handlers: HTTP Layer
	authHandler := handler.NewAuthHandler(authService, validate)
	parkingZoneHandler := handler.NewParkingZoneHandler(parkingZoneService, validate)
	reservationHandler := handler.NewReservationHandler(reservationService, validate)

	// ═══════════════════════════════════════════════════════════════════════
	// 🌐 ECHO SERVER SETUP
	// ═══════════════════════════════════════════════════════════════════════

	e := echo.New()

	// Middleware
	e.Use(echoMiddleware.RequestLogger())
	e.Use(echoMiddleware.Recover())

	// CORS Middleware
	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: []string{"*"}, // Configure for production
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
	}))

	// ═══════════════════════════════════════════════════════════════════════
	// 🛣️ API ROUTES
	// ═══════════════════════════════════════════════════════════════════════

	routes.SetupRoutes(routes.RouterConfig{
		Echo:           e,
		AuthHandler:    authHandler,
		ZoneHandler:    parkingZoneHandler,
		ReserveHandler: reservationHandler,
		Config:         cfg,
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server starting on port %s\n", port)
	if err := e.Start(":" + port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

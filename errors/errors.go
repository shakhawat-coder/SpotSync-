package errors

type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"`
}

func (e *AppError) Error() string {
	return e.Message
}

// Centralized Error Definitions

var (
	ErrInvalidEmail          = &AppError{Code: "INVALID_EMAIL", Message: "Invalid email format", Status: 400}
	ErrWeakPassword          = &AppError{Code: "WEAK_PASSWORD", Message: "Password must be at least 8 characters", Status: 400}
	ErrUserExists            = &AppError{Code: "USER_EXISTS", Message: "User with this email already exists", Status: 400}
	ErrUserNotFound          = &AppError{Code: "USER_NOT_FOUND", Message: "User not found", Status: 404}
	ErrInvalidCredentials    = &AppError{Code: "INVALID_CREDENTIALS", Message: "Invalid email or password", Status: 401}
	ErrUnauthorized          = &AppError{Code: "UNAUTHORIZED", Message: "Missing or invalid token", Status: 401}
	ErrForbidden             = &AppError{Code: "FORBIDDEN", Message: "Insufficient permissions", Status: 403}
	ErrZoneNotFound          = &AppError{Code: "ZONE_NOT_FOUND", Message: "Parking zone not found", Status: 404}
	ErrZoneFull              = &AppError{Code: "ZONE_FULL", Message: "No available spots in this zone", Status: 409}
	ErrReservationNotFound   = &AppError{Code: "RESERVATION_NOT_FOUND", Message: "Reservation not found", Status: 404}
	ErrInvalidRole           = &AppError{Code: "INVALID_ROLE", Message: "Role must be 'driver' or 'admin'", Status: 400}
	ErrDuplicateLicensePlate = &AppError{Code: "DUPLICATE_LICENSE", Message: "License plate already has an active reservation", Status: 409}
	ErrDatabaseError         = &AppError{Code: "DATABASE_ERROR", Message: "Database operation failed", Status: 500}
	ErrInternalServer        = &AppError{Code: "INTERNAL_ERROR", Message: "Internal server error", Status: 500}
)

func NewAppError(code, message string, status int) *AppError {
	return &AppError{Code: code, Message: message, Status: status}
}

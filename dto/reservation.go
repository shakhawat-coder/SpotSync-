package dto
 
import "time"
 
type CreateReservationRequest struct {
	ZoneID       uint   `json:"zone_id" validate:"required,gt=0"`
	LicensePlate string `json:"license_plate" validate:"required,max=15"`
}
 
type ZoneDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}
 
type ReservationResponse struct {
	ID           uint      `json:"id"`
	UserID       uint      `json:"user_id"`
	ZoneID       uint      `json:"zone_id"`
	LicensePlate string    `json:"license_plate"`
	Status       string    `json:"status"`
	Zone         ZoneDTO   `json:"zone,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
 
type ReservationListResponse struct {
	Success bool                      `json:"success"`
	Message string                    `json:"message"`
	Data    []ReservationResponse     `json:"data"`
}
 
type ReservationSingleResponse struct {
	Success bool                    `json:"success"`
	Message string                  `json:"message"`
	Data    ReservationResponse     `json:"data"`
}
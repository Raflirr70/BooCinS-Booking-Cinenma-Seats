package dto

type CreateBookingRequest struct {
	ScheduleID uint   `json:"schedule_id" binding:"required"`
	SeatIDs    []uint `json:"seat_ids" binding:"required,min=1,max=10"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
}

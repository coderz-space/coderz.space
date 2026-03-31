package mentorship

type RequestRoleRequest struct {
	Role string `json:"role" validate:"required,oneof=mentor mentee"`
}

type UpdateStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=active inactive"`
}

type MenteeRequestDTO struct {
	ID         string `json:"id"`
	FirstName  string `json:"firstName"`
	Email      string `json:"email"`
	Status     string `json:"status"`
	SignedUpAt string `json:"signedUpAt"`
}

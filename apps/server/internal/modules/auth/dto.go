package auth

type SignInRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=50"`
	Name     string `json:"name" validate:"required,min=2,max=50"`
}

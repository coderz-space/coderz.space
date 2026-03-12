package auth

import "github.com/DSAwithGautam/Coderz.space/internal/config"

type Handler struct {
	service *Service
	config  *config.Config
}

func NewHandler(service *Service, config *config.Config) *Handler {
	return &Handler{
		service: service,
		config:  config,
	}
}

// func (h *Handler) SignIn(c *echo.Context) error {
// 	var req LoginRequest
// 	if err := c.Bind(&req); err != nil {
// 		return response.NewResponse(c, http.StatusBadRequest, "Invalid request", nil, err, "error")
// 	}
// 	return response.NewResponse(c, http.StatusOK, "Login successful", nil, nil, "success")
// }

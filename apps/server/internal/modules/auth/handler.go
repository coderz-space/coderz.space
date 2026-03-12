package auth

import (
	"github.com/labstack/echo/v5"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) SignIn(c *echo.Context, body SignInRequest) error {
	return nil
}

func (h *Handler) Register(c *echo.Context, body RegisterRequest) error {
	return nil
}

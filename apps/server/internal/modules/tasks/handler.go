package tasks

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v5"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ListMenteeQuestions(c *echo.Context) error {
	id := (*c).Param("id")
	res, err := h.service.ListMenteeQuestions(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, res)
}

func (h *Handler) AssignQuestionToMentee(c *echo.Context, req AssignQuestionRequest) error {
	id := (*c).Param("id")
	
	// Assuming the admin/mentor is identifying themselves too, omitted for MVP
	
	res, err := h.service.AssignQuestionToMentee(c.Request().Context(), id, req)
	if err != nil {
		log.Printf("Failed to assign question: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, res)
}

func (h *Handler) GetMenteeQuestion(c *echo.Context) error {
	id := (*c).Param("id")
	questionId := (*c).Param("questionId")
	
	res, err := h.service.GetMenteeQuestion(c.Request().Context(), id, questionId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, res)
}

func (h *Handler) UpdateQuestionProgress(c *echo.Context, req UpdateProgressRequest) error {
	id := (*c).Param("id")
	questionId := (*c).Param("questionId")
	
	err := h.service.UpdateQuestionProgress(c.Request().Context(), id, questionId, req.ProgressStatus)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]bool{"success": true})
}

func (h *Handler) UpdateQuestionDetails(c *echo.Context, req UpdateDetailsRequest) error {
	id := (*c).Param("id")
	questionId := (*c).Param("questionId")
	
	err := h.service.UpdateQuestionDetails(c.Request().Context(), id, questionId, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]bool{"success": true})
}

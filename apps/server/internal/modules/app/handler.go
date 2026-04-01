package app

import (
	"net/http"

	authmw "github.com/coderz-space/coderz.space/internal/common/middleware/auth"
	"github.com/coderz-space/coderz.space/internal/common/response"
	"github.com/coderz-space/coderz.space/internal/common/utils"
	"github.com/coderz-space/coderz.space/internal/common/validator"
	"github.com/labstack/echo/v5"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetContext(c *echo.Context) error {
	userID, err := currentUserID(c)
	if err != nil {
		return unauthorizedResponse(c, err)
	}

	data, err := h.service.GetContext(c.Request().Context(), userID)
	if err != nil {
		return handleAppError(c, err)
	}

	return response.NewResponse(c, http.StatusOK, "OK", "APP_CONTEXT_RETRIEVED", data, nil)
}

func (h *Handler) MenteeSignup(c *echo.Context) error {
	var body MenteeSignupRequest
	if err := (&echo.DefaultBinder{}).Bind(c, &body); err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_REQUEST_BODY", nil, err)
	}
	if err := validator.NewValidator().ValidateStruct(body); err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "VALIDATION_FAILED", nil, err)
	}

	data, err := h.service.MenteeSignup(c.Request().Context(), body)
	if err != nil {
		return handleAppError(c, err)
	}

	return response.NewResponse(c, http.StatusCreated, "CREATED", "MENTEE_SIGNUP_REQUEST_CREATED", data, nil)
}

func (h *Handler) ListMenteeRequests(c *echo.Context) error {
	userID, err := currentUserID(c)
	if err != nil {
		return unauthorizedResponse(c, err)
	}

	data, err := h.service.ListMenteeRequests(c.Request().Context(), userID)
	if err != nil {
		return handleAppError(c, err)
	}

	return response.NewResponse(c, http.StatusOK, "OK", "MENTEE_REQUESTS_RETRIEVED", data, nil)
}

func (h *Handler) ReviewMenteeRequest(c *echo.Context) error {
	userID, err := currentUserID(c)
	if err != nil {
		return unauthorizedResponse(c, err)
	}

	var body ReviewMenteeRequest
	if err := (&echo.DefaultBinder{}).Bind(c, &body); err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_REQUEST_BODY", nil, err)
	}
	if err := validator.NewValidator().ValidateStruct(body); err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "VALIDATION_FAILED", nil, err)
	}

	data, err := h.service.ReviewMenteeRequest(c.Request().Context(), userID, (*c).Param("requestId"), body)
	if err != nil {
		return handleAppError(c, err)
	}

	return response.NewResponse(c, http.StatusOK, "OK", "MENTEE_REQUEST_UPDATED", data, nil)
}

func (h *Handler) ListSheets(c *echo.Context) error {
	return response.NewResponse(c, http.StatusOK, "OK", "SHEETS_RETRIEVED", h.service.ListSheets(), nil)
}

func (h *Handler) GetDayAssignments(c *echo.Context) error {
	userID, err := currentUserID(c)
	if err != nil {
		return unauthorizedResponse(c, err)
	}

	data, err := h.service.GetDayAssignments(c.Request().Context(), userID, (*c).Param("day"))
	if err != nil {
		return handleAppError(c, err)
	}

	return response.NewResponse(c, http.StatusOK, "OK", "DAY_ASSIGNMENTS_RETRIEVED", data, nil)
}

func (h *Handler) UpdateDayAssignments(c *echo.Context) error {
	userID, err := currentUserID(c)
	if err != nil {
		return unauthorizedResponse(c, err)
	}

	var body UpdateDayAssignmentsRequest
	if err := (&echo.DefaultBinder{}).Bind(c, &body); err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_REQUEST_BODY", nil, err)
	}
	if err := validator.NewValidator().ValidateStruct(body); err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "VALIDATION_FAILED", nil, err)
	}

	data, err := h.service.UpdateDayAssignments(c.Request().Context(), userID, (*c).Param("day"), body)
	if err != nil {
		return handleAppError(c, err)
	}

	return response.NewResponse(c, http.StatusOK, "OK", "DAY_ASSIGNMENTS_UPDATED", data, nil)
}

func (h *Handler) CreateAssignments(c *echo.Context) error {
	userID, err := currentUserID(c)
	if err != nil {
		return unauthorizedResponse(c, err)
	}

	var body CreateAssignmentsRequest
	if err := (&echo.DefaultBinder{}).Bind(c, &body); err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_REQUEST_BODY", nil, err)
	}
	if err := validator.NewValidator().ValidateStruct(body); err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "VALIDATION_FAILED", nil, err)
	}

	data, err := h.service.CreateAssignments(c.Request().Context(), userID, body)
	if err != nil {
		return handleAppError(c, err)
	}

	return response.NewResponse(c, http.StatusCreated, "CREATED", "ASSIGNMENTS_CREATED", data, nil)
}

func (h *Handler) ListMenteeQuestions(c *echo.Context) error {
	userID, err := currentUserID(c)
	if err != nil {
		return unauthorizedResponse(c, err)
	}

	data, err := h.service.ListMenteeQuestions(c.Request().Context(), userID, (*c).Param("username"))
	if err != nil {
		return handleAppError(c, err)
	}

	return response.NewResponse(c, http.StatusOK, "OK", "QUESTIONS_RETRIEVED", data, nil)
}

func (h *Handler) GetMenteeQuestion(c *echo.Context) error {
	userID, err := currentUserID(c)
	if err != nil {
		return unauthorizedResponse(c, err)
	}

	data, err := h.service.GetMenteeQuestion(c.Request().Context(), userID, (*c).Param("username"), (*c).Param("assignmentProblemId"))
	if err != nil {
		return handleAppError(c, err)
	}

	return response.NewResponse(c, http.StatusOK, "OK", "QUESTION_RETRIEVED", data, nil)
}

func (h *Handler) UpdateMenteeQuestion(c *echo.Context) error {
	userID, err := currentUserID(c)
	if err != nil {
		return unauthorizedResponse(c, err)
	}

	var body UpdateQuestionRequest
	if err := (&echo.DefaultBinder{}).Bind(c, &body); err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_REQUEST_BODY", nil, err)
	}
	if err := validator.NewValidator().ValidateStruct(body); err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "VALIDATION_FAILED", nil, err)
	}

	data, err := h.service.UpdateMenteeQuestion(c.Request().Context(), userID, (*c).Param("username"), (*c).Param("assignmentProblemId"), body)
	if err != nil {
		return handleAppError(c, err)
	}

	return response.NewResponse(c, http.StatusOK, "OK", "QUESTION_UPDATED", data, nil)
}

func (h *Handler) GetMenteeProfile(c *echo.Context) error {
	userID, err := currentUserID(c)
	if err != nil {
		return unauthorizedResponse(c, err)
	}

	data, err := h.service.GetMenteeProfile(c.Request().Context(), userID, (*c).Param("username"))
	if err != nil {
		return handleAppError(c, err)
	}

	return response.NewResponse(c, http.StatusOK, "OK", "PROFILE_RETRIEVED", data, nil)
}

func (h *Handler) GetMyProfile(c *echo.Context) error {
	userID, err := currentUserID(c)
	if err != nil {
		return unauthorizedResponse(c, err)
	}

	data, err := h.service.GetMyProfile(c.Request().Context(), userID)
	if err != nil {
		return handleAppError(c, err)
	}

	return response.NewResponse(c, http.StatusOK, "OK", "PROFILE_RETRIEVED", data, nil)
}

func (h *Handler) UpdateMyProfile(c *echo.Context) error {
	userID, err := currentUserID(c)
	if err != nil {
		return unauthorizedResponse(c, err)
	}

	var body UpdateProfileRequest
	if err := (&echo.DefaultBinder{}).Bind(c, &body); err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_REQUEST_BODY", nil, err)
	}
	if err := validator.NewValidator().ValidateStruct(body); err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "VALIDATION_FAILED", nil, err)
	}

	data, err := h.service.UpdateMyProfile(c.Request().Context(), userID, body)
	if err != nil {
		return handleAppError(c, err)
	}

	return response.NewResponse(c, http.StatusOK, "OK", "PROFILE_UPDATED", data, nil)
}

func (h *Handler) UpdateMyPassword(c *echo.Context) error {
	userID, err := currentUserID(c)
	if err != nil {
		return unauthorizedResponse(c, err)
	}

	var body UpdatePasswordRequest
	if err := (&echo.DefaultBinder{}).Bind(c, &body); err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", "INVALID_REQUEST_BODY", nil, err)
	}
	if err := validator.NewValidator().ValidateStruct(body); err != nil {
		return response.NewResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "VALIDATION_FAILED", nil, err)
	}

	if err := h.service.UpdateMyPassword(c.Request().Context(), userID, body); err != nil {
		return handleAppError(c, err)
	}

	return response.NewResponse(c, http.StatusOK, "OK", "PASSWORD_UPDATED", map[string]any{}, nil)
}

func (h *Handler) GetLeaderboard(c *echo.Context) error {
	userID, err := currentUserID(c)
	if err != nil {
		return unauthorizedResponse(c, err)
	}

	data, err := h.service.GetLeaderboard(c.Request().Context(), userID)
	if err != nil {
		return handleAppError(c, err)
	}

	return response.NewResponse(c, http.StatusOK, "OK", "LEADERBOARD_RETRIEVED", data, nil)
}

func currentUserID(c *echo.Context) (string, error) {
	claims, ok := (*c).Get(authmw.ClaimsKey).(*utils.TokenPayload)
	if !ok {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "INVALID_TOKEN_CLAIMS")
	}
	return claims.UserID, nil
}

func unauthorizedResponse(c *echo.Context, err error) error {
	return response.NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error(), nil, nil)
}

func handleAppError(c *echo.Context, err error) error {
	switch err.Error() {
	case "ACCESS_DENIED":
		return response.NewResponse(c, http.StatusForbidden, "FORBIDDEN", "ACCESS_DENIED", nil, nil)
	case "USER_NOT_FOUND", "MENTEE_NOT_FOUND", "QUESTION_NOT_FOUND", "REQUEST_NOT_FOUND", "SHEET_NOT_FOUND":
		return response.NewResponse(c, http.StatusNotFound, "NOT_FOUND", err.Error(), nil, nil)
	case "EMAIL_ALREADY_EXISTS", "USERNAME_ALREADY_EXISTS":
		return response.NewResponse(c, http.StatusConflict, "CONFLICT", err.Error(), nil, nil)
	case "INVALID_USERNAME", "PASSWORD_MUST_CONTAIN_LETTER_AND_NUMBER", "SHEET_REQUIRED", "QUESTION_IDS_REQUIRED", "MENTEES_REQUIRED", "NO_FIELDS_TO_UPDATE", "BOOTCAMP_NOT_CONFIGURED", "INVALID_CURRENT_PASSWORD", "PASSWORD_LOGIN_NOT_AVAILABLE":
		return response.NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", err.Error(), nil, nil)
	default:
		return response.NewResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil, err)
	}
}

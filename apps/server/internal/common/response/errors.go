package response

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Success bool        `json:"success" example:"false"`
	Error   ErrorDetail `json:"error"`
}

// ErrorDetail contains error information
type ErrorDetail struct {
	Status  string `json:"status" example:"BAD_REQUEST"`
	Code    string `json:"code" example:"VALIDATION_ERROR"`
	Message string `json:"message" example:"Validation failed"`
}

// ValidationError returns a standardized validation error response
func ValidationError(c *echo.Context, code string, err error) error {
	return NewResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", code, nil, err)
}

// AuthorizationError returns a standardized authorization error response
func AuthorizationError(c *echo.Context, code, message string) error {
	if message == "" {
		message = "Access denied"
	}
	return NewResponse(c, http.StatusForbidden, "FORBIDDEN", message, nil, nil)
}

// AuthenticationError returns a standardized authentication error response
func AuthenticationError(c *echo.Context, code, message string) error {
	if message == "" {
		message = "Authentication required"
	}
	return NewResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", message, nil, nil)
}

// ConflictError returns a standardized conflict error response
func ConflictError(c *echo.Context, code, message string) error {
	if message == "" {
		message = "Resource conflict"
	}
	return NewResponse(c, http.StatusConflict, "CONFLICT", message, nil, nil)
}

// NotFoundError returns a standardized not found error response
func NotFoundError(c *echo.Context, code, message string) error {
	if message == "" {
		message = "Resource not found"
	}
	return NewResponse(c, http.StatusNotFound, "NOT_FOUND", message, nil, nil)
}

// BadRequestError returns a standardized bad request error response
func BadRequestError(c *echo.Context, code, message string) error {
	if message == "" {
		message = "Bad request"
	}
	return NewResponse(c, http.StatusBadRequest, "BAD_REQUEST", message, nil, nil)
}

// InternalServerError returns a standardized internal server error response
func InternalServerError(c *echo.Context, code string, err error) error {
	return NewResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", code, nil, err)
}

// UnprocessableEntityError returns a standardized unprocessable entity error response
func UnprocessableEntityError(c *echo.Context, code, message string) error {
	if message == "" {
		message = "Unprocessable entity"
	}
	return NewResponse(c, http.StatusUnprocessableEntity, "UNPROCESSABLE_ENTITY", message, nil, nil)
}

// TooManyRequestsError returns a standardized rate limit error response
func TooManyRequestsError(c *echo.Context, code, message string) error {
	if message == "" {
		message = "Too many requests"
	}
	return NewResponse(c, http.StatusTooManyRequests, "TOO_MANY_REQUESTS", message, nil, nil)
}

// HandleServiceError maps service layer errors to appropriate HTTP responses
func HandleServiceError(c *echo.Context, err error) error {
	if err == nil {
		return nil
	}

	errMsg := err.Error()

	// Authentication errors
	switch errMsg {
	case "INVALID_TOKEN", "TOKEN_EXPIRED", "INVALID_TOKEN_CLAIMS":
		return AuthenticationError(c, errMsg, "")
	}

	// Authorization errors
	switch errMsg {
	case "ACCESS_DENIED", "FORBIDDEN", "NOT_MEMBER_OF_ORGANIZATION", "NOT_ENROLLED_IN_BOOTCAMP",
		"ADMIN_REQUIRED", "MENTOR_REQUIRED", "ADMIN_OR_MENTOR_REQUIRED", "SUPER_ADMIN_REQUIRED",
		"ONLY_MENTORS_ADMINS_CAN_RESOLVE", "MENTEES_CANNOT_DELETE_DOUBTS",
		"ONLY_MENTEES_CAN_VOTE", "MENTEES_CANNOT_ACCESS_RESULTS", "MENTEES_CANNOT_ACCESS_VOTES":
		return AuthorizationError(c, errMsg, "")
	}

	// Not found errors
	switch errMsg {
	case "USER_NOT_FOUND", "ORGANIZATION_NOT_FOUND", "BOOTCAMP_NOT_FOUND", "PROBLEM_NOT_FOUND",
		"ASSIGNMENT_NOT_FOUND", "ASSIGNMENT_GROUP_NOT_FOUND", "DOUBT_NOT_FOUND", "POLL_NOT_FOUND",
		"ENROLLMENT_NOT_FOUND", "MEMBER_NOT_FOUND", "TAG_NOT_FOUND", "RESOURCE_NOT_FOUND",
		"ENTRY_NOT_FOUND", "ASSIGNMENT_PROBLEM_NOT_FOUND":
		return NotFoundError(c, errMsg, "")
	}

	// Conflict errors
	switch errMsg {
	case "DUPLICATE_ENTRY", "SLUG_ALREADY_EXISTS", "EMAIL_ALREADY_EXISTS",
		"ORGANIZATION_NOT_APPROVED", "BOOTCAMP_NOT_ACTIVE", "CROSS_ORG_VIOLATION",
		"CROSS_BOOTCAMP_VIOLATION", "PROBLEM_IN_USE", "TAG_IN_USE",
		"ASSIGNMENT_GROUP_HAS_ASSIGNMENTS", "DUPLICATE_ENROLLMENT", "DUPLICATE_ASSIGNMENT":
		return ConflictError(c, errMsg, "")
	}

	// Validation errors
	switch errMsg {
	case "INVALID_UUID", "INVALID_EMAIL", "INVALID_PASSWORD", "INVALID_ROLE",
		"INVALID_STATUS", "INVALID_DATE_RANGE", "INVALID_PROBLEM_ID", "INVALID_BOOTCAMP_ID",
		"INVALID_ORGANIZATION_ID", "INVALID_ENROLLMENT_ID", "INVALID_ASSIGNMENT_PROBLEM_ID",
		"INVALID_USER_ID", "INVALID_DOUBT_ID", "INVALID_POLL_ID", "INVALID_TAG_ID",
		"VALIDATION_FAILED", "NO_FIELDS_PROVIDED":
		return BadRequestError(c, errMsg, "")
	}

	// Default to internal server error
	return InternalServerError(c, "INTERNAL_ERROR", err)
}

// ErrorCode constants for common error scenarios
const (
	// Authentication errors
	ErrInvalidToken       = "INVALID_TOKEN"
	ErrTokenExpired       = "TOKEN_EXPIRED"
	ErrInvalidCredentials = "INVALID_CREDENTIALS"

	// Authorization errors
	ErrAccessDenied   = "ACCESS_DENIED"
	ErrForbidden      = "FORBIDDEN"
	ErrAdminRequired  = "ADMIN_REQUIRED"
	ErrMentorRequired = "MENTOR_REQUIRED"

	// Not found errors
	ErrNotFound             = "NOT_FOUND"
	ErrUserNotFound         = "USER_NOT_FOUND"
	ErrOrganizationNotFound = "ORGANIZATION_NOT_FOUND"
	ErrBootcampNotFound     = "BOOTCAMP_NOT_FOUND"
	ErrProblemNotFound      = "PROBLEM_NOT_FOUND"

	// Conflict errors
	ErrDuplicateEntry    = "DUPLICATE_ENTRY"
	ErrSlugExists        = "SLUG_ALREADY_EXISTS"
	ErrEmailExists       = "EMAIL_ALREADY_EXISTS"
	ErrCrossOrgViolation = "CROSS_ORG_VIOLATION"

	// Validation errors
	ErrValidationFailed = "VALIDATION_FAILED"
	ErrInvalidUUID      = "INVALID_UUID"
	ErrInvalidEmail     = "INVALID_EMAIL"
	ErrInvalidPassword  = "INVALID_PASSWORD"
	ErrNoFieldsProvided = "NO_FIELDS_PROVIDED"
)

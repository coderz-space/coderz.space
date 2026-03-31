package response

import "github.com/labstack/echo/v5"

type apiResponse struct {
	Data    any       `json:"data,omitempty"`
	Message string    `json:"message,omitempty"`
	Status  string    `json:"status,omitempty"`
	Error   *apiError `json:"error,omitempty"`
	Success bool      `json:"success,omitempty"`
}

type apiError struct {
	Message string `json:"message,omitempty"`
	Code    int    `json:"code,omitempty"`
}

func NewResponse(c *echo.Context, statusCode int, status, message string, data any, err error) error {
	res := &apiResponse{
		Message: message,
		Data:    data,
		Success: true,
		Status:  status,
	}

	if statusCode < 200 || statusCode >= 300 || err != nil {
		res.Success = false
		res.Error = &apiError{
			Code:    statusCode,
			Message: status,
		}
	}

	return c.JSON(statusCode, res)
}

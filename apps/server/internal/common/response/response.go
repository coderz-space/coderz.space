package response

import "github.com/labstack/echo/v5"

type apiResponse struct {
	Success bool      `json:"success,omitempty"`
	Message string    `json:"message,omitempty"`
	Data    any       `json:"data,omitempty"`
	Error   *apiError `json:"error,omitempty"`
	Status  string    `json:"status,omitempty"`
}

type apiError struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func NewResponse(c *echo.Context, statusCode int, status string, message string, data any, err error) error {
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

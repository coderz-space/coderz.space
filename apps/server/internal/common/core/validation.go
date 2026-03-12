package core

import (
	"github.com/DSAwithGautam/Coderz.space/internal/common/validator"
	"github.com/labstack/echo/v5"
)

// WithBody decorator for JSON body parsing
// Example usage:
//
//	type CreateUserRequest struct {
//	    Name  string `json:"name" validate:"required"`
//	    Email string `json:"email" validate:"required,email"`
//	}
//	e.POST("/users", WithBody(func(c echo.Context, body CreateUserRequest) error {
//	    return c.JSON(201, body)
//	}))
func WithBody[T any](f func(*echo.Context, T) error) echo.HandlerFunc {
	return func(c *echo.Context) error {
		var body T

		// bind the request body to the generic type
		if err := (&echo.DefaultBinder{}).Bind(c, &body); err != nil {
			return err
		}

		// validate the request body
		if err := validator.NewValidator().ValidateStruct(body); err != nil {
			return err
		}
		return f(c, body)
	}
}


// WithParams decorator for URL path parameters
// Example usage:
//
//	type UserParams struct {
//	    ID string `param:"id"`
//	}
//	e.GET("/users/:id", WithParams(func(c echo.Context, params UserParams) error {
//	    return c.JSON(200, map[string]string{"id": params.ID})
//	}))
func WithParams[T any](f func(*echo.Context, T) error) echo.HandlerFunc {
	return func(c *echo.Context) error {
		var params T

		// bind the request params to the generic type
		if err := (&echo.DefaultBinder{}).Bind(c, &params); err != nil {
			return err
		}

		// validate the request params
		if err := validator.NewValidator().ValidateStruct(params); err != nil {
			return err
		}
		return f(c, params)
	}
}


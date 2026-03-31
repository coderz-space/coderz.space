package core

import (
	"net/http"

	"github.com/coderz-space/coderz.space/internal/common/response"
	"github.com/coderz-space/coderz.space/internal/common/validator"
	"github.com/labstack/echo/v5"
)

// WithBody decorator for JSON body parsing
// Example usage:
//
//	type CreateUserRequest struct {
//	    Name  string `json:"name" validate:"required"`
//	    Email string `json:"email" validate:"required,email"`
//	}
//	e.POST("/users", WithBody(func(c *echo.Context, body CreateUserRequest) error {
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
//	e.GET("/users/:id", WithParams(func(c *echo.Context, params UserParams) error {
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

// WithQuery decorator for query parameters
// Example usage:
//
//	type UserQuery struct {
//	    Page  int `query:"page"`
//	    Limit int `query:"limit"`
//	}
//	e.GET("/users", WithQuery(func(c *echo.Context, query UserQuery) error {
//	    return c.JSON(200, query)
//	}))
func WithQuery[Q any](handler func(*echo.Context, Q) error) echo.HandlerFunc {
	return func(c *echo.Context) error {
		var query Q
		if err := (&echo.DefaultBinder{}).Bind(c, &query); err != nil {
			return response.NewResponse(c, http.StatusBadRequest, "INVALID_QUERY_PARAMETERS", "Failed to bind query parameters", nil, err)
		}

		// Validate the bound query parameters
		if err := validator.NewValidator().ValidateStruct(query); err != nil {
			return response.NewResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "Query validation failed", nil, err)
		}

		return handler(c, query)
	}
}

// WithBodyAndParams combines body and URL parameters validation
// Example usage:
//
//	type UpdateUserRequest struct {
//	    Name  string `json:"name"`
//	    Email string `json:"email"`
//	}
//	type UserParams struct {
//	    ID string `param:"id"`
//	}
//	e.PUT("/users/:id", WithBodyAndParams(func(c *echo.Context, body UpdateUserRequest, params UserParams) error {
//	    return c.JSON(200, map[string]any{"id": params.ID, "user": body})
//	}))
func WithBodyAndParams[B any, P any](handler func(*echo.Context, B, P) error) echo.HandlerFunc {
	return func(c *echo.Context) error {
		// Bind and validate body
		var body B
		if err := (&echo.DefaultBinder{}).Bind(c, &body); err != nil {
			return response.NewResponse(c, http.StatusBadRequest, "INVALID_REQUEST_BODY", "Failed to bind request body", nil, err)
		}

		if err := validator.NewValidator().ValidateStruct(body); err != nil {
			return response.NewResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "Body validation failed", nil, err)
		}

		// Bind and validate path parameters
		var params P
		if err := (&echo.DefaultBinder{}).Bind(c, &params); err != nil {
			return response.NewResponse(c, http.StatusBadRequest, "INVALID_URL_PARAMETERS", "Failed to bind path parameters", nil, err)
		}

		if err := validator.NewValidator().ValidateStruct(params); err != nil {
			return response.NewResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "Parameters validation failed", nil, err)
		}

		return handler(c, body, params)
	}
}

// WithParamsAndQuery combines URL parameters and query parameters validation
// Example usage:
//
//	type UserParams struct {
//	    ID string `param:"id"`
//	}
//	type PostQuery struct {
//	    Page  int `query:"page"`
//	    Limit int `query:"limit"`
//	}
//	e.GET("/users/:id/posts", WithParamsAndQuery(func(c *echo.Context, params UserParams, query PostQuery) error {
//	    return c.JSON(200, map[string]any{"userId": params.ID, "page": query.Page})
//	}))
func WithParamsAndQuery[P any, Q any](handler func(*echo.Context, P, Q) error) echo.HandlerFunc {
	return func(c *echo.Context) error {
		// Bind and validate path parameters
		var params P
		if err := (&echo.DefaultBinder{}).Bind(c, &params); err != nil {
			return response.NewResponse(c, http.StatusBadRequest, "INVALID_URL_PARAMETERS", "Failed to bind path parameters", nil, err)
		}

		if err := validator.NewValidator().ValidateStruct(params); err != nil {
			return response.NewResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "Parameters validation failed", nil, err)
		}

		// Bind and validate query parameters
		var query Q
		if err := (&echo.DefaultBinder{}).Bind(c, &query); err != nil {
			return response.NewResponse(c, http.StatusBadRequest, "INVALID_QUERY_PARAMETERS", "Failed to bind query parameters", nil, err)
		}

		if err := validator.NewValidator().ValidateStruct(query); err != nil {
			return response.NewResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "Query validation failed", nil, err)
		}

		return handler(c, params, query)
	}
}

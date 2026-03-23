package controller

import "github.com/labstack/echo/v4"

type Response struct {
	Success bool      `json:"success"`
	Data    any       `json:"data,omitempty"`
	Error   *APIError `json:"error,omitempty"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func Success(ctx echo.Context, status int, data any) error {
	return ctx.JSON(status, Response{
		Success: true,
		Data:    data,
	})
}

func Error(ctx echo.Context, status int, code, message string) error {
	return ctx.JSON(status, Response{
		Success: false,
		Error: &APIError{
			Code:    code,
			Message: message,
		},
	})
}

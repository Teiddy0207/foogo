package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"fooder-backend/core/controller"
	corecontroller "fooder-backend/core/controller"
	"fooder-backend/internal/dto"
	"fooder-backend/internal/service"
)

type AuthController struct {
	service *service.AuthService
}

func NewAuthController(service *service.AuthService) *AuthController {
	return &AuthController{service: service}
}

func (c *AuthController) HandleHealth(ctx echo.Context) error {
	return corecontroller.Success(ctx, http.StatusOK, map[string]any{"status": "ok"})
}

func (c *AuthController) HandleLogin(ctx echo.Context) error {
	var req dto.LoginRequest
	if err := ctx.Bind(&req); err != nil {
		return corecontroller.Error(ctx, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body")
	}

	result, err := c.service.Login(req)
	if err != nil {
		return corecontroller.Error(ctx, http.StatusUnauthorized, "AUTH_LOGIN_FAILED", err.Error())
	}

	return controller.Controller().SuccessResponse(ctx, result, "login successful")
}

func (c *AuthController) HandleLogout(ctx echo.Context) error {
	return corecontroller.Success(ctx, http.StatusOK, map[string]string{
		"message": "logout success on stateless auth",
	})
}

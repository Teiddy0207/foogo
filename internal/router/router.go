package router

import (
	"github.com/labstack/echo/v4"

	"fooder-backend/internal/controller"
)

type Controllers struct {
	Auth *controller.AuthController
}

func RegisterRoutes(e *echo.Echo, c Controllers) {
	e.GET("/health", c.Auth.HandleHealth)

	auth := e.Group("/auth")
	auth.POST("/login", c.Auth.HandleLogin)
	auth.POST("/logout", c.Auth.HandleLogout)
}

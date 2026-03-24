package main

import (
	"errors"
	"net/http"

	"fooder-backend/core/config"
	"fooder-backend/core/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"fooder-backend/internal/controller"
	"fooder-backend/internal/repository"
	"fooder-backend/internal/router"
	"fooder-backend/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	if err := logger.Init(logger.LogConfig{
		Level:         logger.LogLevel(cfg.Log.Level),
		JSONFormat:    cfg.Log.JSON,
		DailyRotation: cfg.Log.DailyRotation,
		EnableFile:    cfg.Log.EnableFile,
	}); err != nil {
		panic(err)
	}
	defer func() { _ = logger.Close() }()

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(logger.RequestLoggerMiddleware())

	authRepo := repository.NewInMemoryAuthRepository()
	authService := service.NewAuthService(authRepo, cfg.Auth.DevTokenPrefix)
	authController := controller.NewAuthController(authService)

	detectService, err := service.NewDetectService(cfg.App.DetectGRPCAddr)
	if err != nil {
		panic(err)
	}
	defer func() { _ = detectService.Close() }()
	detectController := controller.NewDetectController(detectService)

	router.RegisterRoutes(e, router.Controllers{
		Auth:   authController,
		Detect: detectController,
	})

	addr := ":" + cfg.App.Port
	logger.Info("server starting", "addr", addr, "app_env", cfg.App.Env, "log_level", cfg.Log.Level)
	if err := e.Start(addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("start server failed", "error", err.Error())
	}
}

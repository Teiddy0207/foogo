package controller

import (
	corecontroller "fooder-backend/core/controller"
	"fooder-backend/core/errors"
	"fooder-backend/internal/dto"
	"fooder-backend/internal/service"

	"github.com/labstack/echo/v4"
)

type DetectController struct {
	service *service.DetectService
}

func NewDetectController(service *service.DetectService) *DetectController {
	return &DetectController{service: service}
}

func (c *DetectController) HandleAnalyzeFood(ctx echo.Context) error {
	var req dto.AnalyzeFoodRequest
	if err := ctx.Bind(&req); err != nil {
		return corecontroller.Controller().BadRequest(errors.ErrInvalidRequestData, "invalid request body")
	}

	resp, appErr := c.service.AnalyzeFood(ctx.Request().Context(), req)
	if appErr != nil {
		switch appErr.Code {
		case errors.ErrInvalidInput:
			return corecontroller.Controller().BadRequest(appErr.Code, appErr.Message)
		default:
			return corecontroller.Controller().InternalServerError(appErr.Code, appErr.Message)
		}
	}

	return corecontroller.Controller().SuccessResponse(ctx, resp, "food analyzed successfully")
}

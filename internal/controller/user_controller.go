package controller

import (
	corecontroller "fooder-backend/core/controller"
	"fooder-backend/core/errors"
	"fooder-backend/core/params"
	"fooder-backend/internal/service"

	"github.com/labstack/echo/v4"
)

type UserController struct {
	service *service.UserService
}

func NewUserController(service *service.UserService) *UserController {
	return &UserController{service: service}
}

func (controller *UserController) HandleGetUser(c echo.Context) error {
	ctx := c.Request().Context()
	q := params.NewQueryParams(c)

	users, err := controller.service.GetUsers(ctx, q)
	if err != nil {
		return corecontroller.Controller().InternalServerError(errors.ErrInternalServer, "failed to fetch users", err.Error())
	}

	return corecontroller.Controller().SuccessResponse(c, users, "users fetched successfully")
}

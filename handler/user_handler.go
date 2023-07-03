package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jayzedx/hotel-reservation/errs"
	"github.com/jayzedx/hotel-reservation/repo"
	"github.com/jayzedx/hotel-reservation/resp"
	"github.com/jayzedx/hotel-reservation/service"
)

type userHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *userHandler {
	return &userHandler{
		userService: userService,
	}
}

func (h *userHandler) HandleGetUser(ctx *fiber.Ctx) error {
	var (
		id = ctx.Params("id")
	)

	data, err := h.userService.GetUserById(id)
	if err != nil {
		appErr, ok := err.(errs.AppError)
		if ok {
			return appErr
		} else {
			return err
		}
	}
	return ctx.Status(http.StatusOK).JSON(resp.Response{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Operation completed successfully",
		Data:    data,
	})
}

/*
func (h *userHandler) HandleGetUsers(ctx *fiber.Ctx) error {
	data, err := h.userService.GetUsers()
	if err != nil {
		appErr, ok := err.(errs.AppError)
		if ok {
			return appErr
		} else {
			return errs.AppError{
				Code:    http.StatusBadRequest,
				Message: "Unexpected error",
			}
		}
	}
	return ctx.Status(http.StatusOK).JSON(resp.Response{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Operation completed successfully",
		Data:    data,
	})
}
*/

func (h *userHandler) HandlePostUser(ctx *fiber.Ctx) error {
	var params service.CreateUserParams
	if err := ctx.BodyParser(&params); err != nil {
		return errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "Invalid data provided",
		}
	}
	data, err := h.userService.CreateUser(params)
	if err != nil {
		appErr, ok := err.(errs.AppError)
		if ok {
			return appErr
		} else {
			return err
		}
	}

	return ctx.Status(http.StatusOK).JSON(resp.Response{
		Code:    http.StatusCreated,
		Status:  "success",
		Message: "Operation completed successfully",
		Data:    data,
	})
}

func (h *userHandler) HandlePutUser(ctx *fiber.Ctx) error {
	var (
		id     = ctx.Params("id")
		params service.UpdateUserParams
	)

	if err := ctx.BodyParser(&params); err != nil {
		return errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "Invalid data provided",
		}
	}
	if err := h.userService.UpdateUser(id, params); err != nil {
		appErr, ok := err.(errs.AppError)
		if ok {
			return appErr
		} else {
			return err
		}
	}

	return ctx.Status(http.StatusOK).JSON(resp.Response{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Operation completed successfully",
		Data:    nil,
	})
}

func (h *userHandler) HandleDeleteUser(ctx *fiber.Ctx) error {
	var (
		id = ctx.Params("id")
	)
	if err := h.userService.DeleteUser(id); err != nil {
		appErr, ok := err.(errs.AppError)
		if ok {
			return appErr
		} else {
			return err
		}
	}
	return ctx.Status(http.StatusOK).JSON(resp.Response{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Operation completed successfully",
		Data:    nil,
	})
}

func (h *userHandler) HandleGetUsersByParams(ctx *fiber.Ctx) error {
	var params repo.User
	if err := ctx.QueryParser(&params); err != nil {
		return errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "Invalid data provided",
		}
	}
	data, err := h.userService.GetUsersByParams(params)
	if err != nil {
		appErr, ok := err.(errs.AppError)
		if ok {
			return appErr
		} else {
			return err
		}
	}
	return ctx.Status(http.StatusOK).JSON(resp.Response{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Operation completed successfully",
		Data:    data,
	})
}

package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jayzedx/hotel-reservation/errs"
	"github.com/jayzedx/hotel-reservation/resp"
	"github.com/jayzedx/hotel-reservation/service"
)

type authHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *authHandler {
	return &authHandler{
		authService: authService,
	}
}

func (h *authHandler) HandlePostAuthen(ctx *fiber.Ctx) error {
	var params service.CreateAuthParams
	if err := ctx.BodyParser(&params); err != nil {
		return errs.ErrBadRequest()
	}
	data, err := h.authService.Authenticate(params)
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

package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jayzedx/hotel-reservation/errs"
	"github.com/jayzedx/hotel-reservation/resp"
)

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	appErr, ok := err.(errs.AppError)
	if ok {
		return ctx.Status(appErr.Code).JSON(resp.Response{
			Code:    appErr.Code,
			Status:  "error",
			Message: err.Error(),
			Data:    nil,
			Errors:  appErr.Errors,
		})
	}
	return ctx.Status(fiber.StatusInternalServerError).JSON(resp.Response{
		Code:    0,
		Status:  "error",
		Message: "Unexpected error : " + err.Error(),
		Data:    nil,
		Errors:  appErr.Errors,
	})
}

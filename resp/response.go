package resp

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jayzedx/hotel-reservation/errs"
)

type Response struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Errors  interface{} `json:"errors"`
}

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	appErr, ok := err.(errs.AppError)
	if ok {
		return ctx.Status(appErr.Code).JSON(Response{
			Code:    appErr.Code,
			Status:  "error",
			Message: err.Error(),
			Data:    nil,
			Errors:  appErr.Errors,
		})
	}
	return ctx.Status(fiber.StatusInternalServerError).JSON(Response{
		Code:    0,
		Status:  "error",
		Message: "Unexpected error : " + err.Error(),
		Data:    nil,
		Errors:  appErr.Errors,
	})
}

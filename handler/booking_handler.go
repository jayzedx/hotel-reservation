package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jayzedx/hotel-reservation/errs"
	"github.com/jayzedx/hotel-reservation/resp"
	"github.com/jayzedx/hotel-reservation/service"
)

type bookingHandler struct {
	bookingService service.BookingService
}

func NewBookingHandler(bookingService service.BookingService) *bookingHandler {
	return &bookingHandler{
		bookingService: bookingService,
	}
}

func (h *bookingHandler) HandlePostBooking(ctx *fiber.Ctx) error {
	var (
		roomId = ctx.Params("id")
		params service.CreateBookingParams
	)
	if err := ctx.BodyParser(&params); err != nil {
		return errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "Invalid data provided",
		}
	}

	booking, err := h.bookingService.CreateBooking(ctx, roomId, params)
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
		Data:    booking,
	})
}

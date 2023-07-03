package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jayzedx/hotel-reservation/errs"
	"github.com/jayzedx/hotel-reservation/repo"
	"github.com/jayzedx/hotel-reservation/resp"
	"github.com/jayzedx/hotel-reservation/service"
)

type hotelHandler struct {
	hotelService service.HotelService
}

func NewHotelHandler(hotelService service.HotelService) *hotelHandler {
	return &hotelHandler{
		hotelService: hotelService,
	}
}

func (h *hotelHandler) HandleGetHotels(ctx *fiber.Ctx) error {

	var params repo.Hotel
	if err := ctx.QueryParser(&params); err != nil {
		return errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "Invalid data provided",
		}
	}

	hotels, err := h.hotelService.GetHotels(params)
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
		Data:    hotels,
	})
}

func (h *hotelHandler) HandleGetHotel(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	hotel, err := h.hotelService.GetHotelRooms(id)
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
		Data:    hotel,
	})
}

func (h *hotelHandler) HandlePostHotel(ctx *fiber.Ctx) error {
	var params service.CreateHotelParams
	if err := ctx.BodyParser(&params); err != nil {
		return errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "Invalid data provided",
		}
	}
	hotel, err := h.hotelService.CreateHotel(params)
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
		Data:    hotel,
	})
}

func (h *hotelHandler) HandlePutHotel(ctx *fiber.Ctx) error {
	var (
		id     = ctx.Params("id")
		params service.UpdateHotelParams
	)
	if err := ctx.BodyParser(&params); err != nil {
		return errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "Invalid data provided",
		}
	}
	if err := h.hotelService.UpdateHotel(id, params); err != nil {
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

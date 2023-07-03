package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jayzedx/hotel-reservation/errs"
	"github.com/jayzedx/hotel-reservation/resp"
	"github.com/jayzedx/hotel-reservation/service"
)

type roomHandler struct {
	roomService service.RoomService
}

func NewRoomHandler(roomService service.RoomService) *roomHandler {
	return &roomHandler{
		roomService: roomService,
	}
}

func (h *roomHandler) HandlePostRoom(ctx *fiber.Ctx) error {
	var params service.CreateRoomParams
	if err := ctx.BodyParser(&params); err != nil {
		return errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "Invalid data provided",
		}
	}

	room, err := h.roomService.CreateRoom(params)
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
		Data:    room,
	})
}

func (h *roomHandler) HandlePutRoom(ctx *fiber.Ctx) error {
	var (
		id     = ctx.Params("id")
		params service.UpdateRoomParams
	)
	if err := ctx.BodyParser(&params); err != nil {
		return errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "Invalid data provided",
		}
	}
	if err := h.roomService.UpdateRoom(id, params); err != nil {
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
func (h *roomHandler) HandleDeleteRoom(ctx *fiber.Ctx) error {
	var (
		id = ctx.Params("id")
	)
	if err := h.roomService.DeleteRoom(id); err != nil {
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

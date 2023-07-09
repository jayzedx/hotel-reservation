package service

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jayzedx/hotel-reservation/errs"
	"github.com/jayzedx/hotel-reservation/logs"
	"github.com/jayzedx/hotel-reservation/repo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type bookingService struct {
	roomRepository    repo.RoomRepository
	bookingRepository repo.BookingRepository
}

func NewBookingService(roomRepository repo.RoomRepository, bookingRepository repo.BookingRepository) *bookingService {
	return &bookingService{
		roomRepository:    roomRepository,
		bookingRepository: bookingRepository,
	}
}

func (s *bookingService) CreateBooking(ctx *fiber.Ctx, roomIdstr string, params CreateBookingParams) (*BookingResponse, error) {
	roomId, err := primitive.ObjectIDFromHex(roomIdstr)
	if err != nil {
		return nil, errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "Invalid room id provided",
		}
	}

	user, ok := ctx.Context().Value("user").(*repo.User)
	if !ok {
		return nil, errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "Invalid token provide",
		}
	}

	// validation
	if errors := params.Validate(); len(errors) > 0 {
		return nil, errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "Validation errors",
			Errors:  errors,
		}
	}

	// check room id is correct
	_, err = s.roomRepository.GetRoomById(roomId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errs.AppError{
				Code:    http.StatusBadRequest,
				Message: "Invalid room id provided",
			}
		} else {
			logs.Error(err)
			return nil, errs.AppError{
				Code:    http.StatusBadRequest,
				Message: "Unexpected error",
			}
		}
	}

	// check room id is available
	ok, err = s.isBookingAvailable(ctx, roomId, params.FromDate, params.TilDate)
	if err != nil {
		logs.Error(err)
		return nil, errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "Unexpected error",
		}
	}
	if !ok {
		return nil, errs.AppError{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("Room id : %s , Already booked", roomId.Hex()),
		}
	}

	// set room id and userId
	params.RoomId = roomId
	params.UserId = user.Id

	booking := CreateBookingFromParams(&params)
	if err := s.bookingRepository.CreateBooking(booking); err != nil {
		return nil, err
	}

	return MapBookingResponse(booking), nil
}

func (s *bookingService) isBookingAvailable(ctx *fiber.Ctx, roomId primitive.ObjectID, fromDate time.Time, tilDate time.Time) (bool, error) {
	// check booking is avaliable
	filter := bson.M{
		"room_id":   roomId,
		"from_date": bson.M{"$gte": fromDate},
		"til_date":  bson.M{"$lte": tilDate},
	}
	bookings, err := s.bookingRepository.GetBookings(filter)
	if err != nil {
		return false, err
	}
	return len(bookings) == 0, nil
}

func (s *bookingService) GetBookings(ctx *fiber.Ctx) ([]*BookingResponse, error) {
	bookings, err := s.bookingRepository.GetBookings(bson.M{})
	if err != nil {
		logs.Error(err)
		return nil, errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "Unexpected error",
		}
	}

	data := []*BookingResponse{}
	for _, booking := range bookings {
		data = append(data, MapBookingResponse(booking))
	}
	return data, nil
}

func (s *bookingService) GetBooking(ctx *fiber.Ctx, id string) (*BookingResponse, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "Invalid id provided",
		}
	}

	booking, err := s.bookingRepository.GetBookingById(oid)
	if err != nil {
		return nil, errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "Your booking id isn't exist",
		}
	}

	user, ok := ctx.Context().UserValue("user").(*repo.User)
	if !ok {
		return nil, errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "Unexpected error",
		}
	}
	if booking.UserId != user.Id {
		return nil, errs.AppError{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
		}
	}

	return MapBookingResponse(booking), nil
}

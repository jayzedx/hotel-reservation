package service

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jayzedx/hotel-reservation/repo"
	"github.com/jayzedx/hotel-reservation/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookingService interface {
	CreateBooking(ctx *fiber.Ctx, roomId string, params CreateBookingParams) (*BookingResponse, error)
	GetBookings(ctx *fiber.Ctx) ([]*BookingResponse, error)
	GetBooking(ctx *fiber.Ctx, id string) (*BookingResponse, error)
	CancelBooking(ctx *fiber.Ctx, id string) error
}

type BookingResponse struct {
	Id           string `json:"booking_id,omitempty" mapstructure:"booking_id"`
	RoomId       string `json:"room_id" mapstructure:"room_id"`
	UserId       string `json:"user_id" mapstructure:"user_id"`
	PersonNumber int    `json:"person_number" mapstructure:"person_number"`
	FromDate     string `json:"from_date,omitempty" mapstructure:"from_date"`
	TilDate      string `json:"til_date,omitempty" mapstructure:"til_date"`
	Canceled     bool   `json:"canceled"`
	CancelDate   string `json:"cancel_date,omitempty" mapstructure:"cancel_date"`
}

// {
// 	"person_number": 3,
// 	"from_date": "2023-07-08T10:30:00.000Z",
// 	"til_date": "2023-07-15T10:30:00.000Z"
// }

type CreateBookingParams struct {
	RoomId       primitive.ObjectID `json:"-"`
	UserId       primitive.ObjectID `json:"-"`
	PersonNumber int                `json:"person_number"`
	FromDate     time.Time          `json:"from_date"`
	TilDate      time.Time          `json:"til_date"`
}

func CreateBookingFromParams(params *CreateBookingParams) *repo.Booking {
	return &repo.Booking{
		RoomId:       params.RoomId,
		UserId:       params.UserId,
		PersonNumber: params.PersonNumber,
		FromDate:     params.FromDate,
		TilDate:      params.TilDate,
	}
}

func MapBookingResponse(booking *repo.Booking) *BookingResponse {
	return &BookingResponse{
		Id:           booking.Id.Hex(),
		RoomId:       booking.RoomId.Hex(),
		UserId:       booking.UserId.Hex(),
		PersonNumber: booking.PersonNumber,
		FromDate:     util.ConvertTimeToString(booking.FromDate),
		TilDate:      util.ConvertTimeToString(booking.TilDate),
		Canceled:     booking.Canceled,
		CancelDate:   util.ConvertTimeToString(booking.CancelDate),
	}
}

func (params *CreateBookingParams) Validate() map[string]string {
	errors := map[string]string{}
	now := time.Now()
	if now.After(params.FromDate) {
		errors["from_date"] = "from_date is invalid"
	}
	if now.After(params.TilDate) {
		errors["til_date"] = "til_date is invalid"
	}

	if params.FromDate.After(params.TilDate) {
		errors["til_date"] = "from_date is invalid"
	}

	return errors
}

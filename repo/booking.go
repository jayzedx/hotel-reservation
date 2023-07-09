package repo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Booking struct {
	Id           primitive.ObjectID `bson:"_id,omitempty"`
	RoomId       primitive.ObjectID `bson:"room_id"`
	UserId       primitive.ObjectID `bson:"user_id"`
	PersonNumber int                `bson:"person_number"`
	FromDate     time.Time          `bson:"from_date"`
	TilDate      time.Time          `bson:"til_date"`
	Canceled     bool               `bson:"canceled"`
	CancelDate   time.Time          `bson:"cancel_date"`
}

type BookingRepository interface {
	CreateBooking(*Booking) error
	GetBookings(filter bson.M) ([]*Booking, error)
	GetBookingById(id primitive.ObjectID) (*Booking, error)
	UpdateBooking(filter bson.M, update interface{}) (int64, error)
}

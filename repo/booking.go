package repo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Booking struct {
	Id           primitive.ObjectID `bson:"_id,omitempty"`
	RoomId       primitive.ObjectID `bson:"room_id,omitempty"`
	UserId       primitive.ObjectID `bson:"user_id,omitempty"`
	PersonNumber int                `bson:"person_number,omitempty"`
	FromDate     time.Time          `bson:"from_date,omitempty"`
	TilDate      time.Time          `bson:"til_date,omitempty"`
}

type BookingRepository interface {
	CreateBooking(*Booking) error
	GetBookings(filter bson.M) ([]*Booking, error)
}

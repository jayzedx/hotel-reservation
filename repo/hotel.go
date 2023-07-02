package repo

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Hotel struct {
	Id       primitive.ObjectID   `bson:"_id,omitempty"`
	Name     string               `bson:"Name"`
	Location string               `bson:"Location"`
	Rating   int                  `bson:"rating" query:"rating"`
	Rooms    []primitive.ObjectID `bson:"rooms"`
}

type HotelRepository interface {
	GetHotels(filter bson.M) ([]*Hotel, error)
	GetHotelById(primitive.ObjectID) (*Hotel, error)
	CreateHotel(*Hotel) error
	UpdateRooms(filter bson.M, params bson.M) error
}
package repo

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Hotel struct {
	Id       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `bson:"name"`
	Location string             `bson:"location"`
	Rating   int                `bson:"rating"`
}

type HotelRepository interface {
	GetHotels(filter bson.M) ([]*Hotel, error)
	GetHotelById(primitive.ObjectID) (*Hotel, error)
	CreateHotel(*Hotel) error
	UpdateHotel(filter bson.M, update interface{}) (int64, error)
}

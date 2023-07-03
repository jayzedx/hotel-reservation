package repo

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Hotel struct {
	Id       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `bson:"name,omitempty"`
	Location string             `bson:"location,omitempty"`
	Rating   int                `bson:"rating,omitempty" query:"rating"`
}

type HotelRepository interface {
	GetHotels(filter bson.M) ([]*Hotel, error)
	GetHotelById(primitive.ObjectID) (*Hotel, error)
	CreateHotel(*Hotel) error
	UpdateHotel(filter bson.M, update interface{}) (int64, error)
}

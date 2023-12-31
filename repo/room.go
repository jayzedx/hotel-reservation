package repo

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoomType int

const (
	_ RoomType = iota
	SingleRoomType
	DoubleRoomType
	SeasideRoomType
	DeluxeRoomType
)

type Room struct {
	Id          primitive.ObjectID `bson:"_id,omitempty"` //omitempty - don't show json when id is empty
	Type        RoomType           `bson:"type"`
	Seaside     bool               `bson:"seaside"`
	Size        string             `bson:"size"`
	Price       float64            `bson:"price"`
	Selected    bool               `bson:"selected"`
	HotelId     primitive.ObjectID `bson:"hotel_id"`
	IsAvailable bool               `bson:"is_available"`
}

type RoomRepository interface {
	GetRoomsByPipeline(pipeline []bson.M) ([]*Room, error)
	GetRooms(filter bson.M) ([]*Room, error)
	GetRoomById(id primitive.ObjectID) (*Room, error)
	CreateRoom(*Room) error
	UpdateRoom(filter bson.M, update interface{}) (int64, error)
	GetRoomIds(hotelId primitive.ObjectID) ([]primitive.ObjectID, error)
	DeleteRoom(id primitive.ObjectID) error
}

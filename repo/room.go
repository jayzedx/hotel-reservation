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
	Id       primitive.ObjectID `bson:"_id,omitempty"` //omitempty - don't show json when id is empty
	Type     RoomType           `bson:"type,omitempty"`
	Seaside  bool               `bson:"seaside,omitempty"`
	Size     string             `bson:"size,omitempty"`
	Price    float64            `bson:"price,omitempty"`
	Selected bool               `bson:"selected,omitempty"`
	HotelId  primitive.ObjectID `bson:"hotel_id,omitempty"`
}

type RoomRepository interface {
	GetRoomsByPipeline(pipeline []bson.M) ([]*Room, error)
	GetRooms(filter bson.M) ([]*Room, error)
	CreateRoom(*Room) error
	UpdateRoom(filter bson.M, update interface{}) (int64, error)
	GetRoomIds(hotelId primitive.ObjectID) ([]primitive.ObjectID, error)
	DeleteRoom(id primitive.ObjectID) error
}

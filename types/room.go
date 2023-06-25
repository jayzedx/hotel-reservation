package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type RoomType int

const (
	_ RoomType = iota
	SingleRoomType
	DoubleRoomType
	SeasideRoomType
	DeluxeRoomType
)

type Room struct {
	Id        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"` //omitempty - don't show json when id is empty
	Type      RoomType           `bson:"type" json:"type"`
	BasePrice float64            `bson:"basePrice" json:"basePrice"`
	Price     float64            `bson:"price" json:"price"`
	HotelId   primitive.ObjectID `bson:"hotelId" json:"hotelId"`
}

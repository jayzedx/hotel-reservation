package service

import (
	"github.com/jayzedx/hotel-reservation/repo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoomService interface {
	CreateRoom(params CreateRoomParams) (*RoomResponse, error)
	UpdateRoom(id string, params UpdateRoomParams) (*RoomResponse, error)
	DeleteRoom(id string) error
}

type RoomResponse struct {
	Id      primitive.ObjectID `json:"id,omitempty"`
	Type    repo.RoomType      `json:"type"`
	Seaside bool               `json:"seaside"`
	Size    string             `json:"size"`
	Price   float64            `json:"price"`
	HotelId primitive.ObjectID `json:"hotel_id"`
}

type CreateRoomParams struct {
	Type    repo.RoomType      `json:"type"`
	Seaside bool               `json:"seaside"`
	Size    string             `json:"size"`
	Price   float64            `json:"price"`
	HotelId primitive.ObjectID `json:"hotel_id"`
}

type UpdateRoomParams struct {
	Type    repo.RoomType `json:"type"`
	Seaside bool          `json:"seaside"`
	Size    string        `json:"size"`
	Price   float64       `json:"price"`
}

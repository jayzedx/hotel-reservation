package service

import (
	"github.com/jayzedx/hotel-reservation/repo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoomService interface {
	GetRooms() ([]*RoomResponse, error)
	CreateRoom(params CreateRoomParams) (*RoomResponse, error)
	UpdateRoom(id string, params UpdateRoomParams) error
	DeleteRoom(id string) error
}

type RoomResponse struct {
	Id          string        `json:"room_id" mapstructure:"room_id"`
	Type        repo.RoomType `json:"type"`
	Seaside     bool          `json:"seaside"`
	Size        string        `json:"size"`
	Price       float64       `json:"price"`
	HotelId     string        `json:"hotel_id" mapstructure:"hotel_id"`
	Selected    bool          `json:"selected"`
	IsAvailable bool          `json:"-"`
}

type CreateRoomParams struct {
	Type        repo.RoomType      `json:"type"`
	Seaside     bool               `json:"seaside"`
	Size        string             `json:"size"`
	Price       float64            `json:"price"`
	HotelId     primitive.ObjectID `json:"hotel_id"`
	Selected    bool               `json:"selected"`
	IsAvailable bool               `json:"-"`
}

type UpdateRoomParams struct {
	Type        repo.RoomType `json:"type"`
	Seaside     bool          `json:"seaside"`
	Size        string        `json:"size"`
	Price       float64       `json:"price"`
	Selected    bool          `json:"selected"`
	IsAvailable bool          `json:"is_available"`
}

func (params *CreateRoomParams) Validate() map[string]string {
	errors := map[string]string{}
	return errors
}

func (params *UpdateRoomParams) Validate() map[string]string {
	errors := map[string]string{}
	return errors
}

func CreateRoomFromParams(params *CreateRoomParams) *repo.Room {
	return &repo.Room{
		Type:        params.Type,
		Seaside:     params.Seaside,
		Size:        params.Size,
		Price:       params.Price,
		HotelId:     params.HotelId,
		Selected:    true,
		IsAvailable: true,
	}
}

func UpdateRoomFromParams(params *UpdateRoomParams) *repo.Room {
	return &repo.Room{
		Type:        params.Type,
		Seaside:     params.Seaside,
		Size:        params.Size,
		Price:       params.Price,
		Selected:    params.Selected,
		IsAvailable: params.IsAvailable,
	}
}

func MapRoomResponse(room *repo.Room) *RoomResponse {
	return &RoomResponse{
		Id:          room.Id.Hex(),
		Type:        room.Type,
		Seaside:     room.Seaside,
		Size:        room.Size,
		Price:       room.Price,
		HotelId:     room.HotelId.Hex(),
		Selected:    room.Selected,
		IsAvailable: room.IsAvailable,
	}
}

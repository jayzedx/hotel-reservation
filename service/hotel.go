package service

import (
	"github.com/jayzedx/hotel-reservation/repo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HotelService interface {
	GetHotelById(id string) (*HotelByIdResponse, error)
	GetHotelRooms(id string) (*HotelResponse, error)
	GetHotels(params repo.Hotel) ([]*HotelResponse, error)
	CreateHotel(params CreateHotelParams) (*HotelResponse, error)
	UpdateHotel(id string, params UpdateHotelParams) (*HotelResponse, error)
	DeleteHotel(id string) error
}

// GetHotels, GetHotelRooms
type HotelResponse struct {
	Id         primitive.ObjectID   `json:"id,omitempty"`
	Name       string               `json:"name"`
	Location   string               `json:"location"`
	Rating     int                  `json:"rating"`
	Rooms      []primitive.ObjectID `json:"rooms"`
	HotelRooms []*repo.Room         `json:"hotel_rooms"`
}

// GetHotelById
type HotelByIdResponse struct {
	Id       primitive.ObjectID `json:"id,omitempty"`
	Name     string             `json:"name"`
	Location string             `json:"location"`
	Rating   int                `json:"rating"`
}

type CreateHotelParams struct {
	Name     string               `json:"name"`
	Location string               `json:"location"`
	Rating   int                  `json:"rating"`
	Rooms    []primitive.ObjectID `json:"rooms"`
}

type UpdateHotelParams struct {
	Name        string               `json:"name"`
	Location    string               `json:"location"`
	Rating      int                  `json:"rating"`
	Rooms       []primitive.ObjectID `json:"rooms"`        // update exists rooms
	CreateRooms []CreateRoomParams   `json:"create_rooms"` // create new rooms
}

func (params CreateHotelParams) Validate() map[string]string {
	errors := map[string]string{}
	return errors
}

func (params UpdateHotelParams) Validate() map[string]string {
	errors := map[string]string{}
	return errors
}

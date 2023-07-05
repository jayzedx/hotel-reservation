package service

import (
	"github.com/jayzedx/hotel-reservation/repo"
)

type HotelService interface {
	GetHotelRooms(id string) (*HotelResponse, error)
	GetHotels(params QueryHotelParams) ([]*HotelResponse, error)
	CreateHotel(params CreateHotelParams) (*HotelResponse, error)
	UpdateHotel(id string, params UpdateHotelParams) error
}

// GetHotels, GetHotelRooms
type HotelResponse struct {
	Id         string          `json:"hotel_id,omitempty" mapstructure:"hotel_id"`
	Name       string          `json:"name"`
	Location   string          `json:"location"`
	Rating     int             `json:"rating"`
	HotelRooms []*RoomResponse `json:"hotel_rooms" mapstructure:"hotel_rooms"`
}

type CreateHotelParams struct {
	Name     string `json:"name"`
	Location string `json:"location"`
	Rating   int    `json:"rating"`
}

type UpdateHotelParams struct {
	Name     string `json:"name"`
	Location string `json:"location"`
	Rating   int    `json:"rating"`
}

type QueryHotelParams struct {
	Rating int `query:"rating"`
}

func (params CreateHotelParams) Validate() map[string]string {
	errors := map[string]string{}
	return errors
}

func (params UpdateHotelParams) Validate() map[string]string {
	errors := map[string]string{}
	return errors
}

func CreateHotelFromParams(params *CreateHotelParams) *repo.Hotel {
	return &repo.Hotel{
		Name:     params.Name,
		Location: params.Location,
		Rating:   params.Rating,
	}
}

func UpdateHotelFromParams(params *UpdateHotelParams) *repo.Hotel {
	return &repo.Hotel{
		Name:     params.Name,
		Location: params.Location,
		Rating:   params.Rating,
	}
}

func MapHotelResponse(hotel *repo.Hotel, roomsResponse []*RoomResponse) *HotelResponse {
	return &HotelResponse{
		Id:         hotel.Id.Hex(),
		Name:       hotel.Name,
		Location:   hotel.Location,
		Rating:     hotel.Rating,
		HotelRooms: roomsResponse,
	}
}

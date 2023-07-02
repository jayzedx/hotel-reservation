package service

import (
	"net/http"

	"github.com/jayzedx/hotel-reservation/errs"
	"github.com/jayzedx/hotel-reservation/logs"
	"github.com/jayzedx/hotel-reservation/repo"
	"github.com/jayzedx/hotel-reservation/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type hotelService struct {
	hotelRepository repo.HotelRepository
	roomRepository  repo.RoomRepository
}

func NewHotelService(hotelRepository repo.HotelRepository, roomRepository repo.RoomRepository) *hotelService {
	return &hotelService{
		hotelRepository: hotelRepository,
		roomRepository:  roomRepository,
	}
}

func (s *hotelService) GetHotels(params repo.Hotel) ([]*HotelResponse, error) {
	var (
		rating = params.Rating
		filter = bson.M{}
	)

	if rating != 0 {
		filter = bson.M{
			"rooms": bson.M{"$exists": true},
			"rating": bson.M{
				"$eq": rating,
			},
		}
	}

	hotels, err := s.hotelRepository.GetHotels(filter)
	if err != nil {
		return nil, errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "Invalid data provided. Please check your input and try again.",
		}
	}

	matchStage := bson.M{
		"$match": bson.M{
			"hotel_id": bson.M{
				"$exists": true,
			},
		},
	}
	pipeline := []bson.M{
		matchStage,
	}
	rooms, err := s.roomRepository.GetRoomsByPipeline(pipeline)
	if err != nil {
		logs.Error(err)
		return nil, errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "Invalid data provided. Please check your input and try again.",
		}
	}

	roomsByHotel := make(map[primitive.ObjectID][]*repo.Room)
	for _, room := range rooms {
		roomsByHotel[room.HotelId] = append(roomsByHotel[room.HotelId], room)
	}

	data := []*HotelResponse{}
	for _, hotel := range hotels {
		hotelRooms := roomsByHotel[hotel.Id]
		if len(hotelRooms) == 0 {
			hotelRooms = []*repo.Room{}
		}

		data = append(data,
			&HotelResponse{
				Id:         hotel.Id,
				Name:       hotel.Name,
				Location:   hotel.Location,
				Rating:     hotel.Rating,
				Rooms:      hotel.Rooms,
				HotelRooms: hotelRooms,
			})
	}
	return data, nil
}

func (s *hotelService) GetHotelById(id string) (*HotelByIdResponse, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "Invalid id provided. Please check your input and try again.",
		}
	}

	hotel, err := s.hotelRepository.GetHotelById(oid)
	if err != nil {
		logs.Error(err)
		return nil, errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "Invalid data provided. Please check your input and try again.",
		}
	}

	data := &HotelByIdResponse{
		Id:       hotel.Id,
		Name:     hotel.Name,
		Location: hotel.Location,
		Rating:   hotel.Rating,
	}
	return data, nil
}

func (s *hotelService) GetHotelRooms(id string) (*HotelResponse, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "Invalid id provided. Please check your input and try again.",
		}
	}
	hotel, err := s.hotelRepository.GetHotelById(oid)
	if err != nil {
		logs.Error(err)
		return nil, errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "Invalid data provided. Please check your input and try again.",
		}
	}
	filter := bson.M{
		"hotel_id": hotel.Id,
	}
	rooms, err := s.roomRepository.GetRooms(filter)
	if err != nil {
		logs.Error(err)
		return nil, errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "Invalid data provided. Please check your input and try again.",
		}
	}
	hotelRooms := []*repo.Room{}
	if len(rooms) > 0 {
		hotelRooms = rooms
	}

	data := &HotelResponse{
		Id:         hotel.Id,
		Name:       hotel.Name,
		Location:   hotel.Location,
		Rating:     hotel.Rating,
		Rooms:      hotel.Rooms,
		HotelRooms: hotelRooms,
	}
	return data, nil
}

func (s *hotelService) CreateHotel(params CreateHotelParams) (*HotelResponse, error) {
	//validation for creating hotel
	hotel := &repo.Hotel{
		Name:     params.Name,
		Location: params.Location,
		Rating:   params.Rating,
		Rooms:    []primitive.ObjectID{},
	}
	if err := s.hotelRepository.CreateHotel(hotel); err != nil {
		return nil, err
	}
	data := &HotelResponse{
		Id:         hotel.Id,
		Name:       hotel.Name,
		Location:   hotel.Location,
		Rating:     hotel.Rating,
		Rooms:      hotel.Rooms,
		HotelRooms: []*repo.Room{},
	}
	return data, nil
}

func (s *hotelService) UpdateHotel(id string, params UpdateHotelParams) (*HotelResponse, error) {

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "Invalid id provided. Please check your input and try again.",
		}
	}

	var (
		errors = map[string]string{}
		filter = bson.M{"_id": oid}
	)
	// 1) create rooms by hotel id
	// 1.1) validate data for creating new room

	// 2) update exist rooms by hotel id
	// 2.1) check exists room id

	// update := bson.M{"$push": bson.M{"rooms": params.Rooms}}
	// if err := s.hotelRepository.UpdateRooms(filter, update); err != nil {
	// 	return nil, err
	// }

	// 3) update hotel by hotel id
	// 3.1) validate data for updating hotel
	if errors = params.Validate(); len(errors) > 0 {
		return nil, errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "Validation errors. Please check your input and try again.",
			Errors:  errors,
		}
	}
	updateHotel := util.ToBSON(params)
	// fmt.Println(updateHotel)
	// fmt.Println(filter)
	if err = s.hotelRepository.UpdateHotel(filter, updateHotel); err != nil {
		logs.Error(err)
		return nil, err
	}

	// 4) query hotel by id
	hotelResponse, err := s.GetHotelRooms(id)
	if err != nil {
		return nil, err
	}

	return hotelResponse, nil
}
func (s *hotelService) DeleteHotel(id string) error {
	return nil
}

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

func (s *hotelService) GetHotels(params QueryHotelParams) ([]*HotelResponse, error) {
	var (
		rating = params.Rating
		filter = bson.M{}
	)

	if rating != 0 {
		filter = bson.M{
			"rating": bson.M{
				"$eq": rating,
			},
		}
	}

	hotels, err := s.hotelRepository.GetHotels(filter)
	if err != nil {
		return nil, errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "Invalid data provided",
		}
	}

	/*
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
	*/
	rooms, err := s.roomRepository.GetRooms(bson.M{})
	if err != nil {
		return nil, errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "Invalid data provided",
		}
	}

	// map rooms by hotel id
	roomsReponse := make(map[primitive.ObjectID][]*RoomResponse)
	for _, room := range rooms {
		roomsReponse[room.HotelId] = append(roomsReponse[room.HotelId], MapRoomResponse(room))
	}

	// building response data
	data := []*HotelResponse{}
	for _, hotel := range hotels {
		data = append(data, MapHotelResponse(hotel, roomsReponse[hotel.Id]))
	}
	return data, nil
}

func (s *hotelService) GetHotelRooms(id string) (*HotelResponse, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "Invalid id provided",
		}
	}

	hotel, err := s.hotelRepository.GetHotelById(oid)
	if err != nil {
		logs.Error(err)
		return nil, errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "Upexpected Error",
		}
	}

	filter := bson.M{"hotel_id": hotel.Id}
	rooms, err := s.roomRepository.GetRooms(filter)
	if err != nil {
		logs.Error(err)
		return nil, errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "Upexpected Error",
		}
	}

	roomsReponse := []*RoomResponse{}
	for _, room := range rooms {
		roomsReponse = append(roomsReponse, MapRoomResponse(room))
	}

	return MapHotelResponse(hotel, roomsReponse), nil
}

func (s *hotelService) CreateHotel(params CreateHotelParams) (*HotelResponse, error) {
	if errors := params.Validate(); len(errors) > 0 {
		return nil, errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "Validation errors",
			Errors:  errors,
		}
	}

	hotel := CreateHotelFromParams(&params)
	if err := s.hotelRepository.CreateHotel(hotel); err != nil {
		return nil, err
	}

	return MapHotelResponse(hotel, []*RoomResponse{}), nil
}

func (s *hotelService) UpdateHotel(id string, params UpdateHotelParams) error {
	hotelId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "Invalid id provided",
		}
	}

	// validation
	if errors := params.Validate(); len(errors) > 0 {
		return errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "Validation errors",
			Errors:  errors,
		}
	}

	dataMap, err := util.ConvertToBsonM(params)
	if err != nil {
		return err
	}

	if len(dataMap) == 0 {
		return errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "No field to update.",
		}
	}

	update := bson.M{
		"$set": dataMap,
	}
	filter := bson.M{"_id": hotelId}
	if _, err = s.hotelRepository.UpdateHotel(filter, update); err != nil {
		logs.Error(err)
		return err
	}

	return nil
}

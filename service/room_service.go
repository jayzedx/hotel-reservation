package service

import (
	"net/http"

	"github.com/jayzedx/hotel-reservation/errs"
	"github.com/jayzedx/hotel-reservation/repo"
	"github.com/jayzedx/hotel-reservation/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type roomService struct {
	roomRepository  repo.RoomRepository
	hotelRepository repo.HotelRepository
}

func NewRoomService(roomRepository repo.RoomRepository, hotelRepository repo.HotelRepository) *roomService {
	return &roomService{
		roomRepository:  roomRepository,
		hotelRepository: hotelRepository,
	}
}

func (s *roomService) CreateRoom(params CreateRoomParams) (*RoomResponse, error) {
	// validation
	if errors := params.Validate(); len(errors) > 0 {
		return nil, errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "Validation errors",
			Errors:  errors,
		}
	}

	// check hotel id is exist
	if _, err := s.hotelRepository.GetHotelById(params.HotelId); err != nil {
		return nil, errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "Your hotel id isn't exist",
		}
	}

	room := CreateRoomFromParams(&params)
	if err := s.roomRepository.CreateRoom(room); err != nil {
		return nil, err
	}

	return MapRoomResponse(room), nil
}

func (s *roomService) UpdateRoom(id string, params UpdateRoomParams) error {
	oid, err := primitive.ObjectIDFromHex(id)
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

	filter := bson.M{"_id": oid}
	update := bson.M{"$set": dataMap}
	if _, err = s.roomRepository.UpdateRoom(filter, update); err != nil {
		return errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "Can't update this room.",
		}
	}

	return nil
}

func (s *roomService) DeleteRoom(id string) error {
	roomId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "Invalid id provided",
		}
	}

	if err = s.roomRepository.DeleteRoom(roomId); err != nil {
		return errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "Can't delete this room.",
		}
	}

	return nil
}

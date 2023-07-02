package service

import (
	"github.com/jayzedx/hotel-reservation/repo"
)

type roomService struct {
	roomRepository repo.RoomRepository
}

func NewRoomService(roomRepository repo.RoomRepository) *roomService {
	return &roomService{
		roomRepository: roomRepository,
	}
}

func (s *roomService) GetRoomById(id string) (*RoomResponse, error) {
	return nil, nil
}
func (s *roomService) CreateRoom(params CreateRoomParams) (*RoomResponse, error) {
	return nil, nil
}
func (s *roomService) UpdateRoom(id string, params UpdateRoomParams) (*RoomResponse, error) {
	return nil, nil
}
func (s *roomService) DeleteRoom(id string) error {
	return nil
}

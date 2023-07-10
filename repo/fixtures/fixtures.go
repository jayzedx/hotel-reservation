package fixtures

import (
	"fmt"

	"github.com/jayzedx/hotel-reservation/repo"
	"github.com/jayzedx/hotel-reservation/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateUser(userRepo repo.UserRepository, email string, firstName string, lastName string, password string, isAdmin bool) *repo.User {
	params := service.CreateUserParams{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Password:  password,
	}
	user, err := service.CreateUserFromParams(&params)
	if isAdmin {
		user.IsAdmin = true
	} else {
		user.IsAdmin = false
	}

	if err != nil {
		panic(err)
	}
	if err := userRepo.CreateUser(user); err != nil {
		panic(err)
	}
	fmt.Println("created user : ", user)
	return user
}

func CreateRoom(roomRepo repo.RoomRepository, roomType repo.RoomType, seaside bool, size string, price float64, hotelId primitive.ObjectID, selected bool) *repo.Room {
	params := service.CreateRoomParams{
		Type:     roomType,
		Seaside:  seaside,
		Size:     size,
		Price:    price,
		HotelId:  hotelId,
		Selected: selected,
	}
	room := service.CreateRoomFromParams(&params)

	if err := roomRepo.CreateRoom(room); err != nil {
		panic("can't create room")
	}
	fmt.Println("created room : ", room)
	return room
}

func CreateHotel(hotelRepo repo.HotelRepository, name string, location string, rating int) *repo.Hotel {
	params := service.CreateHotelParams{
		Name:     name,
		Location: location,
		Rating:   rating,
	}
	hotel := service.CreateHotelFromParams(&params)

	if err := hotelRepo.CreateHotel(hotel); err != nil {
		panic("can't create hotel")
	}
	fmt.Println("created hotel : ", hotel)
	return hotel
}

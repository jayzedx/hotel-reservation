package main

import (
	"fmt"

	"github.com/jayzedx/hotel-reservation/repo"
	"github.com/jayzedx/hotel-reservation/test"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	client    *mongo.Client
	dbName    string
	roomRepo  repo.RoomRepository
	hotelRepo repo.HotelRepository
	userRepo  repo.UserRepository
)

func main() {
	testApp := test.NewTestApp()
	client = testApp.GetClient()
	dbName = testApp.GetDatabaseName()
	initRepo()

	fmt.Println("=== seeding the database ===")

	//### hotel 1 ###
	hotelId := createHotel("Bellucia", "France", 4)
	createRoom(repo.SingleRoomType, false, "small", 99.9, hotelId, true)
	createRoom(repo.DoubleRoomType, false, "normal", 120, hotelId, true)

	//### user 1 ###
	createUser()

	fmt.Println("============================")

}

func initRepo() {
	hotelRepo = repo.NewHotelRepository(client, dbName)
	roomRepo = repo.NewRoomRepository(client, dbName)
	userRepo = repo.NewUserRepository(client, dbName)
}

func createHotel(name string, location string, rating int) primitive.ObjectID {
	hotel := repo.Hotel{
		Name:     name,
		Location: location,
		Rating:   rating,
	}
	if err := hotelRepo.CreateHotel(&hotel); err != nil {
		panic("can't create hotel")
	}
	fmt.Println("created hotel : ", hotel)
	return hotel.Id
}

func createRoom(roomType repo.RoomType, seaside bool, size string, price float64, hotelId primitive.ObjectID, selected bool) {
	room := repo.Room{
		Type:     roomType,
		Seaside:  seaside,
		Size:     size,
		Price:    price,
		HotelId:  hotelId,
		Selected: selected,
	}
	if err := roomRepo.CreateRoom(&room); err != nil {
		panic("can't create room")
	}
	fmt.Println("created room : ", room)
}

func createUser() {

}

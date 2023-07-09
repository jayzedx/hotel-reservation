package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jayzedx/hotel-reservation/repo"
	"github.com/jayzedx/hotel-reservation/service"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client    *mongo.Client
	dbName    string
	dbUri     string
	roomRepo  repo.RoomRepository
	hotelRepo repo.HotelRepository
	userRepo  repo.UserRepository
)

func main() {
	initConfig()
	initDatabase()
	initRepo()

	fmt.Println("=== seeding the database ===")
	//### hotel ###
	hotelId := createHotel("Bellucia", "France", 4)
	createRoom(repo.SingleRoomType, false, "small", 99.9, hotelId, true)
	createRoom(repo.DoubleRoomType, false, "normal", 120, hotelId, true)

	//### user ###
	createUser("jay@mail.com", "jay", "layman", "1234567", false)
	createUser("mail@mail.com", "Chalermpong", "Dejnarong", "1234567", true)
	fmt.Println("============================")

}

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../")
	// APP_PORT=3000 go run .
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}

func initDatabase() {
	var err error
	mode := os.Getenv("SEED_MODE")
	if strings.TrimSpace(strings.ToUpper(mode)) == "TEST" {
		fmt.Println(mode)

		dbUri = viper.GetString("test_db.uri")
		dbName = viper.GetString("test_db.database")
	} else {
		dbUri = viper.GetString("db.uri")
		dbName = viper.GetString("db.database")
	}

	fmt.Println("dbUri : ", dbUri)
	fmt.Println("dbName : ", dbName)
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(dbUri))
	if err != nil {
		panic(err)
	}
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

func createUser(email string, firstName string, lastName string, password string, isAdmin bool) {
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
}

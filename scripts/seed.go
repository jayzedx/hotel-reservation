package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/jayzedx/hotel-reservation/repo"
	"github.com/jayzedx/hotel-reservation/repo/fixtures"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	initConfig()
	uri := viper.GetString("test_db.uri")
	dbName := viper.GetString("test_db.database")

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic("Error to connect mongo db")
	}

	hotelRepo := repo.NewHotelRepository(client, dbName)
	roomRepo := repo.NewRoomRepository(client, dbName)
	userRepo := repo.NewUserRepository(client, dbName)
	// bookingRepo := testApp.Repo.Booking
	// authRepo := testApp.Repo.Auth

	fmt.Println("=== seeding the database ===")
	// hotel
	hotel := fixtures.CreateHotel(hotelRepo, "Bellucia", "France", 4)

	// room
	fixtures.CreateRoom(roomRepo, repo.SingleRoomType, false, "small", 99.9, hotel.Id, true)

	// user
	fixtures.CreateUser(userRepo, "jay@mail.com", "jay", "layman", "1234567", false)
	fixtures.CreateUser(userRepo, "mail@mail.com", "Chalermpong", "Dejnarong", "1234567", true)
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

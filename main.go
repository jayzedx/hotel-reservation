package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jayzedx/hotel-reservation/handler"
	"github.com/jayzedx/hotel-reservation/logs"
	"github.com/jayzedx/hotel-reservation/repo"
	"github.com/jayzedx/hotel-reservation/service"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
	ErrorHandler: handler.ErrorHandler,
}

var (
	PORT    string
	DB_NAME string
	DB_URI  string
)

func main() {
	initConfig()
	initTimezone()
	client := initDatabase()

	// listenAddr := flag.String("listenAddr", ":5000", "The listen address of the API server")
	// flag.Parse()

	var (
		userRepository  = repo.NewUserRepository(client, DB_NAME)
		roomRepository  = repo.NewRoomRepository(client, DB_NAME)
		hotelRepository = repo.NewHotelRepository(client, DB_NAME)
		authRepository  = repo.NewAuthRepository(client, DB_NAME)

		userService  = service.NewUserService(userRepository)
		hotelService = service.NewHotelService(hotelRepository, roomRepository)
		roomService  = service.NewRoomService(roomRepository, hotelRepository)
		authService  = service.NewAuthService(userRepository, authRepository)

		userHandler  = handler.NewUserHandler(userService)
		hotelHandler = handler.NewHotelHandler(hotelService)
		roomHandler  = handler.NewRoomHandler(roomService)
		authHandler  = handler.NewAuthHandler(authService)

		app   = fiber.New(config)
		auth  = app.Group("/api")
		apiv1 = app.Group("/api/v1")
		// apiv1 = app.Group("/api/v1", middleware.JWTAuthentication)
	)

	// auth handlers
	auth.Post("/auth", authHandler.HandlePostAuthen)

	// user handlers
	apiv1.Get("/users/**", userHandler.HandleGetUsersByParams)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	apiv1.Post("/user/**", userHandler.HandlePostUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Put("/user/:id", userHandler.HandlePutUser)

	// hotel handlers
	apiv1.Get("/hotels/**", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotel/:id/*", hotelHandler.HandleGetHotel)
	apiv1.Post("/hotel", hotelHandler.HandlePostHotel)
	apiv1.Put("/hotel/:id", hotelHandler.HandlePutHotel)
	// apiv1.Delete("/hotel/:id", roomHandler.HandleDeleteHotel)

	// room handlers
	apiv1.Post("/room", roomHandler.HandlePostRoom)
	apiv1.Put("/room/:id", roomHandler.HandlePutRoom)
	apiv1.Delete("/room/:id", roomHandler.HandleDeleteRoom)

	logs.Info("App service start at port " + viper.GetString("app.port"))
	if err := app.Listen(PORT); err != nil {
		logs.Error(err)
	}

}

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	// APP_PORT=3000 go run .
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	PORT = fmt.Sprintf(":%v", viper.GetInt("app.port"))
	DB_NAME = viper.GetString("db.database")
	DB_URI = viper.GetString("db.uri")

}

func initTimezone() {
	ict, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		panic(err)
	}
	time.Local = ict
}

func initDatabase() *mongo.Client {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(DB_URI))
	if err != nil {
		log.Fatal(err)
	}
	return client
}

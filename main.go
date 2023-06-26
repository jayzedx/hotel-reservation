package main

import (
	"context"
	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jayzedx/hotel-reservation/api"
	"github.com/jayzedx/hotel-reservation/api/middleware"
	"github.com/jayzedx/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
	ErrorHandler: func(ctx *fiber.Ctx, err error) error {
		return ctx.JSON(map[string]string{"error": err.Error()})
	}}

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	listenAddr := flag.String("listenAddr", ":5000", "The listen address of the API server")
	flag.Parse()

	//handlers initialization
	var (
		hotelStore = db.NewMongoHotelStore(client, db.DBNAME)
		roomStore  = db.NewMongoRoomStore(client, db.DBNAME, hotelStore)
		userStore  = db.NewMongoUserStore(client, db.DBNAME)

		store = &db.Store{
			Hotel: hotelStore,
			Room:  roomStore,
			User:  userStore,
		}
		userHandler  = api.NewUserHandler(store)
		hotelHandler = api.NewHotelHandler(store)
		authHandler  = api.NewAuthHandler(userStore)

		// userHandler = api.NewUserHandler(db.NewMongoUserStore(client, db.DBNAME))
		// hotelStore   = db.NewMongoHotelStore(client, db.DBNAME)
		// roomStore    = db.NewMongoRoomStore(client, db.DBNAME, hotelStore)
		// hotelHandler = api.NewHotelHandler(hotelStore, roomStore)

		app   = fiber.New(config)
		auth  = app.Group("/api")
		apiv1 = app.Group("/api/v1", middleware.JWTAuthentication)
	)

	// auth
	auth.Post("/auth", authHandler.HandleAuthenticate)

	// user handlers
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Put("/user/:id", userHandler.HandlePutUser)

	// hotel handlers
	apiv1.Get("/hotel", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotel/:id/room", hotelHandler.HandleGetHotel)
	apiv1.Get("/hotel/:id", hotelHandler.HandleGetHotel)

	app.Listen(*listenAddr)
}

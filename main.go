package main

import (
	"context"
	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jayzedx/hotel-reservation/api"
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

	app := fiber.New(config)
	apiv1 := app.Group("/api/v1")

	app.Get("/foo", handleFoo)

	//handlers initialization
	var (
		userHandler = api.NewUserHandler(db.NewMongoUserStore(client, db.DBNAME))

		hotelStore   = db.NewMongoHotelStore(client, db.DBNAME)
		roomStore    = db.NewMongoRoomStore(client, db.DBNAME, hotelStore)
		hotelHandler = api.NewHotelHandler(hotelStore, roomStore)
	)

	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Put("/user/:id", userHandler.HandlePutUser)

	//hotel
	apiv1.Get("/hotel", hotelHandler.HandleGetHotels)

	app.Listen(*listenAddr)
}

func handleFoo(c *fiber.Ctx) error {
	return c.JSON(map[string]string{
		"msg": "Working just fine!",
	})
}

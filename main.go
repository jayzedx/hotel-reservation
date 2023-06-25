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

const dburi = "mongodb://localhost:27017"
const dbname = "hotel-reservation"
const userColl = "users"

var config = fiber.Config{
	ErrorHandler: func(ctx *fiber.Ctx, err error) error {
		return ctx.JSON(map[string]string{"error": err.Error()})
	}}

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	if err != nil {
		log.Fatal(err)
	}

	/*
		ctx := context.Background()
		coll := client.Database(dbname).Collection(userColl)
		user := types.User{
			FirstName: "Jay",
			LastName:  "Layman ",
		}
		_, err = coll.InsertOne(ctx, user)
		if err != nil {
			log.Fatal(err)
		}

		var jay types.User
		if err := coll.FindOne(ctx, bson.M{}).Decode(&jay); err != nil {
			log.Fatal(err)
		}
		fmt.Println(jay)
	*/

	listenAddr := flag.String("listenAddr", ":5000", "The listen address of the API server")
	flag.Parse()

	app := fiber.New(config)
	apiv1 := app.Group("/api/v1")

	app.Get("/foo", handleFoo)

	//handlers initialization
	userHandler := api.NewUserHandler(db.NewMongoUserStore(client))

	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	apiv1.Post("/user", userHandler.HandlePostUser)

	app.Listen(*listenAddr)
}

func handleFoo(c *fiber.Ctx) error {
	return c.JSON(map[string]string{
		"msg": "Working just fine!",
	})
}

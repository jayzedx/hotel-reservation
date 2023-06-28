package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jayzedx/hotel-reservation/db"
	"github.com/jayzedx/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	roomStore  db.RoomStore
	hotelStore db.HotelStore
	userStore  db.UserStore
	ctx        = context.Background()
)

func init() {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	if err = client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	hotelStore = db.NewMongoHotelStore(client, db.DBNAME)
	roomStore = db.NewMongoRoomStore(client, db.DBNAME, hotelStore)
	userStore = db.NewMongoUserStore(client, db.DBNAME)
}

func seedHotel(name string, location string, rating int) {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}
	rooms := []types.Room{
		{
			Type:  types.SingleRoomType,
			Price: 99.9,
			Size:  "small",
		},
		{
			Type:  types.DoubleRoomType,
			Price: 120.0,
			Size:  "normal",
		},
		{
			Type:  types.SeasideRoomType,
			Price: 159.5,
			Size:  "kingsize",
		},
		{
			Type:  types.DeluxeRoomType,
			Price: 199.9,
			Size:  "kingsize",
		},
	}

	fmt.Println("# insert hotel")
	insertedHotel, err := hotelStore.CreateHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(insertedHotel)

	for _, room := range rooms {
		fmt.Println("# insert room")
		room.HotelId = insertedHotel.Id
		insertedRoom, err := roomStore.CreateRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(insertedRoom)

	}
	fmt.Println("============================")
}

func seedUser(firstName string, lastName string, email string, password string) {
	fmt.Println("# insert user")
	params := types.CreateUserParams{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  password,
	}
	if err := params.Validate(); len(err) > 0 {
		log.Fatal(err)
	}
	user, err := types.NewUserFromParams(params)
	if err != nil {
		log.Fatal(err)
	}
	insertedUser, err := userStore.CreateUser(ctx, user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(insertedUser)
	fmt.Println("============================")

}

func main() {
	fmt.Println("=== seeding the database ===")
	seedHotel("Bellucia", "France", 4)
	seedHotel("The cozy hotel", "Netherland", 1)
	seedHotel("BB hotel", "London", 5)
	seedUser("Jay", "Layman", "mail@mail.com", "1234567")
}

package db

import (
	"context"

	"github.com/jayzedx/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const hotelColl = "hotels"

// interface
type HotelStore interface {
	CreateHotel(context.Context, *types.Hotel) (*types.Hotel, error)
	UpdateRoom(ctx context.Context, filter bson.M, params bson.M) error
}

//struct

type MongoHotelStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

// function
func NewMongoHotelStore(client *mongo.Client, dbname string) *MongoHotelStore {
	return &MongoHotelStore{
		client: client,
		coll:   client.Database(dbname).Collection(hotelColl),
	}
}

// implement
func (s *MongoHotelStore) CreateHotel(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error) {
	resp, err := s.coll.InsertOne(ctx, hotel)
	if err != nil {
		return nil, err
	}
	hotel.Id = resp.InsertedID.(primitive.ObjectID)
	return hotel, nil
}

func (s *MongoHotelStore) UpdateRoom(ctx context.Context, filter bson.M, params bson.M) error {
	_, err := s.coll.UpdateOne(ctx, filter, params)
	return err
}

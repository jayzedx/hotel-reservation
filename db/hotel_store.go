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
	CreateHotel(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error)
	UpdateRoom(ctx context.Context, filter bson.M, params bson.M) error
	GetHotels(ctx context.Context, filter bson.M) ([]*types.Hotel, error)
	GetHotelById(ctx context.Context, id primitive.ObjectID) (*types.Hotel, error)
}

//struct

type mongoHotelStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

// function
func NewMongoHotelStore(client *mongo.Client, dbname string) *mongoHotelStore {
	return &mongoHotelStore{
		client: client,
		coll:   client.Database(dbname).Collection(hotelColl),
	}
}

// implement
func (s *mongoHotelStore) CreateHotel(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error) {
	resp, err := s.coll.InsertOne(ctx, hotel)
	if err != nil {
		return nil, err
	}
	hotel.Id = resp.InsertedID.(primitive.ObjectID)
	return hotel, nil
}

func (s *mongoHotelStore) UpdateRoom(ctx context.Context, filter bson.M, params bson.M) error {
	_, err := s.coll.UpdateOne(ctx, filter, params)
	return err
}

func (s *mongoHotelStore) GetHotels(ctx context.Context, filter bson.M) ([]*types.Hotel, error) {
	cur, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var hotels []*types.Hotel
	if err := cur.All(ctx, &hotels); err != nil {
		return nil, err
	}
	if len(hotels) == 0 {
		hotels = []*types.Hotel{}
	}
	return hotels, nil
}

func (s *mongoHotelStore) GetHotelById(ctx context.Context, id primitive.ObjectID) (*types.Hotel, error) {
	var hotel types.Hotel
	if err := s.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&hotel); err != nil {
		return nil, err
	}
	return &hotel, nil
}
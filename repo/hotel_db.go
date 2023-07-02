package repo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const hotelColl = "hotels"

type hotelRepository struct {
	client *mongo.Client
	coll   *mongo.Collection
	ctx    context.Context
}

func NewHotelRepository(client *mongo.Client, dbname string) *hotelRepository {
	return &hotelRepository{
		client: client,
		coll:   client.Database(dbname).Collection(hotelColl),
		ctx:    context.Background(),
	}
}

func (r *hotelRepository) GetHotels(filter bson.M) ([]*Hotel, error) {
	cur, err := r.coll.Find(r.ctx, filter)
	if err != nil {
		return nil, err
	}
	var hotels []*Hotel
	if err := cur.All(r.ctx, &hotels); err != nil {
		return nil, err
	}
	return hotels, nil
}

func (r *hotelRepository) GetHotelById(id primitive.ObjectID) (*Hotel, error) {
	var hotel Hotel
	if err := r.coll.FindOne(r.ctx, bson.M{"_id": id}).Decode(&hotel); err != nil {
		return nil, err
	}
	return &hotel, nil
}

func (r *hotelRepository) CreateHotel(hotel *Hotel) error {
	resp, err := r.coll.InsertOne(r.ctx, hotel)
	if err != nil {
		return err
	}
	hotel.Id = resp.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *hotelRepository) UpdateRooms(filter bson.M, params bson.M) error {
	_, err := r.coll.UpdateOne(r.ctx, filter, params)
	return err
}

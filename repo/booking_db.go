package repo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const bookingColl = "booking"

type bookingRepository struct {
	client *mongo.Client
	coll   *mongo.Collection
	ctx    context.Context
}

func NewBookingRepository(client *mongo.Client, dbname string) *bookingRepository {
	return &bookingRepository{
		client: client,
		coll:   client.Database(dbname).Collection(bookingColl),
		ctx:    context.Background(),
	}
}

func (r *bookingRepository) CreateBooking(booking *Booking) error {
	resp, err := r.coll.InsertOne(r.ctx, booking)
	if err != nil {
		return err
	}
	booking.Id = resp.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *bookingRepository) GetBookings(filter bson.M) ([]*Booking, error) {
	cur, err := r.coll.Find(r.ctx, filter)
	if err != nil {
		return nil, err
	}
	var bookings []*Booking
	if err := cur.All(r.ctx, &bookings); err != nil {
		return nil, err
	}
	return bookings, nil
}

func (r *bookingRepository) GetBookingById(id primitive.ObjectID) (*Booking, error) {
	var booking Booking
	if err := r.coll.FindOne(r.ctx, bson.M{"_id": id}).Decode(&booking); err != nil {
		return nil, err
	}
	return &booking, nil
}

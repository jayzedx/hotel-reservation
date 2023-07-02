package repo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const roomColl = "rooms"

type roomRepository struct {
	client *mongo.Client
	coll   *mongo.Collection
	ctx    context.Context
}

func NewRoomRepository(client *mongo.Client, dbname string) *roomRepository {
	return &roomRepository{
		client: client,
		coll:   client.Database(dbname).Collection(roomColl),
		ctx:    context.Background(),
	}
}

func (r *roomRepository) GetRoomsByPipeline(pipeline []bson.M) ([]*Room, error) {
	cur, err := r.coll.Aggregate(r.ctx, pipeline)
	if err != nil {
		return nil, err
	}
	var rooms []*Room
	if err := cur.All(r.ctx, &rooms); err != nil {
		return nil, err
	}
	return rooms, nil
}

func (r *roomRepository) GetRooms(filter bson.M) ([]*Room, error) {
	resp, err := r.coll.Find(r.ctx, filter)
	if err != nil {
		return nil, err
	}

	var rooms []*Room
	if err := resp.All(r.ctx, &rooms); err != nil {
		return nil, err
	}
	return rooms, nil
}

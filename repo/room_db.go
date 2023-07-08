package repo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (r *roomRepository) GetRoomById(id primitive.ObjectID) (*Room, error) {
	var room Room
	if err := r.coll.FindOne(r.ctx, bson.M{"_id": id}).Decode(&room); err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *roomRepository) CreateRoom(room *Room) error {
	resp, err := r.coll.InsertOne(r.ctx, room)
	if err != nil {
		return err
	}
	room.Id = resp.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *roomRepository) UpdateRoom(filter bson.M, update interface{}) (int64, error) {
	res, err := r.coll.UpdateOne(r.ctx, filter, update)
	if err != nil {
		return 0, err
	}
	return res.ModifiedCount, nil
}

func (r *roomRepository) GetRoomIds(hotelId primitive.ObjectID) ([]primitive.ObjectID, error) {
	pipeline := []bson.M{
		{
			"$match": bson.M{"hotel_id": hotelId},
		},
		{
			"$group": bson.M{
				"_id":      nil,
				"room_ids": bson.M{"$push": "$_id"},
			},
		},
		{
			"$project": bson.M{
				"_id":      0,
				"room_ids": 1,
			},
		}}

	cur, err := r.coll.Aggregate(r.ctx, pipeline)
	if err != nil {
		return nil, err
	}

	var results struct {
		RoomIds []primitive.ObjectID `bson:"room_ids"`
	}
	if cur.Next(r.ctx) {
		err = cur.Decode(&results)
		if err != nil {
			return nil, err
		}
	}
	// dataType := reflect.TypeOf(data)
	// fmt.Println("Type:", dataType)
	return results.RoomIds, nil
}

func (r *roomRepository) DeleteRoom(id primitive.ObjectID) error {
	res, err := r.coll.DeleteOne(r.ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const userColl = "users"

type userRepository struct {
	client *mongo.Client
	coll   *mongo.Collection
	ctx    context.Context
}

func NewUserRepository(client *mongo.Client, dbname string) *userRepository {
	return &userRepository{
		client: client,
		coll:   client.Database(dbname).Collection(userColl),
		ctx:    context.Background(),
	}
}

func (r *userRepository) GetUserById(id string) (*User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var user User
	if err := r.coll.FindOne(r.ctx, bson.M{"_id": oid}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUsers() ([]*User, error) {
	cur, err := r.coll.Find(r.ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var users []*User
	if err := cur.All(r.ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) CreateUser(user *User) error {
	res, err := r.coll.InsertOne(r.ctx, user)
	if err != nil {
		return err
	}
	user.Id = res.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *userRepository) UpdateUser(filter bson.M, params bson.M) error {
	update := bson.D{
		{
			Key:   "$set",
			Value: params,
		},
	}
	_, err := r.coll.UpdateOne(r.ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepository) DeleteUser(id primitive.ObjectID) error {
	res, err := r.coll.DeleteOne(r.ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return mongo.ErrNoDocuments //fmt.Errorf("user not found")
	}
	return nil
}

func (r *userRepository) GetUserByEmail(email string) (*User, error) {
	var user User
	if err := r.coll.FindOne(r.ctx, bson.M{"email": email}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Drop() error {
	return r.coll.Drop(r.ctx)
}

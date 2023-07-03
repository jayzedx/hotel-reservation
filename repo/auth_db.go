package repo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const authColl = "auths"

type authRepository struct {
	client *mongo.Client
	coll   *mongo.Collection
	ctx    context.Context
}

func NewAuthRepository(client *mongo.Client, dbname string) *authRepository {
	return &authRepository{
		client: client,
		coll:   client.Database(dbname).Collection(authColl),
		ctx:    context.Background(),
	}
}

func (r *authRepository) CreateAuth(auth *Auth) error {
	res, err := r.coll.InsertOne(r.ctx, auth)
	if err != nil {
		return err
	}
	auth.Id = res.InsertedID.(primitive.ObjectID)
	return nil
}

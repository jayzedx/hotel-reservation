package repo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Auth struct {
	Id     primitive.ObjectID `bson:"_id,omitempty"`
	UserId primitive.ObjectID `bson:"user_id"`
	Email  string             `bson:"email"`
	// Token    string             `bson:"token,omitempty"`
	Expires  int64 `bson:"expires"`
	CreateAt int64 `bson:"create_at"`
}

type AuthRepository interface {
	CreateAuth(*Auth) error
}

package repo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Auth struct {
	Id     primitive.ObjectID `bson:"_id,omitempty"`
	UserId primitive.ObjectID `bson:"user_id,omitempty"`
	Email  string             `bson:"email,omitempty"`
	// Token    string             `bson:"token,omitempty"`
	Expires  int64 `bson:"expires,omitempty"`
	CreateAt int64 `bson:"create_at,omitempty"`
}

type AuthRepository interface {
	CreateAuth(*Auth) error
}

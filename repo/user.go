package repo

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id                primitive.ObjectID `bson:"_id,omitempty" query:"id"` //omitempty - don't show json when id is empty
	FirstName         string             `bson:"firstName" query:"firstName"`
	LastName          string             `bson:"lastName" query:"lastName"`
	Email             string             `bson:"email" query:"-"`
	EncryptedPassword string             `bson:"encryptedPassword" query:"-"`
}

type UserRepository interface {
	GetUserById(string) (*User, error)
	GetUsers(filter bson.M) ([]*User, error)
	CreateUser(*User) error
	UpdateUser(filter bson.M, values bson.M) error
	DeleteUser(primitive.ObjectID) error
	Drop() error
}

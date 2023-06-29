package repository

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"` //omitempty - don't show json when id is empty
	FirstName         string             `bson:"firstName" json:"firstName"`
	LastName          string             `bson:"lastName" json:"lastName"`
	Email             string             `bson:"email" json:"email"`
	EncryptedPassword string             `bson:"EncryptedPassword" json:"-"`
}

type UserRepository interface {
	GetUserById(string) (*User, error)
	GetUsers() ([]*User, error)
	CreateUser(*User) error
	UpdateUser(filter bson.M, values bson.M) error
	DeleteUser(primitive.ObjectID) error
	GetUserByEmail(string) (*User, error)
	Drop() error
}

package repo

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id                primitive.ObjectID `bson:"_id,omitempty" query:"id"` //omitempty - don't show json when id is empty
	FirstName         string             `bson:"firstname,omitempty" query:"firstname"`
	LastName          string             `bson:"lastname,omitempty" query:"lastname"`
	Email             string             `bson:"email,omitempty" query:"-"`
	EncryptedPassword string             `bson:"encrypted_password,omitempty" query:"-"`
}

type UserRepository interface {
	GetUserById(string) (*User, error)
	GetUsers(filter bson.M) ([]*User, error)
	CreateUser(*User) error
	UpdateUser(filter bson.M, update bson.M) (int64, error)
	DeleteUser(primitive.ObjectID) error
	Drop() error
}

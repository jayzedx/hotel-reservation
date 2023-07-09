package repo

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id                primitive.ObjectID `bson:"_id,omitempty" query:"id"` //omitempty - don't show json when id is empty
	FirstName         string             `bson:"firstname" query:"firstname"`
	LastName          string             `bson:"lastname" query:"lastname"`
	Email             string             `bson:"email" query:"-"`
	EncryptedPassword string             `bson:"encrypted_password" query:"-"`
	IsAdmin           bool               `bson:"is_admin" query:"-"`
}

type UserRepository interface {
	GetUserById(primitive.ObjectID) (*User, error)
	GetUsers(filter bson.M) ([]*User, error)
	GetUserByEmail(email string) (*User, error)
	CreateUser(*User) error
	UpdateUser(filter bson.M, update bson.M) (int64, error)
	DeleteUser(primitive.ObjectID) error
	Drop() error
}

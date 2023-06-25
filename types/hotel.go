package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Hotel struct {
	Id       primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"` //omitempty - don't show json when id is empty
	Name     string               `bson:"Name" json:"name"`
	Location string               `bson:"Location" json:"location"`
	Rooms    []primitive.ObjectID `bson:"rooms" json:"rooms"`
}

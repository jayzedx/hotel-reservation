package types

type User struct {
	Id        string `bson:"_id,omitempty" json:"id,omitempty"` //omitempty - don't show json when id is empty
	FirstName string `bson:"firstName" json:"firstName"`
	LastName  string `bson:"lastName" json:"lastName"`
}

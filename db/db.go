package db

const (
	DBNAME      = "hotel-reservation"
	TEST_DBNAME = "hotel-reservation-test"
	DBURI       = "mongodb://localhost:27017"
)

// func ToObjectID(id string) primitive.ObjectID {
// 	oid, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return oid
// }

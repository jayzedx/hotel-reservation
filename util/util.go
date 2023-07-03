package util

import (
	"go.mongodb.org/mongo-driver/bson"
)

func ConvertToBsonM(data interface{}) (bson.M, error) {
	bsonData, err := bson.Marshal(data)
	if err != nil {
		return nil, err
	}
	var bsonM bson.M
	err = bson.Unmarshal(bsonData, &bsonM)
	if err != nil {
		return nil, err
	}

	return bsonM, nil
}

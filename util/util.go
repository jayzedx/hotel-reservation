package util

import (
	"time"

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

func ConvertTimeToString(dt time.Time) string {
	if dt.IsZero() {
		return string([]byte{})
	}
	return dt.Format(time.RFC3339)
}

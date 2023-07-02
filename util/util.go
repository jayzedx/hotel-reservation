package util

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
)

func ToBSON(data interface{}) bson.M {
	m := bson.M{}
	dataMap := make(map[string]interface{})

	jsonData, _ := json.Marshal(data)
	json.Unmarshal(jsonData, &dataMap)

	for key, value := range dataMap {
		// fmt.Printf("Key: %s, Value: %v\n", key, value)
		if value != nil && value != "" {
			m[key] = value
		}
	}
	return m
}

package test

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type dbConfig struct {
	uri  string
	name string
}

type TestApp struct {
	client *mongo.Client
	db     dbConfig
}

func (t *TestApp) GetClient() *mongo.Client {
	return t.client
}

func (t *TestApp) GetDatabaseName() string {
	return t.db.name
}

func NewTestApp() *TestApp {
	initConfig()
	var (
		TEST_DB_URI  = viper.GetString("test_db.uri")
		TEST_DB_NAME = viper.GetString("test_db.database")

		testApp TestApp
		err     error
	)
	fmt.Println("DB URI : ", TEST_DB_URI)
	fmt.Println("DB NAME : ", TEST_DB_NAME)
	testApp.client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(TEST_DB_URI))
	testApp.db.name = TEST_DB_NAME
	testApp.db.uri = TEST_DB_URI

	if err != nil {
		log.Fatal(err)
	}
	return &testApp
}

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../")
	// APP_PORT=3000 go run .
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	configPath := os.Getenv("CONFIG_PATH")
	if configPath != "" {
		viper.AddConfigPath(configPath)
	}

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}

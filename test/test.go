package test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/jayzedx/hotel-reservation/repo"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testApp struct {
	client   *mongo.Client
	database dbConfig
	repo     *store
}

type dbConfig struct {
	uri  string
	name string
}

type store struct {
	user    repo.UserRepository
	hotel   repo.HotelRepository
	room    repo.RoomRepository
	auth    repo.AuthRepository
	booking repo.BookingRepository
}

var (
	uri  string
	name string
)

func NewTestApp(config ...dbConfig) *testApp {
	if len(config) > 0 {
		uri = config[0].uri
		name = config[0].name
	} else {
		initConfig()
		uri = viper.GetString("test_db.uri")
		name = viper.GetString("test_db.database")
	}

	fmt.Println("DB URI : ", uri)
	fmt.Println("DB NAME : ", name)

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic("Error to connect mongo db")
	}

	var testApp = &testApp{
		client: client,
		database: dbConfig{
			uri:  uri,
			name: name,
		},
		repo: &store{
			user:    repo.NewUserRepository(client, name),
			hotel:   repo.NewHotelRepository(client, name),
			room:    repo.NewRoomRepository(client, name),
			auth:    repo.NewAuthRepository(client, name),
			booking: repo.NewBookingRepository(client, name),
		},
	}

	return testApp
}

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../")
	// APP_PORT=3000 go run .
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}

func (app *testApp) teardown(t *testing.T) {
	if err := app.client.Database(app.database.name).Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

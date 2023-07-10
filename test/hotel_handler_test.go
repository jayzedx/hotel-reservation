package test

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/jayzedx/hotel-reservation/handler"
	"github.com/jayzedx/hotel-reservation/repo"
	"github.com/jayzedx/hotel-reservation/repo/fixtures"
	"github.com/jayzedx/hotel-reservation/resp"
	"github.com/jayzedx/hotel-reservation/service"
	"github.com/mitchellh/mapstructure"
)

func TestHandleGetHotels(t *testing.T) {

	app := fiber.New()
	testApp := NewTestApp()

	hotelRepo := testApp.repo.hotel
	roomRepo := testApp.repo.room
	hotelService := service.NewHotelService(hotelRepo, roomRepo)
	hotelHandler := handler.NewHotelHandler(hotelService)

	defer testApp.teardown(t)

	// hotel
	hotel := fixtures.CreateHotel(hotelRepo, "Bellucia", "France", 4)
	// room
	fixtures.CreateRoom(roomRepo, repo.SingleRoomType, false, "small", 99.9, hotel.Id, true)

	app.Get("/hotels", hotelHandler.HandleGetHotels)
	values := url.Values{}
	values.Add("rating", "1")

	req := httptest.NewRequest("GET", "/hotels", nil)
	req.URL.RawQuery = values.Encode()

	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	// bodyBytes, _ := ioutil.ReadAll(res.Body)
	// jsonStr := string(bodyBytes)
	// fmt.Println(jsonStr)

	var response resp.Response
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}
	res.Body.Close()

	var hotels []service.HotelResponse
	err = mapstructure.Decode(response.Data, &hotels)
	if err != nil {
		t.Fatal(err)
	}

	if len(hotels) == 0 {
		t.Fatalf("expecting hotels but got 0 hotels")
	}
	if len(hotels[0].HotelRooms) == 0 {
		t.Fatalf("expecting rooms but got 0 rooms")
	}
	if hotels[0].Id == "" {
		fmt.Println(hotels[0].Id)
		t.Fatal("expecting a hotel id to be set")
	}
	if hotels[0].Name == "" {
		t.Fatal("expecting a hotel name to be set")
	}
	if hotels[0].Location == "" {
		t.Fatal("expecting a locationto be set")
	}
	if hotels[0].Rating == 0 {
		t.Fatal("expecting a rating be set")
	}

}

func SetupGetHotels() {

}

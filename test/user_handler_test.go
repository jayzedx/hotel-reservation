package test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/jayzedx/hotel-reservation/handler"
	"github.com/jayzedx/hotel-reservation/repo"
	"github.com/jayzedx/hotel-reservation/resp"
	"github.com/jayzedx/hotel-reservation/service"
	"github.com/mitchellh/mapstructure"
)

func TestHandlePostUser(t *testing.T) {

	app := fiber.New()
	testApp := NewTestApp()

	userRepo := repo.NewUserRepository(testApp.client, testApp.db.name)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	defer userTeardown(t, userRepo)

	app.Post("/user", userHandler.HandlePostUser)
	params := service.CreateUserParams{
		Email:     "foo@mail.com",
		FirstName: "Mars",
		LastName:  "Fullmaker",
		Password:  "12345678",
	}

	byteValue, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/user", bytes.NewReader(byteValue))
	req.Header.Add("Content-Type", "application/json")

	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	var response resp.Response
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}
	res.Body.Close()

	var user service.UserResponse
	err = mapstructure.Decode(response.Data, &user)
	if err != nil {
		t.Fatal(err)
	}

	if user.Id == "" {
		t.Fatal("expecting a user id to be set")
	}
	if user.FirstName != params.FirstName {
		t.Fatalf("expected first name %s but got %s", params.FirstName, user.FirstName)
	}
	if user.LastName != params.LastName {
		t.Fatalf("expected last name %s but got %s", params.LastName, user.LastName)
	}
	if user.Email != params.Email {
		t.Fatalf("expected email %s but got %s", params.Email, user.Email)
	}

}

func userTeardown(t *testing.T, userRepo repo.UserRepository) {
	if err := userRepo.Drop(); err != nil {
		t.Fatal(err)
	}
}

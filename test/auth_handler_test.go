package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/jayzedx/hotel-reservation/handler"
	"github.com/jayzedx/hotel-reservation/repo"
	"github.com/jayzedx/hotel-reservation/resp"
	"github.com/jayzedx/hotel-reservation/service"
	"github.com/mitchellh/mapstructure"
)

func TestAuthenSuccess(t *testing.T) {
	app := fiber.New()
	testApp := NewTestApp()

	userRepository := repo.NewUserRepository(testApp.client, testApp.db.name)
	authRepository := repo.NewAuthRepository(testApp.client, testApp.db.name)
	authService := service.NewAuthService(userRepository, authRepository)
	authHandler := handler.NewAuthHandler(authService)

	_, err := SetupAuthen(userRepository)
	if err != nil {
		t.Fatal("Error from creating user")
	}

	// authen section
	app.Post("/auth", authHandler.HandlePostAuthen)
	params := service.CreateAuthParams{
		Email:    "mail@mail.com",
		Password: "1234567",
	}

	byteValue, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(byteValue))
	req.Header.Add("Content-Type", "application/json")

	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	var response resp.Response
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		t.Fatal("unauthorize")
	}
	res.Body.Close()

	var auth service.AuthResponse
	err = mapstructure.Decode(response.Data, &auth)
	if err != nil {
		t.Fatal(err)
	}

	if auth.Email != params.Email {
		t.Fatalf("expected email %s but got %s", params.Email, auth.Email)
	}
	if auth.Token == "" {
		t.Fatal("expecting a token to be set")
	}

	defer userTeardown(t, userRepository)
}

func SetupAuthen(userRepo repo.UserRepository) (*repo.User, error) {
	params := service.CreateUserParams{
		FirstName: "Jay",
		LastName:  "Layman",
		Email:     "mail@mail.com",
		Password:  "1234567",
	}

	user, err := service.CreateUserFromParams(&params)
	if err != nil {
		return nil, errors.New("Creating user error")
	}

	if err = userRepo.CreateUser(user); err != nil {
		return nil, errors.New("Creating user error")
	}

	return user, nil
}

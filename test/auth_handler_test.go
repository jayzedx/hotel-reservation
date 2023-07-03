package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/jayzedx/hotel-reservation/api/response"
	"github.com/jayzedx/hotel-reservation/db"
	"github.com/jayzedx/hotel-reservation/types"
)

func TestAuthenticateSuccess(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	insertedUser := insertTestAuthen(t, tdb.store.User)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.store.User)
	app.Post("/auth", authHandler.HandleAuthenticate)

	params := AuthParams{
		Email:    "mail@mail.com",
		Password: "1234567",
	}
	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		//defer resp.Body.Close()
		// body, err := ioutil.ReadAll(resp.Body)
		// if err != nil {
		// 	fmt.Println("Failed to read response body:", err)
		// 	return
		// }
		// fmt.Println("Response Status:", resp.Status)
		// fmt.Println("Response Body:", string(body))
		t.Fatalf("expected http status of 200 but got %d", resp.StatusCode)
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		t.Fatal(err)
	}
	if authResp.Token == "" {
		t.Fatal("expected the JWT token to be present in the auth reponse")
	}

	// Set the encrypted password to an empty string, because we do NOT return that any response
	insertedUser.EncryptedPassword = ""
	// fmt.Println(insertedUser)
	// fmt.Println(authResp.User)
	if !reflect.DeepEqual(insertedUser, authResp.User) {
		t.Fatalf("expected the user to be present in the inserted user reponse")
	}
}

func TestAuthenticateWithWrongPassword(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	insertTestAuthen(t, tdb.store.User)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.store.User)
	app.Post("/auth", authHandler.HandleAuthenticate)

	params := AuthParams{
		Email:    "mail@mail.com",
		Password: "passwordisnotcorrect",
	}
	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected http status of 400 but got %d", resp.StatusCode)
	}

	var headResp response.HeadResponse
	if err := json.NewDecoder(resp.Body).Decode(&headResp); err != nil {
		t.Fatal(err)
	}
	if headResp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected response status code to be error but got %d", headResp.StatusCode)
	}
	if headResp.Status != "error" {
		t.Fatalf("expected response status to be error but got %s", headResp.Status)
	}
	if headResp.Msg != "invalid credentials" {
		t.Fatalf("expected response msg to be <invalid credentials> but got %s", headResp.Msg)
	}
}

func insertTestAuthen(t *testing.T, userStore db.UserStore) *types.User {
	params := types.CreateUserParams{
		FirstName: "Jay",
		LastName:  "Layman",
		Email:     "mail@mail.com",
		Password:  "1234567",
	}
	if err := params.Validate(); len(err) > 0 {
		t.Fatal(err)
	}
	user, err := types.NewUserFromParams(params)
	if err != nil {
		t.Fatal(err)
	}
	_, err = userStore.CreateUser(context.TODO(), user)
	if err != nil {
		t.Fatal(err)
	}
	return user
}

package service

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/jayzedx/hotel-reservation/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	GetUserById(id string) (*UserResponse, error)
	GetUsers() ([]*UserResponse, error)
	CreateUser(params CreateUserParams) (*UserResponse, error)
	UpdateUser(id string, params UpdateUserParams) error
	DeleteUser(id string) error
	GetUserByEmail(params UserQueryParams) (*UserResponse, error)
	Drop() error
}

type CreateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type UpdateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type UserQueryParams struct {
	Email string `json:"email"`
}

type UserResponse struct {
	Id        primitive.ObjectID `json:"id,omitempty"` //omitempty - don't show json when id is empty
	FirstName string             `json:"firstName,omitempty"`
	LastName  string             `json:"lastName,omitempty"`
	Email     string             `json:"email,omitempty"`
}

func MapToUserResponse(user *repository.User) *UserResponse {
	return &UserResponse{
		Id:        user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}
}

// validation
const (
	bcryptCost      = 12
	minFirstNameLen = 2
	minLastNameLen  = 2
	minPasswordLen  = 7
)

// validation for create
func (params CreateUserParams) Validate() map[string]string {
	errors := map[string]string{}
	if len(params.FirstName) < minFirstNameLen {
		errors["firstName"] = fmt.Sprintf("firstName length should be at least %d characters", minFirstNameLen)
	}
	if len(params.LastName) < minFirstNameLen {
		errors["lastName"] = fmt.Sprintf("lastName length should be at least %d characters", minLastNameLen)
	}
	if len(params.Password) < minPasswordLen {
		errors["password"] = fmt.Sprintf("password length should be at least %d characters", minPasswordLen)
	}
	if !isEmailValid(params.Email) {
		errors["email"] = fmt.Sprintf("Email is invalid")
	}
	return errors
}

func isEmailValid(e string) bool {
	emailRegx := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegx.MatchString(e)
}

func IsValidPassword(encpw string, pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(encpw), []byte(pw)) == nil
}

func ToBSON(data interface{}) bson.M {
	m := bson.M{}
	dataMap := make(map[string]interface{})

	jsonData, _ := json.Marshal(data)
	json.Unmarshal(jsonData, &dataMap)

	// Iterate over the map using for range loop
	for key, value := range dataMap {
		// fmt.Printf("Key: %s, Value: %v\n", key, value)
		m[key] = value
	}
	return m
}

func (params UpdateUserParams) Validate() map[string]string {
	errors := map[string]string{}
	if len(params.FirstName) < minFirstNameLen {
		errors["firstName"] = fmt.Sprintf("firstName length should be at least %d characters", minFirstNameLen)
	}
	if len(params.LastName) < minFirstNameLen {
		errors["lastName"] = fmt.Sprintf("lastName length should be at least %d characters", minLastNameLen)
	}
	return errors
}

func NewUserFromParams(params CreateUserParams) (*repository.User, error) {
	encp, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)
	if err != nil {
		return nil, err
	}
	return &repository.User{
		FirstName:         params.FirstName,
		LastName:          params.LastName,
		Email:             params.Email,
		EncryptedPassword: string(encp),
	}, nil
}

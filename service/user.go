package service

import (
	"fmt"
	"regexp"

	"github.com/jayzedx/hotel-reservation/repo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	GetUserById(id string) (*UserResponse, error)
	GetUsers() ([]*UserResponse, error)
	CreateUser(params CreateUserParams) (*UserResponse, error)
	UpdateUser(id string, params UpdateUserParams) error
	DeleteUser(id string) error
	GetUsersByParams(params repo.User) ([]*UserResponse, error)
	Drop() error
}

type UserResponse struct {
	Id        primitive.ObjectID `json:"id,omitempty"`
	FirstName string             `json:"firstname"`
	LastName  string             `json:"lastname"`
	Email     string             `json:"email"`
}

type CreateUserParams struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type UpdateUserParams struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

const (
	bcryptCost      = 12
	minFirstNameLen = 2
	minLastNameLen  = 2
	minPasswordLen  = 7
)

func (params CreateUserParams) Validate() map[string]string {
	errors := map[string]string{}
	if len(params.FirstName) < minFirstNameLen {
		errors["firstname"] = fmt.Sprintf("firstname length should be at least %d characters", minFirstNameLen)
	}
	if len(params.LastName) < minFirstNameLen {
		errors["lastname"] = fmt.Sprintf("lastname length should be at least %d characters", minLastNameLen)
	}
	if len(params.Password) < minPasswordLen {
		errors["password"] = fmt.Sprintf("password length should be at least %d characters", minPasswordLen)
	}
	if !isEmailValid(params.Email) {
		errors["email"] = fmt.Sprintf("email is invalid")
	}
	return errors
}

func (params UpdateUserParams) Validate() map[string]string {
	errors := map[string]string{}
	if len(params.FirstName) > 0 && len(params.FirstName) < minFirstNameLen {
		errors["firstname"] = fmt.Sprintf("firstname length should be at least %d characters", minFirstNameLen)
	}

	if len(params.LastName) > 0 && len(params.LastName) < minFirstNameLen {
		errors["lastname"] = fmt.Sprintf("lastname length should be at least %d characters", minLastNameLen)
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

func CreateUserFromParams(params *CreateUserParams) (*repo.User, error) {
	encp, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)
	if err != nil {
		return nil, err
	}
	return &repo.User{
		FirstName:         params.FirstName,
		LastName:          params.LastName,
		Email:             params.Email,
		EncryptedPassword: string(encp),
	}, nil
}

func UpdateUserFromParams(params *UpdateUserParams) *repo.User {
	return &repo.User{
		FirstName: params.FirstName,
		LastName:  params.LastName,
	}
}

func MapUserResponse(user *repo.User) *UserResponse {
	return &UserResponse{
		Id:        user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}
}

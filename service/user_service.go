package service

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/jayzedx/hotel-reservation/errs"
	"github.com/jayzedx/hotel-reservation/logs"
	"github.com/jayzedx/hotel-reservation/repo"
	"github.com/jayzedx/hotel-reservation/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	userRepository repo.UserRepository
}

func NewUserService(userRepository repo.UserRepository) *userService {
	return &userService{
		userRepository: userRepository,
	}
}

func (s *userService) GetUserById(id string) (*UserResponse, error) {
	user, err := s.userRepository.GetUserById(id)
	if err != nil {

		if err == mongo.ErrNoDocuments {
			return nil, errs.AppError{
				Code:    http.StatusBadRequest,
				Message: "Data not found.",
			}
		}

		logs.Error(err)
		return nil, err
	}

	data := &UserResponse{
		Id:        user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}
	return data, nil
}

func (s *userService) GetUsers() ([]*UserResponse, error) {
	users, err := s.userRepository.GetUsers(bson.M{})
	if err != nil {
		logs.Error(err)
		return nil, err
	}

	data := []*UserResponse{}
	for _, user := range users {
		data = append(data, &UserResponse{
			Id:        user.Id,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		})
	}
	return data, nil
}

func (s *userService) CreateUser(params CreateUserParams) (*UserResponse, error) {
	if errors := params.Validate(); len(errors) > 0 {
		return nil, errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "Validation errors. Please check your input and try again.",
			Errors:  errors,
		}
	}

	user, err := CreateUserFromParams(params)
	if err != nil {
		logs.Error(err)
		return nil, err
	}

	if err = s.userRepository.CreateUser(user); err != nil {
		logs.Error(err)
		return nil, err
	}

	data := &UserResponse{
		Id:        user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}
	return data, nil
}

func (s *userService) UpdateUser(userId string, params UpdateUserParams) error {
	oid, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "Invalid id provided. Please check your input and try again.",
		}
	}

	if errors := params.Validate(); len(errors) > 0 {
		return errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "Validation errors. Please check your input and try again.",
			Errors:  errors,
		}
	}

	filter := bson.M{"_id": oid}
	values := util.ToBSON(params)

	if length := len(values); length <= 0 {
		return errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "No field to update.",
		}
	}
	if err := s.userRepository.UpdateUser(filter, values); err != nil {
		logs.Error(err)
		return err
	}
	return nil
}

func (s *userService) DeleteUser(userId string) error {
	oid, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "Invalid id provided. Please check your input and try again.",
		}
	}

	if err := s.userRepository.DeleteUser(oid); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errs.AppError{
				Code:    http.StatusBadRequest,
				Message: "No user to delete.",
			}
		}
		logs.Error(err)
		return err
	}
	return nil
}

// seach with regex by params
func (s *userService) GetUsersByParams(params repo.User) ([]*UserResponse, error) {

	filter := createSearchFilter(params)

	users, err := s.userRepository.GetUsers(filter)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			logs.Error(err)
			return nil, err
		} else {
			return nil, nil
		}
	}

	data := []*UserResponse{}
	for _, user := range users {
		data = append(data, &UserResponse{
			Id:        user.Id,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		})
	}
	return data, nil
}

func (s *userService) Drop() error {
	if err := s.userRepository.Drop(); err != nil {
		logs.Error(err)
		return err
	}
	return nil
}

func createSearchFilter(params repo.User) bson.M {
	var filter = bson.M{}

	val := reflect.ValueOf(params)
	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tagName := field.Tag.Get("bson")
		fieldValue := val.Field(i).Interface()
		if fieldValue == "" {
			continue
		}

		id, ok := fieldValue.(primitive.ObjectID)
		if ok {
			if id.IsZero() {
				continue
			} else {
				objId, err := primitive.ObjectIDFromHex(id.Hex())
				if err != nil {
					continue
				}
				filter[tagName] = bson.M{
					"$eq": objId,
				}
				continue
			}
		}

		filter[tagName] = bson.M{
			"$regex": fmt.Sprintf(".*%s.*", fieldValue),
		}
	}
	return filter
}

func CreateUserFromParams(params CreateUserParams) (*repo.User, error) {
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

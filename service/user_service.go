package service

import (
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
				Message: "User not found.",
			}
		} else {
			logs.Error(err)
			return nil, errs.AppError{
				Code:    http.StatusBadRequest,
				Message: "Unexpected error",
			}
		}
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
		return nil, errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "Unexpected error",
		}
	}

	data := []*UserResponse{}
	for _, user := range users {
		data = append(data, MapUserResponse(user))
	}
	return data, nil
}

func (s *userService) CreateUser(params CreateUserParams) (*UserResponse, error) {
	if errors := params.Validate(); len(errors) > 0 {
		return nil, errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "Validation errors",
			Errors:  errors,
		}
	}

	user, err := CreateUserFromParams(&params)
	if err != nil {
		logs.Error(err)
		return nil, errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "Unexpected error",
		}
	}

	if err = s.userRepository.CreateUser(user); err != nil {
		logs.Error(err)
		return nil, errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "Can't create user",
		}
	}
	return MapUserResponse(user), nil
}

func (s *userService) UpdateUser(userId string, params UpdateUserParams) error {
	oid, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "Invalid data provided",
		}
	}
	// validation
	if errors := params.Validate(); len(errors) > 0 {
		return errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "Validation errors",
			Errors:  errors,
		}
	}

	user := UpdateUserFromParams(&params)
	updateUser, err := util.ConvertToBsonM(user)
	if err != nil {
		return err
	}

	if len(updateUser) == 0 {
		return errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "No field to update.",
		}
	}

	filter := bson.M{"_id": oid}
	update := bson.M{"$set": updateUser}

	if _, err := s.userRepository.UpdateUser(filter, update); err != nil {
		logs.Error(err)
		return errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "Can't update this user",
		}
	}
	return nil
}

func (s *userService) DeleteUser(userId string) error {
	oid, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "Invalid data provided",
		}
	}

	if err := s.userRepository.DeleteUser(oid); err != nil {
		logs.Error(err)
		return errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "Can't delete user",
		}
	}
	return nil
}

// seach with regex by params
func (s *userService) GetUsersByParams(params repo.User) ([]*UserResponse, error) {

	filter := createSearchFilter(params)

	users, err := s.userRepository.GetUsers(filter)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return nil, errs.AppError{
				Code:    http.StatusBadRequest,
				Message: "users not found",
			}
		} else {
			logs.Error(err)
			return nil, errs.AppError{
				Code:    http.StatusBadRequest,
				Message: "Upexpected Error",
			}
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
		return errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "Upexpected Error",
		}
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

package service

import (
	"errors"
	"net/http"

	"github.com/jayzedx/hotel-reservation/errs"
	"github.com/jayzedx/hotel-reservation/logs"
	"github.com/jayzedx/hotel-reservation/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) *userService {
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
	return MapToUserResponse(user), nil
}

func (s *userService) GetUsers() ([]*UserResponse, error) {
	users, err := s.userRepository.GetUsers()
	if err != nil {
		logs.Error(err)
		return nil, err
	}

	data := make([]*UserResponse, 0)
	for _, user := range users {
		data = append(data, MapToUserResponse(user))
	}
	return data, nil
}

func (s *userService) CreateUser(params CreateUserParams) (*UserResponse, error) {
	if errors := params.Validate(); len(errors) > 0 {
		return nil, errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "Invalid data provided.",
			Errors:  errors,
		}
	}

	user, err := NewUserFromParams(params)
	if err != nil {
		logs.Error(err)
		return nil, err
	}

	if err = s.userRepository.CreateUser(user); err != nil {
		logs.Error(err)
		return nil, err
	}

	return MapToUserResponse(&repository.User{
		Id: user.Id,
	}), nil
}

func (s *userService) UpdateUser(userId string, params UpdateUserParams) error {
	oid, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "Invalid user id provided.",
		}
	}

	if errors := params.Validate(); len(errors) > 0 {
		return errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "Invalid data provided.",
			Errors:  errors,
		}
	}

	filter := bson.M{"_id": oid}
	values := ToBSON(params)

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
			Message: "Invalid data provided. Please check your input and try again.",
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

func (s *userService) GetUserByEmail(params UserQueryParams) (*UserResponse, error) {
	var email = params.Email
	if email == "" {
		return nil, errs.AppError{
			Code:    http.StatusBadRequest,
			Message: "Your email is empty.",
		}
	}

	user, err := s.userRepository.GetUserByEmail(email)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			logs.Error(err)
			return nil, err
		} else {
			return nil, nil
		}
	}

	return MapToUserResponse(user), nil
}

func (s *userService) Drop() error {
	if err := s.userRepository.Drop(); err != nil {
		logs.Error(err)
		return err
	}
	return nil
}

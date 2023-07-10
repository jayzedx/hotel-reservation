package service

import (
	"errors"
	"net/http"

	"github.com/jayzedx/hotel-reservation/errs"
	"github.com/jayzedx/hotel-reservation/logs"
	"github.com/jayzedx/hotel-reservation/repo"
	"go.mongodb.org/mongo-driver/mongo"
)

type authService struct {
	userRepository repo.UserRepository
	authRepository repo.AuthRepository
}

func NewAuthService(userRepository repo.UserRepository, authRepository repo.AuthRepository) *authService {
	return &authService{
		userRepository: userRepository,
		authRepository: authRepository,
	}
}

func (s *authService) Authenticate(params CreateAuthParams) (*AuthResponse, error) {

	user, err := s.userRepository.GetUserByEmail(params.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errs.AppError{
				Code:    http.StatusUnauthorized,
				Message: "email or password is incorrect.",
			}
		} else {
			logs.Error(err)
			return nil, errs.ErrUnexpected()
		}
	}

	if !IsValidPassword(user.EncryptedPassword, params.Password) {
		return nil, errs.AppError{
			Code:    http.StatusUnauthorized,
			Message: "email or password is incorrect",
		}
	}

	auth := CreateAuthFromUser(user)

	if err := s.authRepository.CreateAuth(auth); err != nil {
		logs.Error(err)
		return nil, errs.ErrUnexpected()
	}
	token, err := createTokenFromAuth(auth)
	if err != nil {
		logs.Error(err)
		return nil, errs.ErrUnexpected()
	}

	return MapAuthResponse(auth, token), nil
}

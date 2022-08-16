package services

import (
	"github.com/mhamm84/pulse-api/internal/data"
	"github.com/mhamm84/pulse-api/internal/jsonlog"
	"time"
)

type userService struct {
	UserRepository data.UserRepository
	TokenService   TokenService
	Logger         *jsonlog.Logger
}

func NewUserService(userRepository data.UserRepository, tokenService TokenService, logger *jsonlog.Logger) UserService {
	return &userService{userRepository, tokenService, logger}
}

func (u *userService) RegisterUser(user *data.User) (*data.Token, error) {

	err := u.UserRepository.Insert(user)
	if err != nil {
		return nil, err
	}
	token, err := u.TokenService.New(user.ID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		return nil, err
	}
	u.Logger.PrintInfo("new token generated", map[string]interface{}{
		"userId": token.UserId,
		"scope":  token.Scope,
	})

	return token, nil
}

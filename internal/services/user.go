package services

import (
	"github.com/mhamm84/pulse-api/internal/data"
	"github.com/mhamm84/pulse-api/internal/jsonlog"
	"github.com/mhamm84/pulse-api/internal/mailer"
	"time"
)

type userService struct {
	UserRepository data.UserRepository
	TokenService   TokenService
	Mailer         *mailer.Mailer
	Logger         *jsonlog.Logger
}

func NewUserService(userRepository data.UserRepository, tokenService TokenService, mailer *mailer.Mailer, logger *jsonlog.Logger) UserService {
	return &userService{userRepository, tokenService, mailer, logger}
}

func (u *userService) GetByEmail(email string) (*data.User, error) {
	return u.UserRepository.GetByEmail(email)
}

func (u *userService) ActivateUser(token string) (*data.User, error) {

	user, err := u.UserRepository.GetUserFromToken(data.ScopeActivation, token)
	if err != nil {
		u.Logger.PrintError(err, map[string]interface{}{"function": "GetUserFromToken"})
		return nil, err
	}

	user.Activated = true
	err = u.UserRepository.Update(user)
	if err != nil {
		u.Logger.PrintError(err, map[string]interface{}{"function": "Update"})
		return nil, err
	}

	err = u.TokenService.DeleteAllForUser(user.ID, data.ScopeActivation)
	if err != nil {
		u.Logger.PrintError(err, map[string]interface{}{"function": "DeleteAllForUser"})
		return nil, err
	}

	return user, nil
}

func (u *userService) RegisterUser(user *data.User) error {

	err := u.UserRepository.Insert(user)
	if err != nil {
		return err
	}

	token, err := u.TokenService.New(user.ID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		return err
	}

	u.Logger.PrintInfo("new token generated", map[string]interface{}{
		"userId": token.UserId,
		"scope":  token.Scope,
	})

	// Send email
	go func() {
		data := map[string]interface{}{
			"activationToken": token.Plaintext,
			"userID":          user.ID}

		mailerErr := u.Mailer.Send(user.Email, "user_welcome.tmpl", data)
		if mailerErr != nil {
			u.Logger.PrintError(mailerErr, nil)
		}
	}()

	return nil
}

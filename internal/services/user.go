package services

import (
	"github.com/mhamm84/pulse-api/internal/data"
	"github.com/mhamm84/pulse-api/internal/mailer"
	"github.com/mhamm84/pulse-api/internal/utils"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"time"
)

type userService struct {
	UserRepository        data.UserRepository
	PermissionsRepository data.PermissionsRepository
	TokenService          TokenService
	Mailer                *mailer.Mailer
}

func NewUserService(userRepository data.UserRepository, permissionsRepository data.PermissionsRepository, tokenService TokenService, mailer *mailer.Mailer) UserService {
	return &userService{userRepository, permissionsRepository, tokenService, mailer}
}

func (u *userService) GetFromToken(tokenScope, tokenplaintext string) (*data.User, error) {
	return u.UserRepository.GetFromToken(tokenScope, tokenplaintext)
}

func (u *userService) GetByEmail(email string) (*data.User, error) {
	return u.UserRepository.GetByEmail(email)
}

func (u *userService) ActivateUser(token string) (*data.User, error) {

	user, err := u.UserRepository.GetFromToken(data.ScopeActivation, token)
	if err != nil {
		msg := "error getting token when getting token in ActivateUser"
		utils.Logger.Error(msg, zap.Error(err), zap.String("token", token))
		return nil, errors.Wrap(err, msg)
	}

	user.Activated = true
	err = u.UserRepository.Update(user)
	if err != nil {
		msg := "error getting token when updating user in ActivateUser"
		utils.Logger.Error(msg, zap.Error(err), zap.Any("user", &user))
		return nil, errors.Wrap(err, msg)
	}

	err = u.TokenService.DeleteAllForUser(user.ID, data.ScopeActivation)
	if err != nil {
		msg := "error getting token when delete all activation tokens for user in ActivateUser"
		utils.Logger.Error(msg, zap.Error(err), zap.Any("user", &user))
		return nil, errors.Wrap(err, msg)
	}

	return user, nil
}

func (u *userService) RegisterUser(user *data.User) error {

	// TODO
	// Need tp create a DB transaction and pass to all repo functions altering DB state

	err := u.UserRepository.Insert(user)
	if err != nil {
		msg := "error inserting user in RegisterUser"
		utils.Logger.Error(msg, zap.Error(err), zap.Any("user", &user))
		return errors.Wrap(err, msg)
	}

	err = u.PermissionsRepository.AddForUser(user.ID, "economic:all")
	if err != nil {
		msg := "error adding economic all permission to user in RegisterUser"
		utils.Logger.Error(msg, zap.Error(err), zap.Any("user", &user))
		return errors.Wrap(err, msg)
	}

	token, err := u.TokenService.New(user.ID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		msg := "error creating new activation token in RegisterUser"
		utils.Logger.Error(msg, zap.Error(err), zap.Any("user", &user))
		return errors.Wrap(err, msg)
	}

	utils.Logger.Info("new token generated",
		zap.Int64("userId", token.UserId),
		zap.String("scope", token.Scope),
	)

	// Send email
	go func() {
		data := map[string]interface{}{
			"activationToken": token.Plaintext,
			"userID":          user.ID}

		mailerErr := u.Mailer.Send(user.Email, "user_welcome.tmpl", data)
		if mailerErr != nil {
			msg := "error sending email in RegisterUser"
			utils.Logger.Error(msg, zap.Error(mailerErr), zap.Any("email", user.Email))
		}
	}()

	return nil
}

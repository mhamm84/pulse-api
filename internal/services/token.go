package services

import (
	"github.com/mhamm84/pulse-api/internal/data"
	"github.com/mhamm84/pulse-api/internal/validator"
	"time"
)

const (
	tokenValidationKey = "token"
)

type tokenService struct {
	TokenRepository data.TokenRepository
}

func NewTokenService(tokenRepo data.TokenRepository) TokenService {
	return &tokenService{tokenRepo}
}

func (s *tokenService) DeleteAllForUser(userId int64, scope string) error {
	return s.TokenRepository.DeleteAllForUser(userId, scope)
}

func (s *tokenService) New(userID int64, ttl time.Duration, scope string) (*data.Token, error) {
	token, err := data.GenerateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}

	err = s.TokenRepository.Insert(token)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func ValidateTokenPlaintext(v *validator.Validator, plaintexttoken string) {
	v.Check(plaintexttoken != "", tokenValidationKey, "must be provided")
	v.Check(len(plaintexttoken) == 26, tokenValidationKey, "token must be 26 bytes long")
}

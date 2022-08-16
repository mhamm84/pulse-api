package services

import (
	"github.com/mhamm84/pulse-api/internal/data"
	"github.com/mhamm84/pulse-api/internal/jsonlog"
	"github.com/mhamm84/pulse-api/internal/validator"
	"time"
)

type tokenService struct {
	TokenRepository data.TokenRepository
	Logger          *jsonlog.Logger
}

func NewTokenService(tokenRepo data.TokenRepository, logger *jsonlog.Logger) TokenService {
	return &tokenService{tokenRepo, logger}
}

const (
	tokenValidationKey = "token"
)

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

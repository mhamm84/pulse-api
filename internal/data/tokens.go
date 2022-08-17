package data

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"time"
)

const (
	ScopeActivation = "activation"
)

type Token struct {
	Plaintext string    `json:"token"`
	Hash      []byte    `json:"-" db:"hash"`
	UserId    int64     `json:"-" db:"user_id"`
	Expiry    time.Time `json:"expiry" db:"expiry"`
	Scope     string    `json:"-" db:"scope"`
}

func GenerateToken(userId int64, ttl time.Duration, scope string) (*Token, error) {
	token := &Token{
		UserId: userId,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]

	return token, nil
}

type TokenRepository interface {
	Insert(token *Token) error
	DeleteAllForUser(userId int64, scope string) error
}

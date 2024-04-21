package security

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var serverId = uuid.New()

const secretLength = 1024

func generateKey() []byte {
	log.Print("Generating key for signing tokens")
	ret := make([]byte, secretLength)
	_, err := rand.Read(ret)
	if err != nil {
		log.Fatal("Failed to generate key for signing tokens")
	}

	// Self check the token
	zeros := 0
	for _, b := range ret {
		if b == 0 {
			zeros++
		}
	}

	if zeros > secretLength/2 {
		log.Fatal("Generated key has too many zeros, it cannot be guaremteed that it is secure")
	}
	return ret
}

var key = generateKey()

type Claims struct {
	GameId    uuid.UUID `json:"gameId"`
	PlayerId  uuid.UUID `json:"playerId"`
	ServerId  uuid.UUID `json:"serverId"`
	IssuedAt  time.Time `json:"iat"`
	ExpiresAt time.Time `json:"exp"`
}

func (c *Claims) GetExpirationTime() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(c.ExpiresAt), nil
}

func (c *Claims) GetIssuedAt() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(c.IssuedAt), nil
}

func (c *Claims) GetNotBefore() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(c.IssuedAt), nil
}

func (c *Claims) GetIssuer() (string, error) {
	return c.ServerId.String(), nil
}

func (c *Claims) GetSubject() (string, error) {
	return c.GameId.String(), nil
}

func (c *Claims) GetAudience() (jwt.ClaimStrings, error) {
	return []string{c.PlayerId.String()}, nil
}

func (c *Claims) Valid() error {
	if c.ServerId != serverId {
		return errors.New("Invalid server")
	}

	if c.ExpiresAt.Before(time.Now()) {
		return errors.New("Token expired")
	}

	if c.IssuedAt.After(time.Now()) {
		return errors.New("Token issued in the future")
	}

	var nilId uuid.UUID
	if c.GameId == nilId {
		return errors.New("Invalid game id")
	}

	if c.PlayerId == nilId {
		return errors.New("Invalid player id")
	}
	return nil
}

func NewToken(gameId uuid.UUID, playerId uuid.UUID) (string, error) {
	claims := Claims{
		GameId:    gameId,
		PlayerId:  playerId,
		ServerId:  serverId,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(time.Hour * 10),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	return token.SignedString(key)
}

func ParseToken(token string) (*Claims, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	if err != nil {
		return nil, err
	}

	err = claims.Valid()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Invalid token: %s", err))
	}
	return claims, nil
}

package security

import (
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type Security interface {
	ComparePassword(hash string, password string) error
	CreateToken(m map[string]any) (string, error)
	GeneratePassword(password string) ([]byte, error)
}

type security struct {
	tokenKey []byte
}

func NewSecurity(tokenKey string) Security {
	return &security{
		tokenKey: []byte(tokenKey),
	}
}

func (s *security) ComparePassword(hash string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func (s *security) CreateToken(m map[string]any) (string, error) {

	mapClaim := jwt.MapClaims(m)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaim)

	return token.SignedString(s.tokenKey)
}

func (s *security) GeneratePassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), 15)
}

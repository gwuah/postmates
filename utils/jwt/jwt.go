package jwt

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gwuah/api/database/models"
)

var minSecretLen = 128

// New generates new JWT service necessary for auth middleware
func New(algo, secret string, ttlMinutes, minSecretLength int) (Service, error) {
	if minSecretLength > 0 {
		minSecretLen = minSecretLength
	}
	if len(secret) < minSecretLen {
		return Service{}, fmt.Errorf("jwt secret length is %v, which is less than required %v", len(secret), minSecretLen)
	}
	signingMethod := jwt.GetSigningMethod(algo)
	if signingMethod == nil {
		return Service{}, fmt.Errorf("invalid jwt signing method: %s", algo)
	}

	return Service{
		key:  []byte(secret),
		algo: signingMethod,
		ttl:  time.Duration(ttlMinutes) * time.Minute,
	}, nil
}

// Service provides a Json-Web-Token authentication implementation
type Service struct {
	// Secret key used for signing.
	key []byte

	// Duration for which the jwt token is valid.
	ttl time.Duration

	// JWT signing algorithm
	algo jwt.SigningMethod
}

// ParseToken parses token from Authorization header
func (s Service) ParseToken(authHeader string) (*jwt.Token, error) {
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return nil, errors.New("token not passed to authorization header")
	}

	return jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
		if s.algo != token.Method {
			return nil, errors.New("failed to parse token")
		}
		return s.key, nil
	})

}

// GenerateToken generates new JWT token and populates it with user data
func (s Service) GenerateToken(customer *models.Customer) (string, error) {
	return jwt.NewWithClaims(s.algo, jwt.MapClaims{
		"phone": customer.Phone,
		"email": customer.Email,
		"exp":   time.Now().Add(s.ttl).Unix(),
	}).SignedString(s.key)

}

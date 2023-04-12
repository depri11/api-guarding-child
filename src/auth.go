package gc

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type claims struct {
	UserID string `json:"userId"`
	jwt.StandardClaims
}

func (g *GC) CreateJWT(userId string) (string, error) {
	claims := &claims{
		UserID: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
}

func (g *GC) ValidateJWT(token string) (*claims, error) {
	tokens, err := jwt.ParseWithClaims(token, &claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil {
		return nil, err
	}

	// claims, ok := tokens.Claims.(*claims)
	claims, ok := tokens.Claims.(*claims)
	if !ok {
		err = errors.New("couldn't parse claims")
		return nil, err
	}

	return claims, nil
}

func (g *GC) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (g *GC) CheckPassword(hashPass, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPass), []byte(password))
	return err == nil
}

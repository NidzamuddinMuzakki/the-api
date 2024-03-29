package configs

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

func GenerateTokenPair(name string) (map[string]string, error) {
	// Create token
	token := jwt.New(jwt.SigningMethodHS512)
	id := uuid.New()
	// Set claims
	// This is the information which frontend can use
	// The backend can also decode the token and get admin etc.
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = 1
	claims["name"] = name
	claims["uuid"] = id.String()
	claims["exp"] = time.Now().Add(time.Minute * 5).Unix()

	// Generate encoded token and send it as response.
	// The signing string should be secret (a generated UUID works too)
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return nil, err
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["sub"] = 1
	rtClaims["name"] = name
	rtClaims["uuid"] = id.String()
	rtClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	rt, err := refreshToken.SignedString([]byte("secret"))
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"token":        t,
		"refreshtoken": rt,
	}, nil
}

package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTData struct {
	Email string
}

type JWT struct {
	Secret string
}

func NewJWT(secret string) *JWT {

	return &JWT{
		Secret: secret,
	}
}

func (j *JWT) Create(email string) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   expirationTime,
	})

	s, err := token.SignedString([]byte(j.Secret))

	if err != nil {
		return "", err
	}

	return s, nil
}

func (j *JWT) Parse(token string) (bool, *JWTData) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.Secret), nil
	})

	if err != nil {
		return false, nil
	}

	email := t.Claims.(jwt.MapClaims)["email"]

	return t.Valid, &JWTData{
		Email: email.(string),
	}

}

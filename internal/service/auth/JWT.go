package auth

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"

	"crypto/rsa"
	"time"
)

type Jwt struct {
	publicKey *rsa.PublicKey
	privateKey *rsa.PrivateKey

	lifetime time.Duration
}

type Claims struct {
	Id string
	jwt.StandardClaims
}

func NewJwt(private []byte, public []byte, lifetime time.Duration) (*Jwt, error) {
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(private)
	if err != nil {
		return nil, err
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(public)
	if err != nil {
		return nil, err
	}

	return &Jwt{
		publicKey: publicKey,
		privateKey: privateKey,
		lifetime: lifetime,
	}, nil
}

func (j Jwt) IssueToken(userId string) (string, error) {
	claims := Claims{
		Id: userId,
		StandardClaims: jwt.StandardClaims {
			ExpiresAt: time.Now().Add(j.lifetime).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(j.privateKey)
}

func (j Jwt) UserIdByToken(rawToken string) (string, error) {
	token, err := jwt.ParseWithClaims(rawToken, &Claims{}, func (token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Bad token signature")
		}
		return j.publicKey, nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return "", errors.New("Invalid token claims")
	}
	return claims.Id, nil
}

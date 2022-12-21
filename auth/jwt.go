package auth

import (
	"errors"
	"hypegym-backend/models/enums"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

type JWTClaim struct {
	Email string     `json:"email"`
	Role  enums.Role `json:"role"`
	GymID uint       `json:"gym_id"`
	jwt.StandardClaims
}

func GenerateJWT(email string, role enums.Role, gym_id uint) (tokenString string, err error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &JWTClaim{
		Email: email,
		Role:  role,
		GymID: gym_id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(jwtKey)
	return
}

func ValidateToken(signedToken string) (err error, jwtClaim *JWTClaim) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)

	if err != nil {
		return
	}

	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		err = errors.New("couldn't parse claims")
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("token expired")
		return
	}

	jwtClaim = claims
	return
}

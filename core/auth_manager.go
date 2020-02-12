package core

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	SigningKeyBytesAmount   = 4096
	TokenExpireTimeDuration = 24 * time.Hour
)

type AuthManager struct {
	signingKey string // For now it is random bytes string represented in Base64
}

type JWTCustomClaims struct {
	EntityID string `json:"entityID"`
	DeviceID string `json:"deviceID"`
	jwt.StandardClaims
}

func NewAuthManager() (*AuthManager, error) {
	am := &AuthManager{}
	bytes, err := GenRandomBytes(SigningKeyBytesAmount)
	if err != nil {
		return nil, err
	}
	am.signingKey = base64.RawStdEncoding.EncodeToString(bytes)
	return am, nil
}

func (am *AuthManager) CreateNewToken(entityID, deviceID string, tokenExpireTimeDuration time.Duration) (string, error) {
	timeNow := time.Now()
	expiringTime := timeNow.Add(tokenExpireTimeDuration)
	claims := JWTCustomClaims{
		entityID,
		deviceID,
		jwt.StandardClaims{
			ExpiresAt: time.Date(
				expiringTime.Year(),
				expiringTime.Month(),
				expiringTime.Day(),
				expiringTime.Hour(),
				expiringTime.Minute(),
				expiringTime.Second(),
				expiringTime.Nanosecond(),
				time.UTC,
			).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(am.signingKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (am *AuthManager) ValidateToken(tokenString string) (bool, string, string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(am.signingKey), nil
	})

	if claims, ok := token.Claims.(*JWTCustomClaims); ok && token.Valid {
		return true, claims.EntityID, claims.DeviceID, nil
	}
	return false, "", "", err
}

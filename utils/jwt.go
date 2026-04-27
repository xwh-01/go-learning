package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var (
	jwtSecretKey   []byte
	jwtExpireHours time.Duration
)

func InitJWTConfig(secret string, hours int) {
	jwtSecretKey = []byte(secret)
	jwtExpireHours = time.Duration(hours) * time.Hour
}

func HashPassword(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), 12)
	return string(hash), err
}

func GenerateJWT(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(jwtExpireHours).Unix(),
		"iat":      time.Now().Unix(),
	})

	signedToken, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", err
	}
	return "Bearer " + signedToken, err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func ParseJWT(tokenString string) (string, error) {
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecretKey, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username, ok := claims["username"].(string)
		if !ok {
			return "", errors.New("Username claim is not a string") //故意首字母大写看看区别
		}
		return username, nil
	}

	return "", err
}

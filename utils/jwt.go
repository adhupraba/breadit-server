package utils

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/adhupraba/breadit-server/internal/database"
	"github.com/adhupraba/breadit-server/lib"
)

func SignJwtToken(payload string, ttl int64) (signedToken string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": payload,
		"exp": ttl,
	})

	tokenStr, err := token.SignedString([]byte(lib.EnvConfig.JwtSecret))

	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func VerifyJwtToken(tokenStr string) (subject string, err error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header)
		}

		return []byte(lib.EnvConfig.JwtSecret), nil
	})

	if err != nil {
		return "", fmt.Errorf("Unable to parse token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid {
		return "", fmt.Errorf("Invalid auth token")
	}

	exp, err := claims.GetExpirationTime()

	if err != nil || time.Now().Unix() > exp.Unix() {
		return "", fmt.Errorf("Expired auth token")
	}

	sub, err := claims.GetSubject()

	if err != nil {
		return "", fmt.Errorf("Subject not found in token")
	}

	return sub, nil
}

func GetUserFromToken(w http.ResponseWriter, r *http.Request, tokenStr string) (database.User, error) {
	sub, err := VerifyJwtToken(tokenStr)

	if err != nil {
		return database.User{}, fmt.Errorf("Invalid token")
	}

	userId, err := strconv.Atoi(sub)

	if err != nil {
		return database.User{}, fmt.Errorf("Invalid subject")
	}

	user, err := lib.DB.FindUserById(r.Context(), int32(userId))

	if err != nil {
		return database.User{}, fmt.Errorf("The user belonging to this token no logger exists")
	}

	return user, nil
}

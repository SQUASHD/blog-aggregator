package auth

import (
	"errors"
	"net/http"
	"strings"
)

// ErrNoAuthHeaderIncluded -
var ErrNoAuthHeaderIncluded = errors.New("no auth header included in request")

func GetApiKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuthHeaderIncluded
	}
	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || strings.ToLower(splitAuth[0]) != "apikey" {
		return "", errors.New("malformed authorization header")
	}
	return splitAuth[1], nil
}

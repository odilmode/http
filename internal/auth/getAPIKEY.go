package auth

import (
	"net/http"
	"strings"
	"fmt"
)

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("Authorization header missing")
	}

	if strings.HasPrefix(authHeader, "ApiKey ") {
		key := strings.TrimPrefix(authHeader, "ApiKey ")
		return key, nil
	} else {
		return "", fmt.Errorf("Unsupported Authorization scheme")
	}
}

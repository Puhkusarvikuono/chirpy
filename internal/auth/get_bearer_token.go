package auth

import (
	"strings"
	"net/http"
	"fmt"
)

func GetBearerToken(headers http.Header) (string, error) {
	v, ok := headers["Authorization"]
	if !ok {
		return "", fmt.Errorf("Authorization header not found")
	}

  for _, value := range v {
		value = strings.TrimSpace(value)
		if strings.HasPrefix(strings.ToLower(value), "bearer") {
			bearerToken := strings.TrimSpace(value[6:])
			if bearerToken == "" {
				return "", fmt.Errorf("No token string")
			}
			return bearerToken, nil
		}
	}

	return "", fmt.Errorf("No valid bearer token")
}

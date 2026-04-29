package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	v, ok := headers["Authorization"]
	if !ok {
		return "", fmt.Errorf("authorization header not found")
	}

	for _, value := range v {
		value = strings.TrimSpace(value)
		if strings.HasPrefix(strings.ToLower(value), "apikey") {
			APIKey := strings.TrimSpace(value[6:])
			if APIKey == "" {
				return "", fmt.Errorf("no APIKey string")
			}
			return APIKey, nil
		}
	}

	return "", fmt.Errorf("no valid APIKey found")
}

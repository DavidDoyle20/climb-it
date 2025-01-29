package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetAuthorizationHeader(headerName string, headers http.Header) (string, error) {
	vals := headers.Values("Authorization")
	tokenString := ""
	for _, v := range vals {
		if strings.Contains(v, headerName) {
			tokenString = strings.TrimSpace(strings.Replace(v, headerName, "", 1))
		}
	}
	if tokenString == "" {
		return "", fmt.Errorf("%s header doesnt exist", headerName)
	}
	return tokenString, nil
}

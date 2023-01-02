package validation

import "strings"

//trimBearer trim extra objects in Authorization:
func TrimBearer(tokenString string) string {
	tokenString = strings.TrimSpace(tokenString)
	if len(tokenString) < 6 {
		return tokenString
	}

	if strings.ToLower(tokenString[:6]) == "bearer" {
		tokenString = tokenString[6:]
		tokenString = strings.TrimSpace(tokenString)
	}
	return tokenString
}

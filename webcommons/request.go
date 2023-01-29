package webcommons

import (
	"net/http"
	"strings"
)

// GetBearerToken retrieves the bearer token from the Authorization header of an http request.
func GetBearerToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	authz := strings.Split(authHeader, " ")
	if len(authz) != 2 || authz[0] != "Bearer" {
		return ""
	}
	return authz[1]
}

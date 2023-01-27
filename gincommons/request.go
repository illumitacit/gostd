package gincommons

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// GetBearerToken retrieves the bearer token from the Authorization header of a gin context.
func GetBearerToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}

	authz := strings.Split(authHeader, " ")
	if len(authz) != 2 || authz[0] != "Bearer" {
		return ""
	}
	return authz[1]
}

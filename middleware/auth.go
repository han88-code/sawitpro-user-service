package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/SawitProRecruitment/UserService/util"

	"github.com/labstack/echo/v4"
)

// AuthenticationMiddleware checks if the user has a valid JWT token
func JWTAuthentication(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		headers := req.Header

		tokenString := headers.Get("Authorization")
		if tokenString == "" {
			return c.JSON(http.StatusUnauthorized, fmt.Errorf("missing authentication token"))
		}

		// The token should be prefixed with "Bearer "
		tokenParts := strings.Split(tokenString, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return c.JSON(http.StatusUnauthorized, fmt.Errorf("invalid authentication token"))
		}

		tokenString = tokenParts[1]

		claims, err := util.VerifyRSAToken(tokenString)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, fmt.Errorf("invalid authentication token"))
		}

		c.Set("user_id", claims["user_id"])
		return next(c)
	}
}

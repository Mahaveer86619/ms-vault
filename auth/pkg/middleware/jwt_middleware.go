package middleware

import (
	"net/http"
	"strings"

	"github.com/Mahaveer86619/ms/auth/pkg/utils"
	"github.com/Mahaveer86619/ms/auth/pkg/views"
	"github.com/labstack/echo/v4"
)

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, views.Failure{
				StatusCode: http.StatusUnauthorized,
				Message:    "Missing Authorization header",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.JSON(http.StatusUnauthorized, views.Failure{
				StatusCode: http.StatusUnauthorized,
				Message:    "Invalid Authorization header format",
			})
		}

		claims, err := utils.ValidateToken(parts[1], "access")
		if err != nil {
			return c.JSON(http.StatusUnauthorized, views.Failure{
				StatusCode: http.StatusUnauthorized,
				Message:    "Invalid or expired token",
			})
		}

		c.Set("user_id", claims.UserID)
		return next(c)
	}
}

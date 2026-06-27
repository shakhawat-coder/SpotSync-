package middleware

import (
	"fmt"
	"strings"

	"spotsync/errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func JWTMiddleware(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			auth := c.Request().Header.Get("Authorization")
			if auth == "" {
				return c.JSON(401, map[string]any{
					"success": false,
					"message": errors.ErrUnauthorized.Message,
					"errors":  "Missing Authorization header",
				})
			}

			const prefix = "Bearer "
			if !strings.HasPrefix(auth, prefix) {
				return c.JSON(401, map[string]any{
					"success": false,
					"message": errors.ErrUnauthorized.Message,
					"errors":  "Invalid token format",
				})
			}

			token := strings.TrimSpace(auth[len(prefix):])
			if token == "" {
				return c.JSON(401, map[string]any{
					"success": false,
					"message": errors.ErrUnauthorized.Message,
					"errors":  "Invalid token format",
				})
			}

			claims := &JWTClaims{}
			parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(secret), nil
			})

			if err != nil || !parsedToken.Valid {
				return c.JSON(401, map[string]any{
					"success": false,
					"message": errors.ErrUnauthorized.Message,
					"errors":  "Invalid or expired token",
				})
			}

			c.Set("user_id", claims.UserID)
			c.Set("email", claims.Email)
			c.Set("role", claims.Role)

			return next(c)
		}
	}
}

func GetUserIDFromContext(c *echo.Context) uint {
	if uid, ok := c.Get("user_id").(uint); ok {
		return uid
	}
	return 0
}

func GetRoleFromContext(c *echo.Context) string {
	if role, ok := c.Get("role").(string); ok {
		return role
	}
	return ""
}

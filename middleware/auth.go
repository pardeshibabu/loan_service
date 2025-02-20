package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"loan_service/models"

	"github.com/gobuffalo/buffalo"
	"github.com/golang-jwt/jwt"
)

// AuthMiddleware handles JWT authentication
func AuthMiddleware(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		tokenString := extractToken(c.Request())
		if tokenString == "" {
			return c.Error(http.StatusUnauthorized, fmt.Errorf("no authorization token provided"))
		}

		claims, err := validateToken(tokenString)
		if err != nil {
			return c.Error(http.StatusUnauthorized, err)
		}

		// Set employee info in context
		c.Set("employee_id", claims.EmployeeID)
		c.Set("role", claims.Role)

		return next(c)
	}
}

// RoleMiddleware checks if the user has required role
func RoleMiddleware(roles ...models.EmployeeRole) buffalo.MiddlewareFunc {
	return func(next buffalo.Handler) buffalo.Handler {
		return func(c buffalo.Context) error {
			userRole := c.Value("role").(models.EmployeeRole)

			for _, role := range roles {
				if userRole == role {
					return next(c)
				}
			}

			return c.Error(http.StatusForbidden, fmt.Errorf("unauthorized access"))
		}
	}
}

// Claims represents JWT claims
type Claims struct {
	EmployeeID string              `json:"employee_id"`
	Role       models.EmployeeRole `json:"role"`
	jwt.StandardClaims
}

func extractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func validateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("your-secret-key"), nil // Use environment variable in production
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

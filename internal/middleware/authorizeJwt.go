package middleware

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"sandbox.go/jwt/internal/base"
	"sandbox.go/jwt/internal/handler"
)

func validateToken(secret string, encodedToken string) (*jwt.Token, error) {
	return jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if _, isvalid := token.Method.(*jwt.SigningMethodHMAC); !isvalid {
			return nil, fmt.Errorf("Invalid token", token.Header["alg"])
		}
		return []byte(secret), nil
	})

}

var AuthorizeJWT gin.HandlerFunc = func(c *gin.Context) {
	const BEARER_SCHEMA = "Bearer "
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
	} else {
		tokenString := authHeader[len(BEARER_SCHEMA):]
		token, err := validateToken(base.JwtConfig.SecretKey, tokenString)
		if token.Valid {
			claims := token.Claims.(jwt.MapClaims)
			userRole := handler.UserRole(claims["userRole"].(float64))
			c.Set("userRole", userRole)
			c.Set("userName", claims["Name"])
			fmt.Println(claims)
		} else {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}

}

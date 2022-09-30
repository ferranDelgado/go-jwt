package handler

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"sandbox.go/jwt/internal/base"
	"time"
)

type UserRole int64

func (u UserRole) String() string {
	switch u {
	case Root:
		return "Root"
	case User:
		return "User"
	}
	return "Undefined"
}

const (
	Undefined UserRole = iota
	Root
	User
)

type authApiClaims struct {
	Id       string
	Name     string
	UserRole UserRole `json:"userRole"`
	jwt.StandardClaims
}

func generateApiToken(id string, name string, role UserRole) (string, error) {
	expirationDate := time.Now().AddDate(0, 0, base.JwtConfig.Days).Unix()
	claims := &authApiClaims{
		id,
		name,
		role,
		jwt.StandardClaims{
			ExpiresAt: expirationDate,
			Issuer:    base.JwtConfig.Issuer,
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(base.JwtConfig.SecretKey))
	if err != nil {
		return "", err
	}
	return t, nil
}

var CreateTokenHandler gin.HandlerFunc = func(ctx *gin.Context) {
	id := uuid.New().String()
	token, err := generateApiToken(id, "name", User)
	if err != nil {
		log.Printf("Error: %s", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal error",
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"token": token,
		})
	}
}

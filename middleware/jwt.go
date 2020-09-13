package middleware

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type TokenParser interface {
	ParseToken(string) (*jwt.Token, error)
}

func JWT(tokenParser TokenParser) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := tokenParser.ParseToken(c.Request.Header.Get("Authorization"))
		if err != nil || !token.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		phone := claims["phone"].(string)
		email := claims["email"].(string)

		c.Set("phone", phone)
		c.Set("email", email)

		c.Next()
	}
}

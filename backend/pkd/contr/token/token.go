package token

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func CheckToken(c *gin.Context) {
	log.Printf("Check token: %v\n", c.Request.Header.Get("Authorization"))

	tokenStr := c.Request.Header.Get("Authorization")
	if tokenStr == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		jwtTokenSecret := os.Getenv("JWT_TOKEN_SECRET")
		return []byte(jwtTokenSecret), nil
	})

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		c.Set("user", claims["sub"])
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	c.Next()
}

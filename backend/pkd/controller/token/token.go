package token

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func CheckToken(c *gin.Context) {
	log.Printf("Check token: %v\n", c.Request.Header.Get("Authorization"))

	tokenStr := c.Request.Header.Get("Authorization")
	tokenStr = strings.Split(tokenStr, " ")[1]
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
		log.Printf("Token error: %v\n", err.Error())
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if len(claims["uuid"].(string)) < 10 || len(claims["auth"].(string)) <= 3 || float64(time.Now().Unix()) > claims["exp"].(float64) || float64(time.Now().Unix()) < claims["lastmsg"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		c.Set("user", claims["sub"])
		c.Set("roles", claims["auth"])
	} else {
		log.Printf("Token Claim check failed.\n")
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	c.Next()
}

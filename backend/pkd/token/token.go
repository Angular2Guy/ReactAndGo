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
	"github.com/google/uuid"
)

type Role string

const (
	UserRole  Role = "USER"
	GuestRole Role = "GUEST"
)

const (
	HeaderAuth     = "Authorization"
	HeaderBearer   = "Bearer"
	JwtTokenSecret = "JWT_TOKEN_SECRET"
	TokenAuth      = "auth"
	TokenSub       = "sub"
	TokenUuid      = "uuid"
	TokenExp       = "exp"
	TokenLastMsg   = "lastmsg"
)

type TokenUser struct {
	Username string
	Roles    []string
}

func CreateToken(tokenUser TokenUser) (string, error) {
	roles := []string{string(GuestRole)}
	var myUuid uuid.UUID
	if myUuid1, err := uuid.NewRandom(); err != nil {
		log.Printf("Uuid creation failed: %v\n", err.Error())
		return "", err
	} else {
		myUuid = myUuid1
		roles = append(roles, string(UserRole))
	}
	tokenTtl := 60
	if len(strings.TrimSpace(os.Getenv("MSG_MESSAGES"))) >= 3 {
		tokenTtl = tokenTtl * 10
	}
	myToken := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		TokenSub:     tokenUser.Username,
		TokenUuid:    myUuid.String(),
		TokenAuth:    strings.Join(roles[:], ","),
		TokenLastMsg: time.Now().Unix(),
		TokenExp:     time.Now().Add(time.Second * time.Duration(tokenTtl)).Unix(),
	})
	result, err := myToken.SignedString([]byte(os.Getenv(JwtTokenSecret)))
	return result, err
}

func CheckToken(c *gin.Context) {
	log.Printf("Check token: %v\n", c.Request.Header.Get(HeaderAuth))

	tokenStr := c.Request.Header.Get(HeaderAuth)
	tokenStr = strings.Split(tokenStr, " ")[1]
	if tokenStr == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv(JwtTokenSecret)), nil
	})

	if err != nil {
		log.Printf("Token error: %v\n", err.Error())
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if len(claims[TokenUuid].(string)) < 10 || len(claims[TokenAuth].(string)) <= 3 || float64(time.Now().Unix()) > claims[TokenExp].(float64) || float64(time.Now().Unix()) < claims[TokenLastMsg].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		c.Set("user", claims[TokenSub])
		c.Set("roles", claims[TokenAuth])
	} else {
		log.Printf("Token Claim check failed.\n")
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	c.Next()
}

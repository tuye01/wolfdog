package verify

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"wolfdog/internal/consts"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var secretKey = []byte("wW1.eS8[iW9*lE2_pD5:iQ4:wD8>kT3?bD4`")

type Claims struct {
	UserID   int64  `json:"user_id"`
	UserName string `json:"user_name"`
	Mobile   string `json:"mobile"`
	jwt.StandardClaims
}

func GenerateToken(c *gin.Context, userID int64, username string, mobile string) error {
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		UserID:   userID,
		UserName: username,
		Mobile:   mobile,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Token 过期时间为 24 小时
		},
	})
	token, err := tokenClaims.SignedString([]byte(secretKey))
	if err != nil {
		return err
	}

	_, err = consts.RedisDB.Do(c, "SET", userID, token, "EX", time.Now().Add(time.Hour*24).Unix()).Result()
	if err != nil {
		return err
	}
	c.Header("Authorization", token)
	return nil
}

func ParseToken(c *gin.Context) (*Claims, error) {
	tokenSting := strings.Replace(c.GetHeader("Authorization"), "Bearer ", "", 1)
	if len(tokenSting) == 0 {
		return nil, fmt.Errorf("invalid token")
	}
	token, _ := jwt.ParseWithClaims(tokenSting, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if token.Claims == nil {
		return nil, fmt.Errorf("invalid token")
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		claim, err := ParseToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		val, err := consts.RedisDB.Exists(c, strconv.FormatInt(claim.UserID, 10)).Result()
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		if val != 1 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		c.Set("user_id", claim.UserID)
		c.Next()
	}
}

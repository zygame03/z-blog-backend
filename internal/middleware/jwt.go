package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("your_secret_key")

func SetJWTKey(key string) {
	if key == "" {
		return
	}
	jwtKey = []byte(key)
}

// 生成token
func GenerateToken(userID int, username string, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(ttl).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从头部取出token
		// Authorization: Bearer yourtoken...
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing token",
			})
			return
		}
		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))

		// jwtkey解token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		// 从存的数中拿信息
		claims, ok := token.Claims.(jwt.MapClaims)
		if ok {
			id, ok := claims["user_id"].(float64)
			if ok {
				c.Set("userID", strconv.Itoa(int(id)))
			}
			username, ok := claims["username"].(string)
			if ok {
				c.Set("username", username)
			}
		}
		c.Next()
	}
}

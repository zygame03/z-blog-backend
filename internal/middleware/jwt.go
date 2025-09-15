package middleware

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey []byte

func SetJWTKey(key string) {
	if key == "" {
		key = os.Getenv("JWT_KEY")
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
		// 从 header 中拿 token
		// Authorization: Bearer yourtoken
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing token",
			})
			return
		}
		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		// 先检查 token 是否过期
		// 再将 user_id 和 username 写入上下文
		claims, ok := token.Claims.(jwt.MapClaims)
		if ok {
			if time.Now().Unix() > int64(claims["exp"].(float64)) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
				return
			}

			id, ok := claims["user_id"].(int)
			if ok {
				c.Set("userID", strconv.Itoa(id))
			}
			username, ok := claims["username"].(string)
			if ok {
				c.Set("username", username)
			}
		}
		c.Next()
	}
}

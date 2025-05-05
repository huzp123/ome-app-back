package middleware

import (
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"ome-app-back/pkg/errcode"
)

var (
	// JWT密钥，实际项目中应从配置读取
	jwtSecret = []byte("your-jwt-secret-key")
)

// JWT认证中间件
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string
		authorization := c.GetHeader("Authorization")
		if authorization != "" {
			parts := strings.SplitN(authorization, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" {
				token = parts[1]
			}
		}

		if token == "" {
			errcode.UnauthorizedTokenError.Response(c)
			c.Abort()
			return
		}

		// 解析token
		claims := &Claims{}
		tokenClaims, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !tokenClaims.Valid {
			errcode.UnauthorizedTokenError.Response(c)
			c.Abort()
			return
		}

		// token过期检查
		if time.Now().Unix() > claims.ExpiresAt {
			errcode.UnauthorizedTokenTimeout.Response(c)
			c.Abort()
			return
		}

		// 将用户ID存入上下文
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}

// Claims 自定义JWT Claims
type Claims struct {
	UserID int64 `json:"user_id"`
	jwt.StandardClaims
}

// GenerateToken 生成JWT令牌
func GenerateToken(userID int64) (string, error) {
	now := time.Now()
	expireTime := now.Add(24 * time.Hour) // 过期时间设为24小时

	claims := Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			IssuedAt:  now.Unix(),
			Issuer:    "ome-app",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

package middleware

import (
	"beep/internal/config"
	"beep/internal/errors"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

func Auth(config *config.Config, redis *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("X-Token")
		if token == "" {
			panic(errors.NewUnauthorizedError("未登录的请求", nil))
		}
		claims := new(jwt.RegisteredClaims)
		tok, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Jwt.Secret), nil
		})
		if err != nil || !tok.Valid {
			panic(errors.NewUnauthorizedError("token无效", nil))
		}
		result, err := redis.Get(c, claims.ID).Result()
		if err != nil || result == "" {
			panic(errors.NewUnauthorizedError("token无效", err))
		}

		c.Next()
	}
}

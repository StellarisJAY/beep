package middleware

import (
	"beep/internal/config"
	"beep/internal/errors"
	"beep/internal/types"
	"context"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func Auth(config *config.Config, redis *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := c.GetHeader("access_token")
		if accessToken == "" {
			panic(errors.NewUnauthorizedError("未登录的请求", nil))
		}
		data, err := redis.Get(c.Request.Context(), "access_token:"+accessToken).Bytes()
		if err != nil || data == nil {
			panic(errors.NewUnauthorizedError("未登录的请求", err))
		}
		var loginInfo types.LoginInfo
		if err := json.Unmarshal(data, &loginInfo); err != nil {
			panic(errors.NewInternalServerError("登录信息解析失败", err))
		}
		// 登录信息缓存到context中
		c.Set(types.AccessTokenContextKey, accessToken)
		c.Set(types.UserIdContextKey, loginInfo.UserId)
		c.Set(types.WorkspaceIdContextKey, loginInfo.WorkspaceId)
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), types.AccessTokenContextKey, accessToken))
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), types.UserIdContextKey, loginInfo.UserId))
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), types.WorkspaceIdContextKey, loginInfo.WorkspaceId))
		c.Next()
	}
}

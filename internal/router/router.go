package router

import (
	"beep/internal/config"
	"beep/internal/middleware"
	"io"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/dig"
)

type Params struct {
	dig.In

	Redis  *redis.Client
	Config *config.Config
	//TODO handlers
}

func InitRouter(params Params) (*gin.Engine, error) {
	mode := os.Getenv(gin.EnvGinMode)
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	} else if mode == gin.DebugMode {
		gin.SetMode(gin.DebugMode)
	} else if mode == gin.TestMode {
		gin.SetMode(gin.TestMode)
	}

	r := gin.New()
	r.ContextWithFallback = true
	r.Use(gin.CustomRecoveryWithWriter(io.Discard, middleware.Recovery())) // 自定义恢复中间件
	r.Use(middleware.CORS())                                               // 跨域中间件

	v1 := r.Group("/api/v1")
	v1.Use(middleware.Auth(params.Config, params.Redis))

	return r, nil
}

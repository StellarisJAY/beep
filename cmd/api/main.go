package main

import (
	"beep/internal/config"
	"beep/internal/container"
	"beep/internal/migration"
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	args := os.Args[1:]
	if err := run(args); err != nil {
		log.Fatal(err)
	}
}

func run(args []string) error {
	c := container.NewContainer()
	// 迁移数据库
	if err := c.Invoke(migration.MigrateDatabase); err != nil {
		return err
	}
	// 启动服务，注入gin router和配置
	err := c.Invoke(func(e *gin.Engine, config *config.Config) error {
		// HTTP服务器
		server := &http.Server{
			Addr:    fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port),
			Handler: e,
		}
		// 监听操作系统信号
		ctx, cancel := context.WithCancel(context.Background())
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		// gracefully shutdown
		go func() {
			sig := <-signals
			slog.Info("API Server shutting down", "signal", sig.String())

			// 服务关闭的30s超时时间
			shutCtx, shutCancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer shutCancel()
			// 关闭HTTP服务
			if err := server.Shutdown(shutCtx); err != nil {
				log.Fatalf("Server forced to shutdown: %v", err)
			}
			slog.Info("API Server gracefully stopped")
			cancel()
		}()

		// 启动HTTP服务器
		slog.Info("API Server started", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("API Server startup failed: %w", err)
		}
		// 等待gracefully shutdown
		<-ctx.Done()
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

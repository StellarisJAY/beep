package router

import (
	"beep/internal/config"
	"beep/internal/handler"
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
	UserHandler          *handler.UserHandler
	WorkspaceHandler     *handler.WorkspaceHandler
	KnowledgeBaseHandler *handler.KnowledgeBaseHandler
	ModelHandler         *handler.ModelHandler
	MCPServerHandler     *handler.MCPServerHandler
	DocumentHandler      *handler.DocumentHandler
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
	{
		v1.POST("/register", params.UserHandler.Register)
		v1.POST("/login", params.UserHandler.Login)
	}
	// Workspace
	w := v1.Group("/workspace")
	{
		w.Use(middleware.Auth(params.Config, params.Redis))
		w.GET("/members", params.WorkspaceHandler.ListMember)
		w.GET("/list", params.WorkspaceHandler.List)
		w.POST("/invite", params.WorkspaceHandler.InviteMember)
		w.POST("/switch/:id", params.WorkspaceHandler.SwitchWorkspace)
		w.POST("/role", params.WorkspaceHandler.SetRole)
	}
	// 知识库
	k := v1.Group("/kb")
	{
		k.Use(middleware.Auth(params.Config, params.Redis))
		k.GET("/list", params.KnowledgeBaseHandler.List)
		k.POST("/create", params.KnowledgeBaseHandler.Create)
		k.PUT("/update", params.KnowledgeBaseHandler.Update)
		k.DELETE("/delete/:id", params.KnowledgeBaseHandler.Delete)

	}
	// 模型
	m := v1.Group("/model")
	{
		m.Use(middleware.Auth(params.Config, params.Redis))
		m.GET("/factory/list", params.ModelHandler.ListModelFactory)
		m.GET("/list", params.ModelHandler.ListModel)
		m.POST("/factory/create", params.ModelHandler.CreateModelFactory)
		m.PUT("/factory/update", params.ModelHandler.UpdateFactory)
	}

	// MCP服务器
	mcp := v1.Group("/mcp")
	{
		mcp.Use(middleware.Auth(params.Config, params.Redis))
		mcp.POST("/create", params.MCPServerHandler.Create)
		mcp.GET("/list", params.MCPServerHandler.List)
		mcp.PUT("/update", params.MCPServerHandler.Update)
		mcp.DELETE("/delete/:id", params.MCPServerHandler.Delete)
	}

	// 文档
	doc := v1.Group("/doc")
	{
		doc.Use(middleware.Auth(params.Config, params.Redis))
		doc.POST("/create", params.DocumentHandler.CreateFromFile)
		doc.GET("/list", params.DocumentHandler.List)
		doc.DELETE("/delete/:id", params.DocumentHandler.Delete)
		doc.GET("/download/:id", params.DocumentHandler.Download)
	}
	return r, nil
}

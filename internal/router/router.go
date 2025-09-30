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

	Redis  *redis.Client  // 注入Redis客户端，用于auth中间件
	Config *config.Config // 注入配置

	UserHandler          *handler.UserHandler          // 用户
	WorkspaceHandler     *handler.WorkspaceHandler     // 工作空间
	KnowledgeBaseHandler *handler.KnowledgeBaseHandler // 知识库
	ModelHandler         *handler.ModelHandler         // 模型
	MCPServerHandler     *handler.MCPServerHandler     // MCP服务器
	DocumentHandler      *handler.DocumentHandler      // 文档
	AgentHandler         *handler.AgentHandler         // 智能体
	ChatHandler          *handler.ChatHandler          // 聊天
}

// InitRouter 初始化路由
func InitRouter(params Params) (*gin.Engine, error) {
	// 设置Gin模式
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
	// 基础信息
	initBasicRouter(v1, params)
	// 智能体
	initAgentRouter(v1, params)
	// 知识库，文档
	initKnowledgeRouter(v1, params)
	// 模型，mcp
	initModelRouter(v1, params)
	return r, nil
}

// initBasicRouter 初始化基础路由
func initBasicRouter(r *gin.RouterGroup, params Params) {
	{
		r.POST("/register", params.UserHandler.Register)
		r.POST("/login", params.UserHandler.Login)
	}
	// Workspace
	w := r.Group("/workspace")
	{
		w.Use(middleware.Auth(params.Config, params.Redis))
		w.GET("/members", params.WorkspaceHandler.ListMember)
		w.GET("/list", params.WorkspaceHandler.List)
		w.POST("/invite", params.WorkspaceHandler.InviteMember)
		w.POST("/switch/:id", params.WorkspaceHandler.SwitchWorkspace)
		w.POST("/role", params.WorkspaceHandler.SetRole)
	}
}

// initKnowledgeRouter 初始化知识库路由
func initKnowledgeRouter(r *gin.RouterGroup, params Params) {
	// 知识库
	k := r.Group("/kb")
	{
		k.Use(middleware.Auth(params.Config, params.Redis))
		k.GET("/list", params.KnowledgeBaseHandler.List)
		k.POST("/create", params.KnowledgeBaseHandler.Create)
		k.PUT("/update", params.KnowledgeBaseHandler.Update)
		k.DELETE("/delete/:id", params.KnowledgeBaseHandler.Delete)

	}

	// 文档
	doc := r.Group("/doc")
	{
		doc.Use(middleware.Auth(params.Config, params.Redis))
		doc.POST("/create", params.DocumentHandler.CreateFromFile)
		doc.GET("/list", params.DocumentHandler.List)
		doc.DELETE("/delete/:id", params.DocumentHandler.Delete)
		doc.GET("/download/:id", params.DocumentHandler.Download)
		doc.POST("/parse/:id", params.DocumentHandler.Parse)
	}
}

func initModelRouter(r *gin.RouterGroup, params Params) {
	// 模型
	m := r.Group("/model")
	{
		m.Use(middleware.Auth(params.Config, params.Redis))
		m.GET("/factory/list", params.ModelHandler.ListModelFactory)
		m.GET("/list", params.ModelHandler.ListModel)
		m.POST("/factory/create", params.ModelHandler.CreateModelFactory)
		m.PUT("/factory/update", params.ModelHandler.UpdateFactory)
	}

	// MCP服务器
	mcp := r.Group("/mcp")
	{
		mcp.Use(middleware.Auth(params.Config, params.Redis))
		mcp.POST("/create", params.MCPServerHandler.Create)
		mcp.GET("/list", params.MCPServerHandler.List)
		mcp.PUT("/update", params.MCPServerHandler.Update)
		mcp.DELETE("/delete/:id", params.MCPServerHandler.Delete)
	}
}

// initAgentRouter 初始化智能体路由
func initAgentRouter(r *gin.RouterGroup, params Params) {
	// 智能体
	agent := r.Group("/agent")
	{
		agent.Use(middleware.Auth(params.Config, params.Redis))
		agent.POST("/create", params.AgentHandler.Create)
		agent.GET("/list", params.AgentHandler.List)
	}
	// 智能体对话
	chat := r.Group("/chat")
	{
		chat.Use(middleware.Auth(params.Config, params.Redis))
		chat.POST("/send", params.ChatHandler.SendMessage)
	}
}

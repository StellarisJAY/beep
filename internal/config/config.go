package config

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

// Config 系统配置
type Config struct {
	DB struct {
		Driver      string `yaml:"driver"`       // 数据库驱动
		AutoMigrate bool   `yaml:"auto_migrate"` // 是否自动迁移数据库
	} `yaml:"db"` // 数据库配置
	Postgres struct {
		DSN string `yaml:"dsn"` // Postgres连接字符串
	} `yaml:"postgres"` // Postgres配置
	MySQL struct {
		DSN string `yaml:"dsn"` // MySQL连接字符串
	} `yaml:"mysql"` // MySQL配置
	Redis struct {
		Host     string `yaml:"host"`     // Redis主机
		Port     uint   `yaml:"port"`     // Redis端口
		DB       int    `yaml:"db"`       // Redis数据库索引
		Password string `yaml:"password"` // Redis密码
	} `yaml:"redis"` // Redis配置
	Server struct {
		Host string `yaml:"host"` // 服务器主机
		Port string `yaml:"port"` // 服务器端口
	} `yaml:"server"` // 服务器配置
	Logger struct {
		Format string `yaml:"format"` // 日志格式
		Level  string `yaml:"level"`  // 日志级别
		Path   string `yaml:"path"`   // 日志文件路径
	} `yaml:"logger"` // 日志配置
	Encrypt struct {
		Secret string `yaml:"secret"` // 加密密钥
	} `yaml:"encrypt"` // 加密配置
	FileStore string `yaml:"fileStore"` // 文件存储路径
	Minio     struct {
		Endpoint  string `yaml:"endpoint"`  // Minio端点
		AccessKey string `yaml:"accessKey"` // Minio访问密钥
		SecretKey string `yaml:"secretKey"` // Minio秘密密钥
		Bucket    string `yaml:"bucket"`    // Minio存储桶名称
	} `yaml:"minio"` // Minio配置
	Worker string `yaml:"worker"` // 工作池名称
	GoPool struct {
		Size    int           `yaml:"size"`    // 工作池大小
		Timeout time.Duration `yaml:"timeout"` // 工作池超时时间
	} `yaml:"gopool"` // 工作池配置
	KnowledgeGraph struct {
		EntityPrompt   string `yaml:"entityPrompt"`   // 实体抽取提示词
		RelationPrompt string `yaml:"relationPrompt"` // 关系抽取提示词
	} `yaml:"knowledge_graph"` // 知识图谱配置
	VectorStore string `yaml:"vectorStore"` // 向量存储路径
	Milvus      struct {
		Address  string `yaml:"address"`  // Milvus地址
		Username string `yaml:"username"` // Milvus用户名
		Password string `yaml:"password"` // Milvus密码
		Database string `yaml:"database"` // Milvus数据库名称
	} `yaml:"milvus"` // Milvus配置
	ConversationTitlePrompt string `yaml:"conversationTitlePrompt"` // 会话标题提示词
}

// LoadConfig 加载配置文件
func LoadConfig() (*Config, error) {
	// 根据GIN_MODE环境变量加载不同的配置文件
	mode := os.Getenv(gin.EnvGinMode)
	switch mode {
	case gin.ReleaseMode, gin.DebugMode, gin.TestMode:
	default:
		mode = gin.DebugMode
	}
	path := fmt.Sprintf("config/config-%s.yml", mode)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var config Config
	if err := yaml.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}
	return &config, nil
}

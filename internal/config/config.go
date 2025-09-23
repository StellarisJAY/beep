package config

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

type Config struct {
	DB struct {
		Driver      string `yaml:"driver"`
		AutoMigrate bool   `yaml:"auto_migrate"`
	} `yaml:"db"`
	Postgres struct {
		DSN string `yaml:"dsn"`
	} `yaml:"postgres"`
	MySQL struct {
		DSN string `yaml:"dsn"`
	} `yaml:"mysql"`
	Redis struct {
		Host     string `yaml:"host"`
		Port     uint   `yaml:"port"`
		DB       int    `yaml:"db"`
		Password string `yaml:"password"`
	} `yaml:"redis"`
	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"server"`
	Logger struct {
		Format string `yaml:"format"`
		Level  string `yaml:"level"`
		Path   string `yaml:"path"`
	} `yaml:"logger"`
	Encrypt struct {
		Secret string `yaml:"secret"`
	} `yaml:"encrypt"`
	FileStore string `yaml:"fileStore"`
	Minio     struct {
		Endpoint  string `yaml:"endpoint"`
		AccessKey string `yaml:"accessKey"`
		SecretKey string `yaml:"secretKey"`
		Bucket    string `yaml:"bucket"`
	} `yaml:"minio"`
	Worker string `yaml:"worker"`
	GoPool struct {
		Size    int           `yaml:"size"`
		Timeout time.Duration `yaml:"timeout"`
	} `yaml:"gopool"`
	KnowledgeGraph struct {
		EntityPrompt   string `yaml:"entityPrompt"`
		RelationPrompt string `yaml:"relationPrompt"`
	} `yaml:"knowledge_graph"`
	VectorStore string `yaml:"vectorStore"`
	Milvus      struct {
		Address  string `yaml:"address"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
	} `yaml:"milvus"`
}

func LoadConfig() (*Config, error) {
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

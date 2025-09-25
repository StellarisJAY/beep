package container

import (
	"beep/internal/application/repository"
	"beep/internal/application/service"
	"beep/internal/application/service/captcha"
	"beep/internal/application/service/encrypt"
	"beep/internal/application/service/file"
	"beep/internal/application/service/parser"
	"beep/internal/application/service/vector"
	"beep/internal/config"
	"beep/internal/handler"
	"beep/internal/router"
	"beep/internal/types/interfaces"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/panjf2000/ants/v2"
	"github.com/redis/go-redis/v9"
	"go.uber.org/dig"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewContainer() *dig.Container {
	container := dig.New()
	// 配置
	must(container.Provide(config.LoadConfig))
	// 数据库
	must(container.Provide(InitDatabase))
	// Redis
	must(container.Provide(InitRedis))
	// 文件存储
	must(container.Provide(InitFileStore))
	// 向量存储
	must(container.Provide(InitVectorStore))
	// 线程池
	must(container.Provide(InitAntsPool))
	// 验证码
	must(container.Provide(captcha.New))
	// 加密
	must(container.Provide(encrypt.NewAesEncryptService))

	// 数据层
	must(container.Provide(repository.NewUserRepo))
	must(container.Provide(repository.NewWorkspaceRepo))
	must(container.Provide(repository.NewUserWorkspaceRepo))
	must(container.Provide(repository.NewKnowledgeBaseRepo))
	must(container.Provide(repository.NewModelRepo))
	must(container.Provide(repository.NewModelFactoryRepo))
	must(container.Provide(repository.NewMCPServerRepo))
	must(container.Provide(repository.NewDocumentRepo))

	// 服务层
	must(container.Provide(service.NewUserService))
	must(container.Provide(service.NewWorkspaceService))
	must(container.Provide(service.NewKnowledgeBaseService))
	must(container.Provide(service.NewModelService))
	must(container.Provide(service.NewMCPServerService))
	must(container.Provide(service.NewDocumentService))
	// 文档解析器
	must(container.Provide(parser.NewDocumentParser))

	// handler
	must(container.Provide(handler.NewUserHandler))
	must(container.Provide(handler.NewWorkspaceHandler))
	must(container.Provide(handler.NewKnowledgeBaseHandler))
	must(container.Provide(handler.NewModelHandler))
	must(container.Provide(handler.NewMCPServerHandler))
	must(container.Provide(handler.NewDocumentHandler))
	// gin engine
	must(container.Provide(router.InitRouter))
	return container
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func InitDatabase(config *config.Config) (*gorm.DB, error) {
	var dialector gorm.Dialector
	switch config.DB.Driver {
	case "mysql":
		dialector = mysql.Open(config.MySQL.DSN)
	case "postgres":
		dialector = postgres.Open(config.Postgres.DSN)
	default:
		return nil, errors.New("unsupported database")
	}
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func InitRedis(config *config.Config) (*redis.Client, error) {
	cli := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port),
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := cli.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	return cli, nil
}

func InitFileStore(config *config.Config) (interfaces.FileStore, error) {
	switch config.FileStore {
	case "minio":
		return file.NewMinio(config)
	default:
		return nil, errors.New("unsupported file store")
	}
}

func InitVectorStore(config *config.Config) (interfaces.VectorStore, error) {
	switch config.VectorStore {
	case "milvus":
		return vector.NewMilvus(config.Milvus.Address, config.Milvus.Username, config.Milvus.Password, config.Milvus.Database)
	default:
		return nil, errors.New("unsupported vector store")
	}
}

func InitAntsPool(config *config.Config) (*ants.Pool, error) {
	return ants.NewPool(config.GoPool.Size, ants.WithPreAlloc(true))
}

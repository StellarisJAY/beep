package types

type ModelFactoryType string

const (
	FactoryOpenAI           ModelFactoryType = "open_ai"
	FactoryOllama           ModelFactoryType = "ollama"
	FactoryDashscope        ModelFactoryType = "dashscope"
	FactoryOpenAICompatible ModelFactoryType = "open_ai_compatible"
)

type ModelFactoryConfig struct {
	BaseUrl   string `json:"base_url"`
	APIKey    string `json:"api_key"`
	ExtConfig string `json:"ext_config"`
}

// ModelFactory 模型供应商
type ModelFactory struct {
	BaseEntity
	Name          string           `json:"name"`
	Type          ModelFactoryType `json:"type"`
	EncryptConfig string           `json:"encrypt_config"`
	WorkspaceId   int64            `json:"workspace_id"`
}

func (ModelFactory) TableName() string {
	return "model_factories"
}

type ModelType string

const (
	ModelTypeLLM       ModelType = "llm"
	ModelTypeEmbedding ModelType = "embedding"
	ModelTypeReranking ModelType = "reranking"
)

// Model 模型
type Model struct {
	BaseEntity
	Name        string    `json:"name" gorm:"not null;type:varchar(64);'"`
	Type        ModelType `json:"type" gorm:"not null;type:varchar(64);'"`
	FactoryId   int64     `json:"factory_id" gorm:"not null;"`
	Config      string    `json:"config" gorm:"not null;type:text;'"`
	WorkspaceId int64     `json:"workspace_id" gorm:"not null;"`
}

func (Model) TableName() string {
	return "models"
}

// WorkspaceDefaultModel 工作空间默认模型
type WorkspaceDefaultModel struct {
	BaseEntity
	WorkspaceId int64     `json:"workspace_id" gorm:"not null;"`
	Name        string    `json:"name" gorm:"not null;type:varchar(64);'"`
	Type        ModelType `json:"type" gorm:"not null;type:varchar(16);'"`
	Config      string    `json:"config" gorm:"not null;type:text;'"`
	FactoryId   int64     `json:"factory_id" gorm:"not null;"`
}

func (WorkspaceDefaultModel) TableName() string {
	return "workspace_default_models"
}

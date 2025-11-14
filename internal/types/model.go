package types

import "gorm.io/gorm"

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
	Name        string           `json:"name" gorm:"type:varchar(64);not null;"`         // 模型供应商名称
	Type        ModelFactoryType `json:"type" gorm:"type:varchar(64);not null;"`         // 模型供应商类型
	BaseUrl     string           `json:"base_url" gorm:"type:varchar(255);not null;"`    // 模型供应商基础url
	APIKey      string           `json:"api_key" gorm:"type:text;not null;"`             // 模型供应商api key
	WorkspaceId string           `json:"workspace_id" gorm:"not null;type:varchar(36);"` // 工作空间id

	Models []*Model `json:"models" gorm:"-"` // 模型供应商下的模型
}

func (*ModelFactory) TableName() string {
	return "model_factories"
}

func (m *ModelFactory) BeforeCreate(tx *gorm.DB) error {
	if err := m.BaseEntity.BeforeCreate(tx); err != nil {
		return err
	}
	if m.WorkspaceId == "" {
		// 从context中获取workspaceId
		if workspaceId, ok := tx.Statement.Context.Value(WorkspaceIdContextKey).(string); ok {
			m.WorkspaceId = workspaceId
		}
	}
	return nil
}

type ModelType string

const (
	ModelTypeChat      ModelType = "chat"
	ModelTypeEmbedding ModelType = "embedding"
	ModelTypeReranking ModelType = "reranking"
)

// Model 模型
type Model struct {
	BaseEntity
	Name         string    `json:"name" gorm:"not null;type:varchar(64);'"`        // 模型名称
	Type         ModelType `json:"type" gorm:"not null;type:varchar(64);'"`        // 模型类型
	Tags         string    `json:"tags" gorm:"not null;type:varchar(255);'"`       // 模型标签
	MaxTokens    int64     `json:"max_tokens" gorm:"not null;type:bigint;'"`       // 模型最大token数
	FunctionCall bool      `json:"function_call" gorm:"not null;'"`                // 模型是否支持函数调用
	FactoryId    string    `json:"factory_id" gorm:"not null;type:varchar(36);"`   // 模型供应商id
	WorkspaceId  string    `json:"workspace_id" gorm:"not null;type:varchar(36);"` // 工作空间id
	Status       bool      `json:"status" gorm:"not null;type:bool;default:true;"` // 模型是否可用
}

func (*Model) TableName() string {
	return "models"
}

func (m *Model) BeforeCreate(tx *gorm.DB) error {
	if err := m.BaseEntity.BeforeCreate(tx); err != nil {
		return err
	}
	if m.WorkspaceId == "" {
		// 从context中获取workspaceId
		if workspaceId, ok := tx.Statement.Context.Value(WorkspaceIdContextKey).(string); ok {
			m.WorkspaceId = workspaceId
		}
	}
	return nil
}

// WorkspaceDefaultModel 工作空间默认模型
type WorkspaceDefaultModel struct {
	BaseEntity
	WorkspaceId string    `json:"workspace_id" gorm:"not null;"`
	Name        string    `json:"name" gorm:"not null;type:varchar(64);'"`
	Type        ModelType `json:"type" gorm:"not null;type:varchar(16);'"`
	Config      string    `json:"config" gorm:"not null;type:text;'"`
	FactoryId   string    `json:"factory_id" gorm:"not null;"`
}

func (WorkspaceDefaultModel) TableName() string {
	return "workspace_default_models"
}

type CreateModelFactoryReq struct {
	Name    string           `json:"name" binding:"required"`
	Type    ModelFactoryType `json:"type" binding:"required"`
	BaseUrl string           `json:"base_url"`
	ApiKey  string           `json:"api_key"`
}

type UpdateModelFactoryReq struct {
	Id      string `json:"id" binding:"required"`
	Name    string `json:"name" binding:"required"`
	BaseUrl string `json:"base_url"`
	ApiKey  string `json:"api_key"`
}

type ListModelQuery struct {
	FactoryId string    `json:"factory_id" form:"factory_id"`
	Type      ModelType `json:"type" form:"type"`
}

type ModelDetail struct {
	Id           string    `json:"id"`
	Name         string    `json:"name"`          // 模型名称
	Type         ModelType `json:"type"`          // 模型类型
	Tags         string    `json:"tags"`          // 模型标签
	MaxTokens    int64     `json:"max_tokens"`    // 模型最大token数
	FunctionCall bool      `json:"function_call"` // 模型是否支持函数调用
	FactoryId    string    `json:"factory_id"`
	WorkspaceId  string    `json:"workspace_id"`

	FactoryType     ModelFactoryType `json:"factory_type"` // 模型供应商类型
	ApiKey          string           `json:"api_key"`      // 模型供应商api key
	ApiKeyDecrypted string           `json:"-"`            // 解密的APIKey
	BaseUrl         string           `json:"base_url"`     // 模型供应商基础url
}

type ModelFactoryTemplate struct {
	Name    string `json:"name"`
	SdkType string `json:"sdk_type"`
	Models  []struct {
		Name         string    `json:"llm_name"`
		Type         ModelType `json:"model_type"`
		Tags         string    `json:"tags"`
		MaxTokens    int64     `json:"max_tokens"`
		FunctionCall bool      `json:"function_call"`
	} `json:"models"`
	DefaultConfig struct {
		ApiKey  string `json:"api_key"`
		BaseUrl string `json:"base_url"`
	} `json:"default_config"`
}

type ChatModelOption struct {
	Temperature      *float32 `json:"temperature"`
	Thinking         *bool    `json:"thinking"`
	TopP             *float32 `json:"top_p,omitempty"`
	Stop             []string `json:"stop,omitempty"`
	PresencePenalty  *float32 `json:"presence_penalty,omitempty"`
	ResponseFormat   string   `json:"response_format,omitempty"`
	Seed             *int     `json:"seed,omitempty"`
	FrequencyPenalty *float32 `json:"frequency_penalty,omitempty"`
}

package types

type AgentType string

const (
	AgentTypeSimple   AgentType = "simple"
	AgentTypeWorkflow AgentType = "workflow"
)

type Agent struct {
	BaseEntity
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Type         AgentType `json:"type"`
	Config       string    `json:"config"`
	WorkspaceId  uint64    `json:"workspace_id"`
	CreateBy     uint64    `json:"create_by"`
	LastUpdateBy uint64    `json:"last_update_by"`
}

func (*Agent) TableName() string {
	return "agents"
}

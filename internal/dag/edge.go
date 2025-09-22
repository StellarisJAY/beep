package dag

type Edge struct {
	Id       string         // 连线ID uuid
	Source   string         // 源节点ID
	Target   string         // 目标节点ID
	Metadata map[string]any // 连线元数据
}

func NewEdge(id, source, target string, metadata map[string]any) *Edge {
	return &Edge{
		Id:       id,
		Source:   source,
		Target:   target,
		Metadata: metadata,
	}
}

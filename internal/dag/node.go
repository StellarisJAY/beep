package dag

import "context"

// Node 节点
type Node struct {
	id       string         // 节点ID uuid
	name     string         // 节点名称
	metadata map[string]any // 节点元数据，前端流程编排的数据结构转换成元数据
	nodeFunc NodeFunc       // 节点的任务函数

	config   any     // 节点配置，不同类型的实现不同
	parents  []*Node // 父节点
	children []*Node // 子节点
	indegree int     // 入度，节点初始入度为父节点的数量，每完成一个父节点减少一个入度，当入度为0时该节点才能执行
}

func NewNode(id, name string, nodeFunc NodeFunc) *Node {
	return &Node{
		id:       id,
		name:     name,
		nodeFunc: nodeFunc,
		metadata: make(map[string]any),
	}
}

type NodeFuncReturn struct {
	Output       map[string]any // 返回值
	Error        error          // 错误
	SkipChildren []string       // 跳过的子节点列表
}

type NodeFunc func(ctx context.Context, node Node) NodeFuncReturn

func (n *Node) Id() string {
	return n.id
}

func (n *Node) Children() []string {
	children := make([]string, len(n.children))
	for i, child := range n.children {
		children[i] = child.Id()
	}
	return children
}

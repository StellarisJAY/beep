package dag

import (
	"errors"
	"slices"
)

var (
	ErrDuplicateNodeId = errors.New("duplicate node id")
	ErrInvalidEdge     = errors.New("invalid edge")
	ErrLoopDetected    = errors.New("loop detected")
)

type Graph struct {
	nodes []*Node
	edges []*Edge

	nodeMap map[string]*Node
	edgeMap map[string]*Edge
	path    []*Node
}

func NewGraph() *Graph {
	return &Graph{
		nodes:   make([]*Node, 0),
		edges:   make([]*Edge, 0),
		nodeMap: make(map[string]*Node),
		edgeMap: make(map[string]*Edge),
	}
}

func (g *Graph) AddNode(nodes ...*Node) *Graph {
	g.nodes = append(g.nodes, nodes...)
	return g
}

func (g *Graph) AddEdge(edges ...*Edge) *Graph {
	g.edges = append(g.edges, edges...)
	return g
}

// Compile 编排节点，构建子节点、父节点关系，检查图是否有环
func (g *Graph) Compile() error {
	// 创建节点id->节点映射
	for _, node := range g.nodes {
		if _, ok := g.nodeMap[node.id]; ok {
			return ErrDuplicateNodeId
		}
		g.nodeMap[node.id] = node
	}
	// 通过连线，构建父子节点关系
	for _, edge := range g.edges {
		source, ok := g.nodeMap[edge.Source]
		if !ok {
			return ErrInvalidEdge
		}
		target, ok := g.nodeMap[edge.Target]
		if !ok {
			return ErrInvalidEdge
		}
		source.children = append(source.children, target)
		target.parents = append(target.parents, source)
	}

	// 设置节点入度
	for _, node := range g.nodes {
		node.indegree = len(node.parents)
	}

	// 检查是否有环
	hasLoop, path := g.kahnCheckLoop()
	if hasLoop {
		return ErrLoopDetected
	}
	g.path = path
	return nil
}

func (g *Graph) SetNodeFunc(id string, nodeFunc NodeFunc) {
	for _, node := range g.nodes {
		if node.id == id {
			node.nodeFunc = nodeFunc
			return
		}
	}
}

func (g *Graph) ForCompiledPath(fn func(n Node)) {
	for _, node := range g.path {
		fn(*node)
	}
}

// kahnCheckLoop kahn算法检查是否有环
func (g *Graph) kahnCheckLoop() (bool, []*Node) {
	// 复制一份边列表，用来在kahn算法中做删除，保留原图的边列表不变
	edges := slices.Clone(g.edges)

	// kahn算法，see: https://oi-wiki.org/graph/topo/#kahn-%E7%AE%97%E6%B3%95
	S := make([]*Node, 0, len(g.nodes)) // 入度0节点列表
	L := make([]*Node, 0, len(g.nodes)) // 拓扑排序结果列表
	// 找到所有入度为0的节点，即开始节点
	for _, node := range g.nodes {
		if node.indegree == 0 {
			S = append(S, node)
		}
	}

	// 遍历直到没有入度0的节点
	for len(S) > 0 {
		n := S[0]
		S = S[1:]
		L = append(L, n)
		// 遍历入度0节点的所有子节点
		for _, child := range n.children {
			// 删除节点n到子节点的连接
			edges = slices.DeleteFunc(edges, func(e *Edge) bool {
				return e.Source == n.id && e.Target == child.id
			})
			// 删除n到子节点连接后，获取子节点的入度
			hasConnection := slices.ContainsFunc(edges, func(e *Edge) bool {
				return e.Target == child.id
			})
			// 子节点入度为0，加入S
			if !hasConnection {
				S = append(S, child)
			}
		}
	}
	// 还有多余的边，表示有环路
	return len(edges) > 0, L
}

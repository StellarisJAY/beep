package dag

import "fmt"

type Chain Graph

func NewChain() *Chain {
	return (*Chain)(NewGraph())
}

func (c *Chain) AddNode(nodes ...*Node) *Chain {
	if len(nodes) == 0 {
		return c
	}
	n := len(c.nodes)
	c.nodes = append(c.nodes, nodes...)
	var lastNode *Node
	if n == 0 {
		lastNode = nodes[0]
	} else {
		lastNode = c.nodes[n-1]
	}
	for _, node := range nodes {
		if lastNode == node {
			continue
		}
		c.edges = append(c.edges, &Edge{
			Id:     fmt.Sprintf("%s->%s", lastNode.id, node.id),
			Source: lastNode.id,
			Target: node.id,
		})
		lastNode = node
	}
	return c
}

func (c *Chain) Compile() error {
	return (*Graph)(c).Compile()
}

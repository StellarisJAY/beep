package dag

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func Test_CompileAndRun(t *testing.T) {
	node0 := NewNode("0", "0", func(ctx context.Context, node Node) NodeFuncReturn {
		time.Sleep(1 * time.Second)
		fmt.Println("exec node0")
		return NodeFuncReturn{
			Output: map[string]any{
				"hello": "world",
			},
			SkipChildren: []string{},
		}
	})
	node1 := NewNode("1", "1", func(ctx context.Context, node Node) NodeFuncReturn {
		time.Sleep(1 * time.Second)
		fmt.Println("exec node1")
		return NodeFuncReturn{
			Output: map[string]any{
				"res1": "value1",
			},
			SkipChildren: make([]string, 0),
		}
	})

	node2 := NewNode("2", "2", func(ctx context.Context, node Node) NodeFuncReturn {
		time.Sleep(1 * time.Second)
		fmt.Println("exec node2")
		return NodeFuncReturn{
			Output: map[string]any{
				"res2": "value2",
			},
			SkipChildren: make([]string, 0),
		}
	})

	node3 := NewNode("3", "3", func(ctx context.Context, node Node) NodeFuncReturn {
		time.Sleep(1 * time.Second)
		fmt.Println("exec node3")
		return NodeFuncReturn{
			Output: map[string]any{
				"res3": "value3",
			},
			SkipChildren: make([]string, 0),
		}
	})

	node4 := NewNode("4", "4", func(ctx context.Context, node Node) NodeFuncReturn {
		time.Sleep(1 * time.Second)
		fmt.Println("exec node4")
		return NodeFuncReturn{
			Output: map[string]any{
				"res4": "value4",
			},
			SkipChildren: make([]string, 0),
		}
	})

	node5 := NewNode("5", "5", func(ctx context.Context, node Node) NodeFuncReturn {
		time.Sleep(1 * time.Second)
		fmt.Println("exec node5")
		return NodeFuncReturn{
			Output: map[string]any{
				"res5": "value5",
			},
			SkipChildren: make([]string, 0),
		}
	})

	n0n1 := NewEdge("0->1", "0", "1", nil)
	n0n3 := NewEdge("0->3", "0", "3", nil)
	n1n2 := NewEdge("1->2", "1", "2", nil)
	n3n4 := NewEdge("3->4", "3", "4", nil)
	n2n5 := NewEdge("2->5", "2", "5", nil)
	n4n5 := NewEdge("4->5", "4", "5", nil)
	g := NewGraph().AddNode(node0, node1, node2, node3, node4, node5).AddEdge(n1n2, n2n5, n3n4, n4n5, n0n1, n0n3)
	if err := g.Compile(); err != nil {
		t.Fatal(err)
	}
	run := NewGraphRun(g)
	if err := run.Run(WithParallelNum(2), WithTimeout(30*time.Second)); err != nil {
		t.Fatal(err)
	}
	for _, res := range run.execResults {
		fmt.Println(res)
	}
}

package dag

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"slices"
	"sync"
	"time"

	"github.com/panjf2000/ants/v2"
)

type EventType byte

const (
	EventTypeStart EventType = iota
	EventTypeNodeEnd
	EventTypeEnd
)

type CallbackEvent struct {
	Type       EventType
	Data       any
	GlobalData map[string]any
}

type GraphRun struct {
	graph   *Graph        // DAG图
	context *GraphContext // 执行上下文
	cancel  context.CancelFunc

	parallelNum  int       // 最大并行度
	panicHandler func(any) // 异常处理
	nonBlocking  bool      // 是否是非阻塞模式

	workers      *ants.Pool    // 并行任务池,goroutine数量=最大并行度
	nodeChan     chan *Node    // 需要执行的节点channel，用来接受并行任务的待执行节点
	doneNodeChan chan ExecInfo // 节点执行结果channel，用来接受并行任务的结果

	execResults       []ExecInfo          // 执行结果集
	skippedNodesMutex sync.RWMutex        // 并行模式下更新跳过节点的锁
	skippedNodes      map[string]struct{} // 跳过节点

	callbackFn func(event CallbackEvent)
}

type ExecInfo struct {
	Id         string
	OutputKeys []string
	Skipped    bool
}

func NewGraphRun(graph *Graph) *GraphRun {
	g := &GraphRun{
		graph:             graph,
		context:           newGraphContext(),
		parallelNum:       1,            // 默认并行度1
		panicHandler:      func(any) {}, // 默认异常处理
		skippedNodes:      make(map[string]struct{}),
		skippedNodesMutex: sync.RWMutex{},
		execResults:       make([]ExecInfo, 0),
	}
	ctx, cancelFunc := context.WithCancel(context.Background())
	g.cancel = cancelFunc
	g.context.ctx = ctx
	return g
}

func NewChainRun(chain *Chain) *GraphRun {
	return NewGraphRun((*Graph)(chain))
}

type RunOption func(run *GraphRun) error

func WithParallelNum(num int) RunOption {
	num = max(num, 1)
	return func(run *GraphRun) error {
		run.parallelNum = num
		return nil
	}
}

func WithTimeout(timeout time.Duration) RunOption {
	return func(run *GraphRun) error {
		ctx, cancel := context.WithTimeout(run.context.ctx, timeout)
		run.cancel = cancel
		run.context.ctx = ctx
		return nil
	}
}

func WithPanicHandler(f func(err error)) RunOption {
	return func(run *GraphRun) error {
		if f == nil {
			run.panicHandler = func(any) {}
			return nil
		}
		run.panicHandler = func(e any) {
			if err, ok := e.(error); ok {
				f(err)
			}
		}
		return nil
	}
}

func WithNonBlocking() RunOption {
	return func(run *GraphRun) error {
		run.nonBlocking = true
		return nil
	}
}

func WithVariables(variables map[string]any) RunOption {
	return func(run *GraphRun) error {
		for k, v := range variables {
			run.context.variables[k] = v
		}
		return nil
	}
}

func WithContext(ctx context.Context) RunOption {
	return func(run *GraphRun) error {
		cancelCtx, cancelFunc := context.WithCancel(ctx)
		run.context.ctx = cancelCtx
		run.cancel = cancelFunc
		return nil
	}
}

func WithCallback(callback func(event CallbackEvent)) RunOption {
	return func(run *GraphRun) error {
		if callback != nil {
			run.callbackFn = callback
		}
		return nil
	}
}

func WithWorkerPool(pool *ants.Pool) RunOption {
	return func(run *GraphRun) error {
		if pool == nil {
			return fmt.Errorf("worker pool is nil")
		}
		run.workers = pool
		return nil
	}
}

func (g *GraphRun) Run(options ...RunOption) error {
	for _, option := range options {
		if err := option(g); err != nil {
			return fmt.Errorf("apply option error: %w", err)
		}
	}

	if g.workers == nil && g.parallelNum > 1 {
		workers, err := ants.NewPool(g.parallelNum, ants.WithPreAlloc(true), ants.WithPanicHandler(g.panicHandler))
		if err != nil {
			return fmt.Errorf("create worker pool error: %w", err)
		}
		g.workers = workers
	}

	// 非阻塞模式，开启goroutine
	if g.nonBlocking {
		go g.run()
	} else {
		g.run()
	}

	return nil
}

func (g *GraphRun) Cancel() {
	g.cancel()
}

func (g *GraphRun) run() {
	g.callbackFn(CallbackEvent{Type: EventTypeStart, GlobalData: g.context.variables})
	if g.parallelNum == 1 {
		g.serialRun()
	} else {
		// 获取所有入度为0的节点，即所有开始节点
		g.nodeChan = make(chan *Node, len(g.graph.nodes))
		g.doneNodeChan = make(chan ExecInfo, len(g.graph.nodes))
		for _, node := range g.graph.nodes {
			if node.indegree == 0 {
				g.nodeChan <- node
			}
		}
		g.parallelRun()
	}
	g.callbackFn(CallbackEvent{Type: EventTypeEnd, GlobalData: g.context.variables})
}

// serialRun 串行执行
func (g *GraphRun) serialRun() {
	defer func() {
		if err := recover(); err != nil {
			g.panicHandler(err)
		}
	}()

	done := func() bool {
		select {
		case <-g.context.Done():
			return true
		default:
			return false
		}
	}

	skipNodes := make(map[string]struct{})
	for _, node := range g.graph.path {
		if done() {
			break
		}
		if node.indegree > 0 {
			continue
		}
		execRes := ExecInfo{Id: node.id}
		// 节点被跳过，所有子节点入度减一
		if _, ok := skipNodes[node.id]; ok {
			execRes.Skipped = true
			for _, child := range node.children {
				child.indegree -= 1
				// 子节点入度为0，跳过
				if child.indegree == 0 {
					skipNodes[child.id] = struct{}{}
				}
			}
			continue
		}
		// 执行节点函数
		res := node.nodeFunc(g.context, *node)
		if res.Error != nil {
			panic(res.Error)
		}
		// 保存节点输出
		g.context.SaveVariables(res.Output)
		execRes.OutputKeys = slices.Collect(maps.Keys(res.Output))
		// 所有子节点的入度减一
		for _, child := range node.children {
			child.indegree -= 1
		}
		// 跳过的子节点
		for _, skip := range res.SkipChildren {
			skipNodes[skip] = struct{}{}
		}

		g.execResults = append(g.execResults, execRes)
		g.callbackFn(CallbackEvent{
			EventTypeNodeEnd,
			execRes,
			g.context.variables,
		})
	}
}

// parallelRun 并行执行
func (g *GraphRun) parallelRun() {
	defer func() {
		if err := recover(); err != nil {
			g.panicHandler(err)
		}
		g.workers.Release()
	}()
	for {
		select {
		case <-g.context.Done():
			err := g.context.Err()
			if errors.Is(err, context.Canceled) {
				slog.Info("dag execution canceled")
			}
			if errors.Is(err, context.DeadlineExceeded) {
				slog.Info("dag execution timeout")
			}
			g.workers.Release()
			return
		case execRes := <-g.doneNodeChan: // 节点执行结束
			g.execResults = append(g.execResults, execRes)
			if len(g.execResults) == len(g.graph.nodes) {
				// 如果所有节点都已结束
				g.cancel()
			}
		case node := <-g.nodeChan:
			// 入度大于0，还有前驱节点没有结束
			if node.indegree > 0 {
				continue
			}
			// 提交到线程池并发执行
			_ = g.workers.Submit(func() {
				g.parallelRunNodeFunc(node)
			})
		}
	}
}

func (g *GraphRun) parallelRunNodeFunc(node *Node) {
	execRes := ExecInfo{
		Id: node.id,
	}

	g.skippedNodesMutex.RLock()
	_, skipped := g.skippedNodes[node.id]
	g.skippedNodesMutex.RUnlock()

	execRes.Skipped = skipped
	// 当前节点没被跳过
	if !skipped {
		// 执行节点函数
		res := node.nodeFunc(g.context, *node)
		if res.Error != nil {
			panic(res.Error)
		}
		// 保存节点输出的变量
		execRes.OutputKeys = slices.Collect(maps.Keys(res.Output))
		g.context.SaveVariables(res.Output)
		// 节点返回跳过某些子节点，更新被跳过的节点列表
		g.skippedNodesMutex.Lock()
		for _, skip := range res.SkipChildren {
			g.skippedNodes[skip] = struct{}{}
		}
		g.skippedNodesMutex.Unlock()
	}
	// 所有子节点的入度减一
	for _, child := range node.children {
		child.indegree -= 1
		// 如果当前节点被跳过，且子节点没有其他前驱（入度为0），跳过子节点
		if skipped && child.indegree == 0 {
			g.skippedNodesMutex.Lock()
			g.skippedNodes[child.id] = struct{}{}
			g.skippedNodesMutex.Unlock()
		}
		// 子节点加入待执行队列
		g.nodeChan <- child
	}
	// 当前节点处理完成
	g.doneNodeChan <- execRes
	g.callbackFn(CallbackEvent{
		EventTypeNodeEnd,
		execRes,
		g.context.variables,
	})
}

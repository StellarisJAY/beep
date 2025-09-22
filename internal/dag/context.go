package dag

import (
	"context"
	"sync"
	"time"
)

type GraphContext struct {
	variables map[string]any // 全局变量列表
	mutex     sync.RWMutex
	ctx       context.Context
}

func newGraphContext() *GraphContext {
	return &GraphContext{
		variables: make(map[string]interface{}),
		mutex:     sync.RWMutex{},
	}
}

func (g *GraphContext) Deadline() (deadline time.Time, ok bool) {
	return g.ctx.Deadline()
}

func (g *GraphContext) Done() <-chan struct{} {
	return g.ctx.Done()
}

func (g *GraphContext) Err() error {
	return g.ctx.Err()
}

func (g *GraphContext) Value(key any) any {
	g.mutex.RLock()
	defer g.mutex.RUnlock()
	val := g.ctx.Value(key)
	if val != nil {
		return val
	}
	keyStr, ok := key.(string)
	if !ok {
		return nil
	}
	return g.variables[keyStr]
}

func (g *GraphContext) SaveVariables(variables map[string]any) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	for k, v := range variables {
		g.variables[k] = v
	}
}

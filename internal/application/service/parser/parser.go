package parser

import (
	"beep/internal/dag"
	"beep/internal/types"
	"beep/internal/types/interfaces"
	"context"
	"sync"

	"github.com/cloudwego/eino/schema"
	"github.com/panjf2000/ants/v2"
)

var runningParsers = sync.Map{}

// DocumentParser 文档解析器
type DocumentParser struct {
	workerPool  *ants.Pool             // goroutine pool
	vectorStore interfaces.VectorStore // 向量数据库
}

const (
	parseInfoContextKey = "parse_info"
)

func NewDocumentParser(workerPool *ants.Pool, vectorStore interfaces.VectorStore) interfaces.ParseService {
	return &DocumentParser{
		workerPool:  workerPool,
		vectorStore: vectorStore,
	}
}

// Parse 异步解析文档
func (d *DocumentParser) Parse(_ context.Context, info types.ParseInfo) error {
	// 构建dag任务链
	chain := dag.NewChain().AddNode(
		dag.NewNode("chunking", "chunking", d.chunking),          // 文档切片
		dag.NewNode("embedding", "embedding", d.embedding),       // 嵌入
		dag.NewNode("summarizing", "summarizing", d.summarizing)) // 文档内容总结
	if info.EnableKnowledgeGraph {
		chain.AddNode(dag.NewNode("knowledge_graph", "knowledge_graph", d.genKnowledgeGraph)) // 知识图谱提取
	}
	_ = chain.Compile()
	runner := dag.NewChainRun(chain)
	variables := map[string]any{
		parseInfoContextKey: &info,
	}
	// 异步执行任务链
	if err := runner.Run(dag.WithWorkerPool(d.workerPool),
		dag.WithVariables(variables),
		dag.WithNonBlocking(),
		dag.WithPanicHandler(d.panicHandler),
		dag.WithCallback(d.nodeCallback),
		dag.WithParallelNum(1)); err != nil {
		return err
	}
	runningParsers.Store(info.DocId, runner)
	return nil
}

func (d *DocumentParser) CancelParse(ctx context.Context, docId int64) error {
	value, ok := runningParsers.LoadAndDelete(docId)
	if !ok {
		return nil
	}
	value.(*dag.GraphRun).Cancel()
	return nil
}

// summarizing 文档内容总结
func (d *DocumentParser) summarizing(ctx context.Context, _ dag.Node) (result dag.NodeFuncReturn) {
	parseInfo := ctx.Value(parseInfoContextKey).(*types.ParseInfo)
	response, err := parseInfo.ChatModel.Generate(ctx, []*schema.Message{
		{
			Role: schema.System,
			// TODO better prompt for summerizing
			Content: "你是一个文档内容总结器，你的任务是根据用户提供的文档内容生成一个总结。要求使用中文，且总结不超过100字。",
		},
		{
			Role:    schema.User,
			Content: string(parseInfo.Content),
		},
	})
	if err != nil {
		result.Error = err
		return
	}
	result.Output = map[string]any{
		"summary": response,
	}
	return result
}

func (d *DocumentParser) panicHandler(err error) {
	// TODO 处理解析任务链中的panic
}

func (d *DocumentParser) nodeCallback(event dag.CallbackEvent) {
	// TODO 处理解析任务链中的节点回调
}

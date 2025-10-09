package common

import (
	"beep/internal/types"
	"fmt"

	"github.com/cloudwego/eino/schema"
)

func ConvertToolsetToSchemaTools(toolsets []*types.MCPToolSet) ([]*schema.ToolInfo, error) {
	tools := make([]*schema.ToolInfo, 0)
	for _, toolset := range toolsets {
		for _, tool := range toolset.Tools {
			tools = append(tools, &schema.ToolInfo{
				Name: fmt.Sprintf("%s:%s", toolset.Name, tool.Name),
				Desc: tool.Description,
				Extra: map[string]any{
					"inputSchema":  tool.InputSchema,
					"outputSchema": tool.OutputSchema,
					"title":        tool.Title,
				},
			})
		}
	}
	return tools, nil
}

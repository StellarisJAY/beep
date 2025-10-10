package chat

import (
	"beep/internal/types"
	"context"
	"io"
)

type FunctionCall struct {
	FunctionName string `json:"function_name"`
	Arguments    string `json:"arguments"`
}

type Message struct {
	Role          string          `json:"role"`
	Content       string          `json:"content"`
	FunctionCalls []*FunctionCall `json:"function_calls"`
}

type Options struct {
	Thinking         bool    `json:"thinking"`
	MaxTokens        int     `json:"max_tokens"`
	Temperature      float64 `json:"temperature"`
	TopP             float64 `json:"top_p"`
	Stop             string  `json:"stop"`
	FrequencyPenalty float64 `json:"frequency_penalty"`

	McpServers []*types.MCPServer `json:"-"`
}

type BaseModel interface {
	Generate(ctx context.Context, messages []*Message, options *Options) (*Message, error)
	Stream(ctx context.Context, messages []*Message, options *Options) (*Stream, error)
}

type Stream struct {
	messageChan chan *Message
}

func (s *Stream) Next() (*Message, error) {
	select {
	case message, ok := <-s.messageChan:
		if !ok {
			return nil, io.EOF
		}
		return message, nil
	}
}

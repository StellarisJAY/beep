package interfaces

import (
	"beep/internal/types"
	"context"
)

type ChatService interface {
	MessageLoop(ctx context.Context, messageChan chan types.Message, errorChan chan error) error
}

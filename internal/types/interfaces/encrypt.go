package interfaces

import "context"

type EncryptService interface {
	Encrypt(ctx context.Context, data string) (string, error)
	Decrypt(ctx context.Context, data string) (string, error)
}

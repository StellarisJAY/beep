package interfaces

import "context"

// EncryptService 加密服务接口
type EncryptService interface {
	// Encrypt 加密
	Encrypt(ctx context.Context, data string) (string, error)
	// Decrypt 解密
	Decrypt(ctx context.Context, data string) (string, error)
}

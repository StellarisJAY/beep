package encrypt

import (
	"beep/internal/config"
	"beep/internal/types/interfaces"
	"context"
	"crypto/aes"
	"encoding/hex"
)

type AesEncryptService struct {
	config *config.Config
}

func (a *AesEncryptService) Encrypt(ctx context.Context, data string) (string, error) {
	secret := a.config.Encrypt.Secret
	cipher, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return "", err
	}
	// 加密
	encrypted := make([]byte, len(data))
	cipher.Encrypt(encrypted, []byte(data))
	return hex.EncodeToString(encrypted), nil
}

func (a *AesEncryptService) Decrypt(ctx context.Context, data string) (string, error) {
	secret := a.config.Encrypt.Secret
	cipher, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return "", err
	}
	// 解密
	decoded, err := hex.DecodeString(data)
	if err != nil {
		return "", err
	}
	decrypted := make([]byte, len(decoded))
	cipher.Decrypt(decrypted, decoded)
	return string(decrypted), nil
}

func NewAesEncryptService(config *config.Config) interfaces.EncryptService {
	return &AesEncryptService{config: config}
}

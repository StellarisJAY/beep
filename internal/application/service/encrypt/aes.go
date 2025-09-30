package encrypt

import (
	"beep/internal/config"
	"beep/internal/types/interfaces"
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
)

type AesEncryptService struct {
	config *config.Config
}

// PKCS#7填充
func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

// PKCS#7去除填充
func pkcs7Unpadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("解密数据为空")
	}

	unpadding := int(data[length-1])
	if unpadding > length {
		return nil, errors.New("无效的填充")
	}

	return data[:(length - unpadding)], nil
}

func (a *AesEncryptService) Encrypt(ctx context.Context, data string) (string, error) {
	secret := a.config.Encrypt.Secret
	block, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return "", err
	}

	// 使用PKCS#7填充确保数据长度是块大小的倍数
	blockSize := block.BlockSize()
	plaintext := pkcs7Padding([]byte(data), blockSize)

	// 使用CBC模式而不是ECB模式
	ciphertext := make([]byte, len(plaintext))
	iv := make([]byte, blockSize)
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, plaintext)

	return hex.EncodeToString(ciphertext), nil
}

func (a *AesEncryptService) Decrypt(ctx context.Context, ct string) (string, error) {
	ciphertext, err := hex.DecodeString(ct)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(a.config.Encrypt.Secret))
	if err != nil {
		return "", err
	}

	// 检查密文长度
	blockSize := block.BlockSize()
	if len(ciphertext)%blockSize != 0 {
		return "", errors.New("密文不是块大小的倍数")
	}

	// 使用CBC模式解密
	iv := make([]byte, blockSize)
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	// 去除PKCS#7填充
	plaintext, err := pkcs7Unpadding(ciphertext)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

func NewAesEncryptService(config *config.Config) interfaces.EncryptService {
	return &AesEncryptService{config: config}
}

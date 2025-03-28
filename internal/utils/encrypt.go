package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

// 使用 AES-CBC 模式加密消息
func EncryptByAES(message, key string) (string, error) {
	keyBytes := []byte(key)
	iv := keyBytes // IV 和密钥相同

	// 创建 AES 加密块
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	// 填充明文到块大小的整数倍
	blockSize := block.BlockSize()
	plainText := pkcs7Padding([]byte(message), blockSize)

	mode := cipher.NewCBCEncrypter(block, iv)

	cipherText := make([]byte, len(plainText))
	mode.CryptBlocks(cipherText, plainText)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// PKCS7 填充
func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

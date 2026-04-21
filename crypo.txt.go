package myDOH

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// EncryptTxt 将明文字符串加密为一个可传输的字符串。
// 输出格式为：base64(nonce + ciphertext)。
// key 长度必须是 16、24、32 之一，分别对应 AES-128、AES-192、AES-256。
func EncryptTxt(text string, key []byte) (string, error) {
	switch len(key) {
	case 16, 24, 32:
	default:
		return "", fmt.Errorf("invalid key length: %d", len(key))
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("new aes cipher failed: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("new gcm failed: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("generate nonce failed: %w", err)
	}

	cipherText := gcm.Seal(nil, nonce, []byte(text), nil)
	out := append(nonce, cipherText...)

	return base64.StdEncoding.EncodeToString(out), nil
}

// DecryptTxt 将加密字符串解密为原始明文字符串。
// 输入格式必须是：base64(nonce + ciphertext)。
// key 长度必须是 16、24、32 之一，且必须与加密时使用的 key 一致。
func DecryptTxt(encText string, key []byte) (string, error) {
	switch len(key) {
	case 16, 24, 32:
	default:
		return "", fmt.Errorf("invalid key length: %d", len(key))
	}

	raw, err := base64.StdEncoding.DecodeString(encText)
	if err != nil {
		return "", fmt.Errorf("base64 decode failed: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("new aes cipher failed: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("new gcm failed: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(raw) <= nonceSize {
		return "", fmt.Errorf("cipher text too short")
	}

	nonce := raw[:nonceSize]
	cipherText := raw[nonceSize:]

	plain, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", fmt.Errorf("decrypt failed: %w", err)
	}

	return string(plain), nil
}

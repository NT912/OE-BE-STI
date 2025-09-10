package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// AES (Advanced Encryption Standard)
type AES struct{}

func NewAES() *AES {
	return &AES{}
}

// EncryptAES mã hóa plaintext bằng AES-GCM, trả về base64 string
func (*AES) Encrypt(key, plaintext string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", fmt.Errorf("new cipher: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("new gcm: %w", err)
	}

	// tạo nonce (IV) 12 byte ngẫu nhiên
	nonce := make([]byte, aesgcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("generate nonce: %w", err)
	}

	// Seal = nonce + ciphertext + tag
	ciphertext := aesgcm.Seal(nil, nonce, []byte(plaintext), nil)

	// Gắn nonce trước ciphertext để dễ decode
	final := append(nonce, ciphertext...)

	// Encode base64 để lưu DB hoặc truyền network
	return base64.StdEncoding.EncodeToString(final), nil
}

// DecryptAES giải mã base64 ciphertext bằng AES-GCM
func (*AES) Decrypt(key, cipherTextBase64 string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(cipherTextBase64)
	if err != nil {
		return "", fmt.Errorf("decode base64: %w", err)
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", fmt.Errorf("new cipher: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("new gcm: %w", err)
	}

	nonceSize := aesgcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("decrypt failed: %w", err)
	}

	return string(plaintext), nil
}

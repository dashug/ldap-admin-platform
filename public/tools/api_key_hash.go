package tools

import (
	"crypto/rand"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
)

const apiKeyBcryptCost = bcrypt.DefaultCost
const apiKeyPrefix = "glak_"
const apiKeyRawLen = 32 // 32 字节随机，hex 后 64 字符

// GenerateRandomApiKey 生成新的 API Key 明文，格式 glak_<64位hex>
func GenerateRandomApiKey() (string, error) {
	b := make([]byte, apiKeyRawLen)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return apiKeyPrefix + hex.EncodeToString(b), nil
}

// GenApiKeyHash 对 API Key 明文做不可逆哈希，用于存储
func GenApiKeyHash(plainKey string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(plainKey), apiKeyBcryptCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// VerifyApiKeyHash 校验明文 Key 与存储的哈希是否一致
func VerifyApiKeyHash(hash, plainKey string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plainKey))
	return err == nil
}

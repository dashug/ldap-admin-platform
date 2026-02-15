package model

import (
	"time"

	"gorm.io/gorm"
)

// ApiKey 第三方调用用 API Key，仅存储哈希与前缀，创建时明文仅返回一次
type ApiKey struct {
	gorm.Model
	Name     string `gorm:"type:varchar(64);not null;comment:'密钥名称/备注'" json:"name"`
	KeyHash  string `gorm:"type:varchar(255);not null;comment:'密钥哈希'" json:"-"`   // 不明文存储
	KeyPrefix string `gorm:"type:varchar(32);not null;index;comment:'密钥前缀用于查找'" json:"keyPrefix"` // 如 glak_xxx 前若干位
}

func (ApiKey) TableName() string {
	return "api_keys"
}

// CreatedAt 转为可读时间
func (k *ApiKey) GetCreatedAt() time.Time {
	return k.CreatedAt
}

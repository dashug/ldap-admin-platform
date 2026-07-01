package model

import (
	"time"

	"gorm.io/gorm"
)

// ApiKey 第三方调用用 API Key，仅存储哈希与前缀，创建时明文仅返回一次
type ApiKey struct {
	gorm.Model
	Name       string     `gorm:"type:varchar(64);not null;comment:'密钥名称/备注'" json:"name"`
	KeyHash    string     `gorm:"type:varchar(255);not null;comment:'密钥哈希'" json:"-"`                       // 不明文存储
	KeyPrefix  string     `gorm:"type:varchar(32);not null;index;comment:'密钥前缀用于查找'" json:"keyPrefix"` // 如 glak_xxx 前若干位
	Enabled    *bool      `gorm:"default:true;comment:'是否启用，可临时停用而不删除'" json:"enabled"`             // 停用后立即失效
	ExpiresAt  *time.Time `gorm:"comment:'过期时间，空表示永不过期'" json:"expiresAt"`                            // 过期后校验直接失败
	LastUsedAt *time.Time `gorm:"comment:'最近一次使用时间'" json:"lastUsedAt"`                                 // 便于审计与吊销闲置密钥
}

// IsUsable 判断该 Key 当前是否可用（已启用且未过期）
func (k *ApiKey) IsUsable() bool {
	if k.Enabled != nil && !*k.Enabled {
		return false
	}
	if k.ExpiresAt != nil && !k.ExpiresAt.IsZero() && k.ExpiresAt.Before(time.Now()) {
		return false
	}
	return true
}

func (ApiKey) TableName() string {
	return "api_keys"
}

// CreatedAt 转为可读时间
func (k *ApiKey) GetCreatedAt() time.Time {
	return k.CreatedAt
}

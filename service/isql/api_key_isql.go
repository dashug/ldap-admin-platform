package isql

import (
	"errors"
	"fmt"

	"github.com/dashug/ldap-admin-platform/model"
	"github.com/dashug/ldap-admin-platform/public/common"
	"github.com/dashug/ldap-admin-platform/public/tools"

	"gorm.io/gorm"
)

type ApiKeyService struct{}

// List 分页列表
func (s ApiKeyService) List(offset, limit int) ([]*model.ApiKey, error) {
	var list []*model.ApiKey
	err := common.DB.Model(&model.ApiKey{}).Order("created_at DESC").Offset(offset).Limit(limit).Find(&list).Error
	return list, err
}

// Count 总数
func (s ApiKeyService) Count() (int64, error) {
	var n int64
	err := common.DB.Model(&model.ApiKey{}).Count(&n).Error
	return n, err
}

// Add 创建；rawKey 为生成的明文（如 glak_xxx），仅此一次返回，调用方需保存
func (s ApiKeyService) Add(name, rawKey string) (*model.ApiKey, error) {
	hash, err := tools.GenApiKeyHash(rawKey)
	if err != nil {
		return nil, fmt.Errorf("生成密钥哈希失败: %w", err)
	}
	prefix := rawKey
	if len(prefix) > 20 {
		prefix = prefix[:20] // 用于查找时缩小范围
	}
	ak := &model.ApiKey{
		Name:      name,
		KeyHash:   hash,
		KeyPrefix: prefix,
	}
	if err := common.DB.Create(ak).Error; err != nil {
		return nil, err
	}
	return ak, nil
}

// FindByPrefix 按前缀查找（用于校验时先按前缀缩小范围）
func (s ApiKeyService) FindByPrefix(prefix string) ([]*model.ApiKey, error) {
	if len(prefix) < 5 {
		return nil, nil
	}
	var list []*model.ApiKey
	err := common.DB.Where("key_prefix = ? OR key_prefix LIKE ?", prefix, prefix+"%").Find(&list).Error
	return list, err
}

// Verify 校验 rawKey：找到前缀匹配的 key 再校验哈希
func (s ApiKeyService) Verify(rawKey string) (*model.ApiKey, bool) {
	if len(rawKey) < 10 {
		return nil, false
	}
	prefix := rawKey
	if len(prefix) > 20 {
		prefix = prefix[:20]
	}
	keys, err := s.FindByPrefix(prefix)
	if err != nil || len(keys) == 0 {
		return nil, false
	}
	for _, k := range keys {
		if tools.VerifyApiKeyHash(k.KeyHash, rawKey) {
			return k, true
		}
	}
	return nil, false
}

// Delete 按 ID 删除
func (s ApiKeyService) Delete(id uint) error {
	var ak model.ApiKey
	if err := common.DB.First(&ak, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("API Key 不存在")
		}
		return err
	}
	return common.DB.Delete(&ak).Error
}

// Find 按 ID 查询
func (s ApiKeyService) Find(id uint) (*model.ApiKey, error) {
	var ak model.ApiKey
	err := common.DB.First(&ak, id).Error
	return &ak, err
}

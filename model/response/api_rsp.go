package response

import (
	"time"

	"github.com/dashug/ldap-admin-platform/model"
)

// ApiKeyCreateRsp 创建 API Key 成功响应，key 仅此一次返回
type ApiKeyCreateRsp struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	KeyPrefix string    `json:"keyPrefix"`
	Key       string    `json:"key"` // 明文密钥，仅创建时返回一次
	CreatedAt time.Time `json:"createdAt"`
}

// ApiKeyItemRsp 列表项（不含明文）
type ApiKeyItemRsp struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	KeyPrefix string    `json:"keyPrefix"`
	CreatedAt time.Time `json:"createdAt"`
}

// ApiKeyListRsp 列表
type ApiKeyListRsp struct {
	Total int64           `json:"total"`
	Items []ApiKeyItemRsp `json:"items"`
}

type ApiTreeRsp struct {
	ID       int          `json:"ID"`
	Remark   string       `json:"remark"`
	Category string       `json:"category"`
	Children []*model.Api `json:"children"`
}

type ApiListRsp struct {
	Total int64       `json:"total"`
	Apis  []model.Api `json:"apis"`
}

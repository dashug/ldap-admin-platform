package response

import "github.com/dashug/ldap-admin-platform/model"

type UserListRsp struct {
	Total int          `json:"total"`
	Users []model.User `json:"users"`
}

// SyncPreviewRsp 同步预览（Dry Run）返回：新增/更新/跳过数量，不落库
type SyncPreviewRsp struct {
	AddCount    int      `json:"addCount"`    // 将新增到 LDAP 的数量
	UpdateCount int      `json:"updateCount"` // LDAP 中已存在、将更新的数量
	SkipCount   int      `json:"skipCount"`   // 跳过的数量（可选）
	AddList     []string `json:"addList"`     // 将新增的用户名/部门名列表（便于核对）
	UpdateList  []string `json:"updateList"`   // 将更新的列表
}

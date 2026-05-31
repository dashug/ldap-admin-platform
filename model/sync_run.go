package model

import (
	"gorm.io/gorm"
)

// SyncRun 同步运行记录：定时/手动从外部源（AD/钉钉/企微/飞书）同步的每次结果
type SyncRun struct {
	gorm.Model
	Source   string `gorm:"type:varchar(20);comment:'来源 ldap/dingtalk/wecom/feishu'" json:"source"`
	Trigger  string `gorm:"type:varchar(10);comment:'触发方式 auto/manual'" json:"trigger"`
	Success  bool   `gorm:"comment:'是否成功'" json:"success"`
	Message  string `gorm:"type:varchar(500);comment:'结果信息'" json:"message"`
	Duration int64  `gorm:"comment:'耗时(ms)'" json:"duration"`
}

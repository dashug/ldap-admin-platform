package model

import (
	"gorm.io/gorm"
)

// WebhookDelivery Webhook 投递记录：每次回调（含重试后的最终结果）落一条，便于排查与审计
type WebhookDelivery struct {
	gorm.Model
	Event      string `gorm:"type:varchar(50);comment:'事件类型'" json:"event"`
	URL        string `gorm:"type:varchar(2048);comment:'目标地址'" json:"url"`
	Payload    string `gorm:"type:text;comment:'请求体'" json:"payload"`
	Success    bool   `gorm:"comment:'是否投递成功(2xx)'" json:"success"`
	StatusCode int    `gorm:"comment:'最后一次响应状态码'" json:"statusCode"`
	Attempts   int    `gorm:"comment:'尝试次数'" json:"attempts"`
	Error      string `gorm:"type:varchar(500);comment:'错误信息'" json:"error"`
}

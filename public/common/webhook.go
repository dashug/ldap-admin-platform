package common

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/dashug/ldap-admin-platform/config"
)

// WebhookEvent 回调事件类型
const (
	EventUserCreated  = "user.created"
	EventUserUpdated  = "user.updated"
	EventGroupCreated = "group.created"
	EventGroupUpdated = "group.updated"
	EventUserSync     = "user.sync"
	EventGroupSync    = "group.sync"
)

// WebhookPayload 回调请求体
type WebhookPayload struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
	Time  string      `json:"time"`
}

// UserWebhookData 回调中的用户数据（不含密码）
type UserWebhookData struct {
	ID uint `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Mail     string `json:"mail"`
	Mobile   string `json:"mobile"`
	UserDN   string `json:"userDn"`
	Status   uint   `json:"status"`
}

// GroupWebhookData 回调中的部门数据
type GroupWebhookData struct {
	ID        uint   `json:"id"`
	GroupName string `json:"groupName"`
	GroupDN   string `json:"groupDn"`
	Remark    string `json:"remark"`
}

// SendWebhook 异步发送 Webhook（不阻塞主流程）
func SendWebhook(event string, data interface{}) {
	url := ""
	if config.Conf.System != nil {
		url = config.Conf.System.WebhookURL
	}
	if url == "" {
		return
	}
	payload := WebhookPayload{
		Event: event,
		Data:  data,
		Time:  time.Now().Format(time.RFC3339),
	}
	body, err := json.Marshal(payload)
	if err != nil {
		Log.Warnf("webhook marshal: %v", err)
		return
	}
	go func() {
		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Post(url, "application/json", bytes.NewReader(body))
		if err != nil {
			Log.Warnf("webhook post %s: %v", url, err)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			Log.Warnf("webhook %s status: %d", url, resp.StatusCode)
		}
	}()
}

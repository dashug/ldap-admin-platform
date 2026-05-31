package common

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dashug/ldap-admin-platform/config"
	"github.com/dashug/ldap-admin-platform/model"
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

// webhookMaxAttempts 投递最大尝试次数（含首次）
const webhookMaxAttempts = 3

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

// signWebhook 用 HMAC-SHA256 对 (timestamp + "." + body) 签名，返回十六进制串。
// 接收端可用同一密钥验签：hex(hmac_sha256(secret, ts + "." + rawBody)) == signature。
func signWebhook(secret string, ts string, body []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(ts))
	mac.Write([]byte("."))
	mac.Write(body)
	return hex.EncodeToString(mac.Sum(nil))
}

// deliverWebhook 同步投递一次（含重试与退避），返回最终状态码、尝试次数与错误。
// 调用方负责记录投递日志 / 决定是否异步执行。
func deliverWebhook(url, event string, body []byte) (statusCode int, attempts int, finalErr error) {
	secret := ""
	if config.Conf.System != nil {
		secret = config.Conf.System.WebhookSecret
	}
	ts := strconv.FormatInt(time.Now().Unix(), 10)

	for attempts = 1; attempts <= webhookMaxAttempts; attempts++ {
		req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
		if err != nil {
			return 0, attempts, err
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Webhook-Event", event)
		req.Header.Set("X-Webhook-Timestamp", ts)
		if secret != "" {
			req.Header.Set("X-Webhook-Signature", "sha256="+signWebhook(secret, ts, body))
		}

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			finalErr = err
		} else {
			statusCode = resp.StatusCode
			resp.Body.Close()
			if statusCode >= 200 && statusCode < 300 {
				return statusCode, attempts, nil
			}
			finalErr = fmt.Errorf("响应状态码 %d", statusCode)
		}
		// 还有重试机会则退避（1s、2s…），最后一次不再等待
		if attempts < webhookMaxAttempts {
			time.Sleep(time.Duration(attempts) * time.Second)
		}
	}
	// 循环结束 attempts 会比上限多 1，修正为实际尝试次数
	attempts = webhookMaxAttempts
	return statusCode, attempts, finalErr
}

// recordDelivery 落一条投递记录（DB 未就绪时静默跳过）
func recordDelivery(event, url string, body []byte, statusCode, attempts int, err error) {
	if DB == nil {
		return
	}
	rec := model.WebhookDelivery{
		Event:      event,
		URL:        url,
		Payload:    string(body),
		Success:    err == nil,
		StatusCode: statusCode,
		Attempts:   attempts,
	}
	if err != nil {
		msg := err.Error()
		if len(msg) > 500 {
			msg = msg[:500]
		}
		rec.Error = msg
	}
	if e := DB.Create(&rec).Error; e != nil {
		Log.Warnf("webhook 投递记录写入失败: %v", e)
	}
}

// SendWebhook 异步发送 Webhook（不阻塞主流程），带签名、重试与投递记录
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
		statusCode, attempts, derr := deliverWebhook(url, event, body)
		if derr != nil {
			Log.Warnf("webhook post %s 失败(%d 次尝试): %v", url, attempts, derr)
		}
		recordDelivery(event, url, body, statusCode, attempts, derr)
	}()
}

// TestWebhook 同步发送一条测试回调并返回结果（用于「测试 Webhook」按钮）。
// url 为空时回落到已保存配置。同样带签名，并记录一条投递日志。
func TestWebhook(url string) error {
	if url == "" && config.Conf.System != nil {
		url = config.Conf.System.WebhookURL
	}
	if url == "" {
		return fmt.Errorf("未配置 Webhook 地址")
	}
	payload := WebhookPayload{
		Event: "test",
		Data:  map[string]string{"message": "这是一条来自 LDAP 管理平台的测试回调"},
		Time:  time.Now().Format(time.RFC3339),
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	statusCode, attempts, derr := deliverWebhook(url, "test", body)
	recordDelivery("test", url, body, statusCode, attempts, derr)
	return derr
}

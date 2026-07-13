package middleware

import (
	"time"

	"github.com/dashug/ldap-admin-platform/model/response"

	"github.com/gin-gonic/gin"
	gocache "github.com/patrickmn/go-cache"
)

// loginAttemptCache 记录各客户端 IP 在窗口内对敏感端点的访问次数。
var loginAttemptCache = gocache.New(15*time.Minute, 5*time.Minute)

const (
	sensitiveMaxAttempts = 10          // 每个 IP 在窗口内允许的最大次数
	sensitiveWindow      = time.Minute // 计数窗口
)

// SensitiveRateLimit 对敏感且无鉴权的端点（登录 / 发送验证码 / 改密）按客户端 IP 限流，
// 缓解口令爆破与验证码轰炸。与全局令牌桶（进程级共享，易被单点打满）不同，这里按 IP 独立
// 计数，单个 IP 触顶不会影响其他用户的正常访问。
func SensitiveRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := "sens:" + c.ClientIP()
		if v, found := loginAttemptCache.Get(key); found {
			if n, ok := v.(int); ok && n >= sensitiveMaxAttempts {
				response.Fail(c, nil, "操作过于频繁，请稍后再试")
				c.Abort()
				return
			}
			// 保持原有窗口过期时间，仅递增计数
			_, _ = loginAttemptCache.IncrementInt(key, 1)
		} else {
			loginAttemptCache.Set(key, 1, sensitiveWindow)
		}
		c.Next()
	}
}

package middleware

import (
	"net/http"

	"github.com/dashug/ldap-admin-platform/config"
	"github.com/gin-gonic/gin"
)

// 允许跨域携带的请求头
const corsAllowHeaders = "Authorization, Content-Length, X-CSRF-Token, Token, session, X_Requested_With, " +
	"Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language, DNT, X-CustomHeader, Keep-Alive, " +
	"User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma, X-API-Key"

// CORS跨域中间件
//
// 安全策略（由 config.system.allow-origins 控制）：
//   - 为空（默认）：单二进制同源部署，不下发任何跨域响应头，浏览器按同源策略拒绝跨站请求；
//   - ["*"]      ：放行所有来源（兼容旧行为，仅在确有需要时使用，且不允许携带 Cookie 凭证）；
//   - 指定来源    ：仅对命中白名单的 Origin 回显该来源，并允许携带凭证。
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		allowed, allowCredentials := resolveAllowOrigin(origin)

		if allowed != "" {
			c.Header("Access-Control-Allow-Origin", allowed)
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
			c.Header("Access-Control-Allow-Headers", corsAllowHeaders)
			c.Header("Access-Control-Expose-Headers", "Content-Length, Content-Type")
			c.Header("Access-Control-Max-Age", "172800")
			if allowCredentials {
				c.Header("Access-Control-Allow-Credentials", "true")
			}
			// 按 Origin 变化缓存，避免 CDN/代理把某一来源的响应头串给其他来源
			c.Header("Vary", "Origin")
		}

		// 预检请求直接放行结束
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// resolveAllowOrigin 根据白名单决定 Access-Control-Allow-Origin 的取值。
// 返回空字符串表示不下发跨域头（同源场景）。
func resolveAllowOrigin(origin string) (allowOrigin string, allowCredentials bool) {
	if origin == "" {
		return "", false
	}
	var allowList []string
	if config.Conf.System != nil {
		allowList = config.Conf.System.AllowOrigins
	}
	for _, o := range allowList {
		if o == "*" {
			// 通配模式：不允许携带凭证（浏览器禁止 * 与 credentials 同时出现）
			return "*", false
		}
		if o == origin {
			return origin, true
		}
	}
	return "", false
}

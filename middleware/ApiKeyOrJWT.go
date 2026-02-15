package middleware

import (
	"net/http"
	"strings"

	"github.com/dashug/ldap-admin-platform/model/response"
	"github.com/dashug/ldap-admin-platform/service/isql"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// ApiKeyOrJWT 返回一个认证中间件：优先校验 X-API-Key，通过则设置身份用户并放行；否则走 JWT
func ApiKeyOrJWT(authMiddleware *jwt.GinJWTMiddleware) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := strings.TrimSpace(c.GetHeader("X-API-Key"))
		if key != "" {
			_, ok := isql.ApiKey.Verify(key)
			if !ok {
				response.Response(c, http.StatusUnauthorized, http.StatusUnauthorized, nil, "API Key 无效或已失效")
				c.Abort()
				return
			}
			apiUser, err := isql.User.GetApiKeyIdentityUser()
			if err != nil {
				response.Response(c, http.StatusInternalServerError, http.StatusInternalServerError, nil, "获取 API 身份用户失败")
				c.Abort()
				return
			}
			c.Set("user", apiUser)
			c.Next()
			return
		}
		authMiddleware.MiddlewareFunc()(c)
	}
}

package routes

import (
	"github.com/dashug/ldap-admin-platform/controller"
	"github.com/dashug/ldap-admin-platform/middleware"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// InitApiKeyRoutes 注册 API 密钥管理路由（需登录或有效 API Key）
func InitApiKeyRoutes(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) gin.IRoutes {
	apiKey := r.Group("/apiKey")
	apiKey.Use(middleware.ApiKeyOrJWT(authMiddleware))
	apiKey.Use(middleware.CasbinMiddleware())
	{
		apiKey.GET("/list", controller.ApiKey.List)
		apiKey.POST("/create", controller.ApiKey.Create)
		apiKey.POST("/delete", controller.ApiKey.Delete)
	}
	return r
}

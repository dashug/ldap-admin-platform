package routes

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/dashug/ldap-admin-platform/controller"
	"github.com/gin-gonic/gin"
)

// LoginHandler
// @Summary 登录接口 (手动加上: Bearer + token(密码加密接口))
// @Description 用户登录
// @Tags 基础管理
// @Accept application/json
// @Produce application/json
// @Param  data body request.RegisterAndLoginReq true "用户登录信息账号和密码"
// @Success 200 {object} response.ResponseBody
// @Router /base/login [post]
func LoginHandler() {}

// LogoutHandler
// @Summary 退出登录
// @Description 用户退出登录
// @Tags 基础管理
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.ResponseBody
// @Router /base/logout [post]
func LogoutHandler() {
}

// RefreshHandler
// @Summary 刷新 Token
// @Description 使用旧的 Token 获取新的 Token
// @Tags 基础管理
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 旧的 Token"
// @Success 200 {object} response.ResponseBody
// @Router /base/refreshToken [post]
func RefreshHandler() {

}

// 注册基础路由
func InitBaseRoutes(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) gin.IRoutes {
	base := r.Group("/base")
	{
		base.GET("ping", controller.Demo)
		// 生成加密密码：属于运维/调试工具，必须登录后才能使用，避免成为离线攻击辅助接口
		base.GET("encryptpwd", authMiddleware.MiddlewareFunc(), controller.Base.EncryptPasswd)
		// 说明：原 decryptpwd（把密文还原为明文）已移除——它是一个无鉴权的明文还原 oracle，
		// 配合数据库中可逆存储的密码会造成明文泄露，且前后端均未使用。
		base.GET("publicKey", controller.Base.GetPublicKey) // 获取登录用 RSA 公钥（无需鉴权）
		base.GET("config", controller.Base.GetConfig)         // 获取系统配置
		base.GET("ldapStatus", controller.Base.GetLDAPStatus)   // LDAP 连接状态
		base.GET("systemInfo", controller.Base.GetSystemInfo)  // 系统信息（版本、运行时长、DB）
		base.GET("version", controller.Base.GetVersion)       // 获取版本信息
		// 登录登出刷新token无需鉴权
		base.POST("/login", authMiddleware.LoginHandler)
		base.POST("/logout", authMiddleware.LogoutHandler)
		base.POST("/refreshToken", authMiddleware.RefreshHandler)
		base.POST("/sendcode", controller.Base.SendCode)   // 给用户邮箱发送验证码
		base.POST("/changePwd", controller.Base.ChangePwd) // 修改用户密码
		base.GET("/dashboard", controller.Base.Dashboard)  // 系统首页展示数据
		base.POST("/directoryConfig", authMiddleware.MiddlewareFunc(), controller.Base.UpdateDirectoryConfig)
		base.POST("/directoryConfig/test", authMiddleware.MiddlewareFunc(), controller.Base.TestDirectoryConfig)
		base.POST("/configImport", authMiddleware.MiddlewareFunc(), controller.Base.ImportConfig)
		base.POST("/thirdPartyConfig", authMiddleware.MiddlewareFunc(), controller.Base.UpdateThirdPartyConfig)
		base.POST("/thirdPartyConfig/test", authMiddleware.MiddlewareFunc(), controller.Base.TestThirdPartyConfig)
		base.POST("/emailConfig", authMiddleware.MiddlewareFunc(), controller.Base.UpdateEmailConfig)
		base.POST("/emailConfig/test", authMiddleware.MiddlewareFunc(), controller.Base.TestNotification)
		base.GET("/webhookDeliveries", authMiddleware.MiddlewareFunc(), controller.Base.ListWebhookDeliveries)
		base.POST("/syncConfig", authMiddleware.MiddlewareFunc(), controller.Base.UpdateSyncConfig)
		base.POST("/syncRun", authMiddleware.MiddlewareFunc(), controller.Base.RunSyncNow)
		base.GET("/syncRuns", authMiddleware.MiddlewareFunc(), controller.Base.ListSyncRuns)
		base.POST("/userBatchImport", authMiddleware.MiddlewareFunc(), controller.Base.UserBatchImport)
		base.GET("/mfa/status", authMiddleware.MiddlewareFunc(), controller.Base.MfaStatus)
		base.POST("/mfa/setup", authMiddleware.MiddlewareFunc(), controller.Base.MfaSetup)
		base.POST("/mfa/verify", authMiddleware.MiddlewareFunc(), controller.Base.MfaVerify)
		base.POST("/mfa/disable", authMiddleware.MiddlewareFunc(), controller.Base.MfaDisable)
	}
	return r
}

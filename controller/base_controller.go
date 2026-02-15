package controller

import (
	"github.com/dashug/ldap-admin-platform/logic"
	"github.com/dashug/ldap-admin-platform/model/request"
	"github.com/dashug/ldap-admin-platform/public/tools"

	"github.com/gin-gonic/gin"
)

type BaseController struct{}

// SendCode 给用户邮箱发送验证码
// @Summary 发送验证码
// @Description 向指定邮箱发送验证码
// @Tags 基础管理
// @Accept application/json
// @Produce application/json
// @Param data body request.BaseSendCodeReq true "发送验证码请求数据"
// @Success 200 {object} response.ResponseBody
// @Router /base/sendcode [post]
func (m *BaseController) SendCode(c *gin.Context) {
	req := new(request.BaseSendCodeReq)
	Run(c, req, func() (any, any) {
		return logic.Base.SendCode(c, req)
	})
}

// ChangePwd 用户通过邮箱修改密码
// @Summary 用户通过邮箱修改密码
// @Description 使用邮箱验证码修改密码
// @Tags 基础管理
// @Accept application/json
// @Produce application/json
// @Param  data body request.BaseChangePwdReq true "发送验证码请求数据"
// @Success 200 {object} response.ResponseBody
// @Router /base/changePwd [post]
func (m *BaseController) ChangePwd(c *gin.Context) {
	req := new(request.BaseChangePwdReq)
	Run(c, req, func() (any, any) {
		return logic.Base.ChangePwd(c, req)
	})
}

// Dashboard 系统首页展示数据
// @Summary 获取仪表盘数据
// @Description 获取系统仪表盘概览数据
// @Tags 基础管理
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.ResponseBody
// @Router /base/dashboard [get]
func (m *BaseController) Dashboard(c *gin.Context) {
	req := new(request.BaseDashboardReq)
	Run(c, req, func() (any, any) {
		return logic.Base.Dashboard(c, req)
	})
}

// EncryptPasswd 密码加密
// @Summary 密码加密
// @Description 将明文密码加密
// @Tags 基础管理
// @Accept application/json
// @Produce application/json
// @Param passwd query string true "需要加密的明文密码"
// @Success 200 {object} response.ResponseBody
// @Router /base/encryptpwd [get]
func (m *BaseController) EncryptPasswd(c *gin.Context) {
	req := new(request.EncryptPasswdReq)
	Run(c, req, func() (any, any) {
		return logic.Base.EncryptPasswd(c, req)
	})
}

// DecryptPasswd 密码解密为明文
// @Summary 密码解密
// @Description 将加密后的密码解密为明文
// @Tags 基础管理
// @Accept application/json
// @Produce application/json
// @Param passwd query string true "需要解密的加密密码"
// @Success 200 {object} response.ResponseBody
// @Router /base/decryptpwd [get]
func (m *BaseController) DecryptPasswd(c *gin.Context) {
	req := new(request.DecryptPasswdReq)
	Run(c, req, func() (any, any) {
		return logic.Base.DecryptPasswd(c, req)
	})
}

// GetPublicKey 获取登录用 RSA 公钥（供前端加密密码，无需鉴权）
// @Summary 获取 RSA 公钥
// @Description 返回用于登录密码加密的 RSA 公钥 PEM，前端未配置 VUE_APP_PUBLIC_KEY 时可由此接口获取以保证与后端一致
// @Tags 基础管理
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.ResponseBody
// @Router /base/publicKey [get]
func (m *BaseController) GetPublicKey(c *gin.Context) {
	data, err := logic.Base.GetPublicKey(c)
	if err != nil {
		tools.Err(c, tools.ReloadErr(err), nil)
		return
	}
	tools.Success(c, data)
}

// GetConfig 获取系统配置
// @Summary 获取系统配置
// @Description 获取系统配置信息，用于前端判断是否显示同步按钮
// @Tags 基础管理
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.ResponseBody
// @Router /base/config [get]
func (m *BaseController) GetConfig(c *gin.Context) {
	req := new(request.BaseConfigReq)
	Run(c, req, func() (any, any) {
		return logic.Base.GetConfig(c, req)
	})
}

// UpdateDirectoryConfig 更新目录配置
// @Summary 更新目录配置
// @Description 更新 OpenLDAP/AD 目录服务连接参数
// @Tags 基础管理
// @Accept application/json
// @Produce application/json
// @Param data body request.BaseUpdateDirectoryConfigReq true "目录配置"
// @Success 200 {object} response.ResponseBody
// @Router /base/directoryConfig [post]
func (m *BaseController) UpdateDirectoryConfig(c *gin.Context) {
	req := new(request.BaseUpdateDirectoryConfigReq)
	Run(c, req, func() (any, any) {
		return logic.Base.UpdateDirectoryConfig(c, req)
	})
}

// ImportConfig 导入配置（目录与同步规则 JSON）
// @Summary 导入配置
// @Description 从导出的 JSON 恢复目录与同步规则
// @Tags 基础管理
// @Accept application/json
// @Produce application/json
// @Param data body request.BaseConfigImportReq true "导出的配置 JSON"
// @Success 200 {object} response.ResponseBody
// @Router /base/configImport [post]
func (m *BaseController) ImportConfig(c *gin.Context) {
	req := new(request.BaseConfigImportReq)
	Run(c, req, func() (any, any) {
		return logic.Base.ImportConfig(c, req)
	})
}

// UpdateThirdPartyConfig 更新第三方平台配置
// @Summary 更新第三方平台配置
// @Description 更新飞书、企微、钉钉对接参数
// @Tags 基础管理
// @Accept application/json
// @Produce application/json
// @Param data body request.BaseThirdPartyConfigReq true "平台配置"
// @Success 200 {object} response.ResponseBody
// @Router /base/thirdPartyConfig [post]
func (m *BaseController) UpdateThirdPartyConfig(c *gin.Context) {
	req := new(request.BaseThirdPartyConfigReq)
	Run(c, req, func() (any, any) {
		return logic.Base.UpdateThirdPartyConfig(c, req)
	})
}

// TestThirdPartyConfig 测试第三方平台配置
// @Summary 测试第三方平台配置
// @Description 测试飞书、企微、钉钉连接是否可用
// @Tags 基础管理
// @Accept application/json
// @Produce application/json
// @Param data body request.BaseThirdPartyConfigReq true "平台配置"
// @Success 200 {object} response.ResponseBody
// @Router /base/thirdPartyConfig/test [post]
func (m *BaseController) TestThirdPartyConfig(c *gin.Context) {
	req := new(request.BaseThirdPartyConfigReq)
	Run(c, req, func() (any, any) {
		return logic.Base.TestThirdPartyConfig(c, req)
	})
}

// UpdateEmailConfig 更新邮件通知配置
// @Summary 更新邮件通知配置
// @Description 如是否在新建/同步用户时发送通知邮件
// @Tags 基础管理
// @Accept application/json
// @Produce application/json
// @Param data body request.BaseUpdateEmailConfigReq true "邮件配置"
// @Success 200 {object} response.ResponseBody
// @Router /base/emailConfig [post]
func (m *BaseController) UpdateEmailConfig(c *gin.Context) {
	req := new(request.BaseUpdateEmailConfigReq)
	Run(c, req, func() (any, any) {
		return logic.Base.UpdateEmailConfig(c, req)
	})
}

// GetVersion 获取版本信息
// @Summary 获取版本信息
// @Description 获取系统版本号、Git提交哈希和构建时间
// @Tags 基础管理
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.ResponseBody
// @Router /base/version [get]
func (m *BaseController) GetVersion(c *gin.Context) {
	req := new(request.BaseVersionReq)
	Run(c, req, func() (any, any) {
		return logic.Base.GetVersion(c, req)
	})
}

// GetLDAPStatus 获取 LDAP 连接状态
// @Summary 获取 LDAP 连接状态
// @Description 使用当前目录配置探测 LDAP 是否可连接
// @Tags 基础管理
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.ResponseBody
// @Router /base/ldapStatus [get]
func (m *BaseController) GetLDAPStatus(c *gin.Context) {
	req := new(request.BaseLDAPStatusReq)
	Run(c, req, func() (any, any) {
		return logic.Base.GetLDAPStatus(c, req)
	})
}

// GetSystemInfo 获取系统信息
// @Summary 获取系统信息
// @Description 版本、运行时长、数据库状态，供管理员查看
// @Tags 基础管理
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.ResponseBody
// @Router /base/systemInfo [get]
func (m *BaseController) GetSystemInfo(c *gin.Context) {
	req := new(request.BaseSystemInfoReq)
	Run(c, req, func() (any, any) {
		return logic.Base.GetSystemInfo(c, req)
	})
}

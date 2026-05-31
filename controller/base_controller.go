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

// TestDirectoryConfig 测试目录（LDAP）连接
// @Summary 测试目录连接
// @Description 使用提交的目录参数拨号并绑定，验证连接是否可用（不保存）
// @Tags 基础管理
// @Accept application/json
// @Produce application/json
// @Param data body request.BaseTestDirectoryConfigReq true "目录连接参数"
// @Success 200 {object} response.ResponseBody
// @Router /base/directoryConfig/test [post]
func (m *BaseController) TestDirectoryConfig(c *gin.Context) {
	req := new(request.BaseTestDirectoryConfigReq)
	Run(c, req, func() (any, any) {
		return logic.Base.TestDirectoryConfig(c, req)
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

// TestNotification 测试通知（邮件 / Webhook）
// @Summary 测试通知配置
// @Description 发送一封测试邮件，或向 Webhook 地址发送一条测试回调，验证配置是否可用（不保存）
// @Tags 基础管理
// @Accept application/json
// @Produce application/json
// @Param data body request.BaseTestNotificationReq true "测试通知参数"
// @Success 200 {object} response.ResponseBody
// @Router /base/emailConfig/test [post]
func (m *BaseController) TestNotification(c *gin.Context) {
	req := new(request.BaseTestNotificationReq)
	Run(c, req, func() (any, any) {
		return logic.Base.TestNotification(c, req)
	})
}

// ListWebhookDeliveries 查询 Webhook 投递记录
// @Summary 查询 Webhook 投递记录
// @Description 分页查询回调投递历史（含签名/重试后的最终结果），供排查
// @Tags 基础管理
// @Accept application/json
// @Produce application/json
// @Param data query request.BaseWebhookDeliveriesReq true "分页参数"
// @Success 200 {object} response.ResponseBody
// @Router /base/webhookDeliveries [get]
func (m *BaseController) ListWebhookDeliveries(c *gin.Context) {
	req := new(request.BaseWebhookDeliveriesReq)
	Run(c, req, func() (any, any) {
		return logic.Base.ListWebhookDeliveries(c, req)
	})
}

// UpdateSyncConfig 更新定时自动同步配置
// @Summary 更新定时自动同步配置
// @Tags 基础管理
// @Accept application/json
// @Produce application/json
// @Param data body request.BaseUpdateSyncConfigReq true "定时同步配置"
// @Success 200 {object} response.ResponseBody
// @Router /base/syncConfig [post]
func (m *BaseController) UpdateSyncConfig(c *gin.Context) {
	req := new(request.BaseUpdateSyncConfigReq)
	Run(c, req, func() (any, any) {
		return logic.Base.UpdateSyncConfig(c, req)
	})
}

// RunSyncNow 立即触发某来源同步
// @Summary 立即同步
// @Tags 基础管理
// @Accept application/json
// @Produce application/json
// @Param data body request.BaseRunSyncReq true "同步来源"
// @Success 200 {object} response.ResponseBody
// @Router /base/syncRun [post]
func (m *BaseController) RunSyncNow(c *gin.Context) {
	req := new(request.BaseRunSyncReq)
	Run(c, req, func() (any, any) {
		return logic.Base.RunSyncNow(c, req)
	})
}

// ListSyncRuns 查询同步运行记录
// @Summary 查询同步运行记录
// @Tags 基础管理
// @Accept application/json
// @Produce application/json
// @Param data query request.BaseSyncRunsReq true "分页参数"
// @Success 200 {object} response.ResponseBody
// @Router /base/syncRuns [get]
func (m *BaseController) ListSyncRuns(c *gin.Context) {
	req := new(request.BaseSyncRunsReq)
	Run(c, req, func() (any, any) {
		return logic.Base.ListSyncRuns(c, req)
	})
}

// UserBatchImport 批量导入用户
// @Summary 批量导入用户
// @Description CSV/Excel 解析后的用户行批量创建，dryRun 仅校验不落库
// @Tags 基础管理
// @Accept application/json
// @Produce application/json
// @Param data body request.UserBatchImportReq true "导入用户列表"
// @Success 200 {object} response.ResponseBody
// @Router /base/userBatchImport [post]
func (m *BaseController) UserBatchImport(c *gin.Context) {
	req := new(request.UserBatchImportReq)
	Run(c, req, func() (any, any) {
		return logic.User.BatchImport(c, req)
	})
}

// MfaStatus 查询当前用户 MFA 状态
// @Tags 基础管理
// @Success 200 {object} response.ResponseBody
// @Router /base/mfa/status [get]
func (m *BaseController) MfaStatus(c *gin.Context) {
	req := new(request.BaseMfaStatusReq)
	Run(c, req, func() (any, any) {
		return logic.Base.MfaStatus(c, req)
	})
}

// MfaSetup 生成 MFA 密钥与二维码
// @Tags 基础管理
// @Success 200 {object} response.ResponseBody
// @Router /base/mfa/setup [post]
func (m *BaseController) MfaSetup(c *gin.Context) {
	req := new(request.BaseMfaSetupReq)
	Run(c, req, func() (any, any) {
		return logic.Base.MfaSetup(c, req)
	})
}

// MfaVerify 校验验证码并启用 MFA
// @Tags 基础管理
// @Param data body request.BaseMfaCodeReq true "验证码"
// @Success 200 {object} response.ResponseBody
// @Router /base/mfa/verify [post]
func (m *BaseController) MfaVerify(c *gin.Context) {
	req := new(request.BaseMfaCodeReq)
	Run(c, req, func() (any, any) {
		return logic.Base.MfaVerify(c, req)
	})
}

// MfaDisable 关闭 MFA
// @Tags 基础管理
// @Param data body request.BaseMfaCodeReq true "验证码"
// @Success 200 {object} response.ResponseBody
// @Router /base/mfa/disable [post]
func (m *BaseController) MfaDisable(c *gin.Context) {
	req := new(request.BaseMfaCodeReq)
	Run(c, req, func() (any, any) {
		return logic.Base.MfaDisable(c, req)
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

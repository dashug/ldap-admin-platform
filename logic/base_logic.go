package logic

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image/png"
	"strings"

	"github.com/chyroc/lark"
	"github.com/pquerna/otp/totp"
	"github.com/dashug/ldap-admin-platform/config"
	"github.com/dashug/ldap-admin-platform/model"
	"github.com/dashug/ldap-admin-platform/model/request"
	"github.com/dashug/ldap-admin-platform/model/response"
	"github.com/dashug/ldap-admin-platform/public/common"
	"github.com/dashug/ldap-admin-platform/public/tools"
	"github.com/dashug/ldap-admin-platform/public/version"
	"github.com/dashug/ldap-admin-platform/service/ildap"
	"github.com/dashug/ldap-admin-platform/service/isql"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	wecomsdk "github.com/wenerme/go-wecom/wecom"
	dingsdk "github.com/zhaoyunxing92/dingtalk/v2"
)

type BaseLogic struct{}

// SendCode 发送验证码
func (l BaseLogic) SendCode(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.BaseSendCodeReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c
	// 判断邮箱是否正确
	user := new(model.User)
	err := isql.User.Find(tools.H{"mail": r.Mail}, user)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("%s", "通过邮箱查询用户失败"+err.Error()))
	}
	if user.Status != 1 || user.SyncState != 1 {
		return nil, tools.NewMySqlError(fmt.Errorf("该用户已离职或者未同步在ldap，无法重置密码，如有疑问，请联系管理员"))
	}
	err = tools.SendCode([]string{r.Mail})
	if err != nil {
		return nil, tools.NewLdapError(fmt.Errorf("%s", "邮件发送失败"+err.Error()))
	}

	return nil, nil
}

// ChangePwd 重置密码
func (l BaseLogic) ChangePwd(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.BaseChangePwdReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c
	// 判断邮箱是否正确
	if !isql.User.Exist(tools.H{"mail": r.Mail}) {
		return nil, tools.NewValidatorError(fmt.Errorf("邮箱不存在,请检查邮箱是否正确"))
	}
	// 判断验证码是否过期
	cacheCode, ok := tools.VerificationCodeCache.Get(r.Mail)
	if !ok {
		return nil, tools.NewValidatorError(fmt.Errorf("对不起，该验证码已超过5分钟有效期，请重新重新密码"))
	}
	// 判断验证码是否正确
	if cacheCode != r.Code {
		return nil, tools.NewValidatorError(fmt.Errorf("验证码错误，请检查邮箱中正确的验证码，如果点击多次发送验证码，请用最后一次生成的验证码来验证"))
	}

	user := new(model.User)
	err := isql.User.Find(tools.H{"mail": r.Mail}, user)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("%s", "通过邮箱查询用户失败"+err.Error()))
	}

	newpass, err := ildap.User.NewPwd(user.Username)
	if err != nil {
		return nil, tools.NewLdapError(fmt.Errorf("%s", "LDAP生成新密码失败"+err.Error()))
	}

	err = tools.SendMail([]string{user.Mail}, newpass)
	if err != nil {
		return nil, tools.NewLdapError(fmt.Errorf("%s", "邮件发送失败"+err.Error()))
	}

	// 更新数据库密码
	err = isql.User.ChangePwd(user.Username, tools.NewGenPasswd(newpass))
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("%s", "在MySQL更新密码失败: "+err.Error()))
	}

	return nil, nil
}

// Dashboard 仪表盘
func (l BaseLogic) Dashboard(c *gin.Context, req any) (data any, rspError any) {
	_, ok := req.(*request.BaseDashboardReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	userCount, err := isql.User.Count()
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取用户总数失败"))
	}
	groupCount, err := isql.Group.Count()
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取分组总数失败"))
	}
	roleCount, err := isql.Role.Count()
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取角色总数失败"))
	}
	menuCount, err := isql.Menu.Count()
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取菜单总数失败"))
	}
	apiCount, err := isql.Api.Count()
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取接口总数失败"))
	}
	logCount, err := isql.OperationLog.Count()
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取日志总数失败"))
	}

	rst := make([]*response.DashboardList, 0)

	rst = append(rst,
		&response.DashboardList{
			DataType:  "user",
			DataName:  "用户",
			DataCount: userCount,
			Icon:      "people",
			Path:      "#/personnel/user",
		},
		&response.DashboardList{
			DataType:  "group",
			DataName:  "分组",
			DataCount: groupCount,
			Icon:      "peoples",
			Path:      "#/personnel/group",
		},
		&response.DashboardList{
			DataType:  "role",
			DataName:  "角色",
			DataCount: roleCount,
			Icon:      "eye-open",
			Path:      "#/system/role",
		},
		&response.DashboardList{
			DataType:  "menu",
			DataName:  "菜单",
			DataCount: menuCount,
			Icon:      "tree-table",
			Path:      "#/system/menu",
		},
		&response.DashboardList{
			DataType:  "api",
			DataName:  "接口",
			DataCount: apiCount,
			Icon:      "tree",
			Path:      "#/system/api",
		},
		&response.DashboardList{
			DataType:  "log",
			DataName:  "日志",
			DataCount: logCount,
			Icon:      "documentation",
			Path:      "#/log/operation-log",
		},
	)

	return rst, nil
}

// EncryptPasswd
func (l BaseLogic) EncryptPasswd(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.EncryptPasswdReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	return tools.NewGenPasswd(r.Passwd), nil
}

// GetPublicKey 返回登录用 RSA 公钥 PEM，供前端加密密码（与后端密钥对一致）
func (l BaseLogic) GetPublicKey(c *gin.Context) (data any, err error) {
	if config.Conf.System == nil || len(config.Conf.System.RSAPublicBytes) == 0 {
		return nil, fmt.Errorf("系统未配置 RSA 公钥")
	}
	return string(config.Conf.System.RSAPublicBytes), nil
}

// GetConfig 获取系统配置
func (l BaseLogic) GetConfig(c *gin.Context, req any) (data any, rspError any) {
	_, ok := req.(*request.BaseConfigReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	// 安全获取配置值，防止配置段缺失导致空指针
	rsp := &response.BaseConfigRsp{}
	if config.Conf.Ldap != nil {
		rsp.LdapEnableSync = config.Conf.Ldap.EnableSync
		rsp.DirectoryType = strings.TrimSpace(config.Conf.Ldap.DirectoryType)
		if rsp.DirectoryType == "" {
			rsp.DirectoryType = "openldap"
		}
		rsp.Url = config.Conf.Ldap.Url
		rsp.BaseDN = config.Conf.Ldap.BaseDN
		rsp.AdminDN = config.Conf.Ldap.AdminDN
		// 出于安全考虑，不回传管理员密码明文
		rsp.AdminPass = ""
		rsp.UserDN = config.Conf.Ldap.UserDN
		rsp.UserInitPassword = config.Conf.Ldap.UserInitPassword
		rsp.DefaultEmailSuffix = config.Conf.Ldap.DefaultEmailSuffix
		rsp.SyncUsernameRule = strings.TrimSpace(config.Conf.Ldap.SyncUsernameRule)
		if rsp.SyncUsernameRule == "" {
			rsp.SyncUsernameRule = "email_prefix"
		}
		rsp.SyncGroupNameRule = strings.TrimSpace(config.Conf.Ldap.SyncGroupNameRule)
		if rsp.SyncGroupNameRule == "" {
			rsp.SyncGroupNameRule = "name"
		}
		rsp.UserRDNAttr = strings.TrimSpace(config.Conf.Ldap.UserRDNAttr)
		if rsp.UserRDNAttr == "" {
			if strings.EqualFold(rsp.DirectoryType, "ad") {
				rsp.UserRDNAttr = "cn"
			} else {
				rsp.UserRDNAttr = "uid"
			}
		}
		rsp.GroupRDNAttr = strings.TrimSpace(config.Conf.Ldap.GroupRDNAttr)
		if rsp.GroupRDNAttr == "" {
			rsp.GroupRDNAttr = "cn"
		}
	}
	if config.Conf.DingTalk != nil {
		rsp.DingTalkEnableSync = config.Conf.DingTalk.EnableSync
		rsp.DingTalkFlag = config.Conf.DingTalk.Flag
		rsp.DingTalkAppKey = config.Conf.DingTalk.AppKey
		// 出于安全考虑，不回传密钥明文
		rsp.DingTalkAppSecret = ""
		rsp.DingTalkAgentID = config.Conf.DingTalk.AgentId
	}
	if config.Conf.FeiShu != nil {
		rsp.FeiShuEnableSync = config.Conf.FeiShu.EnableSync
		rsp.FeiShuFlag = config.Conf.FeiShu.Flag
		rsp.FeiShuAppID = config.Conf.FeiShu.AppID
		// 出于安全考虑，不回传密钥明文
		rsp.FeiShuAppSecret = ""
	}
	if config.Conf.WeCom != nil {
		rsp.WeComEnableSync = config.Conf.WeCom.EnableSync
		rsp.WeComFlag = config.Conf.WeCom.Flag
		rsp.WeComCorpID = config.Conf.WeCom.CorpID
		// 出于安全考虑，不回传密钥明文
		rsp.WeComCorpSecret = ""
		rsp.WeComAgentID = config.Conf.WeCom.AgentID
	}
	if config.Conf.Email != nil {
		rsp.SendUserCreationMail = config.Conf.Email.SendUserCreationMail
		rsp.SmtpHost = config.Conf.Email.Host
		rsp.SmtpPort = config.Conf.Email.Port
		rsp.SmtpUser = config.Conf.Email.User
		rsp.SmtpFrom = config.Conf.Email.From
	}
	if config.Conf.System != nil {
		rsp.WebhookURL = config.Conf.System.WebhookURL
		rsp.WebhookSecretSet = strings.TrimSpace(config.Conf.System.WebhookSecret) != ""
		rsp.AutoSyncEnabled = config.Conf.System.AutoSyncEnabled
		rsp.AutoSyncCron = config.Conf.System.AutoSyncCron
	}

	return rsp, nil
}

// ImportConfig 从 JSON 导入配置（目录与同步规则），仅更新提供的字段
func (l BaseLogic) ImportConfig(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.BaseConfigImportReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c
	dirType := strings.ToLower(strings.TrimSpace(r.DirectoryType))
	if dirType == "" {
		dirType = "openldap"
	}
	if dirType != "openldap" && dirType != "ad" {
		return nil, tools.NewValidatorError(fmt.Errorf("directoryType 仅支持 openldap 或 ad"))
	}
	if config.Conf.Ldap != nil {
		config.Conf.Ldap.DirectoryType = dirType
		if r.Url != "" {
			config.Conf.Ldap.Url = strings.TrimSpace(r.Url)
		}
		if r.BaseDN != "" {
			config.Conf.Ldap.BaseDN = strings.TrimSpace(r.BaseDN)
		}
		if r.AdminDN != "" {
			config.Conf.Ldap.AdminDN = strings.TrimSpace(r.AdminDN)
		}
		if strings.TrimSpace(r.AdminPass) != "" {
			config.Conf.Ldap.AdminPass = r.AdminPass
		}
		if r.UserDN != "" {
			config.Conf.Ldap.UserDN = strings.TrimSpace(r.UserDN)
		}
		if r.UserInitPassword != "" {
			config.Conf.Ldap.UserInitPassword = r.UserInitPassword
		}
		if r.DefaultEmailSuffix != "" {
			config.Conf.Ldap.DefaultEmailSuffix = strings.TrimSpace(r.DefaultEmailSuffix)
		}
		config.Conf.Ldap.EnableSync = r.LdapEnableSync
		if r.SyncUsernameRule != "" {
			config.Conf.Ldap.SyncUsernameRule = strings.TrimSpace(r.SyncUsernameRule)
		}
		if r.SyncGroupNameRule != "" {
			config.Conf.Ldap.SyncGroupNameRule = strings.TrimSpace(r.SyncGroupNameRule)
		}
		if r.UserRDNAttr != "" {
			config.Conf.Ldap.UserRDNAttr = strings.TrimSpace(r.UserRDNAttr)
		}
		if r.GroupRDNAttr != "" {
			config.Conf.Ldap.GroupRDNAttr = strings.TrimSpace(r.GroupRDNAttr)
		}
	}
	viper.Set("ldap.directory-type", dirType)
	if r.Url != "" {
		viper.Set("ldap.url", strings.TrimSpace(r.Url))
	}
	if r.BaseDN != "" {
		viper.Set("ldap.base-dn", strings.TrimSpace(r.BaseDN))
	}
	if r.AdminDN != "" {
		viper.Set("ldap.admin-dn", strings.TrimSpace(r.AdminDN))
	}
	if strings.TrimSpace(r.AdminPass) != "" {
		viper.Set("ldap.admin-pass", r.AdminPass)
	}
	if r.UserDN != "" {
		viper.Set("ldap.user-dn", strings.TrimSpace(r.UserDN))
	}
	if r.UserInitPassword != "" {
		viper.Set("ldap.user-init-password", r.UserInitPassword)
	}
	if r.DefaultEmailSuffix != "" {
		viper.Set("ldap.default-email-suffix", strings.TrimSpace(r.DefaultEmailSuffix))
	}
	viper.Set("ldap.enable-sync", r.LdapEnableSync)
	if r.SyncUsernameRule != "" {
		viper.Set("ldap.sync-username-rule", strings.TrimSpace(r.SyncUsernameRule))
	}
	if r.SyncGroupNameRule != "" {
		viper.Set("ldap.sync-group-name-rule", strings.TrimSpace(r.SyncGroupNameRule))
	}
	if r.UserRDNAttr != "" {
		viper.Set("ldap.user-rdn-attr", strings.TrimSpace(r.UserRDNAttr))
	}
	if r.GroupRDNAttr != "" {
		viper.Set("ldap.group-rdn-attr", strings.TrimSpace(r.GroupRDNAttr))
	}
	if err := viper.WriteConfig(); err != nil {
		return nil, tools.NewOperationError(fmt.Errorf("保存配置文件失败: %s", err.Error()))
	}
	return nil, nil
}

// UpdateDirectoryConfig 更新目录服务配置
func (l BaseLogic) UpdateDirectoryConfig(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.BaseUpdateDirectoryConfigReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	dirType := strings.ToLower(strings.TrimSpace(r.DirectoryType))
	if dirType == "" {
		dirType = "openldap"
	}
	if dirType != "openldap" && dirType != "ad" {
		return nil, tools.NewValidatorError(fmt.Errorf("directoryType 仅支持 openldap 或 ad"))
	}

	// 更新运行时配置
	if config.Conf.Ldap != nil {
		config.Conf.Ldap.DirectoryType = dirType
		config.Conf.Ldap.Url = strings.TrimSpace(r.Url)
		config.Conf.Ldap.BaseDN = strings.TrimSpace(r.BaseDN)
		config.Conf.Ldap.AdminDN = strings.TrimSpace(r.AdminDN)
		if strings.TrimSpace(r.AdminPass) != "" {
			config.Conf.Ldap.AdminPass = r.AdminPass
		}
		config.Conf.Ldap.UserDN = strings.TrimSpace(r.UserDN)
		config.Conf.Ldap.UserInitPassword = strings.TrimSpace(r.UserInitPassword)
		config.Conf.Ldap.DefaultEmailSuffix = strings.TrimSpace(r.DefaultEmailSuffix)
		config.Conf.Ldap.EnableSync = r.LdapEnableSync
		syncUserRule := strings.TrimSpace(r.SyncUsernameRule)
		if syncUserRule != "" {
			config.Conf.Ldap.SyncUsernameRule = syncUserRule
		}
		syncGroupRule := strings.TrimSpace(r.SyncGroupNameRule)
		if syncGroupRule != "" {
			config.Conf.Ldap.SyncGroupNameRule = syncGroupRule
		}
		userRdn := strings.TrimSpace(r.UserRDNAttr)
		if userRdn != "" {
			config.Conf.Ldap.UserRDNAttr = userRdn
		}
		groupRdn := strings.TrimSpace(r.GroupRDNAttr)
		if groupRdn != "" {
			config.Conf.Ldap.GroupRDNAttr = groupRdn
		}
	}

	// 更新配置文件
	viper.Set("ldap.directory-type", dirType)
	viper.Set("ldap.url", strings.TrimSpace(r.Url))
	viper.Set("ldap.base-dn", strings.TrimSpace(r.BaseDN))
	viper.Set("ldap.admin-dn", strings.TrimSpace(r.AdminDN))
	if strings.TrimSpace(r.AdminPass) != "" {
		viper.Set("ldap.admin-pass", r.AdminPass)
	}
	viper.Set("ldap.user-dn", strings.TrimSpace(r.UserDN))
	viper.Set("ldap.user-init-password", strings.TrimSpace(r.UserInitPassword))
	viper.Set("ldap.default-email-suffix", strings.TrimSpace(r.DefaultEmailSuffix))
	viper.Set("ldap.enable-sync", r.LdapEnableSync)
	if strings.TrimSpace(r.SyncUsernameRule) != "" {
		viper.Set("ldap.sync-username-rule", strings.TrimSpace(r.SyncUsernameRule))
	}
	if strings.TrimSpace(r.SyncGroupNameRule) != "" {
		viper.Set("ldap.sync-group-name-rule", strings.TrimSpace(r.SyncGroupNameRule))
	}
	if strings.TrimSpace(r.UserRDNAttr) != "" {
		viper.Set("ldap.user-rdn-attr", strings.TrimSpace(r.UserRDNAttr))
	}
	if strings.TrimSpace(r.GroupRDNAttr) != "" {
		viper.Set("ldap.group-rdn-attr", strings.TrimSpace(r.GroupRDNAttr))
	}

	if err := viper.WriteConfig(); err != nil {
		return nil, tools.NewOperationError(fmt.Errorf("保存配置文件失败: %s", err.Error()))
	}

	return nil, nil
}

// TestDirectoryConfig 测试目录（LDAP）连接：使用提交的参数拨号并绑定，不修改任何配置。
// url/adminDN/adminPass 留空时回落到已保存配置（与「保存」时 adminPass 留空不修改的语义一致）。
func (l BaseLogic) TestDirectoryConfig(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.BaseTestDirectoryConfigReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	url := strings.TrimSpace(r.Url)
	adminDN := strings.TrimSpace(r.AdminDN)
	adminPass := r.AdminPass
	if config.Conf.Ldap != nil {
		url = firstNonEmpty(url, config.Conf.Ldap.Url)
		adminDN = firstNonEmpty(adminDN, config.Conf.Ldap.AdminDN)
		adminPass = firstNonEmpty(adminPass, config.Conf.Ldap.AdminPass)
	}
	if url == "" {
		return nil, tools.NewValidatorError(fmt.Errorf("LDAP 地址不能为空"))
	}

	connected, message := common.ProbeLDAPConnectionWith(url, adminDN, adminPass)
	if !connected {
		return nil, tools.NewOperationError(fmt.Errorf("%s", message))
	}
	return tools.H{"ok": true, "message": message}, nil
}

// UpdateThirdPartyConfig 更新第三方平台配置
func (l BaseLogic) UpdateThirdPartyConfig(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.BaseThirdPartyConfigReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	platform := strings.ToLower(strings.TrimSpace(r.Platform))
	switch platform {
	case "dingtalk":
		if config.Conf.DingTalk != nil {
			if strings.TrimSpace(r.Flag) != "" {
				config.Conf.DingTalk.Flag = strings.TrimSpace(r.Flag)
				viper.Set("dingtalk.flag", strings.TrimSpace(r.Flag))
			}
			if strings.TrimSpace(r.AppKey) != "" {
				config.Conf.DingTalk.AppKey = strings.TrimSpace(r.AppKey)
				viper.Set("dingtalk.app-key", strings.TrimSpace(r.AppKey))
			}
			if strings.TrimSpace(r.AppSecret) != "" {
				config.Conf.DingTalk.AppSecret = r.AppSecret
				viper.Set("dingtalk.app-secret", r.AppSecret)
			}
			if strings.TrimSpace(r.AgentID) != "" {
				config.Conf.DingTalk.AgentId = strings.TrimSpace(r.AgentID)
				viper.Set("dingtalk.agent-id", strings.TrimSpace(r.AgentID))
			}
			config.Conf.DingTalk.EnableSync = r.EnableSync
			viper.Set("dingtalk.enable-sync", r.EnableSync)
		}
	case "wecom":
		if config.Conf.WeCom != nil {
			if strings.TrimSpace(r.Flag) != "" {
				config.Conf.WeCom.Flag = strings.TrimSpace(r.Flag)
				viper.Set("wecom.flag", strings.TrimSpace(r.Flag))
			}
			if strings.TrimSpace(r.CorpID) != "" {
				config.Conf.WeCom.CorpID = strings.TrimSpace(r.CorpID)
				viper.Set("wecom.corp-id", strings.TrimSpace(r.CorpID))
			}
			if strings.TrimSpace(r.CorpSecret) != "" {
				config.Conf.WeCom.CorpSecret = r.CorpSecret
				viper.Set("wecom.corp-secret", r.CorpSecret)
			}
			if r.WeComAgentID > 0 {
				config.Conf.WeCom.AgentID = r.WeComAgentID
				viper.Set("wecom.agent-id", r.WeComAgentID)
			}
			config.Conf.WeCom.EnableSync = r.EnableSync
			viper.Set("wecom.enable-sync", r.EnableSync)
		}
	case "feishu":
		if config.Conf.FeiShu != nil {
			if strings.TrimSpace(r.Flag) != "" {
				config.Conf.FeiShu.Flag = strings.TrimSpace(r.Flag)
				viper.Set("feishu.flag", strings.TrimSpace(r.Flag))
			}
			if strings.TrimSpace(r.AppID) != "" {
				config.Conf.FeiShu.AppID = strings.TrimSpace(r.AppID)
				viper.Set("feishu.app-id", strings.TrimSpace(r.AppID))
			}
			if strings.TrimSpace(r.AppSecret) != "" {
				config.Conf.FeiShu.AppSecret = r.AppSecret
				viper.Set("feishu.app-secret", r.AppSecret)
			}
			config.Conf.FeiShu.EnableSync = r.EnableSync
			viper.Set("feishu.enable-sync", r.EnableSync)
		}
	default:
		return nil, tools.NewValidatorError(fmt.Errorf("platform 仅支持 dingtalk/wecom/feishu"))
	}

	if err := viper.WriteConfig(); err != nil {
		return nil, tools.NewOperationError(fmt.Errorf("保存配置文件失败: %s", err.Error()))
	}
	return nil, nil
}

// TestThirdPartyConfig 测试第三方平台连接
func (l BaseLogic) TestThirdPartyConfig(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.BaseThirdPartyConfigReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	platform := strings.ToLower(strings.TrimSpace(r.Platform))
	switch platform {
	case "dingtalk":
		appKey := firstNonEmpty(strings.TrimSpace(r.AppKey), config.Conf.DingTalk.AppKey)
		appSecret := firstNonEmpty(r.AppSecret, config.Conf.DingTalk.AppSecret)
		if appKey == "" || appSecret == "" {
			return nil, tools.NewValidatorError(fmt.Errorf("钉钉 appKey/appSecret 不能为空"))
		}
		client, err := dingsdk.NewClient(appKey, appSecret)
		if err != nil {
			return nil, tools.NewOperationError(fmt.Errorf("钉钉初始化失败: %s", err.Error()))
		}
		if _, err = client.FetchDeptList(1, false, "zh_CN"); err != nil {
			return nil, tools.NewOperationError(fmt.Errorf("钉钉连接测试失败: %s", err.Error()))
		}
	case "wecom":
		corpID := firstNonEmpty(strings.TrimSpace(r.CorpID), config.Conf.WeCom.CorpID)
		corpSecret := firstNonEmpty(r.CorpSecret, config.Conf.WeCom.CorpSecret)
		agentID := r.WeComAgentID
		if agentID <= 0 {
			agentID = config.Conf.WeCom.AgentID
		}
		if corpID == "" || corpSecret == "" || agentID <= 0 {
			return nil, tools.NewValidatorError(fmt.Errorf("企微 corpId/corpSecret/agentId 不能为空"))
		}
		client := wecomsdk.NewClient(wecomsdk.Conf{
			CorpID:     corpID,
			AgentID:    agentID,
			CorpSecret: corpSecret,
		})
		if _, err := client.ListDepartment(&wecomsdk.ListDepartmentRequest{}); err != nil {
			return nil, tools.NewOperationError(fmt.Errorf("企微连接测试失败: %s", err.Error()))
		}
	case "feishu":
		appID := firstNonEmpty(strings.TrimSpace(r.AppID), config.Conf.FeiShu.AppID)
		appSecret := firstNonEmpty(r.AppSecret, config.Conf.FeiShu.AppSecret)
		if appID == "" || appSecret == "" {
			return nil, tools.NewValidatorError(fmt.Errorf("飞书 appId/appSecret 不能为空"))
		}
		client := lark.New(lark.WithAppCredential(appID, appSecret))
		pageSize := int64(1)
		pageToken := ""
		fetchChild := false
		reqTmp := lark.GetDepartmentListReq{
			PageToken:    &pageToken,
			FetchChild:   &fetchChild,
			PageSize:     &pageSize,
			DepartmentID: "0",
		}
		if _, _, err := client.Contact.GetDepartmentList(context.Background(), &reqTmp); err != nil {
			return nil, tools.NewOperationError(fmt.Errorf("飞书连接测试失败: %s", err.Error()))
		}
	default:
		return nil, tools.NewValidatorError(fmt.Errorf("platform 仅支持 dingtalk/wecom/feishu"))
	}

	return tools.H{"platform": platform, "ok": true}, nil
}

func firstNonEmpty(current, fallback string) string {
	if strings.TrimSpace(current) != "" {
		return current
	}
	return fallback
}

// TestNotification 测试通知：发送一封测试邮件，或向 Webhook 地址发送一条测试回调。
// SMTP / Webhook 字段留空时回落到已保存配置，可在保存前验证可用性。
func (l BaseLogic) TestNotification(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.BaseTestNotificationReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	switch strings.ToLower(strings.TrimSpace(r.Target)) {
	case "email":
		mailTo := strings.TrimSpace(r.Mail)
		if mailTo == "" {
			return nil, tools.NewValidatorError(fmt.Errorf("请填写测试收件邮箱"))
		}
		if err := tools.SendTestMailWith(r.SmtpHost, r.SmtpPort, r.SmtpUser, r.SmtpPass, r.SmtpFrom, mailTo); err != nil {
			return nil, tools.NewOperationError(fmt.Errorf("测试邮件发送失败: %s", err.Error()))
		}
		return tools.H{"ok": true, "target": "email"}, nil
	case "webhook":
		if err := common.TestWebhook(strings.TrimSpace(r.WebhookURL)); err != nil {
			return nil, tools.NewOperationError(fmt.Errorf("Webhook 测试失败: %s", err.Error()))
		}
		return tools.H{"ok": true, "target": "webhook"}, nil
	default:
		return nil, tools.NewValidatorError(fmt.Errorf("target 仅支持 email/webhook"))
	}
}

// ListWebhookDeliveries 分页查询 Webhook 投递记录（最新在前）
func (l BaseLogic) ListWebhookDeliveries(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.BaseWebhookDeliveriesReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c
	pageNum := r.PageNum
	if pageNum <= 0 {
		pageNum = 1
	}
	pageSize := r.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	var list []model.WebhookDelivery
	var total int64
	common.DB.Model(&model.WebhookDelivery{}).Count(&total)
	if err := common.DB.Order("created_at DESC").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("查询投递记录失败: %s", err.Error()))
	}
	return tools.H{"list": list, "total": total}, nil
}

// UpdateSyncConfig 更新定时自动同步配置（开关 + cron），保存后热加载调度
func (l BaseLogic) UpdateSyncConfig(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.BaseUpdateSyncConfigReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	spec := strings.TrimSpace(r.AutoSyncCron)
	if r.AutoSyncEnabled {
		if spec == "" {
			return nil, tools.NewValidatorError(fmt.Errorf("启用定时同步时必须填写 cron 表达式"))
		}
		if err := ValidateCron(spec); err != nil {
			return nil, tools.NewValidatorError(fmt.Errorf("cron 表达式无效（6 段，含秒）: %s", err.Error()))
		}
	}

	if config.Conf.System != nil {
		config.Conf.System.AutoSyncEnabled = r.AutoSyncEnabled
		config.Conf.System.AutoSyncCron = spec
	}
	viper.Set("system.auto-sync-enabled", r.AutoSyncEnabled)
	viper.Set("system.auto-sync-cron", spec)
	if err := viper.WriteConfig(); err != nil {
		return nil, tools.NewOperationError(fmt.Errorf("保存配置文件失败: %s", err.Error()))
	}
	if err := ReloadAutoSync(); err != nil {
		return nil, tools.NewOperationError(fmt.Errorf("重新加载定时任务失败: %s", err.Error()))
	}
	return nil, nil
}

// RunSyncNow 立即异步触发某来源的同步（部门+用户），结果记入同步运行记录
func (l BaseLogic) RunSyncNow(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.BaseRunSyncReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c
	source := strings.TrimSpace(r.Source)
	go RunSourceSync(source, "manual")
	return tools.H{"triggered": true, "source": source}, nil
}

// ListSyncRuns 分页查询同步运行记录（最新在前）
func (l BaseLogic) ListSyncRuns(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.BaseSyncRunsReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c
	pageNum := r.PageNum
	if pageNum <= 0 {
		pageNum = 1
	}
	pageSize := r.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	var list []model.SyncRun
	var total int64
	common.DB.Model(&model.SyncRun{}).Count(&total)
	if err := common.DB.Order("created_at DESC").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("查询同步记录失败: %s", err.Error()))
	}
	return tools.H{"list": list, "total": total}, nil
}

// currentUserFromCtx 从上下文取出当前登录用户（由 AuthMiddleware 的 authorizator 写入）
func currentUserFromCtx(c *gin.Context) (model.User, error) {
	v, ok := c.Get("user")
	if !ok {
		return model.User{}, fmt.Errorf("未获取到登录用户")
	}
	u, ok := v.(model.User)
	if !ok {
		return model.User{}, fmt.Errorf("登录用户信息异常")
	}
	return u, nil
}

// MfaStatus 查询当前用户 MFA 启用状态
func (l BaseLogic) MfaStatus(c *gin.Context, req any) (data any, rspError any) {
	if _, ok := req.(*request.BaseMfaStatusReq); !ok {
		return nil, ReqAssertErr
	}
	ctxUser, err := currentUserFromCtx(c)
	if err != nil {
		return nil, tools.NewValidatorError(err)
	}
	var user model.User
	if err := common.DB.First(&user, ctxUser.ID).Error; err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("查询用户失败: %s", err.Error()))
	}
	return tools.H{"enabled": user.MfaEnabled}, nil
}

// MfaSetup 生成一个新的 TOTP 密钥并返回二维码（此时仅为待验证状态，需校验后才启用）
func (l BaseLogic) MfaSetup(c *gin.Context, req any) (data any, rspError any) {
	if _, ok := req.(*request.BaseMfaSetupReq); !ok {
		return nil, ReqAssertErr
	}
	ctxUser, err := currentUserFromCtx(c)
	if err != nil {
		return nil, tools.NewValidatorError(err)
	}
	var user model.User
	if err := common.DB.First(&user, ctxUser.ID).Error; err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("查询用户失败: %s", err.Error()))
	}
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "LDAP 管理平台",
		AccountName: user.Username,
	})
	if err != nil {
		return nil, tools.NewOperationError(fmt.Errorf("生成 MFA 密钥失败: %s", err.Error()))
	}
	// 暂存待验证密钥（mfa_enabled 保持 false，校验通过后才真正启用）
	if err := common.DB.Model(&model.User{}).Where("id = ?", user.ID).
		Updates(map[string]any{"otp_secret": key.Secret(), "mfa_enabled": false}).Error; err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("保存 MFA 密钥失败: %s", err.Error()))
	}
	// 生成二维码 PNG（base64 内嵌，前端 <img> 直接展示）
	img, err := key.Image(220, 220)
	if err != nil {
		return nil, tools.NewOperationError(fmt.Errorf("生成二维码失败: %s", err.Error()))
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, tools.NewOperationError(fmt.Errorf("编码二维码失败: %s", err.Error()))
	}
	qr := "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())
	return tools.H{"secret": key.Secret(), "qr": qr, "otpauthUrl": key.URL()}, nil
}

// MfaVerify 校验动态验证码并启用 MFA
func (l BaseLogic) MfaVerify(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.BaseMfaCodeReq)
	if !ok {
		return nil, ReqAssertErr
	}
	ctxUser, err := currentUserFromCtx(c)
	if err != nil {
		return nil, tools.NewValidatorError(err)
	}
	var user model.User
	if err := common.DB.First(&user, ctxUser.ID).Error; err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("查询用户失败: %s", err.Error()))
	}
	if strings.TrimSpace(user.OtpSecret) == "" {
		return nil, tools.NewValidatorError(fmt.Errorf("请先生成 MFA 密钥"))
	}
	if !totp.Validate(strings.TrimSpace(r.Code), user.OtpSecret) {
		return nil, tools.NewValidatorError(fmt.Errorf("验证码错误，请重试"))
	}
	if err := common.DB.Model(&model.User{}).Where("id = ?", user.ID).Update("mfa_enabled", true).Error; err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("启用 MFA 失败: %s", err.Error()))
	}
	return tools.H{"enabled": true}, nil
}

// MfaDisable 关闭 MFA（需提供有效验证码以确认本人操作）
func (l BaseLogic) MfaDisable(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.BaseMfaCodeReq)
	if !ok {
		return nil, ReqAssertErr
	}
	ctxUser, err := currentUserFromCtx(c)
	if err != nil {
		return nil, tools.NewValidatorError(err)
	}
	var user model.User
	if err := common.DB.First(&user, ctxUser.ID).Error; err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("查询用户失败: %s", err.Error()))
	}
	if !user.MfaEnabled {
		return tools.H{"enabled": false}, nil
	}
	if !totp.Validate(strings.TrimSpace(r.Code), user.OtpSecret) {
		return nil, tools.NewValidatorError(fmt.Errorf("验证码错误，无法关闭"))
	}
	if err := common.DB.Model(&model.User{}).Where("id = ?", user.ID).
		Updates(map[string]any{"mfa_enabled": false, "otp_secret": ""}).Error; err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("关闭 MFA 失败: %s", err.Error()))
	}
	return tools.H{"enabled": false}, nil
}

// UpdateEmailConfig 更新邮件通知配置（开关 + 邮件服务器）
func (l BaseLogic) UpdateEmailConfig(c *gin.Context, req any) (data any, rspError any) {
	r, ok := req.(*request.BaseUpdateEmailConfigReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	if config.Conf.Email != nil {
		config.Conf.Email.SendUserCreationMail = r.SendUserCreationMail
		if strings.TrimSpace(r.SmtpHost) != "" {
			config.Conf.Email.Host = strings.TrimSpace(r.SmtpHost)
			viper.Set("email.host", config.Conf.Email.Host)
		}
		if strings.TrimSpace(r.SmtpPort) != "" {
			config.Conf.Email.Port = strings.TrimSpace(r.SmtpPort)
			viper.Set("email.port", config.Conf.Email.Port)
		}
		if strings.TrimSpace(r.SmtpUser) != "" {
			config.Conf.Email.User = strings.TrimSpace(r.SmtpUser)
			viper.Set("email.user", config.Conf.Email.User)
		}
		if strings.TrimSpace(r.SmtpFrom) != "" {
			config.Conf.Email.From = strings.TrimSpace(r.SmtpFrom)
			viper.Set("email.from", config.Conf.Email.From)
		}
		if strings.TrimSpace(r.SmtpPass) != "" {
			config.Conf.Email.Pass = r.SmtpPass
			viper.Set("email.pass", r.SmtpPass)
		}
	}
	viper.Set("email.send-user-creation-mail", r.SendUserCreationMail)
	if config.Conf.System != nil {
		config.Conf.System.WebhookURL = strings.TrimSpace(r.WebhookURL)
		viper.Set("system.webhook-url", config.Conf.System.WebhookURL)
		// 签名密钥：留空表示不修改（与各类密码字段一致）
		if strings.TrimSpace(r.WebhookSecret) != "" {
			config.Conf.System.WebhookSecret = strings.TrimSpace(r.WebhookSecret)
			viper.Set("system.webhook-secret", config.Conf.System.WebhookSecret)
		}
	}
	if err := viper.WriteConfig(); err != nil {
		return nil, tools.NewOperationError(fmt.Errorf("保存配置文件失败: %s", err.Error()))
	}
	return nil, nil
}

// GetVersion 获取版本信息
func (l BaseLogic) GetVersion(c *gin.Context, req any) (data any, rspError any) {
	_, ok := req.(*request.BaseVersionReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	return version.GetVersion(), nil
}

// GetLDAPStatus 获取 LDAP 连接状态（探测当前配置是否可连）
func (l BaseLogic) GetLDAPStatus(c *gin.Context, req any) (data any, rspError any) {
	_, ok := req.(*request.BaseLDAPStatusReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	connected, message := common.ProbeLDAPConnection()
	return response.LDAPStatusRsp{Connected: connected, Message: message}, nil
}

// GetSystemInfo 获取系统信息（版本、运行时长、数据库状态）
func (l BaseLogic) GetSystemInfo(c *gin.Context, req any) (data any, rspError any) {
	_, ok := req.(*request.BaseSystemInfoReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	versionInfo := version.GetVersion()
	uptime := common.GetAppUptime().String()
	driver, err := common.GetDBStatus()
	dbStatus := "正常"
	dbMessage := ""
	if err != nil {
		dbStatus = "异常"
		dbMessage = err.Error()
	}
	inactiveDays := 0
	if config.Conf.System != nil {
		inactiveDays = config.Conf.System.InactiveDays
	}
	return response.SystemInfoRsp{
		Version:      versionInfo,
		Uptime:       uptime,
		DBDriver:     driver,
		DBStatus:     dbStatus,
		DBMessage:    dbMessage,
		InactiveDays: inactiveDays,
	}, nil
}

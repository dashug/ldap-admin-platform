package request

// BaseSendCodeReq 发送验证码
type BaseSendCodeReq struct {
	Mail string `json:"mail" validate:"required,min=0,max=100"`
}

// BaseChangePwdReq 修改密码结构体
type BaseChangePwdReq struct {
	Mail string `json:"mail" validate:"required,min=0,max=100"`
	Code string `json:"code" validate:"required,len=6"`
}

// BaseDashboardReq  系统首页展示数据结构体
type BaseDashboardReq struct {
}

// EncryptPasswdReq
type EncryptPasswdReq struct {
	Passwd string `json:"passwd" form:"passwd" validate:"required"`
}

// BaseConfigReq 获取系统配置结构体
type BaseConfigReq struct {
}

// BaseVersionReq 获取版本信息结构体
type BaseVersionReq struct {
}

// BaseLDAPStatusReq 获取 LDAP 连接状态（无请求体）
type BaseLDAPStatusReq struct {
}

// BaseSystemInfoReq 获取系统信息（无请求体）
type BaseSystemInfoReq struct {
}

// BaseUpdateDirectoryConfigReq 更新目录服务配置
type BaseUpdateDirectoryConfigReq struct {
	DirectoryType      string `json:"directoryType" validate:"required,oneof=openldap ad"`
	Url                string `json:"url" validate:"required,min=1,max=255"`
	BaseDN             string `json:"baseDN" validate:"required,min=1,max=255"`
	AdminDN            string `json:"adminDN" validate:"required,min=1,max=255"`
	AdminPass          string `json:"adminPass" validate:"omitempty,min=0,max=255"`
	UserDN             string `json:"userDN" validate:"required,min=1,max=255"`
	UserInitPassword   string `json:"userInitPassword" validate:"required,min=1,max=255"`
	DefaultEmailSuffix string `json:"defaultEmailSuffix" validate:"required,min=1,max=100"`
	LdapEnableSync     bool   `json:"ldapEnableSync"`
	// 同步与 DN 规则（预制可选）
	SyncUsernameRule  string `json:"syncUsernameRule" validate:"omitempty,oneof=email_prefix pinyin job_number field_relation"`
	SyncGroupNameRule string `json:"syncGroupNameRule" validate:"omitempty,oneof=pinyin name"`
	UserRDNAttr       string `json:"userRDNAttr" validate:"omitempty,oneof=uid cn"`
	GroupRDNAttr      string `json:"groupRDNAttr" validate:"omitempty,oneof=ou cn"`
}

// BaseThirdPartyConfigReq 第三方平台配置
type BaseThirdPartyConfigReq struct {
	Platform   string `json:"platform" validate:"required,oneof=dingtalk wecom feishu"`
	Flag       string `json:"flag" validate:"omitempty,min=1,max=50"`
	EnableSync bool   `json:"enableSync"`

	AppKey    string `json:"appKey" validate:"omitempty,min=1,max=255"`
	AppSecret string `json:"appSecret" validate:"omitempty,min=0,max=255"`
	AgentID   string `json:"agentId" validate:"omitempty,min=0,max=50"`

	CorpID       string `json:"corpId" validate:"omitempty,min=1,max=255"`
	CorpSecret   string `json:"corpSecret" validate:"omitempty,min=0,max=255"`
	WeComAgentID int    `json:"weComAgentId" validate:"omitempty"`

	AppID string `json:"appId" validate:"omitempty,min=1,max=255"`
}

// BaseConfigExportReq 配置导出（无请求体，导出当前目录与同步规则等为 JSON）
type BaseConfigExportReq struct {
}

// BaseConfigImportReq 配置导入（与导出的 JSON 结构一致，用于恢复/迁移）
type BaseConfigImportReq struct {
	DirectoryType      string `json:"directoryType"`
	Url                string `json:"url"`
	BaseDN             string `json:"baseDN"`
	AdminDN            string `json:"adminDN"`
	AdminPass          string `json:"adminPass"`
	UserDN             string `json:"userDN"`
	UserInitPassword   string `json:"userInitPassword"`
	DefaultEmailSuffix string `json:"defaultEmailSuffix"`
	LdapEnableSync     bool   `json:"ldapEnableSync"`
	SyncUsernameRule   string `json:"syncUsernameRule"`
	SyncGroupNameRule  string `json:"syncGroupNameRule"`
	UserRDNAttr        string `json:"userRDNAttr"`
	GroupRDNAttr       string `json:"groupRDNAttr"`
}

// BaseTestDirectoryConfigReq 测试目录（LDAP）连接：仅校验连接所需字段，
// 留空字段在 logic 中回落到已保存配置（adminPass 留空 = 用已保存密码）。
type BaseTestDirectoryConfigReq struct {
	Url       string `json:"url" validate:"omitempty,max=255"`
	AdminDN   string `json:"adminDN" validate:"omitempty,max=255"`
	AdminPass string `json:"adminPass" validate:"omitempty,max=255"`
}

// BaseTestNotificationReq 测试通知：发送测试邮件或测试 Webhook。
// SMTP / Webhook 字段留空时回落到已保存配置。
type BaseTestNotificationReq struct {
	Target     string `json:"target" validate:"required,oneof=email webhook"`
	Mail       string `json:"mail" validate:"omitempty,max=255"` // 测试邮件收件人
	SmtpHost   string `json:"smtpHost" validate:"omitempty,max=255"`
	SmtpPort   string `json:"smtpPort" validate:"omitempty,max=20"`
	SmtpUser   string `json:"smtpUser" validate:"omitempty,max=255"`
	SmtpPass   string `json:"smtpPass" validate:"omitempty,max=255"`
	SmtpFrom   string `json:"smtpFrom" validate:"omitempty,max=255"`
	WebhookURL string `json:"webhookUrl" validate:"omitempty,max=2048"`
}

// BaseUpdateEmailConfigReq 更新邮件通知配置（开关 + 邮件服务器）+ Webhook 回调地址
type BaseUpdateEmailConfigReq struct {
	SendUserCreationMail bool   `json:"sendUserCreationMail"`
	SmtpHost             string `json:"smtpHost" validate:"omitempty,max=255"`
	SmtpPort             string `json:"smtpPort" validate:"omitempty,max=20"`
	SmtpUser             string `json:"smtpUser" validate:"omitempty,max=255"`
	SmtpPass             string `json:"smtpPass" validate:"omitempty,max=255"` // 留空表示不修改
	SmtpFrom             string `json:"smtpFrom" validate:"omitempty,max=255"`
	WebhookURL           string `json:"webhookUrl" validate:"omitempty,max=2048"` // 用户/部门创建或同步后 HTTP 回调地址
	WebhookSecret        string `json:"webhookSecret" validate:"omitempty,max=255"` // HMAC 签名密钥，留空表示不修改
}

// BaseWebhookDeliveriesReq Webhook 投递记录分页查询
type BaseWebhookDeliveriesReq struct {
	PageNum  int `json:"pageNum" form:"pageNum"`
	PageSize int `json:"pageSize" form:"pageSize"`
}

// BaseUpdateSyncConfigReq 更新定时自动同步配置
type BaseUpdateSyncConfigReq struct {
	AutoSyncEnabled bool   `json:"autoSyncEnabled"`
	AutoSyncCron    string `json:"autoSyncCron" validate:"omitempty,max=100"`
}

// BaseRunSyncReq 立即触发某来源同步
type BaseRunSyncReq struct {
	Source string `json:"source" validate:"required,oneof=ldap dingtalk wecom feishu"`
}

// BaseSyncRunsReq 同步运行记录分页查询
type BaseSyncRunsReq struct {
	PageNum  int `json:"pageNum" form:"pageNum"`
	PageSize int `json:"pageSize" form:"pageSize"`
}

// BaseMfaSetupReq 生成 MFA 密钥（无请求体）
type BaseMfaSetupReq struct{}

// BaseMfaStatusReq 查询当前用户 MFA 状态（无请求体）
type BaseMfaStatusReq struct{}

// BaseMfaCodeReq 校验/启用/关闭 MFA 时提交的动态验证码
type BaseMfaCodeReq struct {
	Code string `json:"code" validate:"required,len=6"`
}

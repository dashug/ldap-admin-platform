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

// DecryptPasswdReq
type DecryptPasswdReq struct {
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

// BaseUpdateEmailConfigReq 更新邮件通知配置（开关 + 邮件服务器）+ Webhook 回调地址
type BaseUpdateEmailConfigReq struct {
	SendUserCreationMail bool   `json:"sendUserCreationMail"`
	SmtpHost             string `json:"smtpHost" validate:"omitempty,max=255"`
	SmtpPort             string `json:"smtpPort" validate:"omitempty,max=20"`
	SmtpUser             string `json:"smtpUser" validate:"omitempty,max=255"`
	SmtpPass             string `json:"smtpPass" validate:"omitempty,max=255"` // 留空表示不修改
	SmtpFrom             string `json:"smtpFrom" validate:"omitempty,max=255"`
	WebhookURL           string `json:"webhookUrl" validate:"omitempty,max=2048"` // 用户/部门创建或同步后 HTTP 回调地址
}

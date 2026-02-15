package response

type DashboardList struct {
	DataType  string `json:"dataType"`
	DataName  string `json:"dataName"`
	DataCount int64  `json:"dataCount"`
	Icon      string `json:"icon"`
	Path      string `json:"path"`
}

type BaseConfigRsp struct {
	LdapEnableSync     bool   `json:"ldapEnableSync"`
	DingTalkEnableSync bool   `json:"dingTalkEnableSync"`
	FeiShuEnableSync   bool   `json:"feiShuEnableSync"`
	WeComEnableSync    bool   `json:"weComEnableSync"`
	DirectoryType      string `json:"directoryType"`
	Url                string `json:"url"`
	BaseDN             string `json:"baseDN"`
	AdminDN            string `json:"adminDN"`
	AdminPass          string `json:"adminPass"`
	UserDN             string `json:"userDN"`
	UserInitPassword   string `json:"userInitPassword"`
	DefaultEmailSuffix string `json:"defaultEmailSuffix"`
	// 同步与 DN 规则
	SyncUsernameRule  string `json:"syncUsernameRule"`
	SyncGroupNameRule string `json:"syncGroupNameRule"`
	UserRDNAttr       string `json:"userRDNAttr"`
	GroupRDNAttr      string `json:"groupRDNAttr"`

	DingTalkFlag      string `json:"dingTalkFlag"`
	DingTalkAppKey    string `json:"dingTalkAppKey"`
	DingTalkAppSecret string `json:"dingTalkAppSecret"`
	DingTalkAgentID   string `json:"dingTalkAgentId"`

	WeComFlag       string `json:"weComFlag"`
	WeComCorpID     string `json:"weComCorpId"`
	WeComCorpSecret string `json:"weComCorpSecret"`
	WeComAgentID    int    `json:"weComAgentId"`

	FeiShuFlag      string `json:"feiShuFlag"`
	FeiShuAppID     string `json:"feiShuAppId"`
	FeiShuAppSecret string `json:"feiShuAppSecret"`

	// 是否发送用户创建成功通知邮件（前端可配置，默认关闭）
	SendUserCreationMail bool `json:"sendUserCreationMail"`
	// 邮件服务器配置（供前端展示与编辑；密码不回传）
	SmtpHost string `json:"smtpHost"`
	SmtpPort string `json:"smtpPort"`
	SmtpUser string `json:"smtpUser"`
	SmtpFrom string `json:"smtpFrom"`
}

// LDAPStatusRsp LDAP 连接状态
type LDAPStatusRsp struct {
	Connected bool   `json:"connected"`
	Message   string `json:"message"`
}

// SystemInfoRsp 系统信息（版本、运行时长、数据库状态）
type SystemInfoRsp struct {
	Version      map[string]string `json:"version"`
	Uptime       string            `json:"uptime"`       // 运行时长
	DBDriver     string            `json:"dbDriver"`    // 数据库驱动
	DBStatus     string            `json:"dbStatus"`   // 数据库状态
	DBMessage    string            `json:"dbMessage"`  // 异常时的错误信息
	InactiveDays int               `json:"inactiveDays"` // N 天未登录自动禁用，0=关闭
}

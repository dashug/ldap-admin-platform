package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
)

// 系统配置，对应yml
// viper内置了mapstructure, yml文件用"-"区分单词, 转为驼峰方便

// 全局配置变量
var Conf = new(config)

// rsaPriv/rsaPub 保存运行期装载的登录加解密 RSA 密钥对（见 rsa_key.go: initRSAKeys）。
// 旧版把私钥 go:embed 进二进制并提交进 git，导致所有部署共用同一把「公开」私钥；
// 现改为运行期按 env/文件加载，缺失则自动生成一把独立密钥。
var (
	rsaPriv []byte
	rsaPub  []byte
)

type config struct {
	System   *SystemConfig `mapstructure:"system" json:"system"`
	Logs     *LogsConfig   `mapstructure:"logs" json:"logs"`
	Database *Database     `mapstructure:"database" json:"database"`
	Mysql    *MysqlConfig  `mapstructure:"mysql" json:"mysql"`
	// Casbin    *CasbinConfig    `mapstructure:"casbin" json:"casbin"`
	Jwt       *JwtConfig       `mapstructure:"jwt" json:"jwt"`
	RateLimit *RateLimitConfig `mapstructure:"rate-limit" json:"rateLimit"`
	Ldap      *LdapConfig      `mapstructure:"ldap" json:"ldap"`
	Email     *EmailConfig     `mapstructure:"email" json:"email"`
	DingTalk  *DingTalkConfig  `mapstructure:"dingtalk" json:"dingTalk"`
	WeCom     *WeComConfig     `mapstructure:"wecom" json:"weCom"`
	FeiShu    *FeiShuConfig    `mapstructure:"feishu" json:"feiShu"`
}

// 设置读取配置信息
func InitConfig() {
	workDir, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("读取应用目录失败:%s", err))
	}
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir + "/")
	// 读取配置信息；未找到 config.yml 时回退到内置模板 config.example.yml（便于首次运行/演示，
	// 也让 config.yml 可以不入库——避免真实密钥被提交）
	err = viper.ReadInConfig()
	if err != nil {
		viper.SetConfigName("config.example")
		if err2 := viper.ReadInConfig(); err2 == nil {
			err = nil
		}
	}

	// 热更新配置
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		// 将读取的配置信息保存至全局变量Conf
		if err := viper.Unmarshal(Conf); err != nil {
			panic(fmt.Errorf("初始化配置文件失败:%s", err))
		}
		// 回填运行期 RSA 密钥（viper.Unmarshal 会重建 Conf，需重新赋值）
		Conf.System.RSAPublicBytes = rsaPub
		Conf.System.RSAPrivateBytes = rsaPriv
		// 热更新后重新应用环境变量覆盖，避免 env 注入的密钥被文件里的占位值静默回退
		applyEnvOverrides()
	})

	if err != nil {
		panic(fmt.Errorf("读取配置文件失败:%s", err))
	}
	// 将读取的配置信息保存至全局变量Conf
	if err := viper.Unmarshal(Conf); err != nil {
		panic(fmt.Errorf("初始化配置文件失败:%s", err))
	}
	// 装载/生成运行期 RSA 密钥（不再依赖 go:embed 的仓库内私钥）
	if err := initRSAKeys(); err != nil {
		panic(fmt.Errorf("初始化 RSA 密钥失败:%s", err))
	}

	// 部分配合通过环境变量加载
	dbDriver := os.Getenv("DB_DRIVER")
	if dbDriver != "" {
		Conf.Database.Driver = dbDriver
	}
	mysqlHost := os.Getenv("MYSQL_HOST")
	if mysqlHost != "" {
		Conf.Mysql.Host = mysqlHost
	}
	mysqlUsername := os.Getenv("MYSQL_USERNAME")
	if mysqlUsername != "" {
		Conf.Mysql.Username = mysqlUsername
	}
	mysqlPassword := os.Getenv("MYSQL_PASSWORD")
	if mysqlPassword != "" {
		Conf.Mysql.Password = mysqlPassword
	}
	mysqlDatabase := os.Getenv("MYSQL_DATABASE")
	if mysqlDatabase != "" {
		Conf.Mysql.Database = mysqlDatabase
	}
	mysqlPort := os.Getenv("MYSQL_PORT")
	if mysqlPort != "" {
		Conf.Mysql.Port, _ = strconv.Atoi(mysqlPort)
	}

	ldapUrl := os.Getenv("LDAP_URL")
	if ldapUrl != "" {
		Conf.Ldap.Url = ldapUrl
	}
	ldapDirectoryType := os.Getenv("LDAP_DIRECTORY_TYPE")
	if ldapDirectoryType != "" {
		Conf.Ldap.DirectoryType = ldapDirectoryType
	}
	ldapBaseDN := os.Getenv("LDAP_BASE_DN")
	if ldapBaseDN != "" {
		Conf.Ldap.BaseDN = ldapBaseDN
	}
	ldapAdminDN := os.Getenv("LDAP_ADMIN_DN")
	if ldapAdminDN != "" {
		Conf.Ldap.AdminDN = ldapAdminDN
	}
	ldapAdminPass := os.Getenv("LDAP_ADMIN_PASS")
	if ldapAdminPass != "" {
		Conf.Ldap.AdminPass = ldapAdminPass
	}
	ldapUserDN := os.Getenv("LDAP_USER_DN")
	if ldapUserDN != "" {
		Conf.Ldap.UserDN = ldapUserDN
	}
	ldapUserInitPassword := os.Getenv("LDAP_USER_INIT_PASSWORD")
	if ldapUserInitPassword != "" {

		Conf.Ldap.UserInitPassword = ldapUserInitPassword
	}
	ldapDefaultEmailSuffix := os.Getenv("LDAP_DEFAULT_EMAIL_SUFFIX")
	if ldapDefaultEmailSuffix != "" {
		Conf.Ldap.DefaultEmailSuffix = ldapDefaultEmailSuffix
	}
	ldapUserPasswordEncryptionType := os.Getenv("LDAP_USER_PASSWORD_ENCRYPTION_TYPE")
	if ldapUserPasswordEncryptionType != "" {
		Conf.Ldap.UserPasswordEncryptionType = ldapUserPasswordEncryptionType
	}

	// 安全/凭据相关配置支持环境变量注入，避免明文写入 config.yml 或被烤进镜像。
	applyEnvOverrides()
}

// applyEnvOverrides 用环境变量覆盖安全/凭据相关配置项。
// 在启动装载与 config.yml 热更新回调中都会调用，确保 env 注入的密钥不会被文件里的占位值回退覆盖。
func applyEnvOverrides() {
	if Conf.Jwt != nil {
		applyEnvOverride(&Conf.Jwt.Key, "JWT_KEY")
	}
	if Conf.System != nil {
		applyEnvOverride(&Conf.System.Mode, "SYSTEM_MODE")
		applyEnvOverride(&Conf.System.WebhookSecret, "WEBHOOK_SECRET")
	}
	if Conf.FeiShu != nil {
		applyEnvOverride(&Conf.FeiShu.AppID, "FEISHU_APP_ID")
		applyEnvOverride(&Conf.FeiShu.AppSecret, "FEISHU_APP_SECRET")
	}
	if Conf.DingTalk != nil {
		applyEnvOverride(&Conf.DingTalk.AppKey, "DINGTALK_APP_KEY")
		applyEnvOverride(&Conf.DingTalk.AppSecret, "DINGTALK_APP_SECRET")
	}
	if Conf.WeCom != nil {
		applyEnvOverride(&Conf.WeCom.CorpID, "WECOM_CORP_ID")
		applyEnvOverride(&Conf.WeCom.CorpSecret, "WECOM_CORP_SECRET")
	}
	if Conf.Email != nil {
		applyEnvOverride(&Conf.Email.User, "EMAIL_USER")
		applyEnvOverride(&Conf.Email.Pass, "EMAIL_PASS")
	}
}

// applyEnvOverride 若指定环境变量非空，则用其覆盖目标配置项。
func applyEnvOverride(target *string, envKey string) {
	if v := os.Getenv(envKey); v != "" {
		*target = v
	}
}

type SystemConfig struct {
	Mode            string `mapstructure:"mode" json:"mode"`
	UrlPathPrefix   string `mapstructure:"url-path-prefix" json:"urlPathPrefix"`
	Port            int    `mapstructure:"port" json:"port"`
	InitData        bool   `mapstructure:"init-data" json:"initData"`
	InactiveDays    int    `mapstructure:"inactive-days" json:"inactiveDays"`   // 账户过期策略：N 天未登录自动禁用，0 表示关闭
	WebhookURL      string `mapstructure:"webhook-url" json:"webhookUrl"`      // 用户/部门创建或同步后 HTTP 回调地址，空则不发
	WebhookSecret   string `mapstructure:"webhook-secret" json:"-"`           // Webhook HMAC-SHA256 签名密钥，空则不签名
	AutoSyncEnabled bool   `mapstructure:"auto-sync-enabled" json:"autoSyncEnabled"` // 是否启用定时自动同步
	AutoSyncCron    string `mapstructure:"auto-sync-cron" json:"autoSyncCron"`       // 定时同步 cron 表达式（robfig/cron，6 段，含秒）
	// 允许跨域的来源白名单。
	//  - 为空（默认）：单二进制同源部署，不下发任何跨域响应头（最安全）
	//  - ["*"]      ：放行所有来源（仅在确有需要时使用）
	//  - 指定来源    ：仅对命中白名单的 Origin 回显该来源
	AllowOrigins    []string `mapstructure:"allow-origins" json:"allowOrigins"`
	RSAPublicBytes  []byte   `mapstructure:"-" json:"-"`
	RSAPrivateBytes []byte   `mapstructure:"-" json:"-"`
}

type LogsConfig struct {
	Level      zapcore.Level `mapstructure:"level" json:"level"`
	Path       string        `mapstructure:"path" json:"path"`
	MaxSize    int           `mapstructure:"max-size" json:"maxSize"`
	MaxBackups int           `mapstructure:"max-backups" json:"maxBackups"`
	MaxAge     int           `mapstructure:"max-age" json:"maxAge"`
	Compress   bool          `mapstructure:"compress" json:"compress"`
}

type Database struct {
	Driver string `mapstructure:"driver" json:"driver"`
	Source string `mapstructure:"source" json:"source"`
}

type MysqlConfig struct {
	Username    string `mapstructure:"username" json:"username"`
	Password    string `mapstructure:"password" json:"password"`
	Database    string `mapstructure:"database" json:"database"`
	Host        string `mapstructure:"host" json:"host"`
	Port        int    `mapstructure:"port" json:"port"`
	Query       string `mapstructure:"query" json:"query"`
	LogMode     bool   `mapstructure:"log-mode" json:"logMode"`
	TablePrefix string `mapstructure:"table-prefix" json:"tablePrefix"`
	Charset     string `mapstructure:"charset" json:"charset"`
	Collation   string `mapstructure:"collation" json:"collation"`
}

// type CasbinConfig struct {
// 	ModelPath string `mapstructure:"model-path" json:"modelPath"`
// }

type JwtConfig struct {
	Realm      string `mapstructure:"realm" json:"realm"`
	Key        string `mapstructure:"key" json:"key"`
	Timeout    int    `mapstructure:"timeout" json:"timeout"`
	MaxRefresh int    `mapstructure:"max-refresh" json:"maxRefresh"`
}

type RateLimitConfig struct {
	FillInterval int64 `mapstructure:"fill-interval" json:"fillInterval"`
	Capacity     int64 `mapstructure:"capacity" json:"capacity"`
}

type LdapConfig struct {
	DirectoryType              string `mapstructure:"directory-type" json:"directoryType"`
	Url                        string `mapstructure:"url" json:"url"`
	MaxConn                    int    `mapstructure:"max-conn" json:"maxConn"`
	BaseDN                     string `mapstructure:"base-dn" json:"baseDN"`
	AdminDN                    string `mapstructure:"admin-dn" json:"adminDN"`
	AdminPass                  string `mapstructure:"admin-pass" json:"adminPass"`
	UserDN                     string `mapstructure:"user-dn" json:"userDN"`
	UserInitPassword           string `mapstructure:"user-init-password" json:"userInitPassword"`
	GroupNameModify            bool   `mapstructure:"group-name-modify" json:"groupNameModify"`
	UserNameModify             bool   `mapstructure:"user-name-modify" json:"userNameModify"`
	DefaultEmailSuffix         string `mapstructure:"default-email-suffix" json:"defaultEmailSuffix"`
	UserPasswordEncryptionType string `mapstructure:"user-password-encryption-type" json:"userPasswordEncryptionType"`
	EnableSync                 bool   `mapstructure:"enable-sync" json:"enableSync"`
	// 同步与 DN 规则（同步前可选）
	SyncUsernameRule  string `mapstructure:"sync-username-rule" json:"syncUsernameRule"`   // 用户名规则: email_prefix=邮箱前段, pinyin=姓名拼音, job_number=工号, field_relation=按字段关联
	SyncGroupNameRule string `mapstructure:"sync-group-name-rule" json:"syncGroupNameRule"` // 部门名规则: pinyin=拼音, name=中文名
	UserRDNAttr       string `mapstructure:"user-rdn-attr" json:"userRDNAttr"`             // 用户 DN 的 RDN 属性: uid / cn
	GroupRDNAttr      string `mapstructure:"group-rdn-attr" json:"groupRDNAttr"`          // 部门 DN 的 RDN 属性: ou / cn
}
type EmailConfig struct {
	Host                 string `mapstructure:"host" json:"host"`
	Port                 string `mapstructure:"port" json:"port"`
	User                 string `mapstructure:"user" json:"user"`
	Pass                 string `mapstructure:"pass" json:"pass"`
	From                 string `mapstructure:"from" json:"from"`
	SendUserCreationMail bool   `mapstructure:"send-user-creation-mail" json:"sendUserCreationMail"` // 是否发送用户创建成功通知邮件，默认 false
}

type DingTalkConfig struct {
	AppKey        string   `mapstructure:"app-key" json:"appKey"`
	AppSecret     string   `mapstructure:"app-secret" json:"appSecret"`
	AgentId       string   `mapstructure:"agent-id" json:"agentId"`
	RootOuName    string   `mapstructure:"root-ou-name" json:"rootOuName"`
	Flag          string   `mapstructure:"flag" json:"flag"`
	EnableSync    bool     `mapstructure:"enable-sync" json:"enableSync"`
	DeptSyncTime  string   `mapstructure:"dept-sync-time" json:"deptSyncTime"`
	UserSyncTime  string   `mapstructure:"user-sync-time" json:"userSyncTime"`
	DeptList      []string `mapstructure:"dept-list" json:"deptList"`
	IsUpdateSyncd bool     `mapstructure:"is-update-syncd" json:"isUpdateSyncd"`
	ULeaveRange   uint     `mapstructure:"user-leave-range" json:"userLevelRange"`
}

type WeComConfig struct {
	Flag          string `mapstructure:"flag" json:"flag"`
	CorpID        string `mapstructure:"corp-id" json:"corpId"`
	AgentID       int    `mapstructure:"agent-id" json:"agentId"`
	CorpSecret    string `mapstructure:"corp-secret" json:"corpSecret"`
	EnableSync    bool   `mapstructure:"enable-sync" json:"enableSync"`
	DeptSyncTime  string `mapstructure:"dept-sync-time" json:"deptSyncTime"`
	UserSyncTime  string `mapstructure:"user-sync-time" json:"userSyncTime"`
	IsUpdateSyncd bool   `mapstructure:"is-update-syncd" json:"isUpdateSyncd"`
}

type FeiShuConfig struct {
	Flag          string   `mapstructure:"flag" json:"flag"`
	AppID         string   `mapstructure:"app-id" json:"appId"`
	AppSecret     string   `mapstructure:"app-secret" json:"appSecret"`
	EnableSync    bool     `mapstructure:"enable-sync" json:"enableSync"`
	DeptSyncTime  string   `mapstructure:"dept-sync-time" json:"deptSyncTime"`
	UserSyncTime  string   `mapstructure:"user-sync-time" json:"userSyncTime"`
	DeptList      []string `mapstructure:"dept-list" json:"deptList"`
	IsUpdateSyncd bool     `mapstructure:"is-update-syncd" json:"isUpdateSyncd"`
}

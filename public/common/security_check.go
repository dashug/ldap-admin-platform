package common

import (
	"strings"

	"github.com/dashug/ldap-admin-platform/config"
)

// 已知的弱默认值，命中即视为未修改的出厂配置
const (
	defaultJwtKey      = "secret key"
	defaultWeakPasswd  = "123456"
	defaultWeakPasswd2 = "your password"
	// config.example.yml（镜像默认 config 即来自它）内置的占位值，也应视为「未修改」
	examplePlaceholderPass = "please-change-me"
)

// weakJwtKeys 汇总所有已知的弱/占位签名密钥（含仓库与镜像模板内置值 + 空值），命中即视为未修改。
// 覆盖 config.example.yml 的占位值，避免「镜像默认配置 + release 模式 + 未注入 JWT_KEY」时
// 带着公开可得的签名密钥启动却未被熔断。
var weakJwtKeys = map[string]bool{
	"":            true, // 空密钥无签名安全
	defaultJwtKey: true, // "secret key"
	"please-change-me-to-a-random-32+bytes-secret": true, // config.example.yml 占位
}

// isWeakAdminPass 判断密码是否为已知弱默认/占位值。
func isWeakAdminPass(p string) bool {
	return p == defaultWeakPasswd || p == examplePlaceholderPass
}

// SecurityCheck 在启动时对关键配置做安全自检：
// 命中弱默认值时打印醒目告警，提醒使用者在生产环境务必修改；
// 若运行在 release 模式下仍使用弱配置，则升级为更高等级的告警。
func SecurityCheck() {
	var issues []string
	var fatal []string // release 模式下会导致系统被直接接管的高危默认值，必须阻断启动

	if c := config.Conf; c != nil {
		if c.Jwt != nil && weakJwtKeys[strings.TrimSpace(c.Jwt.Key)] {
			msg := "jwt.key 仍为默认/占位值，攻击者可据此离线伪造任意登录 token（含 admin），请改为随机强密钥（可用环境变量 JWT_KEY 注入）"
			issues = append(issues, msg)
			fatal = append(fatal, msg)
		}
		if c.Ldap != nil {
			if isWeakAdminPass(c.Ldap.AdminPass) {
				msg := "ldap.admin-pass 仍为弱默认/占位值，admin 登录口令过弱，请修改（可用环境变量 LDAP_ADMIN_PASS 注入）"
				issues = append(issues, msg)
				fatal = append(fatal, msg)
			}
			if isWeakAdminPass(c.Ldap.UserInitPassword) {
				issues = append(issues, "ldap.user-init-password 仍为弱默认/占位值，新建用户初始密码过于简单")
			}
		}
		// 仅当真正使用 MySQL 时才校验其密码
		if c.Database != nil && c.Database.Driver == "mysql" && c.Mysql != nil && isWeakAdminPass(c.Mysql.Password) {
			issues = append(issues, "mysql.password 仍为弱默认/占位值，请修改数据库密码")
		}
		if c.Email != nil && c.Email.Pass == defaultWeakPasswd2 {
			issues = append(issues, "email.pass 仍为占位值 \"your password\"，如启用邮件通知请填写真实凭据")
		}
	}

	if len(issues) == 0 {
		Log.Info("安全自检通过：未发现弱默认配置")
		return
	}

	isRelease := config.Conf != nil && config.Conf.System != nil && config.Conf.System.Mode == "release"

	banner := "==================== 安全配置告警 ===================="
	Log.Warn(banner)
	for _, it := range issues {
		if isRelease {
			// 生产模式下弱配置危害更大，提升到 Error 级别以引起重视
			Log.Errorf("[release 模式] %s", it)
		} else {
			Log.Warnf("%s", it)
		}
	}
	Log.Warn("上述项请在 config.yml 中修改，或通过环境变量覆盖后再部署到生产环境")
	Log.Warn(strings.Repeat("=", len(banner)))

	// release 模式下若仍使用会导致系统被直接接管的高危默认值（jwt.key / admin 密码），拒绝启动，避免带病上线。
	// debug/test 模式仅告警，方便本地开发与演示。
	if isRelease && len(fatal) > 0 {
		Log.Fatalf("release 模式检测到 %d 项高危默认配置，已阻断启动；请修正后重启（可用环境变量 JWT_KEY、LDAP_ADMIN_PASS 等注入）。", len(fatal))
	}
}

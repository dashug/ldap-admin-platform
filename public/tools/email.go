package tools

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/dashug/ldap-admin-platform/config"
	"github.com/patrickmn/go-cache"

	"strconv"

	"gopkg.in/gomail.v2"
)

func firstNonEmptyStr(a, b string) string {
	if strings.TrimSpace(a) != "" {
		return a
	}
	return b
}

// 验证码放到缓存当中
var VerificationCodeCache = cache.New(24*time.Hour, 48*time.Hour)

func email(mailTo []string, subject string, body string) error {
	mailConn := map[string]string{
		"user": config.Conf.Email.User,
		"pass": config.Conf.Email.Pass,
		"host": config.Conf.Email.Host,
		"port": config.Conf.Email.Port,
	}
	port, _ := strconv.Atoi(mailConn["port"]) //转换端口类型为int

	newmail := gomail.NewMessage()

	newmail.SetHeader("From", newmail.FormatAddress(mailConn["user"], config.Conf.Email.From))
	newmail.SetHeader("To", mailTo...)    //发送给多个用户
	newmail.SetHeader("Subject", subject) //设置邮件主题
	newmail.SetBody("text/html", body)    //设置邮件正文

	do := gomail.NewDialer(mailConn["host"], port, mailConn["user"], mailConn["pass"])
	return do.DialAndSend(newmail)
}

// SendTestMailWith 用指定 SMTP 参数发送一封测试邮件，留空字段回落到已保存配置。
// 用于通知设置页「发送测试邮件」：可在保存前验证 SMTP 是否真的可用。
func SendTestMailWith(host, port, user, pass, from, mailTo string) error {
	if config.Conf.Email != nil {
		host = firstNonEmptyStr(host, config.Conf.Email.Host)
		port = firstNonEmptyStr(port, config.Conf.Email.Port)
		user = firstNonEmptyStr(user, config.Conf.Email.User)
		pass = firstNonEmptyStr(pass, config.Conf.Email.Pass)
		from = firstNonEmptyStr(from, config.Conf.Email.From)
	}
	if host == "" || port == "" || user == "" {
		return fmt.Errorf("SMTP 服务器/端口/发件人邮箱不完整")
	}
	p, err := strconv.Atoi(strings.TrimSpace(port))
	if err != nil {
		return fmt.Errorf("SMTP 端口无效: %s", port)
	}
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(user, from))
	m.SetHeader("To", mailTo)
	m.SetHeader("Subject", "LDAP 管理平台 - 测试邮件")
	m.SetBody("text/html", "<p>这是一封来自 LDAP 管理平台的测试邮件。</p><p>如果你收到了它，说明当前 SMTP 通知配置可用。</p>")
	d := gomail.NewDialer(host, p, user, pass)
	return d.DialAndSend(m)
}

func SendMail(sendto []string, pass string) error {
	subject := "重置LDAP密码成功"
	// 邮件正文
	body := fmt.Sprintf("<li><a>更改之后的密码为: %s </a></li>", pass)
	return email(sendto, subject, body)
}

// SendCode 发送验证码
func SendCode(sendto []string) error {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	vcode := fmt.Sprintf("%06v", rnd.Int31n(1000000))
	// 把验证码信息放到cache，以便于验证时拿到
	VerificationCodeCache.Set(sendto[0], vcode, time.Minute*5)
	subject := "验证码-重置密码"
	//发送的内容
	body := fmt.Sprintf(`<div>
        <div>
            尊敬的用户，您好！
        </div>
        <div style="padding: 8px 40px 8px 50px;">
            <p>你本次的验证码为 %s ,为了保证账号安全，验证码有效期为5分钟。请确认为本人操作，切勿向他人泄露，感谢您的理解与使用。</p>
        </div>
        <div>
            <p>此邮箱为系统邮箱，请勿回复。</p>
        </div>
    </div>`, vcode)
	return email(sendto, subject, body)
}

// SendUserCreationNotification 发送用户创建成功通知邮件
func SendUserCreationNotification(username, nickname, mail, password string) error {
	subject := "LDAP账户创建成功通知"
	// 邮件正文
	body := fmt.Sprintf(`<div>
        <div>
            尊敬的%s，您好！
        </div>
        <div style="padding: 8px 40px 8px 50px;">
            <p>您的LDAP账户已创建成功，以下是您的账户信息：</p>
            <ul>
                <li>用户名：%s</li>
                <li>昵称：%s</li>
                <li>初始密码：%s</li>
            </ul>
            <p style="color: #ff6600;">请妥善保管您的账户信息，建议首次登录后及时修改密码。</p>
        </div>
        <div>
            <p>此邮箱为系统邮箱，请勿回复。</p>
        </div>
    </div>`, nickname, username, nickname, password)
	return email([]string{mail}, subject, body)
}

// SendPasswordResetNotification 发送密码重置成功通知邮件
func SendPasswordResetNotification(username, nickname, mail, newPassword string) error {
	subject := "LDAP密码重置成功通知"
	// 邮件正文
	body := fmt.Sprintf(`<div>
        <div>
            尊敬的%s，您好！
        </div>
        <div style="padding: 8px 40px 8px 50px;">
            <p>您的LDAP账户密码已成功重置，以下是您的新密码信息：</p>
            <ul>
                <li>用户名：%s</li>
                <li>新密码：%s</li>
            </ul>
            <p style="color: #ff6600;">为了您的账户安全，请尽快登录并修改为您自己的密码。</p>
            <p style="color: #ff0000;">请妥善保管您的账户信息，切勿泄露给他人。</p>
        </div>
        <div>
            <p>此邮箱为系统邮箱，请勿回复。</p>
        </div>
    </div>`, nickname, username, newPassword)
	return email([]string{mail}, subject, body)
}

package tools

import (
	"crypto/rand"
	"math/big"
	"strings"

	"github.com/dashug/ldap-admin-platform/config"
	"golang.org/x/crypto/bcrypt"
)

// 说明：
//   - 同步到 LDAP/AD 的用户必须保留“可逆”密码（NewGenPasswd/NewParPasswd，RSA 加解密），
//     因为创建/改密时需要把明文推送到目录服务。
//   - 不需要下发到 LDAP 的本地账号（如平台管理员 admin），其密码以 bcrypt 不可逆哈希存储，
//     即便数据库与 RSA 私钥同时泄露也无法还原明文。
//   - 登录与改密校验统一走 VerifyPassword，自动识别两种存储格式，平滑兼容历史数据。

// 密码加密（RSA 可逆，供需要同步到 LDAP 的用户使用）
func NewGenPasswd(passwd string) string {
	pass, _ := RSAEncrypt([]byte(passwd), config.Conf.System.RSAPublicBytes)
	return string(pass)
}

// 密码解密（RSA，还原明文以便推送到 LDAP）
func NewParPasswd(passwd string) string {
	pass, _ := RSADecrypt([]byte(passwd), config.Conf.System.RSAPrivateBytes)
	return string(pass)
}

// HashPassword 使用 bcrypt 生成不可逆哈希，供本地账号存储密码
func HashPassword(plaintext string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintext), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(hash)
}

// IsBcryptHash 判断字符串是否为 bcrypt 哈希（$2a$/$2b$/$2y$ 前缀，固定 60 字符）
func IsBcryptHash(stored string) bool {
	if len(stored) != 60 {
		return false
	}
	return strings.HasPrefix(stored, "$2a$") ||
		strings.HasPrefix(stored, "$2b$") ||
		strings.HasPrefix(stored, "$2y$")
}

// VerifyPassword 校验明文与数据库中存储的密码是否匹配。
// 自动兼容两种存储格式：bcrypt 哈希走 bcrypt 比对；否则按历史的 RSA 可逆密文比对。
func VerifyPassword(stored, plaintext string) bool {
	if IsBcryptHash(stored) {
		return bcrypt.CompareHashAndPassword([]byte(stored), []byte(plaintext)) == nil
	}
	return NewParPasswd(stored) == plaintext
}

const (
	passwordLength = 8
	letters        = "abcdefghijklmnopqrstu@vwxyzABCDEFGHIJKL#MNOP*QRSTUVWXYZ0123456789"
	lettersLength  = len(letters)
)

// 生成随机密码
func GenerateRandomPassword() string {
	password := make([]byte, passwordLength)

	for i := range password {
		index, _ := rand.Int(rand.Reader, big.NewInt(int64(lettersLength)))
		password[i] = letters[index.Int64()]
	}

	return string(password)
}

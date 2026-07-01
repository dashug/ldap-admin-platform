package tools

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"strings"
	"testing"

	"github.com/dashug/ldap-admin-platform/config"
)

// TestMain 初始化 config.Conf 的 RSA 密钥对，供依赖它的 NewGenPasswd/NewParPasswd
// （以及既有的 TestGenPass）在无完整启动流程时也能运行。
func TestMain(m *testing.M) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	privPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	pubDER, _ := x509.MarshalPKIXPublicKey(&key.PublicKey)
	pubPEM := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubDER})
	config.Conf.System = &config.SystemConfig{RSAPublicBytes: pubPEM, RSAPrivateBytes: privPEM}
	os.Exit(m.Run())
}

func TestSSHAEncodeMatches(t *testing.T) {
	enc := EncodePass([]byte("directory-pass"))
	if !strings.HasPrefix(enc, "{SSHA}") {
		t.Fatalf("SSHA 编码应以 {SSHA} 开头，得到: %q", enc)
	}
	if !Matches([]byte(enc), []byte("directory-pass")) {
		t.Error("正确密码应与 SSHA 匹配")
	}
	if Matches([]byte(enc), []byte("nope")) {
		t.Error("错误密码不应与 SSHA 匹配")
	}
}

// TestRSARoundTrip 验证 config/rsa_key.go 生成的密钥格式（PKCS1 私钥 + PKIX 公钥）
// 与 tools.RSAEncrypt/RSADecrypt 期望的格式一致。
func TestRSARoundTrip(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("生成 RSA 密钥失败: %v", err)
	}
	privPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	pubDER, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		t.Fatalf("序列化公钥失败: %v", err)
	}
	pubPEM := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubDER})

	enc, err := RSAEncrypt([]byte("login-password"), pubPEM)
	if err != nil {
		t.Fatalf("RSAEncrypt 失败: %v", err)
	}
	dec, err := RSADecrypt(enc, privPEM)
	if err != nil {
		t.Fatalf("RSADecrypt 失败: %v", err)
	}
	if string(dec) != "login-password" {
		t.Errorf("RSA 往返不一致，得到: %q", string(dec))
	}
}

// TestReversiblePasswordRoundTrip 覆盖同步到 LDAP 的「可逆密码」路径与 VerifyPassword 的非 bcrypt 分支。
func TestReversiblePasswordRoundTrip(t *testing.T) {
	enc := NewGenPasswd("directory-secret")
	if enc == "" {
		t.Fatal("NewGenPasswd 返回空")
	}
	if got := NewParPasswd(enc); got != "directory-secret" {
		t.Errorf("RSA 可逆密码往返不一致: %q", got)
	}
	if !VerifyPassword(enc, "directory-secret") {
		t.Error("VerifyPassword 应识别 RSA 可逆密文并校验通过")
	}
}

func TestApiKeyHash(t *testing.T) {
	raw, err := GenerateRandomApiKey()
	if err != nil {
		t.Fatalf("生成 API Key 失败: %v", err)
	}
	if !strings.HasPrefix(raw, "glak_") || len(raw) != len("glak_")+64 {
		t.Fatalf("API Key 格式应为 glak_+64hex，得到: %q (len=%d)", raw, len(raw))
	}
	hash, err := GenApiKeyHash(raw)
	if err != nil {
		t.Fatalf("哈希 API Key 失败: %v", err)
	}
	if !VerifyApiKeyHash(hash, raw) {
		t.Error("正确 API Key 应校验通过")
	}
	if VerifyApiKeyHash(hash, raw+"x") {
		t.Error("被篡改的 API Key 不应校验通过")
	}
}

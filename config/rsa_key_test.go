package config

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"testing"
)

// TestParseRSAPrivateKey 验证注入私钥兼容 PKCS1 与 PKCS8（openssl 默认输出 PKCS8）。
func TestParseRSAPrivateKey(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("生成密钥失败: %v", err)
	}

	// PKCS1
	if got, err := parseRSAPrivateKey(x509.MarshalPKCS1PrivateKey(key)); err != nil {
		t.Errorf("PKCS1 解析失败: %v", err)
	} else if got.N.Cmp(key.N) != 0 {
		t.Error("PKCS1 解析出的私钥与原始不一致")
	}

	// PKCS8
	der, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatalf("PKCS8 编码失败: %v", err)
	}
	if got, err := parseRSAPrivateKey(der); err != nil {
		t.Errorf("PKCS8 解析失败: %v", err)
	} else if got.N.Cmp(key.N) != 0 {
		t.Error("PKCS8 解析出的私钥与原始不一致")
	}

	// 非法输入应报错
	if _, err := parseRSAPrivateKey([]byte("not-a-key")); err == nil {
		t.Error("非法私钥应解析失败")
	}
}

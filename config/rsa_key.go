package config

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
)

const (
	// 运行期私钥的持久化位置（位于数据目录，已被 .gitignore/.dockerignore 排除）
	defaultRSAKeyPath = "data/rsa_private.pem"
	// 历史遗留：旧版 go:embed 的仓库内私钥。仅用于兼容旧本地部署做一次性迁移，之后不再依赖。
	legacyRSAKeyPath = "config/go-ldap-admin-priv.pem"
)

// initRSAKeys 装载登录密码加解密所用的 RSA 密钥对。
//
// 与旧版把私钥 go:embed 进二进制（所有部署共用同一把、且私钥已提交进 git，等同公开）不同，
// 这里在【运行期】按以下优先级装载或自动生成，确保每个部署拥有独立、非公开的密钥：
//  1. 环境变量 RSA_PRIVATE_KEY（PEM 文本）
//  2. 环境变量 RSA_PRIVATE_KEY_FILE 指向的文件
//  3. data/rsa_private.pem（此前生成/迁移的私钥）
//  4. 历史遗留 config/go-ldap-admin-priv.pem（存在则迁移到 data/ 后使用）
//  5. 以上都没有则生成一把全新的 2048 位 RSA 私钥并持久化到 data/
//
// 公钥由私钥派生（PKIX/SPKI PEM，与 tools.RSAEncrypt 期望的格式一致），无需单独文件。
func initRSAKeys() error {
	var privBytes []byte

	switch {
	case os.Getenv("RSA_PRIVATE_KEY") != "":
		privBytes = []byte(os.Getenv("RSA_PRIVATE_KEY"))
	case os.Getenv("RSA_PRIVATE_KEY_FILE") != "":
		b, err := os.ReadFile(os.Getenv("RSA_PRIVATE_KEY_FILE"))
		if err != nil {
			return fmt.Errorf("读取 RSA_PRIVATE_KEY_FILE 失败: %w", err)
		}
		privBytes = b
	case fileExists(defaultRSAKeyPath):
		b, err := os.ReadFile(defaultRSAKeyPath)
		if err != nil {
			return err
		}
		privBytes = b
	case fileExists(legacyRSAKeyPath):
		// 兼容旧本地部署：迁移历史私钥到 data 目录，之后不再依赖仓库内的 pem。
		b, err := os.ReadFile(legacyRSAKeyPath)
		if err != nil {
			return err
		}
		privBytes = b
		if err := persistPrivateKey(defaultRSAKeyPath, b); err != nil {
			fmt.Printf("警告: 迁移历史 RSA 私钥到 %s 失败: %v\n", defaultRSAKeyPath, err)
		}
	default:
		key, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return fmt.Errorf("生成 RSA 密钥失败: %w", err)
		}
		privBytes = pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		})
		if err := persistPrivateKey(defaultRSAKeyPath, privBytes); err != nil {
			// 持久化失败不致命：本次运行内存里仍有可用密钥，但重启会换新钥（会使已加密数据失效）。
			fmt.Printf("警告: 持久化 RSA 私钥到 %s 失败, 重启后将重新生成: %v\n", defaultRSAKeyPath, err)
		}
	}

	block, _ := pem.Decode(privBytes)
	if block == nil {
		return fmt.Errorf("RSA 私钥 PEM 解析失败")
	}
	priKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("RSA 私钥解析失败(需 PKCS1 格式): %w", err)
	}
	pubDER, err := x509.MarshalPKIXPublicKey(&priKey.PublicKey)
	if err != nil {
		return fmt.Errorf("派生 RSA 公钥失败: %w", err)
	}
	pubBytes := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubDER})

	rsaPriv = privBytes
	rsaPub = pubBytes
	Conf.System.RSAPrivateBytes = privBytes
	Conf.System.RSAPublicBytes = pubBytes
	return nil
}

func fileExists(p string) bool {
	info, err := os.Stat(p)
	return err == nil && !info.IsDir()
}

func persistPrivateKey(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

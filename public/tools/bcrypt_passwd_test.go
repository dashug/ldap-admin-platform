package tools

import "testing"

func TestHashAndVerifyPassword_Bcrypt(t *testing.T) {
	plain := "S3cure#Pass!"
	hash := HashPassword(plain)
	if hash == "" {
		t.Fatal("HashPassword 返回空")
	}
	if !IsBcryptHash(hash) {
		t.Fatalf("生成的哈希未被识别为 bcrypt: %q", hash)
	}
	if !VerifyPassword(hash, plain) {
		t.Error("正确密码校验应通过")
	}
	if VerifyPassword(hash, "wrong-password") {
		t.Error("错误密码校验应失败")
	}
}

func TestIsBcryptHash(t *testing.T) {
	// 由 HashPassword 生成的真实 bcrypt 哈希，长度恒为 60
	realHash := HashPassword("any-password")
	cases := map[string]bool{
		realHash: true,
		"$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy": true,
		// 长度不足 60，应判为非 bcrypt（会回退到 RSA 比对路径）
		"$2a$short": false,
		// 典型的 RSA 密文（二进制串）不应被误判为 bcrypt
		"some-random-rsa-ciphertext-bytes": false,
		"": false,
	}
	for in, want := range cases {
		if got := IsBcryptHash(in); got != want {
			t.Errorf("IsBcryptHash(%q) = %v, want %v", in, got, want)
		}
	}
}

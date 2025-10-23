package utils

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	// 初始化测试配置
	InitConfig()
	
	password := "testpassword"
	hash, err := HashPassword(password)
	if err != nil {
		t.Errorf("HashPassword() error = %v", err)
		return
	}
	
	// 验证哈希不为空
	if hash == "" {
		t.Error("HashPassword() returned empty hash")
	}
	
	// 验证哈希不等于原始密码
	if hash == password {
		t.Error("HashPassword() returned unhashed password")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	// 初始化测试配置
	InitConfig()
	
	password := "testpassword"
	hash, err := HashPassword(password)
	if err != nil {
		t.Errorf("HashPassword() error = %v", err)
		return
	}
	
	// 测试正确密码
	if !CheckPasswordHash(password, hash) {
		t.Error("CheckPasswordHash() failed to validate correct password")
	}
	
	// 测试错误密码
	if CheckPasswordHash("wrongpassword", hash) {
		t.Error("CheckPasswordHash() validated wrong password")
	}
}
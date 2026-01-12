package utils

import (
	"crypto/md5"
	"encoding/hex"
)

const (
	// SALT 密码加密盐值
	SALT = "yupi"
)

// EncryptPassword 加密密码
func EncryptPassword(password string) string {
	hasher := md5.New()
	hasher.Write([]byte(SALT + password))
	return hex.EncodeToString(hasher.Sum(nil))
}

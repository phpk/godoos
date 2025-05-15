package utils

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	r "math/rand"
	"regexp"
	"time"
)

// generateSalt 生成随机盐
func GenerateSalt(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// GenerateRandomSixDigitNumber 生成一个6位的随机正整数
func GenerateRandomSixDigitNumber() string {
	return fmt.Sprintf("%06v", r.New(r.NewSource(time.Now().UnixNano())).Int31n(1000000))
}

// hashPassword 生成哈希密码
func HashPassword(password, salt string) string {
	hash := sha256.New()
	hash.Write([]byte(password + salt))
	return base64.URLEncoding.EncodeToString(hash.Sum(nil))
}
func StringToMD5(src string) string {
	m := md5.New()
	m.Write([]byte(src))
	res := hex.EncodeToString(m.Sum(nil))
	return res
}

// IsEmail 校验邮箱格式
func IsEmail(email string) bool {
	// 正则表达式匹配邮箱格式
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

// HashAndMD5 先使用SHA-256对输入字符串进行哈希，然后将结果进行MD5加密
func HashAndMD5(input string) (string, error) {
	// 使用SHA-256进行哈希
	sha256Hash := sha256.New()
	if _, err := sha256Hash.Write([]byte(input)); err != nil {
		return "", err
	}
	sha256Result := sha256Hash.Sum(nil)

	// 将SHA-256的结果转换为字符串形式
	sha256String := hex.EncodeToString(sha256Result)

	// 对SHA-256的结果再次使用MD5进行哈希
	md5Hash := md5.New()
	if _, err := md5Hash.Write([]byte(sha256String)); err != nil {
		return "", err
	}
	md5Result := md5Hash.Sum(nil)

	// 返回最终的MD5哈希值
	return hex.EncodeToString(md5Result), nil
}

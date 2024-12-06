package libs

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
)

func IsEncryptedFile(data string) bool {
	if len(data) < 2 {
		return false
	}
	// 检查是否以 @ 开头
	if !strings.HasPrefix(data, "@") {
		return false
	}
	// 去掉开头的 @
	data = data[1:]
	// 分割加密的私钥、二级密码和加密的文本
	parts := strings.SplitN(data, "@", 3)
	if len(parts) != 3 {
		return false
	}

	hexEncodedPrivateKey := parts[0]
	encryptedSecondaryPassword := parts[1]
	encryptedText := parts[2]

	// 检查各个部分是否都不为空
	if hexEncodedPrivateKey == "" || encryptedSecondaryPassword == "" || encryptedText == "" {
		return false
	}

	// 检查十六进制字符串是否有效
	base64Str, err := hex.DecodeString(hexEncodedPrivateKey)
	if err != nil {
		return false
	}
	// 尝试将 Base64 字符串解码为字节切片
	_, err = base64.URLEncoding.DecodeString(string(base64Str))
	if err != nil {
		return false
	}

	// 检查二级密码和加密文本是否能被正确解码
	_, err = base64.URLEncoding.DecodeString(encryptedSecondaryPassword)
	if err != nil {
		return false
	}

	_, err = base64.URLEncoding.DecodeString(encryptedText)
	return err == nil
}

// EncodeFile 加密文件
func EncodeFile(password string, longText string) (string, error) {
	privateKey, publicKey, err := GenerateRSAKeyPair(2048) // 使用2048位密钥更安全
	if err != nil {
		fmt.Println("生成RSA密钥对失败:", err)
		return "", err
	}

	// 使用密码加密私钥
	pwd, err := hashAndMD5(password)
	if err != nil {
		return "", fmt.Errorf("加密密码失败:%v", err)
	}
	encryptedPrivateKey, err := EncryptPrivateKey(privateKey, pwd)
	if err != nil {
		fmt.Println("加密私钥失败:", err)
		return "", fmt.Errorf("加密私钥失败")
	}

	// 对 encryptedPrivateKey 进行 Base64 编码
	base64EncryptedPrivateKey := base64.URLEncoding.EncodeToString([]byte(encryptedPrivateKey))

	// 将 Base64 编码后的字符串转换为十六进制字符串
	hexEncodedPrivateKey := hex.EncodeToString([]byte(base64EncryptedPrivateKey))

	// 生成二级密码
	secondaryPassword, err := GenerateSecondaryPassword(32) // 32字节的二级密码
	if err != nil {
		return "", err
	}

	// 使用公钥加密二级密码
	encryptedSecondaryPassword, err := EncryptSecondaryPassword(secondaryPassword, publicKey)
	if err != nil {
		return "", fmt.Errorf("加密二级密码失败:%v", err)
	}

	// 使用二级密码加密数据
	encryptedText, err := EncryptDataWithCBC(longText, secondaryPassword)
	if err != nil {
		return "", fmt.Errorf("加密数据失败:%v", err)
	}

	return "@" + hexEncodedPrivateKey + "@" + encryptedSecondaryPassword + "@" + encryptedText, nil
}

// DecodeFile 解密文件
func DecodeFile(password string, encryptedData string) (string, error) {
	// 去掉开头的@
	if !strings.HasPrefix(encryptedData, "@") {
		return "", fmt.Errorf("无效的加密数据格式")
	}
	encryptedData = encryptedData[1:]

	// 分割加密的私钥、加密的二级密码和加密的文本
	parts := strings.SplitN(encryptedData, "@", 3)
	if len(parts) != 3 {
		return "", fmt.Errorf("无效的加密数据格式")
	}

	hexEncodedPrivateKey := parts[0]
	encryptedSecondaryPassword := parts[1]
	encryptedText := parts[2]

	// 将十六进制字符串转换回 Base64 编码字符串
	base64DecodedPrivateKey, err := hex.DecodeString(hexEncodedPrivateKey)
	if err != nil {
		return "", fmt.Errorf("十六进制字符串解码失败:%v", err)
	}

	// 将 Base64 编码字符串解码回原始的 encryptedPrivateKey
	encryptedPrivateKey, err := base64.URLEncoding.DecodeString(string(base64DecodedPrivateKey))
	if err != nil {
		return "", fmt.Errorf("Base64字符串解码失败:%v", err)
	}

	pwd, err := hashAndMD5(password)
	if err != nil {
		return "", fmt.Errorf("加密密码失败:%v", err)
	}

	// 解密私钥
	privateKey, err := DecryptPrivateKey(string(encryptedPrivateKey), pwd)
	if err != nil {
		return "", fmt.Errorf("解密私钥失败:%v", err)
	}

	// 解密二级密码
	secondaryPassword, err := DecryptSecondaryPassword(encryptedSecondaryPassword, privateKey)
	if err != nil {
		return "", fmt.Errorf("解密二级密码失败:%v", err)
	}

	// 使用二级密码解密数据
	decryptedText, err := DecryptDataWithCBC(encryptedText, secondaryPassword)
	if err != nil {
		return "", fmt.Errorf("解密文本失败:%v", err)
	}

	return decryptedText, nil
}

// GenerateRSAKeyPair 生成RSA密钥对
func GenerateRSAKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}

	publicKey := &privateKey.PublicKey

	return privateKey, publicKey, nil
}
func GenertePrivateKey(password string, privateKey *rsa.PrivateKey) (*rsa.PrivateKey, error) {
	encryptedPrivateKey, err := EncryptPrivateKey(privateKey, password)
	if err != nil {
		fmt.Println("加密私钥失败:", err)
		return nil, fmt.Errorf("加密私钥失败")
	}
	decryptedPrivateKey, err := DecryptPrivateKey(encryptedPrivateKey, password)
	if err != nil {
		fmt.Println("解密私钥失败:", err)
		return nil, fmt.Errorf("加密私钥失败")
	}
	return decryptedPrivateKey, nil
}

// EncryptWithPublicKey 使用公钥加密数据
func EncryptWithPublicKey(data []byte, publicKey *rsa.PublicKey) (string, error) {
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, data, nil)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// DecryptWithPrivateKey 使用私钥解密数据
func DecryptWithPrivateKey(encryptedData string, privateKey *rsa.PrivateKey) ([]byte, error) {
	ciphertext, err := base64.URLEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, err
	}

	plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

// hashAndMD5 使用 MD5 哈希密码
func hashAndMD5(password string) (string, error) {
	hasher := md5.New()
	_, err := hasher.Write([]byte(password))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

// EncryptPrivateKey 使用AES加密RSA私钥
func EncryptPrivateKey(privateKey *rsa.PrivateKey, password string) (string, error) {
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	key := []byte(password)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, privateKeyBytes, nil)
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// DecryptPrivateKey 使用AES解密RSA私钥
func DecryptPrivateKey(encryptedPrivateKey string, password string) (*rsa.PrivateKey, error) {
	ciphertext, err := base64.URLEncoding.DecodeString(encryptedPrivateKey)
	if err != nil {
		return nil, err
	}

	key := []byte(password)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	privateKeyBytes, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

// DecryptSecondaryPassword 使用私钥解密二级密码
func DecryptSecondaryPassword(encryptedSecondaryPassword string, privateKey *rsa.PrivateKey) (string, error) {
	ciphertext, err := base64.URLEncoding.DecodeString(encryptedSecondaryPassword)
	if err != nil {
		return "", err
	}

	plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

// DecodeFileWithSecondaryPassword 解密文件
func DecodeFileWithSecondaryPassword(password string, encryptedData string) (string, error) {
	// 去掉开头的@
	if !strings.HasPrefix(encryptedData, "@") {
		return "", fmt.Errorf("无效的加密数据格式")
	}
	encryptedData = encryptedData[1:]

	// 分割加密的二级密码和加密的文本
	parts := strings.SplitN(encryptedData, "@", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("无效的加密数据格式")
	}

	encryptedSecondaryPassword := parts[0]
	encryptedText := parts[1]

	// 使用私钥解密二级密码
	privateKey, err := DecryptPrivateKey(password, password)
	if err != nil {
		return "", fmt.Errorf("解密私钥失败:%v", err)
	}

	// 使用私钥解密二级密码
	secondaryPassword, err := DecryptSecondaryPassword(encryptedSecondaryPassword, privateKey)
	if err != nil {
		return "", fmt.Errorf("解密二级密码失败:%v", err)
	}

	// 使用二级密码解密文本
	decryptedText, err := DecryptDataWithCBC(encryptedText, secondaryPassword)
	if err != nil {
		return "", fmt.Errorf("解密文本失败:%v", err)
	}

	return decryptedText, nil
}

// GenerateSecondaryPassword 生成一个随机的二级密码
func GenerateSecondaryPassword(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// EncryptWithPublicKey 使用公钥加密二级密码
func EncryptSecondaryPassword(secondaryPassword string, publicKey *rsa.PublicKey) (string, error) {
	data := []byte(secondaryPassword)
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, data, nil)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// EncryptDataWithCBC 使用二级密码和CBC模式加密数据
func EncryptDataWithCBC(data, secondaryPassword string) (string, error) {
	// 确保 secondaryPassword 长度为 32 字节
	if len(secondaryPassword) < 32 {
		return "", fmt.Errorf("secondary password too short")
	}
	key := []byte(secondaryPassword[:32])

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(data), nil)
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// DecryptDataWithCBC 使用二级密码和CBC模式解密数据
func DecryptDataWithCBC(encryptedData, secondaryPassword string) (string, error) {
	// 确保 secondaryPassword 长度为 32 字节
	if len(secondaryPassword) < 32 {
		return "", fmt.Errorf("secondary password too short")
	}
	key := []byte(secondaryPassword[:32])

	ciphertext, err := base64.URLEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

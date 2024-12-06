package libs

import (
	"bytes"
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
	"log"
	"strings"
)

func GetEncryptedCode(data string) (string, error) {
	// 检查是否以 @ 开头
	if !strings.HasPrefix(data, "@") {
		return "", fmt.Errorf("invalid input format")
	}
	// 去掉开头的 @
	data = data[1:]
	// 分割加密的私钥和加密的文本
	parts := strings.SplitN(data, "@", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid input format")
	}
	return parts[0], nil
}
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
	// 分割加密的私钥和加密的文本
	parts := strings.SplitN(data, "@", 2)
	if len(parts) != 2 {
		return false
	}
	// 检查加密私钥和加密文本是否都不为空
	hexEncodedPrivateKey := parts[0]
	// 检查十六进制字符串是否有效
	base64Str, err := hex.DecodeString(hexEncodedPrivateKey)
	if err != nil {
		return false
	}
	// 尝试将 Base64 字符串解码为字节切片
	_, err = base64.URLEncoding.DecodeString(string(base64Str))
	return err == nil
}

// EncodeFile 加密文件
func EncodeFile(password string, longText string) (string, error) {
	privateKey, publicKey, err := GenerateRSAKeyPair(1024)
	if err != nil {
		fmt.Println("生成RSA密钥对失败:", err)
		return "", err
	}
	// 使用公钥加密
	encryptedText := ""
	if len(longText) > 0 {
		encryptedText, err = EncryptLongText(longText, publicKey)
		if err != nil {
			return "", fmt.Errorf("加密失败:%v", err)
		}
	}

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

	return "@" + hexEncodedPrivateKey + "@" + encryptedText, nil
}

// DecodeFile 解密文件
func DecodeFile(password string, encryptedData string) (string, error) {
	// 去掉开头的@
	// log.Printf("encryptedData: %s", encryptedData)
	if !strings.HasPrefix(encryptedData, "@") {
		return "", fmt.Errorf("无效的加密数据格式")
	}
	encryptedData = encryptedData[1:]

	// 分割加密的私钥和加密的文本
	parts := strings.SplitN(encryptedData, "@", 2)
	log.Printf("parts:%v", parts)
	if len(parts) == 1 {
		return "", nil
	}
	if len(parts) != 2 {
		return "", fmt.Errorf("无效的加密数据格式")
	}

	hexEncodedPrivateKey := parts[0]
	encryptedText := parts[1]
	if len(encryptedText) == 0 {
		return "", nil
	}

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

	// 解密文本
	decryptedText, err := DecryptLongText(encryptedText, privateKey)
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

// EncryptLongText 使用公钥加密长文本
func EncryptLongText(longText string, publicKey *rsa.PublicKey) (string, error) {
	// 分块加密
	blockSize := publicKey.N.BitLen()/8 - 2*sha256.Size - 2
	chunks := splitIntoChunks([]byte(longText), blockSize)

	var encryptedChunks []string
	for _, chunk := range chunks {
		encryptedChunk, err := EncryptWithPublicKey(chunk, publicKey)
		if err != nil {
			return "", err
		}
		encryptedChunks = append(encryptedChunks, encryptedChunk)
	}

	return strings.Join(encryptedChunks, ":"), nil
}

// DecryptLongText 使用私钥解密长文本
func DecryptLongText(encryptedLongText string, privateKey *rsa.PrivateKey) (string, error) {
	// 分块解密
	encryptedChunks := strings.Split(encryptedLongText, ":")

	var decryptedChunks [][]byte
	for _, encryptedChunk := range encryptedChunks {
		decryptedChunk, err := DecryptWithPrivateKey(encryptedChunk, privateKey)
		if err != nil {
			return "", err
		}
		decryptedChunks = append(decryptedChunks, decryptedChunk)
	}

	return string(bytes.Join(decryptedChunks, nil)), nil
}

// splitIntoChunks 将数据分割成指定大小的块
func splitIntoChunks(data []byte, chunkSize int) [][]byte {
	var chunks [][]byte
	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize
		if end > len(data) {
			end = len(data)
		}
		chunks = append(chunks, data[i:end])
	}
	return chunks
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

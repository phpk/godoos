package vector

import (
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

const metadataFileName = "00000000"

// hash2hex 将字符串转换为 SHA256 哈希并返回前 8 位的十六进制表示。
func hash2hex(name string) string {
	hash := sha256.Sum256([]byte(name))
	return hex.EncodeToString(hash[:4])
}

// persistToFile 将对象持久化到文件。支持 Gob 序列化、Gzip 压缩和 AES-GCM 加密。
func persistToFile(filePath string, obj any, compress bool, encryptionKey string) error {
	if filePath == "" {
		return fmt.Errorf("文件路径为空")
	}
	if encryptionKey != "" && len(encryptionKey) != 32 {
		return errors.New("加密密钥必须是 32 字节长")
	}

	// 确保父目录存在
	if err := os.MkdirAll(filepath.Dir(filePath), 0o700); err != nil {
		return fmt.Errorf("无法创建父目录: %w", err)
	}

	// 打开或创建文件
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("无法创建文件: %w", err)
	}
	defer f.Close()

	return persistToWriter(f, obj, compress, encryptionKey)
}

// persistToWriter 将对象持久化到 io.Writer。支持 Gob 序列化、Gzip 压缩和 AES-GCM 加密。
func persistToWriter(w io.Writer, obj any, compress bool, encryptionKey string) error {
	if encryptionKey != "" && len(encryptionKey) != 32 {
		return errors.New("加密密钥必须是 32 字节长")
	}

	var chainedWriter io.Writer
	if encryptionKey == "" {
		chainedWriter = w
	} else {
		chainedWriter = &bytes.Buffer{}
	}

	var gzw *gzip.Writer
	var enc *gob.Encoder
	if compress {
		gzw = gzip.NewWriter(chainedWriter)
		enc = gob.NewEncoder(gzw)
	} else {
		enc = gob.NewEncoder(chainedWriter)
	}

	if err := enc.Encode(obj); err != nil {
		return fmt.Errorf("无法编码或写入对象: %w", err)
	}

	if compress {
		if err := gzw.Close(); err != nil {
			return fmt.Errorf("无法关闭 Gzip 写入器: %w", err)
		}
	}

	if encryptionKey == "" {
		return nil
	}

	block, err := aes.NewCipher([]byte(encryptionKey))
	if err != nil {
		return fmt.Errorf("无法创建 AES 密码: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("无法创建 GCM 包装器: %w", err)
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return fmt.Errorf("无法读取随机字节作为 nonce: %w", err)
	}

	buf := chainedWriter.(*bytes.Buffer)
	encrypted := gcm.Seal(nonce, nonce, buf.Bytes(), nil)
	if _, err := w.Write(encrypted); err != nil {
		return fmt.Errorf("无法写入加密数据: %w", err)
	}

	return nil
}

// readFromFile 从文件中读取对象。支持 Gob 反序列化、Gzip 解压和 AES-GCM 解密。
func readFromFile(filePath string, obj any, encryptionKey string) error {
	if filePath == "" {
		return fmt.Errorf("文件路径为空")
	}
	if encryptionKey != "" && len(encryptionKey) != 32 {
		return errors.New("加密密钥必须是 32 字节长")
	}

	r, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("无法打开文件: %w", err)
	}
	defer r.Close()

	return readFromReader(r, obj, encryptionKey)
}

// readFromReader 从 io.Reader 中读取对象。支持 Gob 反序列化、Gzip 解压和 AES-GCM 解密。
func readFromReader(r io.ReadSeeker, obj any, encryptionKey string) error {
	if encryptionKey != "" && len(encryptionKey) != 32 {
		return errors.New("加密密钥必须是 32 字节长")
	}

	var chainedReader io.Reader
	if encryptionKey != "" {
		encrypted, err := io.ReadAll(r)
		if err != nil {
			return fmt.Errorf("无法读取数据: %w", err)
		}
		block, err := aes.NewCipher([]byte(encryptionKey))
		if err != nil {
			return fmt.Errorf("无法创建 AES 密码: %w", err)
		}
		gcm, err := cipher.NewGCM(block)
		if err != nil {
			return fmt.Errorf("无法创建 GCM 包装器: %w", err)
		}
		nonceSize := gcm.NonceSize()
		if len(encrypted) < nonceSize {
			return fmt.Errorf("加密数据太短")
		}
		nonce, ciphertext := encrypted[:nonceSize], encrypted[nonceSize:]
		data, err := gcm.Open(nil, nonce, ciphertext, nil)
		if err != nil {
			return fmt.Errorf("无法解密数据: %w", err)
		}
		chainedReader = bytes.NewReader(data)
	} else {
		chainedReader = r
	}

	magicNumber := make([]byte, 2)
	_, err := chainedReader.Read(magicNumber)
	if err != nil {
		return fmt.Errorf("无法读取魔数以确定是否压缩: %w", err)
	}
	compressed := magicNumber[0] == 0x1f && magicNumber[1] == 0x8b

	// 重置读取器位置
	if s, ok := chainedReader.(io.Seeker); !ok {
		return fmt.Errorf("读取器不支持寻址")
	} else {
		_, err := s.Seek(0, 0)
		if err != nil {
			return fmt.Errorf("无法重置读取器: %w", err)
		}
	}

	if compressed {
		gzr, err := gzip.NewReader(chainedReader)
		if err != nil {
			return fmt.Errorf("无法创建 Gzip 读取器: %w", err)
		}
		defer gzr.Close()
		chainedReader = gzr
	}

	dec := gob.NewDecoder(chainedReader)
	if err := dec.Decode(obj); err != nil {
		return fmt.Errorf("无法解码对象: %w", err)
	}

	return nil
}

// removeFile 删除指定路径的文件。如果文件不存在，则无操作。
func removeFile(filePath string) error {
	if filePath == "" {
		return fmt.Errorf("文件路径为空")
	}

	err := os.Remove(filePath)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("无法删除文件 %q: %w", filePath, err)
	}

	return nil
}

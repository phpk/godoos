/*
 * GodoOS - A lightweight cloud desktop
 * Copyright (C) 2024 https://godoos.com
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 2.1 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package files

import (
	"fmt"
	"godo/libs"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Common response structure

type OsFileInfo struct {
	IsFile     bool        `json:"isFile"`
	IsDir      bool        `json:"isDirectory"`
	IsSymlink  bool        `json:"isSymlink"`
	Size       int64       `json:"size"`
	ModTime    time.Time   `json:"modTime"`
	AccessTime time.Time   `json:"atime"`
	CreateTime time.Time   `json:"birthtime"`
	Mode       os.FileMode `json:"mode"`
	Name       string      `json:"name"`         // 文件名
	Path       string      `json:"path"`         // 文件路径
	OldPath    string      `json:"oldPath"`      // 旧的文件路径
	ParentPath string      `json:"parentPath"`   // 父目录路径
	Content    string      `json:"content"`      // 文件内容
	Ext        string      `json:"ext"`          // 文件扩展名
	Title      string      `json:"title"`        // 文件名（不包含扩展名）
	ID         int         `json:"id,omitempty"` // 文件ID（可选）
	IsPwd      bool        `json:"isPwd"`        // 是否加密
}

// validateFilePath 验证路径不为空
func validateFilePath(path string) error {
	if path == "" {
		return fmt.Errorf("empty path")
	}
	return nil
}

// validateFilePathPair 验证源路径和目标路径不为空
func validateFilePathPair(src, dst string) error {
	if err := validateFilePath(src); err != nil {
		return err
	}
	if err := validateFilePath(dst); err != nil {
		return err
	}
	return nil
}

// ReadDir reads a directory and returns a list of files
func ReadDir(basePath, path string) ([]os.DirEntry, error) {
	fullPath := filepath.Join(basePath, path)
	return os.ReadDir(fullPath)
}

// Stat retrieves file information
func Stat(basePath, path string) (os.FileInfo, error) {
	fullPath := filepath.Join(basePath, path)
	return os.Stat(fullPath)
}

// Exists checks if a file or directory exists
func Exists(basePath, path string) bool {
	fullPath := filepath.Join(basePath, path)
	// _, err := os.Stat(fullPath)
	// return err == nil
	return libs.PathExists(fullPath)
}

// ReadFile reads a file's content
func ReadFile(basePath, path string) ([]byte, error) {
	fullPath := filepath.Join(basePath, path)
	return os.ReadFile(fullPath)
}

// Unlink removes a file
func Unlink(basePath, path string) error {
	fullPath := filepath.Join(basePath, path)
	return os.Remove(fullPath)
}

// Clear removes the entire filesystem
func Clear(basePath string) error {
	return os.RemoveAll(basePath)
}

// Rename renames a file
func Rename(basePath, oldPath, newPath string) error {
	return os.Rename(filepath.Join(basePath, oldPath), filepath.Join(basePath, newPath))
}

// Mkdir creates a directory
func Mkdir(basePath, dirPath string) error {
	fullPath := filepath.Join(basePath, dirPath)
	return os.MkdirAll(fullPath, 0755)
}

// Rmdir removes a directory
func Rmdir(basePath, dirPath string) error {
	fullPath := filepath.Join(basePath, dirPath)
	return os.RemoveAll(fullPath)
}
func CopyResource(src, dst string) error {
	// 检查源路径是否为目录或文件
	info, err := os.Stat(src)
	if err != nil {
		return err
	}

	// 如果源路径是一个目录，则递归复制整个目录
	if info.IsDir() {
		return CopyDirectory(src, dst)
	}

	// 如果源路径是一个文件，则直接复制文件
	return CopyFile(src, dst)
}

func CopyDirectory(src, dst string) error {
	// 创建目标目录
	if err := os.MkdirAll(dst, os.ModePerm); err != nil {
		return err
	}

	// 遍历源目录中的所有文件和子目录
	files, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, file := range files {
		// 构建源文件和目标文件的完整路径
		srcPath := filepath.Join(src, file.Name())
		dstPath := filepath.Join(dst, file.Name())

		// 递归复制子目录或复制文件
		if file.IsDir() {
			if err := CopyDirectory(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := CopyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// CopyFile copies a file
func CopyFile(src, dst string) error {
	sf, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sf.Close()

	df, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer df.Close()

	_, err = io.Copy(df, sf)
	if err != nil {
		return err
	}

	err = df.Sync()
	if err != nil {
		return err
	}

	return df.Close()
}

// WriteFile writes content to a file
func WriteFile(filePath string, content []byte, perm os.FileMode) error {
	return os.WriteFile(filePath, content, perm)
}

// AppendToFile appends content to a file
// AppendToFile appends the given content to the specified file.
// It assumes that 'content' is an io.Reader representing the data to append.
func AppendToFile(filePath string, content io.Reader) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		// 文件不存在或打开时发生其他错误
		return err
	}
	defer file.Close()

	// 将内容写入文件
	if _, err := io.Copy(file, content); err != nil {
		return err
	}

	return nil
}
func Chmod(basePath string, path string, mode os.FileMode) error {
	fullPath := filepath.Join(basePath, path)
	return os.Chmod(fullPath, mode)
}

// GetFileInfo 从fs.DirEntry或路径获取文件信息
func GetFileInfo(entry interface{}, basePath, parentPath string) (*OsFileInfo, error) {
	var fileInfo os.FileInfo
	var err error
	var filePath string

	switch v := entry.(type) {
	case fs.DirEntry:
		// 如果传入的是fs.DirEntry
		fileInfo, err = v.Info()
		if err != nil {
			return nil, fmt.Errorf("failed to get file info from DirEntry: %v", err)
		}
		filePath = filepath.Join(basePath, parentPath, fileInfo.Name())
	case string:
		// 如果传入的是路径字符串
		filePath = filepath.Join(basePath, parentPath, v)
		fileInfo, err = os.Stat(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to stat file from path: %v", err)
		}
	default:
		return nil, fmt.Errorf("unsupported type provided to GetFileInfo, expected fs.DirEntry or string")
	}

	osFileInfo := &OsFileInfo{
		Name: fileInfo.Name(),
		Path: filePath, // 根据实际情况调整，可能需要使用之前构造的完整路径filePath
	}

	// 初始化其他字段，如IsFile, IsDir, Size等
	osFileInfo.ID = 0
	osFileInfo.IsFile = !fileInfo.IsDir()
	osFileInfo.IsDir = fileInfo.IsDir()
	osFileInfo.Size = fileInfo.Size()
	osFileInfo.ModTime = fileInfo.ModTime()
	osFileInfo.AccessTime = fileInfo.ModTime() // 注意：os.FileInfo没有直接提供访问时间，这里用修改时间代替
	osFileInfo.CreateTime = fileInfo.ModTime() // 同上，仅作示例
	osFileInfo.Mode = fileInfo.Mode()

	osFileInfo.Path = strings.TrimPrefix(filePath, basePath)
	osFileInfo.OldPath = osFileInfo.Path
	osFileInfo.ParentPath = parentPath
	osFileInfo.Title = strings.TrimSuffix(osFileInfo.Name, filepath.Ext(osFileInfo.Name))
	osFileInfo.Ext = strings.TrimPrefix(filepath.Ext(osFileInfo.Name), ".")

	return osFileInfo, nil
}
func GetDesktopPath() string {
	return filepath.Join("C", "Users", "Desktop")
}

// RemoveFirstSlashOrBackslash 检查字符串的第一个字符是否为 '/' 或 '\'
// 如果是，则返回去掉该字符的剩余部分；如果不是，则返回原字符串。
func RemoveFirstSlashOrBackslash(s string) string {
	if len(s) > 0 && (s[0] == '/' || s[0] == '\\') {
		return s[1:]
	}
	return s
}
func CheckAddDesktop(filePath string) error {
	checkPath := RemoveFirstSlashOrBackslash(filepath.Dir(filePath))
	if checkPath != GetDesktopPath() {
		return fmt.Errorf("文件不在桌面目录下")
	}
	basePath, err := libs.GetOsDir()
	if err != nil {
		return fmt.Errorf("获取系统目录失败")
	}
	if !Exists(basePath, filePath) {
		return fmt.Errorf("文件不存在")
	}
	osFileInfo, err := GetFileInfo(filePath, basePath, "")
	if err != nil {
		return fmt.Errorf("获取文件信息失败")
	}
	err = AddDesktop(*osFileInfo, "Desktop")
	if err != nil {
		return fmt.Errorf("添加桌面失败")
	}
	return nil
}
func CheckDeleteDesktop(filePath string) error {
	checkPath := RemoveFirstSlashOrBackslash(filepath.Dir(filePath))
	if checkPath != GetDesktopPath() {
		return fmt.Errorf("文件不在桌面目录下")
	}
	basePath, err := libs.GetOsDir()
	if err != nil {
		return fmt.Errorf("获取系统目录失败")
	}
	if !Exists(basePath, filePath) {
		return fmt.Errorf("文件不存在")
	}
	osFileInfo, err := GetFileInfo(filePath, basePath, "")
	if err != nil {
		return fmt.Errorf("获取文件信息失败")
	}
	err = DeleteDesktop(osFileInfo.Name)
	if err != nil {
		return fmt.Errorf("添加桌面失败")
	}
	return nil
}

// 校验文件密码
func CheckFilePwd(fpwd, salt string) bool {
	// 1. 对输入的密码进行哈希
	pwd := libs.HashPassword(fpwd, salt)
	// 2. 获取存储的密码哈希
	oldPwd, isHas := libs.GetConfig("filePwd")
	if !isHas {
		slog.Error("文件密码获取失败")
		return false
	}
	// 3. 比对密码哈希
	return oldPwd == pwd
}

func IsHavePwd(pwd string) bool {
	if len(pwd) > 0 {
		return true
	} else {
		return false
	}
}

// salt值优先从server端获取，如果没有则从header获取
func GetSalt(r *http.Request) (string, error) {
	data, ishas := libs.GetConfig("salt")
	if ishas {
		// 断言成功则返回
		if salt, ok := data.(string); ok {
			return salt, nil
		}
		return "", fmt.Errorf("类型断言失败，期望类型为 string")
	}
	salt := r.Header.Get("salt")
	if salt != "" {
		return salt, nil
	}
	return "", fmt.Errorf("无法获取salt")
}

// 该函数用于获取文件是否加密的标志
func GetPwdFlag() (bool, error) {
	// 从配置中获取加密标志
	isPwd, has := libs.GetConfig("isPwd")
	if !has {
		// 如果没有设置加密标志，默认设置为不加密
		err := libs.SetConfigByName("isPwd", false)
		if err != nil {
			return false, err // 返回错误信息
		}
		return false, nil // 返回默认值和无错误
	}

	// 尝试将获取的值转换为布尔类型
	if pwdFlag, ok := isPwd.(bool); ok {
		return pwdFlag, nil // 返回加密标志和无错误
	}

	// 如果类型断言失败，返回默认值和错误信息
	return false, fmt.Errorf("类型断言失败，期望类型为 bool")
}

// 检查文件是否加密
func IsHaveHiddenFile(basePath, filePath string) bool {
	// 通过查找同名隐藏文件判断
	// 获取文件名和目录
	dir := filepath.Dir(filePath)
	fileName := filepath.Base(filePath)
	hiddenFilePath := filepath.Join(basePath, dir, "."+fileName)
	_, err := os.Stat(hiddenFilePath)
	return err == nil
}

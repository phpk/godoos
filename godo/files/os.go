package files

import (
	"fmt"
	"godo/libs"
	"io"
	"io/fs"
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

package files

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// 压缩文件到指定的目录
func zipCompress(folderPath, zipFilePath string) error {
	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	info, err := os.Stat(folderPath)
	if err != nil {
		return err
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(folderPath)
	}

	filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(folderPath, path)
		if err != nil {
			return err
		}

		if baseDir != "" {
			header.Name = filepath.Join(baseDir, relPath)
		} else {
			header.Name = relPath
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(writer, file)
		return err
	})

	return nil
}

// IsSupportedCompressionFormat 判断文件后缀是否为支持的压缩格式
func IsSupportedFormat(filename string) bool {
	supportedFormats := map[string]bool{
		".zip": true,
		".tar": true,
		".gz":  true,
	}

	for format := range supportedFormats {
		if strings.HasSuffix(strings.ToLower(filename), format) {
			return true
		}
	}
	return false
}

// CompressFileOrFolder 压缩文件或文件夹到指定的压缩文件
func CompressFileOrFolder(sourcePath, compressedFilePath string) error {
	if IsSupportedFormat(compressedFilePath) {
		return Encompress(sourcePath, compressedFilePath)
	}
	// 如果压缩文件格式不受支持，则返回错误或采取其他措施
	return fmt.Errorf("unsupported compression format for file: %s", compressedFilePath)
}

// Encompress 压缩文件或文件夹到指定的压缩格式
func Encompress(sourcePath, compressedFilePath string) error {
	switch filepath.Ext(compressedFilePath) {
	case ".zip":
		return zipCompress(sourcePath, compressedFilePath)
	case ".tar":
		return tarCompress(sourcePath, compressedFilePath)
	case ".gz":
		return tarGzCompress(sourcePath, compressedFilePath)
	default:
		return fmt.Errorf("unsupported compression format: %s", filepath.Ext(compressedFilePath))
	}
}

// tarCompress 使用 tar 格式压缩文件或文件夹
func tarCompress(folderPath, tarFilePath string) error {
	// 打开目标文件
	tarFile, err := os.Create(tarFilePath)
	if err != nil {
		return err
	}
	defer tarFile.Close()

	// 创建 tar.Writer
	tarWriter := tar.NewWriter(tarFile)
	defer tarWriter.Close()

	// 遍历文件夹并添加到 tar 文件
	err = filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}

		if err := tarWriter.WriteHeader(header); err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(tarWriter, file)
			if err != nil {
				return err
			}
		}
		return nil
	})

	return err
}

// tarGzCompress 使用 tar.gz 格式压缩文件或文件夹
func tarGzCompress(folderPath, tarGzFilePath string) error {
	// 打开目标文件
	gzFile, err := os.Create(tarGzFilePath)
	if err != nil {
		return err
	}
	defer gzFile.Close()

	// 创建 gzip.Writer
	gzWriter := gzip.NewWriter(gzFile)
	defer gzWriter.Close()

	// 创建 tar.Writer
	tarWriter := tar.NewWriter(gzWriter)
	defer tarWriter.Close()

	// 遍历文件夹并添加到 tar 文件
	err = filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}

		if err := tarWriter.WriteHeader(header); err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(tarWriter, file)
			if err != nil {
				return err
			}
		}
		return nil
	})

	return err
}

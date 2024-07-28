package files

import (
	"compress/bzip2"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"archive/tar"
	"archive/zip"
)

// Decompress 解压文件到指定目录，支持多种压缩格式
func Decompress(compressedFilePath, destPath string) error {
	ext := filepath.Ext(compressedFilePath)
	switch ext {
	case ".zip":
		return Unzip(compressedFilePath, destPath)
	case ".tar":
		tarFile, err := os.Open(compressedFilePath)
		if err != nil {
			return err
		}
		defer tarFile.Close()
		return Untar(tarFile, destPath)
	case ".gz":
		return Untargz(compressedFilePath, destPath)
	case ".bz2":
		return Untarbz2(compressedFilePath, destPath)
	default:
		return fmt.Errorf("unsupported compression format: %s", ext)
	}
}

// Untar 解压.tar文件到指定目录
func Untar(reader io.Reader, destPath string) error {
	tarReader := tar.NewReader(reader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		target := filepath.Join(destPath, header.Name)
		if header.Typeflag == tar.TypeDir {
			if err := os.MkdirAll(target, os.FileMode(header.Mode)); err != nil {
				return err
			}
			continue
		}

		outFile, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(header.Mode))
		if err != nil {
			return err
		}
		defer outFile.Close()

		if _, err := io.Copy(outFile, tarReader); err != nil {
			return err
		}
	}
	return nil
}

// Untargz 解压.tar.gz文件到指定目录
func Untargz(targzFilePath, destPath string) error {
	gzFile, err := os.Open(targzFilePath)
	if err != nil {
		return err
	}
	defer gzFile.Close()

	gzReader, err := gzip.NewReader(gzFile)
	if err != nil {
		return err
	}
	defer gzReader.Close()

	return Untar(gzReader, destPath) // 传递gzReader给Untar
}

// Untarbz2 解压.tar.bz2文件到指定目录
func Untarbz2(tarbz2FilePath, destPath string) error {
	bz2File, err := os.Open(tarbz2FilePath)
	if err != nil {
		return err
	}
	defer bz2File.Close()

	bz2Reader := bzip2.NewReader(bz2File)
	return Untar(bz2Reader, destPath) // 传递bz2Reader给Untar
}

// DecompressZip 解压zip文件到指定目录
func Unzip(zipFilePath, destPath string) error {
	r, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		path := filepath.Join(destPath, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
			continue
		}

		os.MkdirAll(filepath.Dir(path), f.Mode())

		outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		defer outFile.Close()

		_, err = io.Copy(outFile, rc)
		if err != nil {
			return err
		}
	}

	return nil
}

// IsSupportedCompressionFormat 判断文件后缀是否为支持的压缩格式
func IsSupportedCompressionFormat(filename string) bool {
	supportedFormats := map[string]bool{
		".zip": true,
		".tar": true,
		".gz":  true,
		".bz2": true,
	}

	for format := range supportedFormats {
		if strings.HasSuffix(strings.ToLower(filename), format) {
			return true
		}
	}
	return false
}

// HandlerFile 根据文件后缀判断是否支持解压，支持则解压，否则移动文件
func HandlerFile(filePath string, destDir string) error {
	if !IsSupportedCompressionFormat(filePath) {
		// 移动文件
		newPath := filepath.Join(destDir, filepath.Base(filePath))
		return os.Rename(filePath, newPath)
	}
	// 解压文件
	return Decompress(filePath, destDir)
}

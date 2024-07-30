package files

import (
	"compress/bzip2"
	"compress/gzip"
	"fmt"
	"godo/libs"
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

		// Ensure that all directories in the path are created.
		if header.Typeflag == tar.TypeDir {
			if !libs.PathExists(target) {
				if err := os.MkdirAll(target, os.FileMode(header.Mode)); err != nil {
					return err
				}
			}

			continue
		}

		// If it's a regular file, ensure its parent directories exist.
		if !libs.PathExists(target) {
			if err := os.MkdirAll(filepath.Dir(target), os.ModePerm); err != nil {
				return err
			}
		}

		// Create the file and write the contents from the tar archive.
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

		// Ensure parent directory exists before creating the file or directory.
		if !libs.PathExists(filepath.Dir(path)) {
			if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
				return err
			}
		}

		if f.FileInfo().IsDir() {
			// Create the directory itself if it doesn't exist yet.
			if !libs.PathExists(path) {
				if err := os.Mkdir(path, f.Mode()); err != nil {
					return err
				}
			}

			continue
		}

		// Open the file for writing, truncating if it already exists.
		outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		defer outFile.Close()

		// Write the contents of the zip file entry to the output file.
		if _, err := io.Copy(outFile, rc); err != nil {
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
	// log.Printf("HandlerFile: %s", filePath)
	// log.Printf("destDir: %s", destDir)
	if IsSupportedCompressionFormat(filePath) {
		// 移动文件
		// newPath := filepath.Join(destDir, filepath.Base(filePath))
		// return os.Rename(filePath, newPath)
		// 解压文件
		return Decompress(filePath, destDir)
	}
	return nil

}

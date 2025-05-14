package libs

import (
	"errors"
	"godocms/common"
	"godocms/libs"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

func GetTempDir(pathname string) (string, error) {
	tempDir, err := os.MkdirTemp("", pathname)
	if err != nil {
		log.Println("Failed to create temporary directory:", err)
		return "./", err
	}

	log.Println("Temporary directory created:", tempDir)
	return tempDir, nil
}
func getConvertDir() string {
	runOsDir := common.GetRunOsDir()
	return filepath.Join(runOsDir, "goconv")
}
func getRapidDir() (string, error) {
	convertDir := getConvertDir()
	var path string
	if runtime.GOOS == "windows" {
		path = filepath.Join(convertDir, "rapid", "RapidOcrOnnx.exe")
	} else {
		path = filepath.Join(convertDir, "rapid", "RapidOcrOnnx")
	}
	if libs.PathExists(path) {
		return path, nil
	} else {
		return "", errors.New("RapidOcrOnnx not found")
	}
}

func getRapidModelDir() (string, error) {
	convertDir := getConvertDir()
	path := filepath.Join(convertDir, "rapid", "models")
	if libs.PathExists(path) {
		return path, nil
	} else {
		return "", errors.New("RapidOcrOnnx model not found")
	}
}
func getXpdfDir(exename string) (string, error) {
	convertDir := getConvertDir()
	var path string
	if runtime.GOOS == "windows" {
		path = filepath.Join(convertDir, "pdf", exename+".exe")
	} else {
		path = filepath.Join(convertDir, "pdf", exename)
	}
	if libs.PathExists(path) {
		return path, nil
	} else {
		return "", errors.New("pdf convert exe not found")
	}
}

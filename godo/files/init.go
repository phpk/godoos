package files

import (
	"encoding/json"
	"fmt"
	"godo/libs"
	"os"
	"path/filepath"
	"time"
)

func RecoverOsSystem() error {
	basePath, err := libs.GetOsDir()
	if err != nil {
		return fmt.Errorf("GetOsDir error: %v", err)
	}

	osCpath := filepath.Join(basePath, "C")
	if libs.PathExists(osCpath) {
		err := os.RemoveAll(osCpath)
		if err != nil {
			return fmt.Errorf("RemoveAll error: %v", err)
		}
	}
	err = InitOsSystem()
	if err != nil {
		return fmt.Errorf("InitOsSystem error: %v", err)
	}
	return nil
}
func InitOsSystem() error {
	basePath, err := libs.GetOsDir()
	if err != nil {
		return fmt.Errorf("GetOsDir error: %v", err)
	}
	if !libs.PathExists(basePath) {
		if err := os.MkdirAll(basePath, 0755); err != nil {
			return fmt.Errorf("Mkdir error: %v", err)
		}
	}
	osCpath := filepath.Join(basePath, "C")
	if !libs.PathExists(osCpath) {
		if err := os.MkdirAll(osCpath, 0755); err != nil {
			return fmt.Errorf("Mkdir error: %v", err)
		}
		baseOsDir := []string{
			"D",
			"E",
			"B",
		}
		for _, dir := range baseOsDir {
			dirPath := filepath.Join(basePath, dir)
			if !libs.PathExists(dirPath) {
				if err := os.MkdirAll(dirPath, 0755); err != nil {
					return fmt.Errorf("Mkdir error: %v", err)
				}
			}
		}
		systemPath := filepath.Join(osCpath, "System")
		if !libs.PathExists(systemPath) {
			if err := os.MkdirAll(systemPath, 0755); err != nil {
				return fmt.Errorf("Mkdir error: %v", err)
			}
		}
		userPath := filepath.Join(osCpath, "Users")
		if !libs.PathExists(userPath) {
			if err := os.MkdirAll(userPath, 0755); err != nil {
				return fmt.Errorf("Mkdir error: %v", err)
			}
			InitPaths := []string{
				"Desktop",
				"Menulist",
				"Documents",
				"Downloads",
				"Music",
				"Pictures",
				"Videos",
				"Schedule",
				"Reciv",
			}
			for _, dir := range InitPaths {
				dirPath := filepath.Join(userPath, dir)
				if !libs.PathExists(dirPath) {
					if err := os.MkdirAll(dirPath, 0755); err != nil {
						return fmt.Errorf("Mkdir error: %v", err)
					}
				}
			}
			InitDocPath := []string{
				"Word",
				"Markdown",
				"PPT",
				"Baiban",
				"Kanban",
				"Execl",
				"Mind",
				"Screenshot",
				"ScreenRecording",
			}
			docpath := filepath.Join(userPath, "Documents")
			for _, dir := range InitDocPath {
				dirPath := filepath.Join(docpath, dir)
				if !libs.PathExists(dirPath) {
					if err := os.MkdirAll(dirPath, 0755); err != nil {
						return fmt.Errorf("Mkdir error: %v", err)
					}
				}
			}

			applist := GetInitRootList()
			desktopPath := filepath.Join(userPath, "Desktop")
			for _, app := range applist.Desktop {
				appPath := filepath.Join(desktopPath, app.Name)
				if !libs.PathExists(appPath) {
					if err := os.WriteFile(appPath, []byte(app.Content), 0644); err != nil {
						return fmt.Errorf("failed to write to file: %w", err)
					}
				}
			}
			menulistPath := filepath.Join(userPath, "Menulist")
			for _, app := range applist.Menulist {
				appPath := filepath.Join(menulistPath, app.Name)
				if !libs.PathExists(appPath) {
					if err := os.WriteFile(appPath, []byte(app.Content), 0644); err != nil {
						return fmt.Errorf("failed to write to file: %w", err)
					}
				}
			}
			content, err := json.MarshalIndent(applist, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal reqBodies to JSON: %w", err)
			}
			appPath, err := libs.GetAppDir()
			if err != nil {
				return fmt.Errorf("GetAppDir error: %v", err)
			}
			desktopFilePath := filepath.Join(appPath, "desktop.json")
			if err := os.WriteFile(desktopFilePath, content, 0644); err != nil {
				return fmt.Errorf("failed to write to file: %w", err)
			}
		}

	}

	return nil
}
func GetDesktop() (RootObject, error) {
	var applist RootObject
	desktoppath, err := GetDeskTopPath()
	if err != nil || !libs.PathExists(desktoppath) {
		return applist, fmt.Errorf("desktop.json not found")

	}
	content, err := os.ReadFile(desktoppath)
	if err != nil {
		return applist, fmt.Errorf("failed to read file: %w", err)
	}

	err = json.Unmarshal(content, &applist)
	if err != nil {
		return applist, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	return applist, nil
}
func GetDeskTopPath() (string, error) {
	appPath, err := libs.GetAppDir()
	if err != nil {
		return "", fmt.Errorf("GetAppDir error: %v", err)
	}
	return filepath.Join(appPath, "desktop.json"), nil
}
func WriteDesktop(rootInfo RootObject) error {
	rootInfo.UpdateTime = time.Now()
	content, err := json.MarshalIndent(rootInfo, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal reqBodies to JSON: %w", err)
	}
	desktoppath, err := GetDeskTopPath()
	if err != nil {
		return fmt.Errorf("GetDeskTopPath error: %v", err)
	}
	return os.WriteFile(desktoppath, content, 0644)
}
func AddDesktop(app OsFileInfo, position string) error {
	applist, err := GetDesktop()
	if err != nil {
		return fmt.Errorf("GetDesktop error: %v", err)
	}
	if position == "Desktop" {
		if !checkHasApp(applist.Desktop, app) {
			applist.Desktop = append(applist.Desktop, app)
		}
	} else if position == "Menulist" {
		if !checkHasApp(applist.Menulist, app) {
			applist.Menulist = append(applist.Menulist, app)
		}
	} else {
		return fmt.Errorf("position error")
	}
	//log.Printf("add app to desktop %v", applist)
	return WriteDesktop(applist)
}
func DeleteDesktop(name string) error {
	rootInfo, err := GetDesktop()
	if err != nil {
		return fmt.Errorf("GetDesktop error: %v", err)
	}
	indexToDelete := -1
	for i, item := range rootInfo.Desktop {
		if item.Name == name {
			indexToDelete = i
			break
		}
	}
	// 如果找到了要删除的元素
	if indexToDelete >= 0 {
		// 使用 append 删除元素
		rootInfo.Desktop = append(rootInfo.Desktop[:indexToDelete], rootInfo.Desktop[indexToDelete+1:]...)
	} else {
		return fmt.Errorf("item with name '%s' not found", name)
	}
	return WriteDesktop(rootInfo)
}
func checkHasApp(list []OsFileInfo, app OsFileInfo) bool {
	for _, item := range list {
		if item.Name == app.Name {
			return true
		}
	}
	return false
}

package store

import (
	"encoding/json"
	"fmt"
	"godo/libs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func InstallHandler(w http.ResponseWriter, r *http.Request) {
	pluginName := r.URL.Query().Get("name")
	if pluginName == "" {
		libs.ErrorMsg(w, "the app name is empty!")
		return
	}
	exePath := GetExePath(pluginName)
	//log.Printf("the app path is %s", exePath)
	if !libs.PathExists(exePath) {
		libs.ErrorMsg(w, "the app path is not exists!")
		return
	}
	installInfo, err := GetInstallInfo(pluginName)
	if err != nil {
		libs.ErrorMsg(w, "the install.json is error:"+err.Error())
		return
	}
	if pluginName != installInfo.Name {
		libs.ErrorMsg(w, "the app name must equal the install.json!")
		return
	}
	if !installInfo.NeedInstall {
		libs.SuccessMsg(w, "success", "the app is installed!")
		return
	}
	storeFile := filepath.Join(exePath, "store.json")
	if !libs.PathExists(storeFile) {
		libs.ErrorMsg(w, "the store.json is not exists!")
		return
	}
	var storeInfo StoreInfo
	content, err := os.ReadFile(storeFile)
	if err != nil {
		libs.ErrorMsg(w, "cant read the store.json!")
		return
	}
	//设置 info.json
	exePath = strings.ReplaceAll(exePath, "\\", "/")
	contentBytes := []byte(strings.ReplaceAll(string(content), "{exePath}", exePath))
	//log.Printf("====content: %s", string(contentBytes))
	err = json.Unmarshal(contentBytes, &storeInfo)
	if err != nil {
		log.Printf("Error during unmarshal: %v", err)
		libs.ErrorMsg(w, "the store.json is error: "+err.Error())
		return
	}
	replacePlaceholdersInCmds(&storeInfo)
	storeInfo.Name = installInfo.Name
	err = SaveInfoFile(storeInfo)
	if err != nil {
		libs.ErrorMsg(w, "the store info.json is error: "+err.Error())
		return
	}
	//设置config
	if libs.PathExists(storeInfo.Setting.ConfPath + ".tpl") {
		err = SaveStoreConfig(storeInfo, exePath)
		if err != nil {
			libs.ErrorMsg(w, "save the config is error!")
			return
		}
	}
	err = InstallStore(storeInfo)
	if err != nil {
		libs.ErrorMsg(w, "install the app is error!")
		return
	}
	var res string
	//复制static目录
	staticPath := filepath.Join(exePath, "static")
	if libs.PathExists(staticPath) {
		staticDir := libs.GetStaticDir()
		targetPath := filepath.Join(staticDir, pluginName)
		if !libs.PathExists(targetPath) {
			err = os.Rename(staticPath, targetPath)
			if err != nil {
				log.Printf("copy the static is error! %v", err)
			}
			iconPath := filepath.Join(targetPath, storeInfo.Icon)
			if libs.PathExists(iconPath) {
				res = "http://localhost:56780/static/" + pluginName + "/" + storeInfo.Icon
			}
		}

	}
	libs.SuccessMsg(w, res, "install the app success!")
}

func UnInstallHandler(w http.ResponseWriter, r *http.Request) {
	pluginName := r.URL.Query().Get("name")
	if pluginName == "" {
		libs.ErrorMsg(w, "the app name is empty!")
		return
	}
	log.Printf("uninstall the app %s", pluginName)
	err := StopCmd(pluginName)
	if err != nil {
		// libs.ErrorMsg(w, "stop the app is error!")
		// return
		log.Printf("stop the app is error! %s", err)
	}
	installInfo, err := GetInstallInfo(pluginName)
	if err != nil {
		libs.ErrorMsg(w, "the install.json is error:"+err.Error())
		return
	}
	if installInfo.IsDev {
		libs.SuccessMsg(w, "success", "uninstall the app success!")
		return
	}
	exePath := GetExePath(pluginName)
	//log.Printf("the app path is %s", exePath)
	if libs.PathExists(exePath) {
		err := os.RemoveAll(exePath)
		if err != nil {
			libs.ErrorMsg(w, "delete the app is error!")
			return
		}
	}
	staticDir := libs.GetStaticDir()
	staticPath := filepath.Join(staticDir, pluginName)
	if libs.PathExists(staticPath) {
		err := os.RemoveAll(staticPath)
		if err != nil {
			libs.ErrorMsg(w, "delete the static is error!")
			return
		}
	}
	libs.SuccessMsg(w, "success", "uninstall the app success!")

}
func GetInstallInfo(pluginName string) (InstallInfo, error) {
	var installInfo InstallInfo
	exePath := GetExePath(pluginName)
	installFile := filepath.Join(exePath, "install.json")
	if !libs.PathExists(installFile) {
		return installInfo, fmt.Errorf("install.json is not exist:%s", installFile)
	}
	content, err := os.ReadFile(installFile)
	if err != nil {
		return installInfo, err
	}
	err = json.Unmarshal(content, &installInfo)
	if err != nil {
		return installInfo, err
	}
	return installInfo, nil
}
func GetExePath(pluginName string) string {
	exeDir := libs.GetRunDir()
	return filepath.Join(exeDir, pluginName)
}

func convertValueToString(value any) string {
	var strValue string
	switch v := value.(type) {
	case string:
		strValue = v
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		strValue = fmt.Sprintf("%v", v)
	default:
		strValue = fmt.Sprintf("%v", v) // 处理其他类型，例如布尔型或更复杂的类型
	}
	return strValue
}

func InstallStore(storeInfo StoreInfo) error {
	err := SetEnvs(storeInfo.Install.Envs)
	if err != nil {
		return fmt.Errorf("failed to set install environment variable %s: %w", storeInfo.Name, err)
	}
	if len(storeInfo.Install.Cmds) > 0 {
		for _, cmdKey := range storeInfo.Install.Cmds {
			if _, ok := storeInfo.Cmds[cmdKey]; ok {
				// 如果命令存在，你可以进一步处理 cmds
				RunCmds(storeInfo, cmdKey)
			}
		}

	}
	return nil
}
func RunCmd(storeInfo StoreInfo, cmd Cmd) error {
	if cmd.Name == "stop" {
		runStop(storeInfo)
	}
	if cmd.Name == "start" {
		runStart(storeInfo)
	}
	if cmd.Name == "restart" {
		runRestart(storeInfo)
	}
	if cmd.Name == "exec" {
		runExec(storeInfo, cmd)
	}
	if cmd.Name == "writeFile" {
		err := WriteFile(cmd)
		if err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}
	}
	if cmd.Name == "changeFile" {
		err := ChangeFile(storeInfo, cmd)
		if err != nil {
			return fmt.Errorf("failed to change to file: %w", err)
		}
	}
	if cmd.Name == "deleteFile" {
		err := DeleteFile(cmd)
		if err != nil {
			return fmt.Errorf("failed to delete file: %w", err)
		}
	}
	if cmd.Name == "unzip" {
		err := Unzip(cmd)
		if err != nil {
			return fmt.Errorf("failed to unzip file: %w", err)
		}
	}
	if cmd.Name == "zip" {
		err := Zip(cmd)
		if err != nil {
			return fmt.Errorf("failed to unzip file: %w", err)
		}
	}
	if cmd.Name == "mkdir" {
		err := MkDir(cmd)
		if err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}
	return nil
}
func RunCmds(storeInfo StoreInfo, cmdKey string) error {
	cmds := storeInfo.Cmds[cmdKey]
	for _, cmd := range cmds {
		if _, ok := storeInfo.Cmds[cmd.Name]; ok {
			RunCmds(storeInfo, cmd.Name)
		} else {
			RunCmd(storeInfo, cmd)
		}
	}
	return nil
}

func SetEnvs(envs []Item) error {
	if len(envs) > 0 {
		for _, item := range envs {
			strValue := convertValueToString(item.Value)
			if strValue != "" {
				err := os.Setenv(item.Name, strValue)
				if err != nil {
					return fmt.Errorf("failed to set environment variable %s: %w", item.Name, err)
				}
			}
		}
	}
	return nil
}
func SaveInfoFile(storeInfo StoreInfo) error {
	exePath := GetExePath(storeInfo.Name)
	infoFile := filepath.Join(exePath, "info.json")
	return SaveStoreInfo(storeInfo, infoFile)
}
func SaveStoreInfo(storeInfo StoreInfo, infoFile string) error {
	// 使用 json.MarshalIndent 直接获取内容的字节切片
	content, err := json.MarshalIndent(storeInfo, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal reqBodies to JSON: %w", err)
	}
	if err := os.WriteFile(infoFile, content, 0644); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}
	return nil
}
func SaveStoreConfig(storeInfo StoreInfo, exePath string) error {
	content, err := os.ReadFile(storeInfo.Setting.ConfPath + ".tpl")
	if err != nil {
		return fmt.Errorf("failed to read tpl file: %w", err)
	}
	contentstr := strings.ReplaceAll(string(content), "{exePath}", exePath)
	contentstr = ChangeConfig(storeInfo, contentstr)
	if err := os.WriteFile(storeInfo.Setting.ConfPath, []byte(contentstr), 0644); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}
	return nil
}
func ChangeConfig(storeInfo StoreInfo, contentstr string) string {
	pattern := regexp.MustCompile(`\{(\w+)\}`)
	for key, value := range storeInfo.Config {
		strValue := convertValueToString(value)
		contentstr = pattern.ReplaceAllStringFunc(contentstr, func(match string) string {
			if strings.Trim(match, "{}") == key {
				return strValue
			}
			return match
		})
	}
	return contentstr
}

func replacePlaceholdersInCmds(storeInfo *StoreInfo) {
	// 创建一个正则表达式模式，用于匹配形如{key}的占位符
	pattern := regexp.MustCompile(`\{(\w+)\}`)

	// 遍历Config中的所有键值对
	for key, value := range storeInfo.Config {
		// 遍历Cmds中的所有命令组
		for cmdGroupName, cmds := range storeInfo.Cmds {
			// 遍历每个命令组中的所有Cmd
			for i, cmd := range cmds {
				strValue := convertValueToString(value)
				// 使用正则表达式替换Content中的占位符
				storeInfo.Cmds[cmdGroupName][i].Content = pattern.ReplaceAllStringFunc(
					cmd.Content,
					func(match string) string {
						// 检查占位符是否与Config中的键匹配
						if strings.Trim(match, "{}") == key {
							return strValue
						}
						return match
					},
				)
			}
		}
	}
}

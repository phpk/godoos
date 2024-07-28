package store

import (
	"encoding/json"
	"fmt"
	"godo/files"
	"godo/libs"
	"godo/progress"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

type StoreInfo struct {
	Name         string           `json:"name"`
	URL          string           `json:"url"`
	NeedDownload bool             `json:"needDownload"`
	Icon         string           `json:"icon"`
	Setting      Setting          `json:"setting"`
	Config       map[string]any   `json:"config"`
	Cmds         map[string][]Cmd `json:"cmds"`
	Install      Install          `json:"install"`
	Start        Install          `json:"start"`
}
type Setting struct {
	BinPath  string `json:"binPath"`
	ConfPath string `json:"confPath"`
	DataDir  string `json:"dataDir"`
	LogDir   string `json:"logDir"`
}

type Item struct {
	Name  string `json:"name"`
	Value any    `json:"value"`
}
type Cmd struct {
	Name     string   `json:"name"`
	FilePath string   `json:"filePath,omitempty"`
	Content  string   `json:"content,omitempty"`
	BinPath  string   `json:"binPath,omitempty"`
	Cmds     []string `json:"cmds,omitempty"`
	Waiting  bool     `json:"waiting"`
	Envs     []Item   `json:"envs"`
}

type Install struct {
	Envs []Item   `json:"envs"`
	Cmds []string `json:"cmds"`
}

func InstallHandler(w http.ResponseWriter, r *http.Request) {
	pluginName := r.URL.Query().Get("name")
	if pluginName == "" {
		libs.ErrorMsg(w, "the app name is empty!")
		return
	}
	exeDir := libs.GetRunDir()
	exePath := filepath.Join(exeDir, pluginName)
	//log.Printf("the app path is %s", exePath)
	if !libs.PathExists(exePath) {
		libs.ErrorMsg(w, "the app path is not exists!")
		return
	}
	installFile := filepath.Join(exePath, "install.json")
	if !libs.PathExists(installFile) {
		libs.ErrorMsg(w, "the install.json is not exists!")
		return
	}
	var storeInfo StoreInfo
	content, err := os.ReadFile(installFile)
	if err != nil {
		libs.ErrorMsg(w, "cant read the install.json!")
		return
	}
	//设置 info.json
	exePath = strings.ReplaceAll(exePath, "\\", "/")
	contentBytes := []byte(strings.ReplaceAll(string(content), "{exePath}", exePath))
	//log.Printf("====content: %s", string(contentBytes))
	err = json.Unmarshal(contentBytes, &storeInfo)
	if err != nil {
		log.Printf("Error during unmarshal: %v", err)
		libs.ErrorMsg(w, "the install.json is error: "+err.Error())
		return
	}
	replacePlaceholdersInCmds(&storeInfo)
	infoFile := filepath.Join(exePath, "info.json")
	err = SaveStoreInfo(storeInfo, infoFile)
	if err != nil {
		libs.ErrorMsg(w, "save the info.json is error!")
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
	//复制static目录
	staticPath := filepath.Join(exePath, "static")
	if libs.PathExists(staticPath) {
		staticDir := libs.GetStaticDir()
		targetPath := filepath.Join(staticDir, pluginName)
		err = os.Rename(staticPath, targetPath)
		if err != nil {
			libs.ErrorMsg(w, "copy the static is error!")
			return
		}
	}
	libs.SuccessMsg(w, "success", "install the app success!")
}
func UnInstallHandler(w http.ResponseWriter, r *http.Request) {
	pluginName := r.URL.Query().Get("name")
	if pluginName == "" {
		libs.ErrorMsg(w, "the app name is empty!")
		return
	}
	err := progress.StopCmd(pluginName)
	if err != nil {
		libs.ErrorMsg(w, "stop the app is error!")
		return
	}
	exeDir := libs.GetRunDir()
	exePath := filepath.Join(exeDir, pluginName)
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
func RunCmds(storeInfo StoreInfo, cmdKey string) error {
	cmds := storeInfo.Cmds[cmdKey]
	for _, cmd := range cmds {

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
	}
	return nil
}
func runStop(storeInfo StoreInfo) error {
	return progress.StopCmd(storeInfo.Name)
}
func runStart(storeInfo StoreInfo) error {
	err := SetEnvs(storeInfo.Start.Envs)
	if err != nil {
		return fmt.Errorf("failed to set start environment variable %s: %w", storeInfo.Name, err)
	}
	if len(storeInfo.Start.Cmds) > 0 {
		if !libs.PathExists(storeInfo.Setting.BinPath) {
			return fmt.Errorf("script file %s does not exist", storeInfo.Setting.BinPath)
		}
		cmd := exec.Command(storeInfo.Setting.BinPath, storeInfo.Start.Cmds...)
		if runtime.GOOS == "windows" {
			// 在Windows上，通过设置CreationFlags来隐藏窗口
			cmd = progress.SetHideConsoleCursor(cmd)
		}
		// 启动脚本命令并返回可能的错误
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("failed to start process %s: %w", storeInfo.Name, err)
		}

		progress.RegisterProcess(storeInfo.Name, cmd)
	}
	return nil
}
func runRestart(storeInfo StoreInfo) error {
	err := runStop(storeInfo)
	if err != nil {
		return fmt.Errorf("failed to stop process %s: %w", storeInfo.Name, err)
	}
	return runStart(storeInfo)
}
func runExec(storeInfo StoreInfo, cmdParam Cmd) error {
	err := SetEnvs(cmdParam.Envs)
	if err != nil {
		return fmt.Errorf("failed to set start environment variable %s: %w", storeInfo.Name, err)
	}
	log.Printf("bin path:%v", cmdParam.BinPath)
	log.Printf("cmds:%v", cmdParam.Cmds)
	cmd := exec.Command(cmdParam.BinPath, cmdParam.Cmds...)
	if runtime.GOOS == "windows" {
		// 在Windows上，通过设置CreationFlags来隐藏窗口
		cmd = progress.SetHideConsoleCursor(cmd)
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to run exec process %s: %w", storeInfo.Name, err)
	}
	if cmdParam.Waiting {
		if err = cmd.Wait(); err != nil {
			return fmt.Errorf("failed to wait for exec process %s: %w", storeInfo.Name, err)
		}
	}
	log.Printf("run exec process %s, name is %s", storeInfo.Name, cmdParam.Name)
	return nil
}
func WriteFile(cmd Cmd) error {
	if cmd.FilePath != "" {
		content := cmd.Content
		if content != "" {
			err := os.WriteFile(cmd.FilePath, []byte(content), 0644)
			if err != nil {
				return fmt.Errorf("failed to write to file: %w", err)
			}
		}
	}
	return nil
}
func DeleteFile(cmd Cmd) error {
	if cmd.FilePath != "" {
		err := os.Remove(cmd.FilePath)
		if err != nil {
			return fmt.Errorf("failed to delete file: %w", err)
		}
	}
	return nil
}
func MkDir(cmd Cmd) error {
	err := os.MkdirAll(cmd.FilePath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to make dir: %w", err)
	}
	return nil
}
func Unzip(cmd Cmd) error {
	return files.Decompress(cmd.FilePath, cmd.Content)
}
func Zip(cmd Cmd) error {
	return files.Encompress(cmd.FilePath, cmd.Content)
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

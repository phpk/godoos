package store

import (
	"encoding/json"
	"fmt"
	"godo/libs"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func GetStoreListHandler(w http.ResponseWriter, r *http.Request) {
	cate := r.URL.Query().Get("cate")
	list, err := GetInstallList(cate)
	if err != nil {
		libs.ErrorMsg(w, err.Error())
		return
	}
	libs.SuccessMsg(w, list, "")
}
func InstallByName(name string, cate string) (InstallInfo, error) {
	var installInfo InstallInfo
	list, err := GetInstallList(cate)

	if err != nil {
		return installInfo, fmt.Errorf("failed to get plugin list: %v", err)
	}
	for _, item := range list {
		if item.Name == name {
			installInfo = item
			break
		}
	}
	if installInfo.Name == "" {
		return installInfo, fmt.Errorf("plugin not found")
	}
	return Installation(installInfo.Name)
}
func GetInstallList(cate string) ([]InstallInfo, error) {
	if cate == "" {
		cate = "hots"
	}
	os := runtime.GOOS
	arch := runtime.GOARCH
	var list []InstallInfo
	pluginUrl := "https://gitee.com/ruitao_admin/godoos-image/raw/master/store/" + os + "/" + arch + "/" + cate + ".json"
	res, err := http.Get(pluginUrl)
	if err != nil {
		return list, fmt.Errorf("failed to get plugin list: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return list, fmt.Errorf("failed to get plugin list")
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return list, fmt.Errorf("failed to get plugin list: %v", err)
	}
	err = json.Unmarshal(body, &list)
	if err != nil {
		return list, fmt.Errorf("failed to get plugin list: %v", err)
	}
	return list, nil

}
func GetInstallInfoHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		libs.ErrorMsg(w, "name is required")
		return
	}
	info, err := GetInstallInfo(name)
	if err != nil {
		libs.ErrorMsg(w, err.Error())
		return
	}
	libs.SuccessMsg(w, info, "")
}
func StoreSettingHandler(w http.ResponseWriter, r *http.Request) {
	var req map[string]any
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		libs.ErrorMsg(w, err.Error())
		return
	}
	if req["name"] == nil {
		libs.ErrorMsg(w, "name is required")
		return
	}
	if req["cmdKey"] == nil {
		libs.ErrorMsg(w, "cmdKey is required")
		return
	}
	name := req["name"].(string)
	cmdKey := req["cmdKey"].(string)
	// 替换字符串类型的值中的反斜杠
	for key, value := range req {
		if strValue, ok := value.(string); ok {
			req[key] = strings.ReplaceAll(strValue, "\\", "/")
		}
	}
	storeInfo, err := GetStoreInfo(name)
	if err != nil {
		libs.ErrorMsg(w, err.Error())
		return
	}
	// 更新storeInfo.Config
	for k, v := range req {
		if k != "name" && k != "cmdKey" {
			storeInfo.Config[k] = v // 如果k存在，则更新；如果不存在，则新增
		}
	}
	replacePlaceholdersInCmds(&storeInfo)
	err = SaveInfoFile(storeInfo)
	if err != nil {
		libs.ErrorMsg(w, "the store info.json is error: "+err.Error())
		return
	}
	_, ok := storeInfo.Commands[cmdKey]
	if !ok {
		libs.ErrorMsg(w, "cmdKey is not found")
		return
	}
	err = RunCmds(storeInfo, cmdKey)
	if err != nil {
		libs.ErrorMsg(w, err.Error())
		return
	}
	libs.SuccessMsg(w, "success", "success")

}
func GetStoreInfo(name string) (StoreInfo, error) {
	var storeInfo StoreInfo
	exeDir := libs.GetRunDir()
	infoPath := filepath.Join(exeDir, name, "info.json")
	if !libs.PathExists(infoPath) {
		return storeInfo, fmt.Errorf("process information for '%s' not found", name)
	}
	return GetStoreInfoByPath(infoPath)
}
func GetStoreInfoByPath(infoPath string) (StoreInfo, error) {
	var storeInfo StoreInfo
	content, err := os.ReadFile(infoPath)
	if err != nil {
		return storeInfo, fmt.Errorf("failed to read info.json: %v", err)
	}
	if err := json.Unmarshal(content, &storeInfo); err != nil {
		return storeInfo, fmt.Errorf("failed to unmarshal info.json: %v", err)
	}
	scriptPath := storeInfo.Setting.BinPath
	if !libs.PathExists(scriptPath) {
		return storeInfo, fmt.Errorf("script file '%s' not found", scriptPath)
	}
	return storeInfo, nil
}
func GetInstalled() []InstallInfo {
	runDir := libs.GetRunDir()
	var list []InstallInfo

	filepath.Walk(runDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			infoPath := filepath.Join(path, "install.json")
			if libs.PathExists(infoPath) {
				installInfo, err := GetInstallInfoByPath(infoPath)
				if err == nil {
					list = append(list, installInfo)
				}
			}
		}
		return nil // 返回nil表示继续遍历
	})
	return list
}

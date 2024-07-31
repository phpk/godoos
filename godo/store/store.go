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
	os := runtime.GOOS
	arch := runtime.GOARCH
	if cate == "" {
		libs.ErrorMsg(w, "cate is required")
		return
	}
	pluginUrl := "https://gitee.com/ruitao_admin/godoos-image/raw/master/store/" + os + "/" + arch + "/" + cate + ".json"
	res, err := http.Get(pluginUrl)
	if err != nil {
		libs.ErrorMsg(w, err.Error())
	}
	defer res.Body.Close()
	if res.StatusCode == 200 {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			libs.ErrorMsg(w, err.Error())
			return
		}
		var info interface{}
		err = json.Unmarshal(body, &info)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			return
		}
		json.NewEncoder(w).Encode(info)

	}
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
	err = SaveInfoFile(storeInfo)
	if err != nil {
		libs.ErrorMsg(w, "the store info.json is error: "+err.Error())
		return
	}
	_, ok := storeInfo.Cmds[cmdKey]
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

package cmd

import (
	"encoding/json"
	"godo/libs"
	"net/http"
	"os"
)

func GetOsPath() string {
	osInfo, _ := libs.GetConfig("osInfo")
	return osInfo.Value
}
func HandleSetConfig(w http.ResponseWriter, r *http.Request) {
	var req libs.ReqBody
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		libs.ErrorMsg(w, "The params is error!")
		return
	}
	if req.Name == "osInfo" && req.Value != "" {
		osInfo, _ := libs.GetConfig("osInfo")
		if osInfo.Value != req.Value {
			if !libs.PathExists(req.Value) {
				libs.ErrorMsg(w, "The Path is not exists!")
				return
			}
			err = os.Chmod(req.Value, 0755)
			if err != nil {
				libs.ErrorMsg(w, "The Path chmod is error!")
				return
			}
			osInfo.Value = req.Value
			osInfo.Type = req.Type
			libs.SetConfig(osInfo)
		}

	}
	if req.Name == "userInfo" ||
		req.Name == "dbInfo" {
		libs.SetConfig(req)
	}
	err = libs.LoadConfig()
	if err != nil {
		libs.ErrorMsg(w, "The config load error!")
		return
	}
	libs.SuccessMsg(w, "success", "The config set success!")
}

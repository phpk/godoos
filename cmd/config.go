package cmd

import (
	"encoding/json"
	"godoos/libs"
	"net/http"
	"os"
)

func HandleSetConfig(w http.ResponseWriter, r *http.Request) {
	var req libs.ReqBody
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		ErrorMsg(w, "The params is error!")
		return
	}
	if req.Name == "osInfo" {
		osInfo, _ := libs.GetConfig("osInfo")
		if req.Value != "" {
			if !libs.PathExists(req.Value) {
				ErrorMsg(w, "The Path is not exists!")
				return
			}
			err = os.Chmod(req.Value, 0755)
			if err != nil {
				ErrorMsg(w, "The Path chmod is error!")
				return
			}
			osInfo.Value = req.Value
		}

		osInfo.Type = req.Type
		libs.SetConfig(osInfo)
	}
	if req.Name == "userInfo" ||
		req.Name == "dbInfo" {
		libs.SetConfig(req)
	}
	err = libs.LoadConfig()
	if err != nil {
		ErrorMsg(w, "The config load error!")
		return
	}
	SuccessMsg(w, "success", "The config set success!")
}

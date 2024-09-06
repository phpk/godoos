package localchat

import (
	"encoding/json"
	"godo/libs"
	"net/http"
	"os"
)

func HandleCheck(w http.ResponseWriter, r *http.Request) {
	hostname, err := os.Hostname()
	if err != nil {
		libs.ErrorMsg(w, "HandleMessage error")
		return
	}
	libs.SuccessMsg(w, hostname, "")
}
func HandleAddr(w http.ResponseWriter, r *http.Request) {
	var ipStr libs.UserChatIpSetting
	err := json.NewDecoder(r.Body).Decode(&ipStr)
	if err != nil {
		libs.HTTPError(w, http.StatusBadRequest, "Failed to parse request body")
		return
	}
	libs.SetConfigByName("ChatIpSetting", ipStr)
	go CheckOnline()
	libs.SuccessMsg(w, nil, "success")
}

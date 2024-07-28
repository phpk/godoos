package store

import (
	"encoding/json"
	"fmt"
	"godo/libs"
	"io"
	"net/http"
	"runtime"
)

func GetStoreInfoHandler(w http.ResponseWriter, r *http.Request) {
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

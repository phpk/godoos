// MIT License
//
// Copyright (c) 2024 godoos.com
// Email: xpbb@qq.com
// GitHub: github.com/phpk/godoos
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

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

func HandleGetAddr(w http.ResponseWriter, r *http.Request) {
	ipInfo := libs.GetChatIpSetting()
	libs.SuccessMsg(w, ipInfo, "success")
}

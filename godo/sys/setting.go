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
package sys

import (
	"encoding/json"
	"fmt"
	"godo/libs"
	"godo/webdav"
	"net"
	"net/http"
	"os"
)

func ConfigHandler(w http.ResponseWriter, r *http.Request) {
	var reqs []libs.ReqBody
	err := json.NewDecoder(r.Body).Decode(&reqs)
	if err != nil {
		libs.ErrorMsg(w, "The params is error!")
		return
	}
	for _, req := range reqs {
		if req.Name == "osPath" {
			reqPath := req.Value.(string)
			osPath, ok := libs.GetConfig("osPath")
			if !ok || osPath != reqPath {
				if !libs.PathExists(reqPath) {
					libs.ErrorMsg(w, "The Path is not exists!")
					return
				}
				err = os.Chmod(reqPath, 0755)
				if err != nil {
					libs.ErrorMsg(w, "The Path chmod is error!")
					return
				}
				libs.SetConfig(req)
			}
		}
		if req.Name == "webdavClient" {
			libs.SetConfig(req)
			err := webdav.InitWebdav()
			if err != nil {
				libs.ErrorMsg(w, "The webdav client init is error："+err.Error())
				return
			}
		}
	}

	libs.SuccessMsg(w, "success", "The config set success!")
}
func SetIplist(req libs.ReqBody) error {
	reqIplist, ok := req.Value.([]string)
	if !ok {
		return fmt.Errorf("unexpected type for iplist:%v", reqIplist)
	}
	if len(reqIplist) > 0 {
		for _, ip := range reqIplist {
			_, _, err := net.ParseCIDR(ip)
			if err != nil {
				return fmt.Errorf("the iplist is error:%s", err.Error())
			}
		}
		libs.SetConfig(req)
	}
	return nil
}

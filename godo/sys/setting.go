/*
 * godoos - A lightweight cloud desktop
 * Copyright (C) 2024 godoos.com
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 2.1 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

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

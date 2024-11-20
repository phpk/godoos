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

package localchat

import (
	"encoding/json"
	"godo/libs"
	"net/http"
	"os"
	"time"
)

func HandlerApplySendFile(w http.ResponseWriter, r *http.Request) {
	var msg UdpMessage
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&msg); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	hostname, err := os.Hostname()
	if err != nil {
		libs.ErrorMsg(w, "HandleMessage error")
		return
	}
	msg.Hostname = hostname
	msg.Time = time.Now()
	msg.Type = "fileSending"
	SendToIP(msg)
	libs.SuccessMsg(w, nil, "请求文件发送成功")
}
func HandlerCannelFile(w http.ResponseWriter, r *http.Request) {
	var msg UdpMessage
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&msg); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	hostname, err := os.Hostname()
	if err != nil {
		libs.ErrorMsg(w, "HandleMessage error")
		return
	}
	msg.Hostname = hostname
	msg.Time = time.Now()
	msg.Type = "fileCannel"
	SendToIP(msg)
	libs.SuccessMsg(w, nil, "请求文件发送成功")
}
func HandlerAccessFile(w http.ResponseWriter, r *http.Request) {
	var msg UdpMessage
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&msg); err != nil {
		libs.ErrorMsg(w, "Invalid request body")
		return
	}
	defer r.Body.Close()
	hostname, err := os.Hostname()
	if err != nil {
		libs.ErrorMsg(w, "HandleMessage error")
		return
	}
	msg.Hostname = hostname
	msg.Time = time.Now()
	msg.Type = "fileAccessed"
	SendToIP(msg)
	path, err := downloadFiles(msg)
	if err != nil {
		libs.ErrorMsg(w, "HandleMessage error")
		return
	}
	res := map[string]interface{}{
		"path": path,
		"msg":  msg.Message,
	}
	libs.SuccessMsg(w, res, "接收文件成功")
}

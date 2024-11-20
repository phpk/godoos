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

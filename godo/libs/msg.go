/*
 * GodoOS - A lightweight cloud desktop
 * Copyright (C) 2024 https://godoos.com
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

package libs

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

func WriteJSONResponse(w http.ResponseWriter, res APIResponse, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(res)
}

// HTTPError 返回带有JSON错误消息的HTTP错误
func HTTPError(w http.ResponseWriter, status int, message string) {
	WriteJSONResponse(w, APIResponse{Message: message, Code: -1}, status)
}
func ErrorMsg(w http.ResponseWriter, message string) {
	WriteJSONResponse(w, APIResponse{Message: message, Code: -1}, 200)
}
func ErrorData(w http.ResponseWriter, data any, message string) {
	WriteJSONResponse(w, APIResponse{Message: message, Data: data, Code: -1}, 200)
}
func Error(w http.ResponseWriter, message string, err string) {
	WriteJSONResponse(w, APIResponse{Message: message, Error: err, Code: -1}, 200)
}
func SuccessMsg(w http.ResponseWriter, data any, message string) {
	WriteJSONResponse(w, APIResponse{Message: message, Data: data, Code: 0}, 200)
}

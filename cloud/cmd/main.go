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

package cmd

import (
	"context"
	"godocloud/deps"
	"io/fs"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const serverAddress = ":56781"

var srv *http.Server

func OsStart() {
	//InitServer()
	router := mux.NewRouter()
	router.Use(recoverMiddleware)
	router.Use(corsMiddleware())
	// 使用带有日志装饰的处理器注册路由
	router.Use(loggingMiddleware{}.Middleware)

	// 注册根路径的处理函数
	distFS, _ := fs.Sub(deps.Frontendassets, "dist")
	fileServer := http.FileServer(http.FS(distFS))

	// 注册根路径的处理函数
	router.PathPrefix("/").Handler(fileServer)

	//serverAddress := ":56781"
	log.Printf("Listening on port: %v", serverAddress)
	srv = &http.Server{Addr: serverAddress, Handler: router}
	Serve(srv)
}
func OsStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server stopped.")
}
func OsRestart() {
	// 停止当前服务
	OsStop()
	// 重新启动服务
	OsStart()
}

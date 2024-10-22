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
package cmd

import (
	"context"
	"godo/deps"
	"godo/files"
	"godo/libs"
	"godo/localchat"
	"godo/store"
	"godo/sys"
	"godo/webdav"
	"io/fs"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const serverAddress = ":56780"

var srv *http.Server

func OsStart() {
	libs.InitServer()
	err := files.InitOsSystem()
	if err != nil {
		log.Fatalf("InitOsSystem error: %v", err)
		return
	}
	webdav.InitWebdav()
	router := mux.NewRouter()
	router.Use(corsMiddleware())
	// 使用带有日志装饰的处理器注册路由
	router.Use(loggingMiddleware{}.Middleware)
	staticDir := libs.GetStaticDir()
	router.PathPrefix("/static").Handler(http.StripPrefix("/static", http.FileServer(http.Dir(staticDir))))
	router.HandleFunc("/ping", store.Ping).Methods(http.MethodGet)

	progressRouter := router.PathPrefix("/store").Subrouter()
	progressRouter.HandleFunc("/start/{name}", store.StartProcess).Methods(http.MethodGet)
	progressRouter.HandleFunc("/stop/{name}", store.StopProcess).Methods(http.MethodGet)
	progressRouter.HandleFunc("/startall", store.StartAll).Methods(http.MethodGet)
	progressRouter.HandleFunc("/stopall", store.StopAll).Methods(http.MethodGet)
	progressRouter.HandleFunc("/restart/{name}", store.ReStartProcess).Methods(http.MethodGet)
	progressRouter.HandleFunc("/listporgress", store.Status).Methods(http.MethodGet)
	progressRouter.HandleFunc("/listport", store.ListAllProcessesHandler).Methods(http.MethodGet)
	progressRouter.HandleFunc("/killport", store.KillProcessByNameHandler).Methods(http.MethodGet)
	progressRouter.HandleFunc("/storelist", store.GetStoreListHandler).Methods(http.MethodGet)
	progressRouter.HandleFunc("/download", store.DownloadHandler).Methods(http.MethodPost)
	progressRouter.HandleFunc("/install", store.InstallHandler).Methods(http.MethodGet)
	progressRouter.HandleFunc("/installInfo", store.GetInstallInfoHandler).Methods(http.MethodGet)
	progressRouter.HandleFunc("/installOut", store.RunOutHandler).Methods(http.MethodGet)
	progressRouter.HandleFunc("/uninstall", store.UnInstallHandler).Methods(http.MethodGet)
	progressRouter.HandleFunc("/setting", store.StoreSettingHandler).Methods(http.MethodPost)
	progressRouter.HandleFunc("/upload", store.UploadHandler).Methods(http.MethodPost)

	router.HandleFunc("/system/message", sys.HandleSystemEvents).Methods(http.MethodGet)
	router.HandleFunc("/system/update", sys.UpdateAppHandler).Methods(http.MethodGet)
	router.HandleFunc("/system/setting", sys.ConfigHandler).Methods(http.MethodPost)
	fileRouter := router.PathPrefix("/file").Subrouter()
	fileRouter.HandleFunc("/desktop", files.HandleDesktop).Methods(http.MethodGet)
	fileRouter.HandleFunc("/info", files.HandleSystemInfo).Methods(http.MethodGet)
	fileRouter.HandleFunc("/read", files.HandleReadDir).Methods(http.MethodGet)
	fileRouter.HandleFunc("/stat", files.HandleStat).Methods(http.MethodGet)
	fileRouter.HandleFunc("/chmod", files.HandleChmod).Methods(http.MethodPost)
	fileRouter.HandleFunc("/exists", files.HandleExists).Methods(http.MethodGet)
	fileRouter.HandleFunc("/readfile", files.HandleReadFile).Methods(http.MethodGet)
	fileRouter.HandleFunc("/unlink", files.HandleUnlink).Methods(http.MethodGet)
	fileRouter.HandleFunc("/clear", files.HandleClear).Methods(http.MethodGet)
	fileRouter.HandleFunc("/rename", files.HandleRename).Methods(http.MethodGet)
	fileRouter.HandleFunc("/mkdir", files.HandleMkdir).Methods(http.MethodPost)
	fileRouter.HandleFunc("/rmdir", files.HandleRmdir).Methods(http.MethodGet)
	fileRouter.HandleFunc("/copyfile", files.HandleCopyFile).Methods(http.MethodGet)
	fileRouter.HandleFunc("/writefile", files.HandleWriteFile).Methods(http.MethodPost)
	fileRouter.HandleFunc("/appendfile", files.HandleAppendFile).Methods(http.MethodPost)
	fileRouter.HandleFunc("/zip", files.HandleZip).Methods(http.MethodGet)
	fileRouter.HandleFunc("/unzip", files.HandleUnZip).Methods(http.MethodGet)
	fileRouter.HandleFunc("/watch", files.WatchHandler).Methods(http.MethodGet)

	localchatRouter := router.PathPrefix("/localchat").Subrouter()
	localchatRouter.HandleFunc("/message", localchat.HandleMessage).Methods(http.MethodPost)
	localchatRouter.HandleFunc("/applyfile", localchat.HandlerApplySendFile).Methods(http.MethodPost)
	localchatRouter.HandleFunc("/cannelfile", localchat.HandlerCannelFile).Methods(http.MethodPost)
	localchatRouter.HandleFunc("/accessfile", localchat.HandlerAccessFile).Methods(http.MethodPost)
	localchatRouter.HandleFunc("/getfiles", localchat.HandleGetFiles).Methods(http.MethodPost)
	localchatRouter.HandleFunc("/servefile", localchat.HandleServeFile).Methods(http.MethodGet)
	localchatRouter.HandleFunc("/sendimage", localchat.HandlerSendImg).Methods(http.MethodPost)
	localchatRouter.HandleFunc("/viewimage", localchat.HandleViewImg).Methods(http.MethodGet)
	localchatRouter.HandleFunc("/setting", localchat.HandleAddr).Methods(http.MethodPost)
	localchatRouter.HandleFunc("/getsetting", localchat.HandleGetAddr).Methods(http.MethodGet)
	// 注册 WebDAV 路由
	webdavRouter := router.PathPrefix("/webdav").Subrouter()
	webdavRouter.HandleFunc("/read", webdav.HandleReadDir).Methods(http.MethodGet)
	webdavRouter.HandleFunc("/stat", webdav.HandleStat).Methods(http.MethodGet)
	webdavRouter.HandleFunc("/chmod", webdav.HandleChmod).Methods(http.MethodPost)
	webdavRouter.HandleFunc("/exists", webdav.HandleExists).Methods(http.MethodGet)
	webdavRouter.HandleFunc("/readfile", webdav.HandleReadFile).Methods(http.MethodGet)
	webdavRouter.HandleFunc("/unlink", webdav.HandleUnlink).Methods(http.MethodGet)
	webdavRouter.HandleFunc("/clear", webdav.HandleClear).Methods(http.MethodGet)
	webdavRouter.HandleFunc("/rename", webdav.HandleRename).Methods(http.MethodGet)
	webdavRouter.HandleFunc("/mkdir", webdav.HandleMkdir).Methods(http.MethodPost)
	webdavRouter.HandleFunc("/rmdir", webdav.HandleRmdir).Methods(http.MethodGet)
	webdavRouter.HandleFunc("/copyfile", webdav.HandleCopyFile).Methods(http.MethodGet)
	webdavRouter.HandleFunc("/writefile", webdav.HandleWriteFile).Methods(http.MethodPost)
	webdavRouter.HandleFunc("/appendfile", webdav.HandleAppendFile).Methods(http.MethodPost)

	distFS, _ := fs.Sub(deps.Frontendassets, "dist")
	fileServer := http.FileServer(http.FS(distFS))
	router.PathPrefix("/").Handler(fileServer)

	go store.CheckActive(context.Background())
	log.Printf("Listening on port: %v", serverAddress)
	srv = &http.Server{Addr: serverAddress, Handler: router}
	Serve(srv)
}
func OsStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err := store.StopAllHandler()
	if err != nil {
		log.Fatalf("Servers forced to shutdown error: %v", err)
	}
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server stopped.")
}

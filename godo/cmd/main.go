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
	model "godo/ai/server"
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
	err = deps.InitDir()
	if err != nil {
		log.Fatalf("Init Dir error: %v", err)
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
	fileRouter.HandleFunc("/setfilepwd", files.HandleSetFilePwd).Methods(http.MethodGet)

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
	// 注册ai路由
	aiRouter := router.PathPrefix("/ai").Subrouter()
	aiRouter.HandleFunc("/download", model.Download).Methods(http.MethodPost)
	aiRouter.HandleFunc("/server", model.DownServerHandler).Methods(http.MethodGet)
	aiRouter.HandleFunc("/delete", model.DeleteFileHandle).Methods(http.MethodPost)
	aiRouter.HandleFunc("/tags", model.Tagshandler).Methods(http.MethodGet)
	aiRouter.HandleFunc("/show", model.ShowHandler).Methods(http.MethodGet)
	aiRouter.HandleFunc("/refreshOllama", model.RefreshOllamaHandler).Methods(http.MethodGet)
	aiRouter.HandleFunc("/chat", model.ChatHandler).Methods(http.MethodPost)
	aiRouter.HandleFunc("/embeddings", model.EmbeddingHandler).Methods(http.MethodPost)
	// router.Handle("/model/uploadimage", http.MethodPost, sd.UploadHandler)
	// router.Handle("/model/image", http.MethodPost, sd.CreateImage)
	// router.Handle("/model/deleteimage", http.MethodPost, sd.DeleteImageHandler)
	// router.Handle("/model/viewimage", http.MethodGet, sd.ServeImage)
	// router.Handle("/model/voice", http.MethodPost, voice.UploadHandler)
	// router.Handle("/model/tts", http.MethodPost, voice.TtsHandler)
	// router.Handle("/model/audio", http.MethodGet, voice.ServeAudio)
	ieRouter := router.PathPrefix("/ie").Subrouter()
	ieRouter.HandleFunc("/navigate", store.HandleNavigate).Methods(http.MethodGet)
	ieRouter.HandleFunc("/back", store.HandleBack).Methods(http.MethodGet)
	ieRouter.HandleFunc("/forward", store.HandleForward).Methods(http.MethodGet)
	ieRouter.HandleFunc("/refresh", store.HandleRefresh).Methods(http.MethodGet)

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

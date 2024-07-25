package cmd

import (
	"context"
	"godo/files"
	"godo/libs"
	"godo/localchat"
	"godo/progress"
	"godo/sys"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const serverAddress = ":56780"

var srv *http.Server

func OsStart() {
	libs.Initdir()
	router := mux.NewRouter()
	router.Use(corsMiddleware())
	// 使用带有日志装饰的处理器注册路由
	router.Use(loggingMiddleware{}.Middleware)
	staticDir := libs.GetStaticDir()
	router.PathPrefix("/static").Handler(http.StripPrefix("/static", http.FileServer(http.Dir(staticDir))))
	router.HandleFunc("/ping", progress.Ping).Methods(http.MethodGet)
	if libs.PathExists("./dist") {
		router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./dist"))))
	} else {
		router.HandleFunc("/", progress.Ping).Methods(http.MethodGet)
	}
	progressRouter := router.PathPrefix("/progress").Subrouter()
	progressRouter.HandleFunc("/start/{name}", progress.StartProcess).Methods(http.MethodGet)
	progressRouter.HandleFunc("/stop/{name}", progress.StopProcess).Methods(http.MethodGet)
	progressRouter.HandleFunc("/startall", progress.StartAll).Methods(http.MethodGet)
	progressRouter.HandleFunc("/stopall", progress.StopAll).Methods(http.MethodGet)
	progressRouter.HandleFunc("/restart/{name}", progress.ReStartProcess).Methods(http.MethodGet)
	progressRouter.HandleFunc("/list", progress.Status).Methods(http.MethodGet)
	progressRouter.HandleFunc("/listport", progress.ListPortsHandler).Methods(http.MethodGet)
	progressRouter.HandleFunc("/killport", progress.KillPortHandler).Methods(http.MethodGet)
	progressRouter.HandleFunc("/app/{name}", progress.ForwardRequest).Methods(http.MethodGet, http.MethodPost)
	progressRouter.HandleFunc("/app/{name}/{subpath:.*}", progress.ForwardRequest).Methods(http.MethodGet, http.MethodPost)

	router.HandleFunc("/system/updateInfo", sys.GetUpdateUrlHandler).Methods(http.MethodGet)
	router.HandleFunc("/system/update", sys.UpdateAppHandler).Methods(http.MethodGet)
	router.HandleFunc("/system/storeList", sys.GetStoreInfoHandler).Methods(http.MethodGet)
	router.HandleFunc("/system/setting", sys.HandleSetConfig).Methods(http.MethodPost)
	router.HandleFunc("/files/info", files.HandleSystemInfo).Methods(http.MethodGet)
	router.HandleFunc("/file/read", files.HandleReadDir).Methods(http.MethodGet)
	router.HandleFunc("/file/stat", files.HandleStat).Methods(http.MethodGet)
	router.HandleFunc("/file/chmod", files.HandleChmod).Methods(http.MethodPost)
	router.HandleFunc("/file/exists", files.HandleExists).Methods(http.MethodGet)
	router.HandleFunc("/file/readfile", files.HandleReadFile).Methods(http.MethodGet)
	router.HandleFunc("/file/unlink", files.HandleUnlink).Methods(http.MethodGet)
	router.HandleFunc("/file/clear", files.HandleClear).Methods(http.MethodGet)
	router.HandleFunc("/file/rename", files.HandleRename).Methods(http.MethodGet)
	router.HandleFunc("/file/mkdir", files.HandleMkdir).Methods(http.MethodPost)
	router.HandleFunc("/file/rmdir", files.HandleRmdir).Methods(http.MethodGet)
	router.HandleFunc("/file/copyfile", files.HandleCopyFile).Methods(http.MethodGet)
	router.HandleFunc("/file/writefile", files.HandleWriteFile).Methods(http.MethodPost)
	router.HandleFunc("/file/appendfile", files.HandleAppendFile).Methods(http.MethodPost)
	router.HandleFunc("/file/watch", files.WatchHandler).Methods(http.MethodGet)
	router.HandleFunc("/localchat/sse", localchat.SseHandler).Methods(http.MethodGet)
	router.HandleFunc("/localchat/message", localchat.HandleMessage).Methods(http.MethodPost)
	router.HandleFunc("/localchat/upload", localchat.MultiUploadHandler).Methods(http.MethodPost)

	go progress.CheckActive(context.Background())
	log.Printf("Listening on port: %v", serverAddress)
	srv = &http.Server{Addr: serverAddress, Handler: router}
	Serve(srv)
}
func OsStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err := progress.StopAllHandler()
	if err != nil {
		log.Fatalf("Servers forced to shutdown error: %v", err)
	}
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server stopped.")
}

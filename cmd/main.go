package cmd

import (
	"context"
	"godoos/libs"
	"godoos/progress"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const serverAddress = ":56710"

var srv *http.Server

func OsStart() {
	libs.Initdir()
	router := mux.NewRouter()
	router.Use(corsMiddleware())
	// 使用带有日志装饰的处理器注册路由
	router.Use(loggingMiddleware{}.Middleware)
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
	router.HandleFunc("/ping", progress.Ping).Methods(http.MethodGet)
	router.HandleFunc("/", progress.Ping).Methods(http.MethodGet)
	router.HandleFunc("/system/info", HandleSystemInfo).Methods(http.MethodGet)
	router.HandleFunc("/system/setting", HandleSetConfig).Methods(http.MethodPost)
	router.HandleFunc("/file/read", HandleReadDir).Methods(http.MethodGet)
	router.HandleFunc("/file/stat", HandleStat).Methods(http.MethodGet)
	router.HandleFunc("/file/chmod", HandleChmod).Methods(http.MethodPost)
	router.HandleFunc("/file/exists", HandleExists).Methods(http.MethodGet)
	router.HandleFunc("/file/readfile", HandleReadFile).Methods(http.MethodGet)
	router.HandleFunc("/file/unlink", HandleUnlink).Methods(http.MethodGet)
	router.HandleFunc("/file/clear", HandleClear).Methods(http.MethodGet)
	router.HandleFunc("/file/rename", HandleRename).Methods(http.MethodGet)
	router.HandleFunc("/file/mkdir", HandleMkdir).Methods(http.MethodPost)
	router.HandleFunc("/file/rmdir", HandleRmdir).Methods(http.MethodGet)
	router.HandleFunc("/file/copyfile", HandleCopyFile).Methods(http.MethodGet)
	router.HandleFunc("/file/writefile", HandleWriteFile).Methods(http.MethodPost)
	router.HandleFunc("/file/appendfile", HandleAppendFile).Methods(http.MethodPost)
	router.HandleFunc("/file/watch", WatchHandler).Methods(http.MethodGet)

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

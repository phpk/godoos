package cmd

import (
	"context"
	"godo/files"
	"godo/libs"
	"godo/localchat"
	"godo/store"
	"godo/sys"
	"godo/webdav"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const serverAddress = ":56780"
const staticServerAddress = ":8215"

var srv *http.Server
var staticSrv *http.Server
var staticRouter *http.ServeMux

func OsStart() {
	libs.InitServer()
	webdav.InitWebdav()
	router := mux.NewRouter()
	router.Use(corsMiddleware())
	// 使用带有日志装饰的处理器注册路由
	router.Use(loggingMiddleware{}.Middleware)
	staticDir := libs.GetStaticDir()
	router.PathPrefix("/static").Handler(http.StripPrefix("/static", http.FileServer(http.Dir(staticDir))))
	router.HandleFunc("/ping", store.Ping).Methods(http.MethodGet)
	if !libs.PathExists("./dist") {
		router.HandleFunc("/", store.Ping).Methods(http.MethodGet)
	} else {
		staticRouter = http.NewServeMux()
		staticRouter.Handle("/", http.FileServer(http.Dir("./dist")))
	}

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

	router.HandleFunc("/system/updateInfo", sys.GetUpdateUrlHandler).Methods(http.MethodGet)
	router.HandleFunc("/system/update", sys.UpdateAppHandler).Methods(http.MethodGet)
	router.HandleFunc("/system/setting", sys.ConfigHandler).Methods(http.MethodPost)
	fileRouter := router.PathPrefix("/file").Subrouter()
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
	localchatRouter.HandleFunc("/sse", localchat.SseHandler).Methods(http.MethodGet)
	localchatRouter.HandleFunc("/message", localchat.HandleMessage).Methods(http.MethodPost)
	localchatRouter.HandleFunc("/upload", localchat.MultiUploadHandler).Methods(http.MethodPost)
	localchatRouter.HandleFunc("/check", localchat.CheckUserHanlder).Methods(http.MethodGet)

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

	// 将静态文件服务放在最后，作为默认处理程序
	router.PathPrefix("/").Handler(http.NotFoundHandler())
	if staticRouter != nil {
		go func() {
			log.Printf("Static Server listening on port: %v", staticServerAddress)
			staticSrv = &http.Server{Addr: staticServerAddress, Handler: staticRouter}
			Serve(staticSrv)
		}()
	}
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
	if staticSrv != nil {
		if err := staticSrv.Shutdown(ctx); err != nil {
			log.Fatalf("Static Server forced to shutdown: %v", err)
		}
	}
	log.Println("Server stopped.")
}

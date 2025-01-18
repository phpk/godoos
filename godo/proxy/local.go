package proxy

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"godo/libs"
	"godo/model"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"
)

// 全局代理服务映射
var proxyServers sync.Map

// ProxyServer 结构体用于存储服务类型和实际服务对象
type ProxyServer struct {
	Type   string
	Server interface{}
}

// FileServer 结构体用于存储文件静态服务的信息
type FileServer struct {
	Port   int
	Server *http.Server
}

// 创建 LocalProxy 的 HTTP 处理函数
func CreateLocalProxyHandler(w http.ResponseWriter, r *http.Request) {
	var lp model.LocalProxy
	err := json.NewDecoder(r.Body).Decode(&lp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// 残留数据id 设置为 0
	lp.ID = 0

	err = model.Db.Create(&lp).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 启动代理服务
	go startProxy(lp)

	libs.SuccessMsg(w, lp, "")
}

// GetLocalProxiesHandler 获取所有 LocalProxy 的 HTTP 处理函数
func GetLocalProxiesHandler(w http.ResponseWriter, r *http.Request) {
	// 获取查询参数 page 和 limit
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	// 定义响应结构体
	type ProxyResponse struct {
		Proxies []model.LocalProxy `json:"proxies"`
		Total   int64              `json:"total"`
	}

	// 修改处理函数
	proxies, total, err := model.GetLocalProxies(page, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	libs.SuccessMsg(w, ProxyResponse{Proxies: proxies, Total: total}, "")
}

// GetLocalProxyHandler 获取单个 LocalProxy 的 HTTP 处理函数
func GetLocalProxyHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	var proxy model.LocalProxy
	err = model.Db.First(&proxy, uint(id)).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	libs.SuccessMsg(w, proxy, "")
}

// UpdateLocalProxyHandler 更新 LocalProxy 的 HTTP 处理函数
func UpdateLocalProxyHandler(w http.ResponseWriter, r *http.Request) {
	var lp model.LocalProxy
	err := json.NewDecoder(r.Body).Decode(&lp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if lp.ID == 0 {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	updata := map[string]interface{}{
		"port":       lp.Port,
		"proxy_type": lp.ProxyType,
		"domain":     lp.Domain,
		"status":     lp.Status,
		// path
	}

	err = model.Db.Model(&model.LocalProxy{}).Where("id = ?", lp.ID).Updates(updata).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 停止旧的代理服务
	stopProxy(lp.ID)

	if !lp.Status {
		// 启动新的代理服务
		go startProxy(lp)
	}

	libs.SuccessMsg(w, lp, "")
}

// DeleteLocalProxyHandler 删除 LocalProxy 的 HTTP 处理函数
func DeleteLocalProxyHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// 停止代理服务
	stopProxy(uint(id))

	err = model.Db.Delete(&model.LocalProxy{}, uint(id)).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	libs.SuccessMsg(w, nil, "delete proxy success")
}

// HandlerSetProxyStatus 设置代理状态
func HandlerSetProxyStatus(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		libs.ErrorMsg(w, "id is empty")
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil || id == 0 {
		libs.ErrorMsg(w, "id is not number")
		return
	}
	var proxy model.LocalProxy
	if err := model.Db.First(&proxy, uint(id)).Error; err != nil {
		libs.ErrorMsg(w, "proxy not found")
		return
	}
	if err := model.Db.Model(&model.LocalProxy{}).Where("id = ?", proxy.ID).Update("status", !proxy.Status).Error; err != nil {
		libs.ErrorMsg(w, "update proxy status failed")
		return
	}
	if !proxy.Status {
		startProxy(proxy)
	} else {
		stopProxy(proxy.ID)
	}
	libs.SuccessMsg(w, nil, "")

}

// 初始化代理处理函数
func InitProxyHandlers() {
	go InitFrpcServer()
	proxies, _, err := model.GetLocalProxies(1, 1000) // 获取所有代理配置
	if err != nil {
		fmt.Println("Failed to get local proxies:", err)
		return
	}

	for _, proxy := range proxies {
		go startProxy(proxy)
	}
}

// 启动代理服务
func startProxy(proxy model.LocalProxy) {
	switch proxy.ProxyType {
	case "http":
		go func(p model.LocalProxy) {
			httpProxyHandler(p)
		}(proxy)
	case "udp":
		go func(p model.LocalProxy) {
			udpProxyHandler(p)
		}(proxy)
	case "file":
		go func(p model.LocalProxy) {
			fileServerHandler(p)
		}(proxy)
	default:
		fmt.Printf("Unknown proxy type: %s\n", proxy.ProxyType)
	}
}

// 停止代理服务
func stopProxy(id uint) {
	if server, ok := proxyServers.Load(id); ok {
		proxyServer, ok := server.(ProxyServer)
		if !ok {
			fmt.Printf("Failed to load proxy server for ID %d\n", id)
			return
		}

		switch proxyServer.Type {
		case "http":
			fmt.Println("closing http proxy...")
			httpServer, ok := proxyServer.Server.(*http.Server)
			if ok {
				// 创建一个上下文，用于传递给 server.Shutdown(ctx)
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				defer httpServer.Close()
				if err := httpServer.Shutdown(ctx); err != nil {
					fmt.Printf("Failed to shutdown HTTP server on port %s: %v\n", httpServer.Addr, err)
				} else {
					fmt.Printf("Stopped HTTP server on port %s\n", httpServer.Addr)
				}
			}
		case "udp":
			udpConn, ok := proxyServer.Server.(*net.UDPConn)
			if ok {
				udpConn.Close()
				fmt.Printf("Stopped UDP server on port %d\n", udpConn.LocalAddr().(*net.UDPAddr).Port)
			}
		case "file":
			fileServer, ok := proxyServer.Server.(*FileServer)
			if ok {
				// 创建一个上下文，用于传递给 server.Shutdown(ctx)
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				if err := fileServer.Server.Shutdown(ctx); err != nil {
					fmt.Printf("Failed to shutdown file server on port %d: %v\n", fileServer.Port, err)
				} else {
					fmt.Printf("Stopped file server on port %d\n", fileServer.Port)
				}
			}
		default:
			fmt.Printf("Unknown proxy type: %s\n", proxyServer.Type)
		}

		proxyServers.Delete(id)
	}
}

// HTTP 代理处理函数
func httpProxyHandler(proxy model.LocalProxy) {
	remote, err := url.Parse(proxy.Domain)
	if err != nil {
		fmt.Printf("Failed to parse remote URL for port %d: %v\n", proxy.Port, err)
		return
	}

	if remote.Scheme == "" {
		fmt.Printf("Remote URL for port %d does not contain a scheme (http/https): %s\n", proxy.Port, proxy.Domain)
		return
	}

	// 创建反向代理
	reverseProxy := httputil.NewSingleHostReverseProxy(remote)

	// 设置请求头（如需要可以启用）
	reverseProxy.Director = func(req *http.Request) {
		req.Header.Set("X-Forwarded-Host", req.Host)
		req.Header.Set("X-Origin-Host", remote.Host)
		req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; curl/7.68.0)") // 设置常见的 User-Agent
		req.Host = remote.Host
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
	}

	reverseProxy.ErrorLog = log.New(os.Stderr, "proxy-debug: ", log.LstdFlags)

	// 启动 HTTP 服务器并监听指定端口
	bindAddress := "127.0.0.1"
	serverAddr := fmt.Sprintf("%s:%d", bindAddress, proxy.Port)
	server := &http.Server{
		Addr:    serverAddr,
		Handler: reverseProxy,
	}

	proxyServers.Store(proxy.ID, ProxyServer{Type: "http", Server: server})

	fmt.Printf("Starting HTTP proxy on LocalHost port %d and forwarding to %s \n", proxy.Port, proxy.Domain)
	if err := server.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("HTTP proxy on port %d stopped gracefully.\n", proxy.Port)
		} else {
			fmt.Printf("Error starting HTTP proxy on port %d: %v\n", proxy.Port, err)
		}
	}
}

// UDP 代理处理函数
func udpProxyHandler(proxy model.LocalProxy) {
	localAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", proxy.Port))
	if err != nil {
		fmt.Printf("Failed to resolve local UDP address for port %d: %v\n", proxy.Port, err)
		return
	}

	remoteAddr, err := net.ResolveUDPAddr("udp", proxy.Domain)
	if err != nil {
		fmt.Printf("Failed to resolve remote UDP address for port %d: %v\n", proxy.Port, err)
		return
	}

	conn, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		fmt.Printf("Failed to listen on UDP port %d: %v\n", proxy.Port, err)
		return
	}
	defer conn.Close()

	proxyServers.Store(proxy.ID, ProxyServer{Type: "udp", Server: conn})

	buffer := make([]byte, 1024)
	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Printf("Failed to read from UDP: %v\n", err)
			continue
		}

		remoteConn, err := net.DialUDP("udp", nil, remoteAddr)
		if err != nil {
			fmt.Printf("Failed to dial remote UDP: %v\n", err)
			continue
		}

		_, err = remoteConn.Write(buffer[:n])
		if err != nil {
			fmt.Printf("Failed to write to remote UDP: %v\n", err)
			continue
		}

		// 读取远程服务器的响应并转发回客户端
		n, err = remoteConn.Read(buffer)
		if err != nil {
			fmt.Printf("Failed to read from remote UDP: %v\n", err)
			continue
		}

		_, err = conn.WriteToUDP(buffer[:n], addr)
		if err != nil {
			fmt.Printf("Failed to write to client UDP: %v\n", err)
			continue
		}
	}
}

// 文件静态服务处理函数
func fileServerHandler(proxy model.LocalProxy) {
	// 创建文件服务器的 HTTP 处理函数
	fileHandler := http.FileServer(http.Dir(proxy.Path))

	// 启动 HTTP 服务器并监听指定端口
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", proxy.Port),
		Handler: fileHandler,
	}

	proxyServers.Store(proxy.ID, ProxyServer{Type: "file", Server: &FileServer{Port: int(proxy.Port), Server: server}})

	fmt.Printf("Starting file server on port %d serving files from %s\n", proxy.Port, proxy.Path)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Printf("Failed to start file server on port %d: %v\n", proxy.Port, err)
	}
}

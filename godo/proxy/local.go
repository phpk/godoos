package proxy

import (
	"encoding/json"
	"fmt"
	"godo/model"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
)

// 创建 LocalProxy 的 HTTP 处理函数
func CreateLocalProxyHandler(w http.ResponseWriter, r *http.Request) {

	var lp model.LocalProxy
	err := json.NewDecoder(r.Body).Decode(&lp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = model.Db.Create(&lp).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(lp)

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

	w.Header().Set("Content-Type", "application/json")
	response := ProxyResponse{
		Proxies: proxies,
		Total:   total,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(proxy)
}

// UpdateLocalProxyHandler 更新 LocalProxy 的 HTTP 处理函数
func UpdateLocalProxyHandler(w http.ResponseWriter, r *http.Request) {

	var lp model.LocalProxy
	err := json.NewDecoder(r.Body).Decode(&lp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = model.Db.Model(&model.LocalProxy{}).Where("id = ?", lp.ID).Updates(lp).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(lp)
}

// DeleteLocalProxyHandler 删除 LocalProxy 的 HTTP 处理函数
func DeleteLocalProxyHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = model.Db.Delete(&model.LocalProxy{}, uint(id)).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
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
		switch proxy.ProxyType {
		case "http":
			go func(p model.LocalProxy) {
				httpProxyHandler(p)
			}(proxy)
		case "udp":
			// UDP 转发的处理逻辑
			go func(p model.LocalProxy) {
				udpProxyHandler(p)
			}(proxy)
		case "file":
			http.Handle(fmt.Sprintf("/file/%d/", proxy.Port), http.StripPrefix(fmt.Sprintf("/file/%d/", proxy.Port), http.FileServer(http.Dir(proxy.Domain))))
			fmt.Printf("Setting up file server for port %d at %s\n", proxy.Port, proxy.Domain)
		default:
			fmt.Printf("Unknown proxy type: %s\n", proxy.ProxyType)
		}
	}
}

// HTTP 代理处理函数
func httpProxyHandler(proxy model.LocalProxy) {
	remote, err := url.Parse(proxy.Domain)
	if err != nil {
		fmt.Printf("Failed to parse remote URL for port %d: %v\n", proxy.Port, err)
		return
	}

	reverseProxy := httputil.NewSingleHostReverseProxy(remote)

	// 启动 HTTP 服务器并监听指定端口
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", proxy.Port),
		Handler: reverseProxy,
	}

	fmt.Printf("Starting HTTP proxy on port %d and forwarding to %s:%d\n", proxy.Port, proxy.Domain, proxy.Port)
	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("Failed to start HTTP proxy on port %d: %v\n", proxy.Port, err)
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

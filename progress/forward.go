package progress

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/gorilla/mux"
)

func ForwardRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	proxyName := vars["name"] // 获取代理名称
	log.Printf("Forwarding request to proxy %s", proxyName)
	log.Printf("processes: %+v", processes)
	cmd, ok := processes[proxyName]
	if !ok || !cmd.Running {
		// Start the process again
		err := ExecuteScript(proxyName)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to restart process %s: %v", proxyName, err))
			return
		}
		time.Sleep(1 * time.Second)
	}
	if cmd.PingURL == "" {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get url %s", proxyName))
		return
	}
	// 构造目标基础URL，假设proxyName直接对应目标服务的基础路径
	targetBaseURLStr := cmd.PingURL
	targetBaseURL, err := url.Parse(targetBaseURLStr)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to parse target base URL")
		return
	}
	// 获取请求的子路径
	subpath := r.URL.Path[len("/"+proxyName):]

	targetURL := &url.URL{
		Scheme:   targetBaseURL.Scheme,
		Host:     targetBaseURL.Host,
		Path:     path.Clean(path.Join(targetBaseURL.Path, subpath)), // 使用path.Clean处理路径
		RawQuery: r.URL.RawQuery,                                     // 正确编码查询参数
	}
	log.Printf("Forwarding request to %s", targetURL)

	// 执行重定向
	http.Redirect(w, r, targetURL.String(), http.StatusTemporaryRedirect)
	/*



		// 保留原请求的路径和查询字符串，构建完整的目标URL
		// 如果没有子路径，直接使用基础URL
		// 如果有子路径，将其附加到基础URL

		client := http.DefaultClient
		req, err := http.NewRequest(r.Method, targetURL.String(), r.Body)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Failed to create request: "+err.Error())
			return
		}

		copyHeader(req.Header, r.Header)

		resp, err := client.Do(req)
		if err != nil {
			respondWithError(w, http.StatusServiceUnavailable, "Failed to forward request: "+err.Error())
			return
		}
		defer resp.Body.Close()

		// 将响应头复制到w.Header()
		for k, vv := range resp.Header {
			for _, v := range vv {
				w.Header().Add(k, v)
			}
		}

		w.WriteHeader(resp.StatusCode)
		if _, err := io.Copy(w, resp.Body); err != nil {
			log.Printf("Error copying response body: %v", err)
		}*/
}

// 复制源Header到目标Header
// func copyHeader(dst, src http.Header) {
// 	for k, vv := range src {
// 		for _, v := range vv {
// 			dst.Add(k, v)
// 		}
// 	}
// }

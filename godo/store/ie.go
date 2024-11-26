package store

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

type PageState struct {
	URL     string `json:"url"`
	History []string
	Index   int
}

var pageState = &PageState{}

func fetchPageContent(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	html := string(body)

	// 插入 <base> 标签
	baseTag := fmt.Sprintf("<base href=\"%s\" />", url)
	html = strings.Replace(html, "<head>", fmt.Sprintf("<head>%s", baseTag), 1)
	// 过滤掉可能引起弹窗的脚本
	re := regexp.MustCompile(`(window\.open)\(.*?\);`)
	html = re.ReplaceAllString(html, "")
	// 移除或修改 a 标签的 target="_blank" 属性
	reA := regexp.MustCompile(`(<a[^>]*?)\s*target\s*=\s*"?_blank"?`)
	html = reA.ReplaceAllString(html, "$1")
	// 移除 X-Frame-Options 和 Content-Security-Policy 头
	resp.Header.Del("X-Frame-Options")
	resp.Header.Del("Content-Security-Policy")
	html += `<script>
        window.onload = function() {
          window.parent.postMessage(window.location.href, '*');
        };
      </script>`
	return html, nil
}

func HandleNavigate(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	if url == "" {
		http.Error(w, "URL parameter is required", http.StatusBadRequest)
		return
	}

	html, err := fetchPageContent(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pageState.History = append(pageState.History, url)
	pageState.Index = len(pageState.History) - 1
	pageState.URL = url

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}

func HandleBack(w http.ResponseWriter, r *http.Request) {
	if pageState.Index > 0 {
		pageState.Index--
		url := pageState.History[pageState.Index]
		html, err := fetchPageContent(url)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pageState.URL = url
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, html)
	} else {
		http.Error(w, "No more history to go back", http.StatusForbidden)
	}
}

func HandleForward(w http.ResponseWriter, r *http.Request) {
	if pageState.Index < len(pageState.History)-1 {
		pageState.Index++
		url := pageState.History[pageState.Index]
		html, err := fetchPageContent(url)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pageState.URL = url
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, html)
	} else {
		http.Error(w, "No more history to go forward", http.StatusForbidden)
	}
}

func HandleRefresh(w http.ResponseWriter, r *http.Request) {
	if pageState.URL != "" {
		html, err := fetchPageContent(pageState.URL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, html)
	} else {
		http.Error(w, "No current URL to refresh", http.StatusForbidden)
	}
}

// func main() {
// 	http.HandleFunc("/navigate", HandleNavigate)
// 	http.HandleFunc("/back", HandleBack)
// 	http.HandleFunc("/forward", HandleForward)
// 	http.HandleFunc("/refresh", HandleRefresh)
//
// 	http.ListenAndServe(":8080", nil)
// }

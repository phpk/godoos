package proxy

import (
	"encoding/json"
	"fmt"
	"godo/libs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/ini.v1"
)

// FrpConfig 结构体用于存储 FRP 配置
type FrpConfig struct {
	ServerAddr                 string // 服务器地址
	ServerPort                 int    // 服务器端口
	AuthMethod                 string // 认证方法，例如 "token" 或 "multiuser"
	AuthToken                  string // 认证令牌，当 AuthMethod 为 "token" 时使用
	User                       string // 用户名，当 AuthMethod 为 "multiuser" 时使用
	MetaToken                  string // 元数据令牌，当 AuthMethod 为 "multiuser" 时使用
	TransportHeartbeatInterval int    // 传输心跳间隔（秒）
	TransportHeartbeatTimeout  int    // 传输心跳超时（秒）
	LogLevel                   string // 日志级别，例如 "info"、"warn"、"error"
	LogMaxDays                 int    // 日志保留天数
	WebPort                    int    // Web 管理界面的端口
	TlsConfigEnable            bool   // 是否启用 TLS 配置
	TlsConfigCertFile          string // TLS 证书文件路径
	TlsConfigKeyFile           string // TLS 密钥文件路径
	TlsConfigTrustedCaFile     string // TLS 信任的 CA 文件路径
	TlsConfigServerName        string // TLS 服务器名称
	ProxyConfigEnable          bool   // 是否启用代理配置
	ProxyConfigProxyUrl        string // 代理 URL
}

// Proxy 结构体用于存储代理配置
type Proxy struct {
	Name              string   `json:"name"`              // 代理名称
	Type              string   `json:"type"`              // 代理类型，例如 "tcp"、"http"、"stcp"
	LocalIp           string   `json:"localIp"`           // 本地 IP 地址
	LocalPort         int      `json:"localPort"`         // 本地端口
	RemotePort        int      `json:"remotePort"`        // 远程端口
	CustomDomains     []string `json:"customDomains"`     // 自定义域名列表
	Subdomain         string   `json:"subdomain"`         // 子域名
	BasicAuth         bool     `json:"basicAuth"`         // 是否启用基本认证
	HttpUser          string   `json:"httpUser"`          // HTTP 基本认证用户名
	HttpPassword      string   `json:"httpPassword"`      // HTTP 基本认证密码
	StcpModel         string   `json:"stcpModel"`         // STCP 模式，例如 "visitors" 或 "visited"
	ServerName        string   `json:"serverName"`        // 服务器名称
	BindAddr          string   `json:"bindAddr"`          // 绑定地址
	BindPort          int      `json:"bindPort"`          // 绑定端口
	FallbackTo        string   `json:"fallbackTo"`        // 回退到的目标
	FallbackTimeoutMs int      `json:"fallbackTimeoutMs"` // 回退超时时间（毫秒）
	SecretKey         string   `json:"secretKey"`         // 密钥，用于加密通信
}

type Response struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Success bool        `json:"success"`
}

func NewResponse(code string, message string, data interface{}) Response {
	return Response{
		Code:    code,
		Message: message,
		Data:    data,
		Success: true,
	}
}

// isRangePort 检查端口是否为范围端口
func isRangePort(proxy Proxy) bool {
	// 这里假设范围端口的判断逻辑
	// 你可以根据实际情况调整
	return strings.Contains(proxy.Name, ":")
}

// GenFrpcIniConfig 生成 FRP 配置文件内容
func GenFrpcIniConfig(config FrpConfig, proxys []Proxy) string {
	var proxyIni []string

	for _, m := range proxys {
		rangePort := isRangePort(m)
		ini := fmt.Sprintf("[%s%s]\ntype = \"%s\"\n",
			func() string {
				if rangePort {
					return "range:"
				}
				return ""
			}(),
			m.Name,
			m.Type,
		)

		switch m.Type {
		case "tcp", "udp":
			ini += fmt.Sprintf(`
localIP = "%s"
localPort = %d
remotePort = %d
`,
				m.LocalIp,
				m.LocalPort,
				m.RemotePort,
			)
		case "http", "https":
			ini += fmt.Sprintf(`
localIP = "%s"
localPort = %d
`,
				m.LocalIp,
				m.LocalPort,
			)

			if len(m.CustomDomains) > 0 {
				ini += fmt.Sprintf(`custom_domains = [%s]
`,
					strings.Join(m.CustomDomains, ","),
				)
			}

			if m.Subdomain != "" {
				ini += fmt.Sprintf(`subdomain="%s"
`,
					m.Subdomain,
				)
			}
			if m.BasicAuth {
				ini += fmt.Sprintf(`
httpUser = "%s"
httpPassword = "%s"
`,
					m.HttpUser,
					m.HttpPassword,
				)
			}
		case "stcp", "xtcp", "sudp":
			if m.StcpModel == "visitors" {
				// 访问者
				ini += fmt.Sprintf(`
role = visitor
serverName = "%s"
bindAddr = "%s"
bindPort = %d
`,
					m.ServerName,
					m.BindAddr,
					m.BindPort,
				)
				if m.FallbackTo != "" {
					ini += fmt.Sprintf(`
fallbackTo = %s
fallbackTimeoutMs = %d
`,
						m.FallbackTo,
						m.FallbackTimeoutMs,
					)
				}
			} else if m.StcpModel == "visited" {
				// 被访问者
				ini += fmt.Sprintf(`
localIP = "%s"
localPort = %d
`,
					m.LocalIp,
					m.LocalPort,
				)
			}
			ini += fmt.Sprintf(`
sk="%s"
`,
				m.SecretKey,
			)
		default:
			// 默认情况不做处理
		}

		proxyIni = append(proxyIni, ini)
	}

	ini := fmt.Sprintf(`[common]
serverAddr = %s
serverPort = %d
%s
%s
%s
%s
logFile = "frpc.log"
logLevel = %s
logMaxDays = %d
adminAddr = 127.0.0.1
adminPort = %d
tlsEnable = %t
%s
%s
%s
%s
%s
`,
		config.ServerAddr,
		config.ServerPort,
		func() string {
			if config.AuthMethod == "token" {
				return fmt.Sprintf(`
authenticationMethod = %s
token = %s
`,
					config.AuthMethod,
					config.AuthToken,
				)
			}
			return ""
		}(),
		func() string {
			if config.AuthMethod == "multiuser" {
				return fmt.Sprintf(`
user = %s
metaToken = %s
`,
					config.User,
					config.MetaToken,
				)
			}
			return ""
		}(),
		func() string {
			if config.TransportHeartbeatInterval > 0 {
				return fmt.Sprintf(`
heartbeatInterval = %d
`,
					config.TransportHeartbeatInterval,
				)
			}
			return ""
		}(),
		func() string {
			if config.TransportHeartbeatTimeout > 0 {
				return fmt.Sprintf(`
heartbeatTimeout = %d
`,
					config.TransportHeartbeatTimeout,
				)
			}
			return ""
		}(),
		config.LogLevel,
		config.LogMaxDays,
		config.WebPort,
		config.TlsConfigEnable,
		func() string {
			if config.TlsConfigEnable && config.TlsConfigCertFile != "" {
				return fmt.Sprintf(`
tlsCertFile = %s
`,
					config.TlsConfigCertFile,
				)
			}
			return ""
		}(),
		func() string {
			if config.TlsConfigEnable && config.TlsConfigKeyFile != "" {
				return fmt.Sprintf(`
tlsKeyFile = %s
`,
					config.TlsConfigKeyFile,
				)
			}
			return ""
		}(),
		func() string {
			if config.TlsConfigEnable && config.TlsConfigTrustedCaFile != "" {
				return fmt.Sprintf(`
tlsTrustedCaFile = %s
`,
					config.TlsConfigTrustedCaFile,
				)
			}
			return ""
		}(),
		func() string {
			if config.TlsConfigEnable && config.TlsConfigServerName != "" {
				return fmt.Sprintf(`
tlsServerName = %s
`,
					config.TlsConfigServerName,
				)
			}
			return ""
		}(),
		func() string {
			if config.ProxyConfigEnable {
				return fmt.Sprintf(`
httpProxy = "%s"
`,
					config.ProxyConfigProxyUrl,
				)
			}
			return ""
		}(),
	)

	ini += strings.Join(proxyIni, "")

	return ini
}

// CreateFrpcHandler 创建frpc配置
func CreateFrpcHandler(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Config  FrpConfig `json:"config"`
		Proxies []Proxy   `json:"proxies"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 读取现有配置
	path := libs.GetAppExecDir()
	filePath := filepath.Join(path, "frpc", "frpc.ini")

	var existingProxies []Proxy
	if iniContent, err := os.ReadFile(filePath); err == nil {
		cfg, err := ini.Load(iniContent)
		if err == nil {
			for _, section := range cfg.Sections() {
				if section.Name() == "common" {
					continue
				}
				proxy := Proxy{
					Name:              section.Name(),
					Type:              section.Key("type").String(),
					LocalIp:           section.Key("localIP").String(),
					LocalPort:         section.Key("localPort").MustInt(),
					RemotePort:        section.Key("remotePort").MustInt(),
					CustomDomains:     strings.Split(section.Key("custom_domains").String(), ","),
					Subdomain:         section.Key("subdomain").String(),
					BasicAuth:         section.Key("httpUser").String() != "",
					HttpUser:          section.Key("httpUser").String(),
					HttpPassword:      section.Key("httpPassword").String(),
					StcpModel:         section.Key("role").String(),
					ServerName:        section.Key("serverName").String(),
					BindAddr:          section.Key("bindAddr").String(),
					BindPort:          section.Key("bindPort").MustInt(),
					FallbackTo:        section.Key("fallbackTo").String(),
					FallbackTimeoutMs: section.Key("fallbackTimeoutMs").MustInt(),
					SecretKey:         section.Key("sk").String(),
				}
				existingProxies = append(existingProxies, proxy)
			}
		}
	}

	// 合并现有代理和新代理
	proxyMap := make(map[string]Proxy)
	for _, proxy := range existingProxies {
		proxyMap[proxy.Name] = proxy
	}

	for _, proxy := range requestData.Proxies {
		proxyMap[proxy.Name] = proxy // 新代理会覆盖同 Name 的旧代理
	}

	var allProxies []Proxy
	for _, proxy := range proxyMap {
		allProxies = append(allProxies, proxy)
	}

	// 生成新的配置内容
	iniContent := GenFrpcIniConfig(requestData.Config, allProxies)

	if _, err := os.Stat(filepath.Join(path, "frpc")); os.IsNotExist(err) {
		os.MkdirAll(filepath.Join(path, "frpc"), 0755)
	}

	err = os.WriteFile(filePath, []byte(iniContent), 0644)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(NewResponse("0", "Configuration created successfully", nil))
}

// GetFrpcConfigHandler 处理获取单个 FRPC 配置的 HTTP 请求
func GetFrpcConfigHandler(w http.ResponseWriter, r *http.Request) {
	// 获取请求参数中的 name
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Name parameter is required", http.StatusBadRequest)
		return
	}

	// 使用 os.ReadFile 读取配置文件内容
	path := libs.GetAppExecDir()
	filePath := filepath.Join(path, "frpc", "frpc.ini")

	iniContent, err := os.ReadFile(filePath)
	if err != nil {
		http.Error(w, "Failed to read configuration file", http.StatusInternalServerError)
		return
	}

	// 解析 ini 文件
	cfg, err := ini.Load(iniContent)
	if err != nil {
		http.Error(w, "Failed to parse configuration file", http.StatusInternalServerError)
		return
	}

	// 提取 [common] 部分
	commonSection := cfg.Section("common")
	frpConfig := FrpConfig{
		ServerAddr:                 commonSection.Key("serverAddr").String(),
		ServerPort:                 commonSection.Key("serverPort").MustInt(),
		AuthMethod:                 commonSection.Key("authenticationMethod").String(),
		AuthToken:                  commonSection.Key("token").String(),
		User:                       commonSection.Key("user").String(),
		MetaToken:                  commonSection.Key("metaToken").String(),
		TransportHeartbeatInterval: commonSection.Key("heartbeatInterval").MustInt(),
		TransportHeartbeatTimeout:  commonSection.Key("heartbeatTimeout").MustInt(),
		LogLevel:                   commonSection.Key("logLevel").String(),
		LogMaxDays:                 commonSection.Key("logMaxDays").MustInt(),
		WebPort:                    commonSection.Key("adminPort").MustInt(),
		TlsConfigEnable:            commonSection.Key("tlsEnable").MustBool(),
		TlsConfigCertFile:          commonSection.Key("tlsCertFile").String(),
		TlsConfigKeyFile:           commonSection.Key("tlsKeyFile").String(),
		TlsConfigTrustedCaFile:     commonSection.Key("tlsTrustedCaFile").String(),
		TlsConfigServerName:        commonSection.Key("tlsServerName").String(),
		ProxyConfigEnable:          commonSection.Key("httpProxy").String() != "",
		ProxyConfigProxyUrl:        commonSection.Key("httpProxy").String(),
	}

	// 查找指定 name 的代理
	var proxy *Proxy
	for _, section := range cfg.Sections() {
		if section.Name() == "common" {
			continue
		}
		if section.Name() == name {
			proxy = &Proxy{
				Name:       section.Name(),
				Type:       section.Key("type").String(),
				LocalIp:    section.Key("localIP").String(),
				LocalPort:  section.Key("localPort").MustInt(),
				RemotePort: section.Key("remotePort").MustInt(),
				// 继续提取其他字段...
			}
			break
		}
	}

	if proxy == nil {
		http.Error(w, "Proxy not found", http.StatusNotFound)
		return
	}

	// 创建响应
	responseData := struct {
		Config FrpConfig `json:"config"`
		Proxy  Proxy     `json:"proxy"`
	}{
		Config: frpConfig,
		Proxy:  *proxy,
	}

	response := NewResponse("0", "Configuration retrieved successfully", responseData)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetFrpcProxiesHandler 获取 FRPC 代理列表
func GetFrpcProxiesHandler(w http.ResponseWriter, r *http.Request) {
	// 使用 os.ReadFile 读取配置文件内容
	path := libs.GetAppExecDir()
	filePath := filepath.Join(path, "frpc", "frpc.ini")

	iniContent, err := os.ReadFile(filePath)
	if err != nil {
		http.Error(w, "Failed to read configuration file", http.StatusInternalServerError)
		return
	}

	// 解析 ini 文件
	cfg, err := ini.Load(iniContent)
	if err != nil {
		http.Error(w, "Failed to parse configuration file", http.StatusInternalServerError)
		return
	}

	// 提取代理部分
	var proxies []Proxy
	for _, section := range cfg.Sections() {
		if section.Name() == "common" {
			continue
		}
		proxy := Proxy{
			Name:              section.Name(),
			Type:              section.Key("type").String(),
			LocalIp:           section.Key("localIP").String(),
			LocalPort:         section.Key("localPort").MustInt(),
			RemotePort:        section.Key("remotePort").MustInt(),
			CustomDomains:     strings.Split(section.Key("custom_domains").String(), ","),
			Subdomain:         section.Key("subdomain").String(),
			BasicAuth:         section.Key("httpUser").String() != "",
			HttpUser:          section.Key("httpUser").String(),
			HttpPassword:      section.Key("httpPassword").String(),
			StcpModel:         section.Key("role").String(),
			ServerName:        section.Key("serverName").String(),
			BindAddr:          section.Key("bindAddr").String(),
			BindPort:          section.Key("bindPort").MustInt(),
			FallbackTo:        section.Key("fallbackTo").String(),
			FallbackTimeoutMs: section.Key("fallbackTimeoutMs").MustInt(),
			SecretKey:         section.Key("sk").String(),
		}
		proxies = append(proxies, proxy)
	}

	// 创建响应
	response := NewResponse("0", "Proxies retrieved successfully", proxies)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func DeleteFrpcConfigHandler(w http.ResponseWriter, r *http.Request) {
	// 获取请求参数中的 name
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Name parameter is required", http.StatusBadRequest)
		return
	}

	// 使用 os.ReadFile 读取配置文件内容
	path := libs.GetAppExecDir()
	filePath := filepath.Join(path, "frpc", "frpc.ini")

	iniContent, err := os.ReadFile(filePath)
	if err != nil {
		http.Error(w, "Failed to read configuration file", http.StatusInternalServerError)
		return
	}

	// 解析 ini 文件
	cfg, err := ini.Load(iniContent)
	if err != nil {
		http.Error(w, "Failed to parse configuration file", http.StatusInternalServerError)
		return
	}

	// 检查是否存在指定 name 的代理
	if cfg.Section(name) == nil {
		http.Error(w, "Proxy not found", http.StatusNotFound)
		return
	}

	// 删除指定 name 的代理
	cfg.DeleteSection(name)

	// 将更新后的配置写回文件
	err = cfg.SaveTo(filePath)
	if err != nil {
		http.Error(w, "Failed to save configuration file", http.StatusInternalServerError)
		return
	}

	// 创建响应
	response := NewResponse("0", "Configuration deleted successfully", nil)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func UpdateFrpcConfigHandler(w http.ResponseWriter, r *http.Request) {
	// 获取请求参数中的 name
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Name parameter is required", http.StatusBadRequest)
		return
	}

	// 解析请求体中的更新数据
	var updateData Proxy
	err := json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 使用 os.ReadFile 读取配置文件内容
	path := libs.GetAppExecDir()
	filePath := filepath.Join(path, "frpc", "frpc.ini")

	iniContent, err := os.ReadFile(filePath)
	if err != nil {
		http.Error(w, "Failed to read configuration file", http.StatusInternalServerError)
		return
	}

	// 解析 ini 文件
	cfg, err := ini.Load(iniContent)
	if err != nil {
		http.Error(w, "Failed to parse configuration file", http.StatusInternalServerError)
		return
	}

	// 查找指定 name 的代理
	section := cfg.Section(name)
	if section == nil {
		http.Error(w, "Proxy not found", http.StatusNotFound)
		return
	}

	// 更新代理配置
	section.Key("type").SetValue(updateData.Type)
	section.Key("localIP").SetValue(updateData.LocalIp)
	section.Key("localPort").SetValue(fmt.Sprintf("%d", updateData.LocalPort))
	section.Key("remotePort").SetValue(fmt.Sprintf("%d", updateData.RemotePort))
	section.Key("custom_domains").SetValue(strings.Join(updateData.CustomDomains, ","))
	section.Key("subdomain").SetValue(updateData.Subdomain)
	section.Key("httpUser").SetValue(updateData.HttpUser)
	section.Key("httpPassword").SetValue(updateData.HttpPassword)
	section.Key("role").SetValue(updateData.StcpModel)
	section.Key("serverName").SetValue(updateData.ServerName)
	section.Key("bindAddr").SetValue(updateData.BindAddr)
	section.Key("bindPort").SetValue(fmt.Sprintf("%d", updateData.BindPort))
	section.Key("fallbackTo").SetValue(updateData.FallbackTo)
	section.Key("fallbackTimeoutMs").SetValue(fmt.Sprintf("%d", updateData.FallbackTimeoutMs))
	section.Key("sk").SetValue(updateData.SecretKey)

	// 将更新后的配置写回文件
	err = cfg.SaveTo(filePath)
	if err != nil {
		http.Error(w, "Failed to save configuration file", http.StatusInternalServerError)
		return
	}

	// 创建响应
	response := NewResponse("0", "Configuration updated successfully", nil)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

package proxy

import (
	"encoding/json"
	"fmt"
	"godo/libs"
	"godo/model"
	"godo/progress"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// FrpConfig 结构体用于存储 FRP 配置
type FrpConfig struct {
	IsOn       bool   `json:"isOn"`
	ServerAddr string `json:"serverAddr"`
	ServerPort int    `json:"serverPort"`
	AuthMethod string `json:"authMethod"`
	AuthToken  string `json:"authToken"`
	User       string `json:"user"`
	MetaToken  string `json:"metaToken"`
}

var FrpcRequest struct {
	Config  FrpConfig         `json:"config"`
	Proxies []model.FrpcProxy `json:"proxies"`
}

func InitFrpcServer() {
	frpcConf, err := GetFrpcConfig()
	if err != nil {
		return
	}
	if frpcConf.IsOn {
		if err := progress.StartCmd("frpc"); err != nil {
			log.Printf("Failed to start frpc: %v", err)
			return
		}
	}
}

// GenFrpcIniConfig 生成 FRP 配置文件内容
func GenFrpcIniConfig(config FrpConfig, proxys []model.FrpcProxy) string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf(`serverAddr = "%s"
serverPort = %d
`,
		config.ServerAddr,
		config.ServerPort,
	))

	if config.AuthMethod == "token" {
		builder.WriteString(fmt.Sprintf(`auth.token = "%s"
`,
			config.AuthToken,
		))
	}

	for _, m := range proxys {
		builder.WriteString(fmt.Sprintf(`[[proxies]]
name="%s"
type = "%s"
`,
			m.Name,
			m.Type,
		))

		switch m.Type {
		case "tcp", "udp":
			if m.LocalIp != "" {
				builder.WriteString(fmt.Sprintf(`localIP = "%s"
`,
					m.LocalIp,
				))
			}
			if m.LocalPort > 0 {
				builder.WriteString(fmt.Sprintf(`localPort = %d
`,
					m.LocalPort,
				))
			}
			if m.RemotePort > 0 {
				builder.WriteString(fmt.Sprintf(`remotePort = %d
`,
					m.RemotePort,
				))
			}
			if m.Type == "tcp" {
				if m.StaticFile && m.LocalPath != "" {
					builder.WriteString(fmt.Sprintf(`[proxies.plugin]
type = "static_file"
localPath = "%s"
`, m.LocalPath))
					if m.StripPrefix != "" {
						builder.WriteString(fmt.Sprintf(`stripPrefix = "%s"
`,
							m.StripPrefix,
						))

					}
				}
			}
		case "http", "https":
			if m.LocalIp != "" {
				builder.WriteString(fmt.Sprintf(`localIP = "%s"
`,
					m.LocalIp,
				))
			}
			builder.WriteString(fmt.Sprintf(`localPort = %d
`,
				m.LocalPort,
			))

			if len(m.CustomDomains) > 0 {
				domains := strings.Split(m.CustomDomains, ",")
				domainArr := []string{}
				for _, domain := range domains {
					domainArr = append(domainArr, fmt.Sprintf(`"%s"`, domain))
				}
				builder.WriteString(fmt.Sprintf(`customDomains = [%s]
`,
					strings.Join(domainArr, ","),
				))
			}

			if m.Subdomain != "" {
				builder.WriteString(fmt.Sprintf(`subdomain = "%s"
`,
					m.Subdomain,
				))
			}
			if m.BasicAuth {
				builder.WriteString(fmt.Sprintf(`httpUser = "%s"
httpPassword = "%s"
`,
					m.HttpUser,
					m.HttpPassword,
				))
			}
			if m.Type == "https" {
				builder.WriteString(fmt.Sprintf(`[proxies.plugin]
type = "https2http"
localAddr = "%s:%d"
`,
					m.LocalIp,
					m.LocalPort,
				))
				if m.Https2httpCaFile != "" && m.Https2httpKeyFile != "" {
					builder.WriteString(fmt.Sprintf(`crtPath = "%s"
	keyPath = "%s"
	`,
						m.Https2httpCaFile,
						m.Https2httpKeyFile,
					))
				}
			}

		case "stcp", "xtcp", "sudp":
			if m.StcpModel == "visitors" {
				// 访问者
				builder.WriteString(fmt.Sprintf(`[[visitors]]
serverName = "%s"
bindAddr = "%s"
bindPort = %d
`,
					m.ServerName,
					m.BindAddr,
					m.BindPort,
				))
				if m.FallbackTo != "" {
					builder.WriteString(fmt.Sprintf(`fallbackTo = %s
fallbackTimeoutMs = %d
`,
						m.FallbackTo,
						m.FallbackTimeoutMs,
					))
				}
			} else if m.StcpModel == "visited" {
				// 被访问者
				builder.WriteString(fmt.Sprintf(`localIP = "%s"
localPort = %d
`,
					m.LocalIp,
					m.LocalPort,
				))
			}
			builder.WriteString(fmt.Sprintf(`sk="%s"
`,
				m.SecretKey,
			))
		default:
			// 默认情况不做处理
		}
	}

	ini := builder.String()

	// 去除多余的空格和换行符
	ini = strings.TrimSpace(ini)
	ini = strings.ReplaceAll(ini, "\n\n", "\n")

	return ini
}
func CheckConfig() error {
	runDir := libs.GetAppRunDir()
	configPath := filepath.Join(runDir, "frpc", "config.json")
	if !libs.PathExists(configPath) {
		return fmt.Errorf("请先配置服务端")
	}
	return nil
}
func GetFrpcConfig() (FrpConfig, error) {
	runDir := libs.GetAppRunDir()
	configPath := filepath.Join(runDir, "frpc", "config.json")
	if !libs.PathExists(configPath) {
		return FrpConfig{}, fmt.Errorf("请先配置服务端")
	}

	configFile, err := os.Open(configPath)
	if err != nil {
		return FrpConfig{}, err
	}
	defer configFile.Close()

	var config FrpConfig
	decoder := json.NewDecoder(configFile)
	if err := decoder.Decode(&config); err != nil {
		return FrpConfig{}, err
	}

	return config, nil
}

func SetFrpcConfig(config FrpConfig) error {
	runDir := libs.GetAppRunDir()
	configPath := filepath.Join(runDir, "frpc", "config.json")

	// 检查配置目录是否存在
	if !libs.PathExists(filepath.Join(runDir, "frpc")) {
		log.Printf("Config file not found at %s, creating new file", configPath)
	}

	// 打开或创建配置文件
	configFile, err := os.Create(configPath)
	if err != nil {
		log.Printf("Failed to create or open config file: %v", err)
		return err
	}
	defer configFile.Close()

	// 将配置写入文件
	encoder := json.NewEncoder(configFile)
	encoder.SetIndent("", "  ") // 设置缩进以便于阅读
	if err := encoder.Encode(config); err != nil {
		log.Printf("Failed to encode config to file: %v", err)
		return err
	}
	return nil
}

func UpdateFrpcIni() error {
	frpcConf, err := GetFrpcConfig()
	if err != nil {
		return err
	}
	var proxyList []model.FrpcProxy
	if err := model.Db.Model(&proxyList).Find(&proxyList).Error; err != nil {
		// 发生了查询错误，返回内部服务器错误
		return fmt.Errorf("查询失败")
	}
	runDir := libs.GetAppRunDir()
	filePath := filepath.Join(runDir, "frpc", "frpc.ini")
	// 生成新的配置内容
	iniContent := GenFrpcIniConfig(frpcConf, proxyList)

	err = os.WriteFile(filePath, []byte(iniContent), 0644)
	if err != nil {
		return fmt.Errorf("写入失败")
	}
	return nil
}

// CreateFrpcHandler 创建frpc配置
func CreateFrpcHandler(w http.ResponseWriter, r *http.Request) {
	var req model.FrpcProxy
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		libs.ErrorMsg(w, "解析结构体错误："+err.Error())
		return
	}
	err = CheckConfig()
	if err != nil {
		libs.ErrorMsg(w, err.Error())
		return
	}

	var existingProxy model.FrpcProxy
	err = model.Db.Model(&existingProxy).Where("name = ?", req.Name).First(&existingProxy).Error

	if err == nil {
		// 记录存在，返回错误
		libs.ErrorMsg(w, "已存在相同名称的配置")
		return
	}
	if err := model.Db.Create(&req).Error; err != nil {
		// 发生了插入错误，返回内部服务器错误
		libs.ErrorMsg(w, "数据库插入错误")
		return
	}
	if err := UpdateFrpcIni(); err != nil {
		libs.ErrorMsg(w, "更新配置错误："+err.Error())
		return
	}
	if err := RestartFrpc(); err != nil {
		log.Printf("Failed to restarted frpc service: %v", err)
	}
	libs.SuccessMsg(w, nil, "配置已创建")
}

// GetFrpcConfigHandler 处理获取单个 FRPC 配置的 HTTP 请求
func GetFrpcHandler(w http.ResponseWriter, r *http.Request) {
	// 获取请求参数中的 id
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Name parameter is required", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	var proxy model.FrpcProxy
	if err := model.Db.First(&proxy, uint(id)).Error; err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	libs.SuccessMsg(w, proxy, "")
}

// GetFrpcProxiesHandler 获取 FRPC 代理列表
func GetFrpcListHandler(w http.ResponseWriter, r *http.Request) {
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
		List  []model.FrpcProxy `json:"list"`
		Total int64             `json:"total"`
	}

	// 修改处理函数
	proxies, total, err := model.GetFrpcList(page, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := ProxyResponse{
		List:  proxies,
		Total: total,
	}
	libs.SuccessMsg(w, response, "")
}

func DeleteFrpcHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Name parameter is required", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := model.Db.Delete(&model.FrpcProxy{}, uint(id)).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := UpdateFrpcIni(); err != nil {
		libs.ErrorMsg(w, err.Error())
		return
	}
	if err := RestartFrpc(); err != nil {
		log.Printf("Failed to restarted frpc service: %v", err)
	}
	libs.SuccessMsg(w, nil, "配置已删除")
}

func UpdateFrpcHandler(w http.ResponseWriter, r *http.Request) {
	// 解析请求体中的更新数据
	var updateData model.FrpcProxy
	err := json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if err := model.Db.Model(&model.FrpcProxy{}).Where("id = ?", updateData.ID).Updates(updateData).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := UpdateFrpcIni(); err != nil {
		libs.ErrorMsg(w, err.Error())
		return
	}
	if err := RestartFrpc(); err != nil {
		log.Printf("Failed to restarted frpc service: %v", err)
	}
	libs.SuccessMsg(w, nil, "配置已更新")
}
func GetFrpcConfigHandler(w http.ResponseWriter, r *http.Request) {
	frpcConf, err := GetFrpcConfig()
	if err != nil {
		libs.ErrorMsg(w, err.Error())
		return
	}
	libs.SuccessMsg(w, frpcConf, "配置已更新")
}
func UpdateFrpcConfigHandler(w http.ResponseWriter, r *http.Request) {
	// 解析请求体中的更新数据
	var updateData FrpConfig
	err := json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if err := SetFrpcConfig(updateData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := UpdateFrpcIni(); err != nil {
		libs.ErrorMsg(w, err.Error())
		return
	}
	if err := RestartFrpc(); err != nil {
		log.Printf("Failed to restarted frpc service: %v", err)
	}
	libs.SuccessMsg(w, nil, "配置已更新")
}
func RestartFrpc() error {
	status := progress.GetCmd("frpc").Running
	if status {
		return progress.RestartCmd("frpc")
	}
	return nil
}
func StatusFrpcHandler(w http.ResponseWriter, r *http.Request) {
	libs.SuccessMsg(w, progress.GetCmd("frpc").Running, "")
}
func StartFrpcHandler(w http.ResponseWriter, r *http.Request) {
	err := CheckConfig()
	if err != nil {
		libs.ErrorMsg(w, err.Error())
		return
	}
	if err := progress.StartCmd("frpc"); err != nil {
		libs.ErrorMsg(w, err.Error())
		return
	}
	libs.SuccessMsg(w, nil, "frpc service started")
}

func StopFrpcHandler(w http.ResponseWriter, r *http.Request) {
	err := CheckConfig()
	if err != nil {
		libs.ErrorMsg(w, err.Error())
		return
	}
	if err := progress.StopCmd("frpc"); err != nil {
		libs.ErrorMsg(w, err.Error())
		return
	}
	libs.SuccessMsg(w, nil, "frpc service stoped")
}

func RestartFrpcHandler(w http.ResponseWriter, r *http.Request) {
	err := CheckConfig()
	if err != nil {
		libs.ErrorMsg(w, err.Error())
		return
	}
	if err := progress.RestartCmd("frpc"); err != nil {
		libs.ErrorMsg(w, err.Error())
		return
	}
	libs.SuccessMsg(w, nil, "frpc service restarted")
}

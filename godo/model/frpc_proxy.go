package model

type FrpcProxy struct {
	BaseModel
	Name              string `gorm:"unique;not null;column:name" json:"name"`             // 代理名称
	Type              string `gorm:"column:type" json:"type"`                             // 代理类型，例如 "tcp"、"http"、"stcp"
	LocalIp           string `gorm:"column:local_ip" json:"localIp"`                      // 本地 IP 地址
	LocalPort         int    `gorm:"column:local_port" json:"localPort"`                  // 本地端口
	RemotePort        int    `gorm:"column:remote_port" json:"remotePort"`                // 远程端口
	CustomDomains     string `gorm:"column:custom_domains" json:"customDomains"`          // 自定义域名列表，可以存储为逗号分隔的字符串
	Subdomain         string `gorm:"column:subdomain" json:"subdomain"`                   // 子域名
	BasicAuth         bool   `gorm:"default:false;column:basic_auth" json:"basicAuth"`    // 是否启用基本认证
	HttpUser          string `gorm:"column:http_user" json:"httpUser"`                    // HTTP 基本认证用户名
	HttpPassword      string `gorm:"column:http_password" json:"httpPassword"`            // HTTP 基本认证密码
	StcpModel         string `gorm:"column:stcp_model" json:"stcpModel"`                  // STCP 模式，例如 "visitors" 或 "visited"
	ServerName        string `gorm:"column:server_name" json:"serverName"`                // 服务器名称
	BindAddr          string `gorm:"column:bind_addr" json:"bindAddr"`                    // 绑定地址
	BindPort          int    `gorm:"column:bind_port" json:"bindPort"`                    // 绑定端口
	FallbackTo        string `gorm:"column:fallback_to" json:"fallbackTo"`                // 回退到的目标
	FallbackTimeoutMs int    `gorm:"column:fallback_timeout_ms" json:"fallbackTimeoutMs"` // 回退超时时间（毫秒）
	SecretKey         string `gorm:"column:secret_key" json:"secretKey"`                  // 密钥，用于加密通信
	//Port              string `gorm:"column:port" json:"port"`                             // 端口
	//Domain            string `gorm:"column:domain" json:"domain"`                         // 域名
	ServerAddr        string `gorm:"column:server_addr" json:"serverAddr"`                // 服务器地址
	ServerPort        int    `gorm:"column:server_port" json:"serverPort"`                // 服务器端口
	Https2http        bool   `gorm:"column:https2http" json:"https2http"`                 // 是否启用 HTTPS 到 HTTP
	Https2httpCaFile  string `gorm:"column:https2http_ca_file" json:"https2httpCaFile"`   // HTTPS 到 HTTP 的 CA 文件
	Https2httpKeyFile string `gorm:"column:https2http_key_file" json:"https2httpKeyFile"` // HTTPS 到 HTTP 的密钥文件
	KeepAlive         bool   `gorm:"column:keep_alive" json:"keepAlive"`                  // 是否保持隧道开启
	VisitedName       string `gorm:"column:visited_name" json:"visitedName"`              // 被访问者代理名称
	StaticFile        bool   `gorm:"column:static_file" json:"staticFile"`
	LocalPath         string `gorm:"column:local_path" json:"localPath"`
	StripPrefix       string `gorm:"column:strip_prefix" json:"stripPrefix"`
}

// TableName 指定表名
func (FrpcProxy) TableName() string {
	return "frpc_proxies"
}

func GetFrpcList(page, limit int) ([]FrpcProxy, int64, error) {
	var proxies []FrpcProxy
	var total int64

	// 先计算总数
	if err := Db.Model(&FrpcProxy{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 再进行分页查询
	offset := (page - 1) * limit
	if err := Db.Limit(limit).Offset(offset).Find(&proxies).Error; err != nil {
		return nil, 0, err
	}

	return proxies, total, nil
}

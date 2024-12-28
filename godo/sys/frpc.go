package sys

import (
	"fmt"
	"strings"
)

// FrpConfig 结构体用于存储 FRP 配置
type FrpConfig struct {
	ServerAddr                 string
	ServerPort                 int
	AuthMethod                 string
	AuthToken                  string
	User                       string
	MetaToken                  string
	TransportHeartbeatInterval int
	TransportHeartbeatTimeout  int
	LogLevel                   string
	LogMaxDays                 int
	WebPort                    int
	TlsConfigEnable            bool
	TlsConfigCertFile          string
	TlsConfigKeyFile           string
	TlsConfigTrustedCaFile     string
	TlsConfigServerName        string
	ProxyConfigEnable          bool
	ProxyConfigProxyUrl        string
}

// Proxy 结构体用于存储代理配置
type Proxy struct {
	Name              string
	Type              string
	LocalIp           string
	LocalPort         int
	RemotePort        int
	CustomDomains     []string
	Subdomain         string
	BasicAuth         bool
	HttpUser          string
	HttpPassword      string
	StcpModel         string
	ServerName        string
	BindAddr          string
	BindPort          int
	FallbackTo        string
	FallbackTimeoutMs int
	SecretKey         string
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

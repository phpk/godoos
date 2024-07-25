package progress

type ProcessInfo struct {
	Port     string
	PingPath string // 增加ping或健康检查路径
}

var ProcessInfoMap = map[string]ProcessInfo{
	"model":     {Port: "56711", PingPath: "ping"},
	"goconv":    {Port: "56712", PingPath: "ping"},
	"file":      {Port: "56713", PingPath: "ping"},
	"knowledge": {Port: "56714", PingPath: "ping"},
	//"llmcpu":    {Port: "56715", PingPath: "healthz"},
	"ollama":    {Port: "56715", PingPath: ""},
	"sd":        {Port: "56716", PingPath: "ping"},
	"voice":     {Port: "56717", PingPath: "ping"},
	"localchat": {Port: "56718", PingPath: "ping"},
	"llmgpu":    {Port: "56719", PingPath: "healthz"},
}

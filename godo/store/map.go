package store

type ProcessInfo struct {
	Port     string
	PingPath string // 增加ping或健康检查路径
}

var ProcessInfoMap = map[string]ProcessInfo{}

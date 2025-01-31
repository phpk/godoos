package model

type LocalProxy struct {
	BaseModel
	Port      uint   `json:"port"`      // 本地端口
	ProxyType string `json:"proxyType"` // 代理类型
	Domain    string `json:"domain"`    // 代理域名
	Path      string `json:"path"`      // 代理路径
	Status    bool   `json:"status"`    // 状态
}

func (*LocalProxy) TableName() string {
	return "local_proxy"
}
func GetLocalProxiesOn() ([]LocalProxy, error) {
	var proxies []LocalProxy
	err := Db.Where("status = ?", true).Find(&proxies).Error
	return proxies, err
}

// GetLocalProxies 获取所有 LocalProxy，支持分页
func GetLocalProxies(page, limit int) ([]LocalProxy, int64, error) {
	var proxies []LocalProxy
	var total int64

	// 先计算总数
	if err := Db.Model(&LocalProxy{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 再进行分页查询
	offset := (page - 1) * limit
	if err := Db.Limit(limit).Offset(offset).Find(&proxies).Error; err != nil {
		return nil, 0, err
	}

	return proxies, total, nil
}

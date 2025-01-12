package model

type LocalProxy struct {
	BaseModel
	Port      uint   `json:"port"`
	ProxyType string `json:"proxy_type"`
	Domain    string `json:"domain"`
}

func (*LocalProxy) TableName() string {
	return "local_proxy"
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

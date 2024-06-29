package sms4go

// SmsConfig Sms4Go 配置内容
type SmsConfig struct {
	IsPrint bool
}

// SupplierConfig 供应商配置规范
type SupplierConfig interface {
	// GetConfigId 用于获取配置对象的唯一标识符
	GetConfigId() string

	// GetSupplier 用于获取配置对象所属的供应商名称或标识符
	GetSupplier() string
}

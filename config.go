package sms4go

// SmsConfig Sms4Go 配置内容
type SmsConfig struct {
}

type SupplierConfig interface {
	GetConfigId() string

	GetSupplier() string
}

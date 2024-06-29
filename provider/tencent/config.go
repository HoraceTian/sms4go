package tencent

import "sms4go/provider"

// 基础常量
const (
	DefaultTerritory   = "ap-guangzhou"
	DefaultRequestUrl  = "sms.tencentcloudapi.com"
	DefaultConnTimeout = 60
	DefaultAction      = "SendSms"
	DefaultVersion     = "2021-01-11"
	DefaultService     = "sms"
)

type Config struct {
	provider.BaseConfig
	territory   string
	connTimeout int
	requestUrl  string
	action      string
	version     string
	service     string
}

func (c *Config) GetConfigId() string {
	return c.ConfigId
}

func (c *Config) GetSupplier() string {
	return provider.Tencent
}

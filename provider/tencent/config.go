package tencent

import "sms4go/provider"

type Config struct {
	base        provider.BaseConfig
	territory   string
	connTimeout int
	requestUrl  string
	action      string
	version     string
	service     string
}

func (c *Config) GetConfigId() string {
	return c.base.ConfigId
}

func (c *Config) GetSupplier() string {
	return provider.Tencent
}

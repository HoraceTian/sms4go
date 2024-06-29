package tencent

import (
	"sms4go"
	"sms4go/provider"
)

type Factory struct {
	instance *Blender
}

func NewFactory() *Factory {
	return &Factory{}
}

func (f Factory) CreateSms(config sms4go.SupplierConfig) sms4go.ISmsBlender {
	// 1. 单例兜底
	if f.instance != nil {
		return f.instance
	}

	// 2. configId 兜底
	configId := provider.ExtractConfigId(config)

	// 3. 初始化 Blender
	cfg := config.(*Config)
	replenishmentDefaultConfig(cfg)
	parent := sms4go.BaseBlender{
		Config:   config,
		ConfigId: configId,
	}

	// 4. 保存单例
	f.instance = &Blender{
		parent: parent,
	}
	return f.instance
}

func replenishmentDefaultConfig(config *Config) {
	if config.territory == "" {
		config.territory = DefaultTerritory
	}
	if config.connTimeout == 0 {
		config.connTimeout = DefaultConnTimeout
	}
	if config.requestUrl == "" {
		config.requestUrl = DefaultRequestUrl
	}
	if config.action == "" {
		config.action = DefaultAction
	}
	if config.version == "" {
		config.version = DefaultVersion
	}
	if config.service == "" {
		config.service = DefaultService
	}
}

func (f Factory) GetSupplier() string {
	return provider.Tencent
}

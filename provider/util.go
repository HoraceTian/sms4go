package provider

import "sms4go"

// ExtractConfigId 提取 ConfigId
func ExtractConfigId(config sms4go.SupplierConfig) string {
	configId := ""
	if config.GetConfigId() != "" {
		configId = config.GetConfigId()
	}
	return configId
}

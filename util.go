package sms4go

import (
	"strings"
)

// ExtractConfigId 提取 ConfigId
func ExtractConfigId(config SupplierConfig) string {
	configId := ""
	if config.GetConfigId() != "" {
		configId = config.GetConfigId()
	}
	return configId
}

func AddPrefixIfNot(str string, prefix string) string {
	if !strings.HasPrefix(str, prefix) {
		str = prefix + str
	}
	return str
}

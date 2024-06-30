package infra

import (
	"sms4go"
	"strings"
)

// ExtractConfigId 提取 ConfigId
func ExtractConfigId(config sms4go.SupplierConfig) string {
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

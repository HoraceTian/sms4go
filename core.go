package sms4go

// ISmsBlender 是短信 Blender 的接口，定义了短信服务的基本操作
type ISmsBlender interface {
	// GetConfigId 获取配置 ID
	GetConfigId() string

	// GetSupplier 获取供应商信息
	GetSupplier() string

	// SendMessage 发送短信并返回响应
	SendMessage(phone, message string) *SmsResponse

	// SendMessageWithParams 发送短信并使用给定参数返回响应
	SendMessageWithParams(phone string, params map[string]string) *SmsResponse

	// SendMessageWithParamsAndTemplate 发送短信并使用给定参数与指定模版返回响应
	SendMessageWithParamsAndTemplate(phone, templateId string, params map[string]string) *SmsResponse

	// SendMessageAsync 异步发送短信并返回响应
	SendMessageAsync(phone, message string, callback Callback) *SmsResponse

	// SendMessageWithParamsAsync 异步发送短信并使用给定参数返回响应
	SendMessageWithParamsAsync(phone string, params map[string]string, callback Callback) *SmsResponse

	// SendMessageWithTemplateAsync 异步发送具有模板 ID 和参数的短信并返回响应
	SendMessageWithTemplateAsync(phone, templateId string, params map[string]string, callback Callback) *SmsResponse
}

// BaseBlender 共有抽象
type BaseBlender struct {
	ConfigId string
	Config   SupplierConfig
}

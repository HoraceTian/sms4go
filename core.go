package sms4go

type ISmsBlender interface {
	GetConfigId() string

	GetSupplier() string

	SendMessage(phone, message string) SmsResponse

	SendMessageWithParams(phone string, params map[string]string) SmsResponse

	SendMessageWithTemplate(phone, templateId string, params map[string]string) SmsResponse

	SendMessageAsync(phone, message string, callback Callback) SmsResponse

	SendMessageWithParamsAsync(phone string, params map[string]string, callback Callback) SmsResponse

	SendMessageWithTemplateAsync(phone, templateId string, params map[string]string, callback Callback) SmsResponse
}

type BaseBlender struct {
	ConfigId string
	Config   SupplierConfig
}

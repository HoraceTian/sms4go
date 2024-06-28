package sms4go

type (
	sms4goOptions struct {
		smsConfig    *SmsConfig
		configMap    map[string]SupplierConfig
		interceptors []Interceptor
	}
)

type Option func(options *sms4goOptions)

func WithInterceptors(interceptors ...Interceptor) Option {
	return func(options *sms4goOptions) {
		options.interceptors = interceptors
	}
}

func WithConfigMap(configMap map[string]SupplierConfig) Option {
	return func(options *sms4goOptions) {
		options.configMap = configMap
	}
}

func WithSmsConfig(smsConfig *SmsConfig) Option {
	return func(options *sms4goOptions) {
		options.smsConfig = smsConfig
	}
}

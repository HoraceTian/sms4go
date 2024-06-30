package sms4go

import "github.com/panjf2000/ants/v2"

type (
	sms4goOptions struct {
		smsConfig    *SmsConfig
		configMap    map[string]SupplierConfig
		interceptors []Interceptor
		factories    []IProviderFactory
		routinePool  *ants.Pool
	}
)

type Option func(options *sms4goOptions)

func WithRoutinePool(pool *ants.Pool) Option {
	return func(options *sms4goOptions) {
		options.routinePool = pool
	}
}

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

func WithProviderFactories(factories ...IProviderFactory) Option {
	return func(options *sms4goOptions) {
		options.factories = factories
	}
}

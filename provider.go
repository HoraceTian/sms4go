package sms4go

// IProviderFactory 基础工厂规范
type IProviderFactory interface {
	CreateSms(config SupplierConfig) ISmsBlender
	GetSupplier() string
}

// Sms Provider 工厂 Holder
type providerFactoryHolder struct {
	factories map[string]IProviderFactory
}

// 注册 Sms Provider 工厂
func (ph *providerFactoryHolder) registerFactory(factory IProviderFactory) {
	if factory != nil {
		if ph.factories == nil {
			ph.factories = make(map[string]IProviderFactory)
		}
		ph.factories[factory.GetSupplier()] = factory
	}
}

// 批量注册 Sms Provider 工厂
func (ph *providerFactoryHolder) registerFactories(factories []IProviderFactory) {
	// 1. 兜底
	if len(factories) < 1 {
		return
	}

	// 2. 循环处理
	for _, factory := range factories {
		if factory == nil {
			continue
		}
		ph.registerFactory(factory)
	}
}

// 根据 Supplier 获取 Provider 工厂
func (ph *providerFactoryHolder) requireForSupplier(supplier string) IProviderFactory {
	return ph.factories[supplier]
}

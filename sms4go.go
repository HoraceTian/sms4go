package sms4go

// NewSmsClient 初始化 Sms 客户端
func NewSmsClient(options ...Option) Client {
	// 1. 初始化 option
	opts := &sms4goOptions{}

	// 2. 遍历用户选项
	for i := range options {
		options[i](opts)
	}

	// 3. 创建客户端
	client := &smsClient{}

	// 4. 初始化 Provider Holder
	initSmsProviderHolder(client, *opts)

	// 5. 初始化厂商 Sms 客户端
	for cfgKey := range opts.configMap {
		cfg := opts.configMap[cfgKey]
		providerFactory := client.providerHolder.factories[cfg.GetSupplier()]
		if providerFactory != nil {
			smsBlender := providerFactory.CreateSms(cfg)
			if client.blends == nil {
				client.blends = make(map[string]ISmsBlender)
			}
			client.blends[smsBlender.GetSupplier()] = smsBlender
		}
	}

	// 4. 返回客户端
	return client
}

// 初始化 Sms Factory
func initSmsProviderHolder(client *smsClient, opts sms4goOptions) {
	// 1. 提取 Sms Provider 列表
	factories := opts.factories

	// 2. 初始化 Holder
	providerHolder := &providerFactoryHolder{}
	providerHolder.registerFactories(factories)

	// 3. 初始化 Sms Provider Holder
	client.providerHolder = providerHolder
}

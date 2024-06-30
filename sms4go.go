package sms4go

import (
	"fmt"
	"github.com/panjf2000/ants/v2"
	"net/http"
	"time"
)

// NewSmsClient 初始化 Sms 客户端
func NewSmsClient(options ...Option) Client {
	// 1. 初始化 option
	opts := &sms4goOptions{}

	// 2. 遍历用户选项
	for _, option := range options {
		option(opts)
	}

	// 3. 创建客户端
	client := &smsClient{}

	// 4. 初始化 Provider Holder
	initSmsProviderHolder(client, opts)

	// 5. 初始化 httpClient
	client.httpClient = initSmsHttpClient(opts)

	// 6. 初始化协程池
	initSmsRoutinePool(client, opts)

	// 7. 初始化厂商 Sms 客户端
	initSmsBlenders(client, opts)

	// 8. 返回客户端
	return client
}

// 初始化 Sms Factory
func initSmsProviderHolder(client *smsClient, opts *sms4goOptions) {
	// 1. 提取 Sms Provider 列表
	factories := opts.factories

	// 2. 初始化 Holder
	providerHolder := &providerFactoryHolder{}
	providerHolder.registerFactories(factories)

	// 3. 初始化 Sms Provider Holder
	client.providerHolder = providerHolder
}

// 初始化 Sms HttpClient
func initSmsHttpClient(opts *sms4goOptions) *http.Client {
	// 1. 提取配置
	cfg := opts.smsConfig

	// 2. 创建客户端
	httpClient := &http.Client{}

	// 3. 处理参数（后续可扩展）
	if cfg.HttpTimeout != 0 {
		httpClient.Timeout = time.Duration(cfg.HttpTimeout) * time.Second // 确保 Timeout 是以秒为单位
	}
	return httpClient
}

// 初始化协程池
func initSmsRoutinePool(client *smsClient, opts *sms4goOptions) {
	// 1. 如果已有 routine pool，直接使用
	if opts.routinePool != nil {
		client.routinePool = opts.routinePool
		return
	}

	// 2. 使用默认值
	routineSize := DefaultRoutinePoolSize
	maxBlockingTasks := DefaultMaxBlockingTasks
	expiryDuration := DefaultExpiryDuration

	// 3. 如果提供了 poolConfig，使用其中的值覆盖默认值
	if poolConfig := opts.smsConfig.PoolConfig; poolConfig != nil {
		if poolConfig.ExpiryDuration != 0 {
			expiryDuration = poolConfig.ExpiryDuration
		}
		if poolConfig.RoutinePoolSize != 0 {
			routineSize = poolConfig.RoutinePoolSize
		}
		if poolConfig.MaxBlockingTasks != 0 {
			maxBlockingTasks = poolConfig.MaxBlockingTasks
		}
	}

	// 4. 初始化 routine pool
	pool, err := ants.NewPool(routineSize,
		ants.WithExpiryDuration(expiryDuration),
		ants.WithMaxBlockingTasks(maxBlockingTasks))

	if err != nil {
		panic(fmt.Sprintf("[sms4go] |- init sms routine pool error: %v", err))
	}

	client.routinePool = pool
}

// 初始化厂商 Sms 客户端
func initSmsBlenders(client *smsClient, opts *sms4goOptions) {
	if client.blends == nil {
		client.blends = make(map[string]ISmsBlender)
	}
	for cfgKey := range opts.configMap {
		cfg := opts.configMap[cfgKey]
		providerFactory := client.providerHolder.factories[cfg.GetSupplier()]
		if providerFactory != nil {
			smsBlender := providerFactory.CreateSms(cfg)
			smsBlender.SetHttpClient(client.httpClient)                 // 填充 HttpClient
			smsBlender.SetRoutinePool(client.routinePool)               // 填充协程池
			smsBlender.SetDelayQueue(NewDelayQueue(client.routinePool)) // 填充延时队列
			client.blends[smsBlender.GetSupplier()] = smsBlender
		}
	}
}

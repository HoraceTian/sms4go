package main

import (
	"sms4go"
	"sms4go/provider"
	"sms4go/provider/tencent"
)

func main() {
	// 1.  提供 SmsConfig
	smsConfig := &sms4go.SmsConfig{
		IsPrint: true,
	}

	// 2. 提供拦截器

	// 3. 提供厂商配置
	configMap := make(map[string]sms4go.SupplierConfig)
	configMap["tencent"] = &tencent.Config{
		BaseConfig: provider.BaseConfig{
			AccessKeyId:     "AKIDTowrXCgO8a1JkAe0CD6sR6nLkN4hHpIb",
			AccessKeySecret: "WtyESqNrX9g8WlwFTVPH8nGbpwwxAJLB",
			Signature:       "田浩然前端技术分享",
			SDKAppId:        "1400626380",
		},
	}

	// 4. 初始化厂商工厂

	client := sms4go.NewSmsClient(sms4go.WithConfigMap(configMap),
		sms4go.WithInterceptors(),
		sms4go.WithSmsConfig(smsConfig),
		sms4go.WithProviderFactories(tencent.NewFactory()),
	)

	client.CreateSmsBlender()

	blender := client.GetBySupplier("tencent")
	blender.SendMessage("13947856739", "aaaa")

}